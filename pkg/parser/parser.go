package parser

import (
	"encoding/xml"
	"fmt"
	"time"
)

type ScheduleRoot struct {
	XMLName xml.Name      `xml:"dataroot"`
	Rows    []ScheduleRow `xml:"My"`
}

type ScheduleRow struct {
	ID        string `xml:"ID"`
	DateRaw   string `xml:"DAT"`
	LessonNum int    `xml:"UR"`
	Teacher   string `xml:"FAMIO"`
	Subject   string `xml:"SPPRED.NAIM"`
	Group     string `xml:"SPGRUP.NAIM"`
	SubGroup  int    `xml:"IDGG"`
	Type      int    `xml:"ZAM"`
}

func ParseFile(data []byte) ([]ScheduleRow, error) {
	var root ScheduleRoot
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling xml: %w", err)
	}
	return root.Rows, nil
}

func ParseXMLDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': %w", dateStr, err)
	}
	return t, nil
}
