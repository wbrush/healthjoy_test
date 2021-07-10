//TODO: NOTICE: this is usually placed in datamodels package of go-common lib!
package datamodels

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/urlstruct"
)

type TemplateStatus string

const (
	TemplateStatusNew     TemplateStatus = "new"
	TemplateStatusActive  TemplateStatus = "active"
	TemplateStatusRemoved TemplateStatus = "removed"
)

func (status TemplateStatus) IsValid() bool {
	s := strings.ToLower(string(status))

	switch TemplateStatus(s) {
	case TemplateStatusNew,
		TemplateStatusActive,
		TemplateStatusRemoved:
		return true
	default:
		return false
	}
}

type (
	Template struct {
		tableName struct{} `sql:"?SHARD.templates"`

		Id           int64          `json:"templateId" pg:"id"`
		Name         string         `json:"name" pg:"name"`
		Status       TemplateStatus `json:"status" pg:"status"`
		Description  string         `json:"description" pg:"description"`
		TemplateSelf string         `json:"self" pg:"template_self"`

		CreatedAt  time.Time  `json:"createdAt" pg:"created_at,default:now()"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" pg:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" pg:"archived_at"`
	}

	TemplateUpdate struct {
		Status      *TemplateStatus `json:"status,omitempty"`
		Description *string         `json:"description,omitempty"`
	}

	TemplateFilter struct {
		tableName struct{} `urlstruct:"template"`

		urlstruct.Pager

		Name        string
		Status      TemplateStatus
		Description string
	}
)

func (t Template) String() string {
	return fmt.Sprintf("Template<%d %s>", t.Id, t.Name)
}

func (t Template) Validate() error {
	if !t.Status.IsValid() {
		return fmt.Errorf("bad status %s provided", t.Status)
	}

	// TODO: add here some expected validation rules

	return nil
}
