// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package echo

import (
	"encoding/json"

	"github.com/schrodit/landscaper-dashboard/server/routes"
)

func Route(data []byte) (*routes.Response, error) {
	return routes.OkResponse(json.RawMessage(data)), nil
}
