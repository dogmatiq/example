package domain_test

import "time"

// The current expected daily debit limit.
const expectedDailyDebitLimit = 900000

var dateTimeNow = time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC)
var businessDateToday = dateTimeNow.Format("2006-01-02")
var startOfBusinessDateTimeToday, _ = time.Parse("2006-01-02", businessDateToday)
var businessDateTomorrow = dateTimeNow.Add(time.Hour * 24).Format("2006-01-02")
var startOfBusinessDateTimeTomorrow, _ = time.Parse("2006-01-02", businessDateTomorrow)
