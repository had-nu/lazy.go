package scaffold

import (
	"fmt"
	"time"

	"github.com/had-nu/lazy.go/pkg/config"
)

// GenerateLicense returns the full license text for the given type.
func GenerateLicense(t config.LicenseType, author string, year int) string {
	if year == 0 {
		year = time.Now().Year()
	}
	switch t {
	case config.LicenseMIT:
		return mitLicense(author, year)
	case config.LicenseApache2:
		return apache2License(author, year)
	case config.LicenseGPL3:
		return gpl3License()
	case config.LicenseProprietary:
		return proprietaryLicense(author, year)
	default:
		return fmt.Sprintf("# LICENSE\n\nCopyright (c) %d %s. All rights reserved.\n", year, author)
	}
}

func mitLicense(author string, year int) string {
	return fmt.Sprintf(`MIT License

Copyright (c) %d %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year, author)
}

func apache2License(author string, year int) string {
	return fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %d %s

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`, year, author)
}

func gpl3License() string {
	return `GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
Everyone is permitted to copy and distribute verbatim copies
of this license document, but changing it is not allowed.

                            PREAMBLE

The GNU General Public License is a free, copyleft license for
software and other kinds of works.

[...Full GPL-3.0 text truncated for brevity — see https://www.gnu.org/licenses/gpl-3.0.txt]

END OF TERMS AND CONDITIONS
`
}

func proprietaryLicense(author string, year int) string {
	return fmt.Sprintf(`PROPRIETARY LICENSE

Copyright (c) %d %s. All rights reserved.

This software and its source code are proprietary and confidential.
Unauthorized copying, distribution, modification, or use of this software,
in whole or in part, is strictly prohibited without prior written permission
from the copyright owner.

For licensing inquiries, contact the author directly.
`, year, author)
}
