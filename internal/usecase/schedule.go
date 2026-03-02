package usecase

import (
	"context"
	"fmt"
	"polytech_timetable/internal/entity"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/model/converter"
	"polytech_timetable/internal/repository"
	"polytech_timetable/pkg/customerrors"
	"polytech_timetable/pkg/parser"
	"polytech_timetable/pkg/polytech"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type ScheduleUseCase struct {
	log          *logrus.Logger
	client       *polytech.Client
	scheduleRepo *repository.ScheduleRepository
	userRepo     *repository.UserRepository
	teacherRepo  *repository.TeacherRepository
	bot          *tgbotapi.BotAPI
}

func NewScheduleUseCase(
	log *logrus.Logger,
	client *polytech.Client,
	scheduleRepo *repository.ScheduleRepository,
	userRepo *repository.UserRepository,
	teacherRepo *repository.TeacherRepository,
	bot *tgbotapi.BotAPI,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		log:          log,
		client:       client,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		teacherRepo:  teacherRepo,
		bot:          bot,
	}
}

type change struct {
	New entity.Lesson
	Old entity.Lesson
}

func (uc *ScheduleUseCase) SyncAllSchedules(ctx context.Context) error {
	uc.log.Info("Starting global schedule sync")

	oldLessonsMap := make(map[string]entity.Lesson)
	var currentLessons []entity.Lesson

	currentLessons, err := uc.scheduleRepo.GetFutureLessons(ctx)
	if err != nil {
		uc.log.Errorf("Failed to get future lessons: %v", err)
	}

	for _, l := range currentLessons {
		key := fmt.Sprintf("%s_%s_%d_%d_%d", l.Group, l.Date.Format("2006-01-02"), l.LessonNum, l.SubGroup, l.Zam)
		oldLessonsMap[key] = l
	}

	var allLessons []entity.Lesson

	var changes []change

	var uniqueTeachers = make(map[string]struct{})

	fileID := 1
	consecutive404 := 0
	for {
		if consecutive404 >= 3 {
			break
		}

		xmlBytes, err := uc.client.DownloadScheduleFile(ctx, fileID)
		if err != nil {
			consecutive404++
			fileID++
			continue
		}
		consecutive404 = 0

		rows, err := parser.ParseFile(xmlBytes)
		if err != nil {
			fileID++
			continue
		}

		for _, row := range rows {
			date, err := parser.ParseXMLDate(row.DateRaw)
			if err != nil || date.Before(time.Now().AddDate(0, 0, -30)) {
				continue
			}

			cleanTeacherName := parser.NormalizeTeacherName(row.Teacher)
			if cleanTeacherName != "" {
				uniqueTeachers[cleanTeacherName] = struct{}{}
			}

			newLesson := entity.Lesson{
				ExternalID: row.ID,
				Subject:    row.Subject,
				Teacher:    cleanTeacherName,
				Group:      row.Group,
				Date:       date,
				LessonNum:  row.LessonNum,
				SubGroup:   row.SubGroup,
				Zam:        row.Type,
			}

			key := fmt.Sprintf("%s_%s_%d_%d_%d", newLesson.Group, newLesson.Date.Format("2006-01-02"), newLesson.LessonNum, newLesson.SubGroup, newLesson.Zam)
			if old, ok := oldLessonsMap[key]; ok {
				if old.Subject != newLesson.Subject || old.Zam != newLesson.Zam || old.Teacher != newLesson.Teacher {
					uc.log.Infof("Change detected: key=%s old=[%s|%v|%s] new=[%s|%v|%s]", key, old.Subject, old.Zam, old.Teacher, newLesson.Subject, newLesson.Zam, newLesson.Teacher)
					changes = append(changes, change{New: newLesson, Old: old})
				}
			}

			allLessons = append(allLessons, newLesson)
		}
		fileID++
		time.Sleep(50 * time.Millisecond)
	}

	var newTeachers []entity.Teacher

	for name := range uniqueTeachers {
		newTeachers = append(newTeachers, entity.Teacher{Name: name})
	}

	if err := uc.teacherRepo.SaveBatchIgnore(ctx, newTeachers); err != nil {
		uc.log.Errorf("Failed to save batch teachers: %+v", err)
	}

	if len(changes) > 0 {
		go uc.sendNotifications(ctx, changes)
	}

	return uc.scheduleRepo.ReplaceAll(ctx, allLessons)
}

func (uc *ScheduleUseCase) sendNotifications(ctx context.Context, changes []change) {
	defer func() {
		if err := recover(); err != nil {
			uc.log.Errorf("Panic sending notifications: %+v", err)
		}
	}()
	byGroup := make(map[string][]change)
	for _, ch := range changes {
		byGroup[ch.New.Group] = append(byGroup[ch.New.Group], ch)
	}

	for group, groupChanges := range byGroup {
		users, err := uc.userRepo.FindByGroupWithNotifications(ctx, group)
		if err != nil {
			uc.log.WithError(err).WithField("group", group).Error("Failed to fetch users for notifications")
			continue
		}

		if len(users) == 0 {
			continue
		}
		var text strings.Builder
		text.WriteString("⚠️ <b>Изменения в расписании</b>\n")
		for _, ch := range groupChanges {
			dateStr := ch.New.Date.Format("02.01")
			fmt.Fprintf(&text, "\n📌 <b>%s • Пара %d</b>\n", dateStr, ch.New.LessonNum)

			if ch.Old.Subject != ch.New.Subject {
				fmt.Fprintf(&text, " • Предмет: <s>%s</s> ⟶ <b>%s</b>\n", ch.Old.Subject, ch.New.Subject)
			} else {
				fmt.Fprintf(&text, " • Предмет: %s\n", ch.New.Subject)
			}
			if ch.Old.Teacher != ch.New.Teacher {
				fmt.Fprintf(&text, " • Преподаватель: <s>%s</s> ⟶ <b>%s</b>\n", ch.Old.Teacher, ch.New.Teacher)
			}
		}

		for _, user := range users {
			msg := tgbotapi.NewMessage(user.TelegramID, text.String())
			msg.ParseMode = "HTML"

			if _, err := uc.bot.Send(msg); err != nil {
				uc.log.WithError(err).WithField("telegram_id", user.TelegramID).Error("Failed to send notification")
			}

			time.Sleep(40 * time.Millisecond)
		}
	}

	uc.log.WithField("changes", len(changes)).Info("Notifications sent")
}

func (uc *ScheduleUseCase) GetLessonsForGroup(ctx context.Context, group string) (*model.GetLessonsForGroupResponse, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -7).Truncate(24 * time.Hour)
	end := now.AddDate(0, 0, 14).Truncate(24 * time.Hour)

	lessons, err := uc.scheduleRepo.GetLessonsByGroup(ctx, group, start, end)
	if err != nil {
		uc.log.Errorf("Failed to get lessons for group: %+v", err)
		return nil, customerrors.ErrInternalServerError
	}

	response := converter.LessonsToResponse(lessons)

	return &model.GetLessonsForGroupResponse{
		Group:   group,
		Lessons: response,
	}, nil
}

func (uc *ScheduleUseCase) GetGroupsList(ctx context.Context) (*model.GetGroupsResponse, error) {
	groups, err := uc.scheduleRepo.GetAllGroups(ctx)
	if err != nil {
		uc.log.Errorf("Failed to get groups: %+v", err)
		return nil, customerrors.ErrInternalServerError
	}

	return &model.GetGroupsResponse{
		Groups: groups,
	}, nil
}

func (uc *ScheduleUseCase) GetTeacherSchedule(ctx context.Context, teacherID int64) (*model.GetTeacherScheduleResponse, error) {
	teacher, err := uc.teacherRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, customerrors.ErrNotFound
	}

	lessons, err := uc.scheduleRepo.SearchTeacher(ctx, teacher.Name)
	if err != nil {
		uc.log.Errorf("Failed to search teacher schedule: %v", err)
		return nil, customerrors.ErrInternalServerError
	}

	if lessons == nil {
		lessons = []entity.Lesson{}
	}

	response := converter.LessonsToResponse(lessons)

	return &model.GetTeacherScheduleResponse{
		Lessons: response,
	}, nil
}
