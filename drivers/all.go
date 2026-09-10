package drivers

import (
	_ "github.com/OpenListTeam/OpenList/v4/drivers/local"
	// drivers/template is intentionally NOT registered: it is a copy-me skeleton
	// kept for developers writing new drivers, and its List/Link return
	// errs.NotImplement. Registering it would only add a non-functional
	// "Template" entry to the driver picker.
)
