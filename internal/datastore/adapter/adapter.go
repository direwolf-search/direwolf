package adapter

import (
	"time"

	"direwolf/internal/datastore/models"
	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	"direwolf/internal/domain/model/starter"
)

type Adapter struct{}

func (a *Adapter) ConvertHostToModel(hostEntity *host.Host) *models.Host {
	return &models.Host{
		H1:        hostEntity.H1,
		Title:     hostEntity.Title,
		Hash:      hostEntity.Hash,
		Text:      hostEntity.Text,
		URL:       hostEntity.URL,
		Status:    hostEntity.Status,
		CreatedAt: time.Now(),
	}
}

func (a *Adapter) ConvertLinkToModel(linkEntity *link.Link) *models.Link {
	return &models.Link{
		From:    linkEntity.From,
		Body:    linkEntity.Body,
		Snippet: linkEntity.Snippet,
		IsV3:    linkEntity.IsV3,
	}
}

func (a *Adapter) ConvertStarterToModel(starterEntity *starter.Starter) *models.Starter {
	return &models.Starter{
		Body:      starterEntity.Body,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}
