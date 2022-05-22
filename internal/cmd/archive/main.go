// SPDX-License-Identifier: AGPL-3.0-only
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>

//nolint:govet
package main

import (
	"archive/zip"
	"context"
	"io"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/bangumi/server/internal/config"
	"github.com/bangumi/server/internal/dal"
	"github.com/bangumi/server/internal/dal/dao"
	"github.com/bangumi/server/internal/dal/query"
	"github.com/bangumi/server/internal/logger"
	"github.com/bangumi/server/internal/metrics"
	"github.com/bangumi/server/internal/model"
)

const defaultStep = 50

func main() {
	if err := mysql.SetLogger(logger.Std()); err != nil {
		logger.Panic("can't replace mysql driver's errLog", zap.Error(err))
	}

	out := pflag.String("out", "archive.zip", "zip file output location")
	if out == nil {
		*out = "archive.zip"
	}

	start(*out)
}

var ctx = context.Background() //nolint:gochecknoglobals

var maxSubjectID model.SubjectIDType     //nolint:gochecknoglobals
var maxCharacterID model.CharacterIDType //nolint:gochecknoglobals
var maxPersonID model.PersonIDType       //nolint:gochecknoglobals

func start(out string) {
	var q *query.Query
	err := fx.New(
		logger.FxLogger(),
		fx.Provide(
			dal.NewConnectionPool, dal.NewDB,

			config.NewAppConfig, logger.Copy, metrics.NewScope,

			query.Use,
		),

		fx.Populate(&q),
	).Err()

	if err != nil {
		logger.Err(err, "failed to fill deps")
	}

	getMaxID(q)

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}

	z := zip.NewWriter(f)

	for _, s := range []struct {
		FileName string
		Fn       func(q *query.Query, w io.Writer)
	}{
		{FileName: "subject.jsonlines", Fn: exportSubjects},
		{FileName: "person.jsonlines", Fn: exportPersons},
		{FileName: "character.jsonlines", Fn: exportCharacters},
		{FileName: "episode.jsonlines", Fn: exportEpisodes},
		{FileName: "subject-relations.jsonlines", Fn: exportSubjectRelations},
		{FileName: "subject-persons.jsonlines", Fn: exportSubjectPersonRelations},
		{FileName: "subject-characters.jsonlines", Fn: exportSubjectCharacterRelations},
		{FileName: "person-characters.jsonlines", Fn: exportPersonCharacterRelations},
	} {
		w, err := z.Create(s.FileName)
		if err != nil {
			panic(err)
		}

		s.Fn(q, w)
	}

	if err = z.Close(); err != nil {
		panic(err)
	}

	if err = f.Close(); err != nil {
		panic(err)
	}
}

func getMaxID(q *query.Query) {
	lastSubject, err := q.WithContext(ctx).Subject.Order(q.Subject.ID.Desc()).First()
	if err != nil {
		panic(err)
	}
	maxSubjectID = lastSubject.ID

	lastCharacter, err := q.WithContext(ctx).Character.Order(q.Character.ID.Desc()).First()
	if err != nil {
		panic(err)
	}
	maxCharacterID = lastCharacter.ID

	lastPerson, err := q.WithContext(ctx).Person.Order(q.Person.ID.Desc()).First()
	if err != nil {
		panic(err)
	}
	maxPersonID = lastPerson.ID
}

type Subject struct {
	ID       model.SubjectIDType `json:"id"`
	Type     model.SubjectType   `json:"type"`
	Name     string              `json:"name"`
	NameCN   string              `json:"name_cn"`
	Infobox  string              `json:"infobox"`
	Platform uint16              `json:"platform"`
	Summary  string              `json:"summary"`
	Nsfw     bool                `json:"nsfw"`
}

func exportSubjects(q *query.Query, w io.Writer) {
	for i := model.SubjectIDType(0); i < maxSubjectID; i += defaultStep {
		subjects, err := q.WithContext(ctx).Subject.
			Where(q.Subject.ID.Gt(i), q.Subject.ID.Lte(i+defaultStep), q.Subject.Ban.Eq(0)).Find()
		if err != nil {
			panic(err)
		}

		for _, subject := range subjects {
			encode(w, Subject{
				ID:       subject.ID,
				Type:     subject.TypeID,
				Name:     subject.Name,
				NameCN:   subject.NameCN,
				Infobox:  subject.Infobox,
				Platform: subject.Platform,
				Summary:  subject.Summary,
				Nsfw:     subject.Nsfw,
			})
		}
	}
}

type Person struct {
	ID      model.PersonIDType `json:"id"`
	Name    string             `json:"name"`
	Type    uint8              `json:"type"`
	Career  []string           `json:"career"`
	Infobox string             `json:"infobox"`
	Summary string             `json:"summary"`
}

func exportPersons(q *query.Query, w io.Writer) {
	for i := model.PersonIDType(0); i < maxPersonID; i += defaultStep {
		persons, err := q.WithContext(context.Background()).Person.
			Where(q.Person.ID.Gt(i), q.Person.ID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, p := range persons {
			encode(w, Person{
				ID:      p.ID,
				Name:    p.Name,
				Type:    p.Type,
				Career:  careers(p),
				Infobox: p.Infobox,
				Summary: p.Summary,
			})
		}
	}
}

func careers(p *dao.Person) []string {
	s := make([]string, 0, 7)

	if p.Writer {
		s = append(s, "writer")
	}

	if p.Producer {
		s = append(s, "producer")
	}

	if p.Mangaka {
		s = append(s, "mangaka")
	}

	if p.Artist {
		s = append(s, "artist")
	}

	if p.Seiyu {
		s = append(s, "seiyu")
	}

	if p.Writer {
		s = append(s, "writer")
	}

	if p.Illustrator {
		s = append(s, "illustrator")
	}

	if p.Actor {
		s = append(s, "actor")
	}

	return s
}

type Character struct {
	ID      model.CharacterIDType `json:"id"`
	Role    uint8                 `json:"name"`
	Infobox string                `json:"infobox"`
	Summary string                `json:"summary"`
}

func exportCharacters(q *query.Query, w io.Writer) {
	for i := model.CharacterIDType(0); i < maxCharacterID; i += defaultStep {
		characters, err := q.WithContext(context.Background()).Character.
			Where(q.Character.ID.Gt(i), q.Character.ID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, c := range characters {
			encode(w, Character{
				ID:      c.ID,
				Role:    c.Role,
				Infobox: c.Infobox,
				Summary: c.Summary,
			})
		}
	}
}

type Episode struct {
	ID          model.EpisodeIDType `json:"id"`
	Name        string              `json:"name"`
	NameCn      string              `json:"name_cn"`
	Description string              `json:"description"`
	AirDate     string              `json:"airdate"`
	Disc        uint8               `json:"disc"`
	SubjectID   model.SubjectIDType `json:"subject_id"`
	Sort        float32             `json:"sort"`
	Type        int8                `json:"type"`
}

func exportEpisodes(q *query.Query, w io.Writer) {
	lastEpisode, err := q.WithContext(ctx).Episode.Order(q.Episode.ID.Desc()).First()
	if err != nil {
		panic(err)
	}
	for i := model.EpisodeIDType(0); i < lastEpisode.ID; i += defaultStep {
		episodes, err := q.WithContext(context.Background()).Episode.
			Where(q.Episode.ID.Gt(i), q.Episode.ID.Lte(i+defaultStep), q.Episode.Ban.Eq(0)).Find()
		if err != nil {
			panic(err)
		}

		for _, e := range episodes {
			encode(w, Episode{
				ID:          e.ID,
				Name:        e.Name,
				NameCn:      e.NameCn,
				Sort:        e.Sort,
				SubjectID:   e.SubjectID,
				Description: e.Desc,
				Type:        e.Type,
				AirDate:     e.Airdate,
				Disc:        e.Disc,
			})
		}
	}
}

type SubjectRelation struct {
	SubjectID        model.SubjectIDType `json:"subject_id"`
	RelationType     uint16              `json:"relation_type"`
	RelatedSubjectID model.SubjectIDType `json:"related_subject_id"`
	Order            uint8               `json:"order"`
}

func exportSubjectRelations(q *query.Query, w io.Writer) {
	for i := model.SubjectIDType(0); i < maxSubjectID; i += defaultStep {
		relations, err := q.WithContext(context.Background()).SubjectRelation.
			Order(q.SubjectRelation.SubjectID, q.SubjectRelation.SubjectID).
			Where(q.SubjectRelation.SubjectID.Gt(i), q.SubjectRelation.SubjectID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, rel := range relations {
			encode(w, SubjectRelation{
				SubjectID:        rel.SubjectID,
				RelationType:     rel.RelationType,
				RelatedSubjectID: rel.RelatedSubjectID,
				Order:            rel.Order,
			})
		}
	}
}

type SubjectPerson struct {
	PersonID  model.PersonIDType  `json:"person_id"`
	SubjectID model.SubjectIDType `json:"subject_id"`
	Position  uint16              `json:"position"`
}

func exportSubjectPersonRelations(q *query.Query, w io.Writer) {
	for i := model.SubjectIDType(0); i < maxSubjectID; i += defaultStep {
		relations, err := q.WithContext(context.Background()).PersonSubjects.
			Order(q.PersonSubjects.SubjectID, q.PersonSubjects.PersonID).
			Where(q.PersonSubjects.SubjectID.Gt(i), q.PersonSubjects.SubjectID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, rel := range relations {
			encode(w, SubjectPerson{
				PersonID:  rel.PersonID,
				SubjectID: rel.SubjectID,
				Position:  rel.PrsnPosition,
			})
		}
	}
}

type SubjectCharacter struct {
	CharacterID model.CharacterIDType `json:"character_id"`
	SubjectID   model.SubjectIDType   `json:"subject_id"`
	Type        uint8                 `json:"type"`
	Order       uint8                 `json:"order"`
}

func exportSubjectCharacterRelations(q *query.Query, w io.Writer) {
	for i := model.SubjectIDType(0); i < maxSubjectID; i += defaultStep {
		relations, err := q.WithContext(context.Background()).CharacterSubjects.
			Order(q.CharacterSubjects.SubjectID, q.CharacterSubjects.CrtOrder).
			Where(q.CharacterSubjects.SubjectID.Gt(i), q.CharacterSubjects.SubjectID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, rel := range relations {
			encode(w, SubjectCharacter{
				CharacterID: rel.CharacterID,
				SubjectID:   rel.SubjectID,
				Type:        rel.CrtType,
				Order:       rel.CrtOrder,
			})
		}
	}
}

type PersonCharacter struct {
	PersonID    model.PersonIDType    `json:"person_id"`
	SubjectID   model.SubjectIDType   `json:"subject_id"`
	CharacterID model.CharacterIDType `json:"character_id"`
	Summary     string                `json:"summary"`
}

func exportPersonCharacterRelations(q *query.Query, w io.Writer) {
	for i := model.PersonIDType(0); i < maxPersonID; i += defaultStep {
		relations, err := q.WithContext(context.Background()).Cast.
			Order(q.Cast.PersonID, q.Cast.CharacterID).
			Where(q.Cast.PersonID.Gt(i), q.Cast.PersonID.Lte(i+defaultStep)).Find()
		if err != nil {
			panic(err)
		}

		for _, rel := range relations {
			encode(w, PersonCharacter{
				PersonID:    rel.PersonID,
				SubjectID:   rel.SubjectID,
				CharacterID: rel.CharacterID,
				Summary:     rel.Summary,
			})
		}
	}
}

func encode(w io.Writer, object interface{}) {
	if err := json.NewEncoder(w).Encode(object); err != nil {
		panic(err)
	}
}
