package main

import (
	"context"
	"encoding/json"
	"strconv"
)

type CourseSchedule struct {
	Module string `json:"module,omitempty"`
	Day    Day    `json:"day,omitempty"`
	Period int    `json:"period,omitempty"`
	Room   string `json:"room,omitempty"`
}

type Course struct {
	ID        string
	Code      string           `json:"code,omitempty"`
	Name      string           `json:"name,omitempty"`
	Methods   []string         `json:"methods,omitempty"`
	Schedules []CourseSchedule `json:"schedules,omitempty"`
}

func GetCourses(ctx context.Context, year int) ([]Course, error) {
	var data []struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
		Course struct {
			ID                string           `json:"id"`
			Year              int              `json:"year"`
			Code              string           `json:"code"`
			Name              string           `json:"name"`
			Instructor        string           `json:"instructor"`
			Credit            json.Number      `json:"credit"`
			Overview          string           `json:"overview"`
			Remarks           string           `json:"remarks"`
			RecommendedGrades []int            `json:"recommendedGrades"`
			Methods           []string         `json:"methods"`
			Schedules         []CourseSchedule `json:"schedules"`
			Isannual          bool             `json:"isAnnual"`
			HasParseError     bool             `json:"hasParseError"`
		} `json:"course"`
		Year       int               `json:"year"`
		Name       *string           `json:"name"`
		Instructor *string           `json:"instructor"`
		Credit     *json.Number      `json:"credit"`
		Methods    *[]string         `json:"methods"`
		Schedules  *[]CourseSchedule `json:"schedules"`
		Memo       string            `json:"memo"`
		Attendance int               `json:"attendance"`
		Absence    int               `json:"absence"`
		Late       int               `json:"late"`
		Tags       []struct {
			ID string `json:"id"`
		} `json:"tags"`
	}
	err := GetAPI(ctx, "/registered-courses?year="+strconv.Itoa(year), &data)
	if err != nil {
		return nil, err
	}

	result := make([]Course, len(data))
	for i, v := range data {
		c := Course{
			ID:        v.ID,
			Code:      v.Course.Code,
			Name:      v.Course.Name,
			Methods:   v.Course.Methods,
			Schedules: v.Course.Schedules,
		}
		if v.Name != nil {
			c.Name = *v.Name
		}
		if v.Methods != nil {
			c.Methods = *v.Methods
		}
		if v.Schedules != nil {
			c.Schedules = *v.Schedules
		}
		result[i] = c
	}
	return result, nil
}
