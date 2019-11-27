// Copyright (c) 2019 Morpheus Data (https://www.morpheusdata.com), All rights reserved.
// morpheus source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package morpheus

import (
	"os"
)

// global stuff here

var (
	USE_FORCE = (os.Getenv("USE_FORCE") == "true")
)