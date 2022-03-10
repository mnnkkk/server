// Copyright (c) 2022 Trim21 <trim21.me@gmail.com>
//
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

package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/uber-go/tally/v4"

	"github.com/bangumi/server/web/handler"
)

func addHandle(
	scope tally.Scope,
	reg func(path string, handlers ...fiber.Handler) fiber.Router,
	path string,
	handler fiber.Handler,
) {
	reqCounter := scope.
		Tagged(map[string]string{"handler": utils.FunctionName(handler)}).
		Counter("request_count")

	reg(path, func(c *fiber.Ctx) error {
		reqCounter.Inc(1)

		return c.Next()
	}, handler)
}

// ResistRouter add all router and default 404 Handler to app.
func ResistRouter(app *fiber.App, h handler.Handler, scope tally.Scope) {
	app.Use(h.MiddlewareAccessUser())

	addHandle(scope, app.Get, "/v0/subjects/:id", h.GetSubject)
	addHandle(scope, app.Get, "/v0/subjects/:id/persons", h.GetSubjectRelatedPersons)
	addHandle(scope, app.Get, "/v0/subjects/:id/subjects", h.GetSubjectRelatedSubjects)
	addHandle(scope, app.Get, "/v0/subjects/:id/characters", h.GetSubjectRelatedCharacters)

	addHandle(scope, app.Get, "/v0/persons/:id", h.GetPerson)
	addHandle(scope, app.Get, "/v0/persons/:id/subjects", h.GetPersonRelatedSubjects)
	addHandle(scope, app.Get, "/v0/persons/:id/characters", h.GetPersonRelatedCharacters)

	addHandle(scope, app.Get, "/v0/characters/:id", h.GetCharacter)
	addHandle(scope, app.Get, "/v0/characters/:id/subjects", h.GetCharacterRelatedSubjects)
	addHandle(scope, app.Get, "/v0/characters/:id/persons", h.GetCharacterRelatedPersons)

	addHandle(scope, app.Get, "/v0/episodes/:id", h.GetEpisode)
	addHandle(scope, app.Get, "/v0/episodes", h.ListEpisode)

	addHandle(scope, app.Get, "/v0/me", h.GetCurrentUser)
	addHandle(scope, app.Get, "/v0/users/:username/collections", h.GetCollection)

	addHandle(scope, app.Get, "/v0/indices/:id", h.GetIndex)
	addHandle(scope, app.Get, "/v0/indices/:id/subjects", h.GetIndexSubjects)

	// default 404 Handler, all router should be added before this router
	app.Use(func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound).
			Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		return c.SendString(`{
  "title": "Not Found",
  "description": "The path you requested doesn't exist",
  "Detail": "This is default 404 response, if you see this response, please check your request path"
}`)
	})
}