package utils

import "time"

const dateDBLayout = "2006-01-02 15:04:05"

var DateUtils dateUtilsInterface = &dateUtilsStruct{}

type dateUtilsStruct struct{}

type dateUtilsInterface interface {
	GetNow() time.Time
	GetNowDBFormat() string
}

func (dateUtilsStruct) GetNow() time.Time { return time.Now().UTC() }

func (d dateUtilsStruct) GetNowDBFormat() string { return d.GetNow().Format(dateDBLayout) }
