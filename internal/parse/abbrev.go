/**
 *
 */

package parse

import "strings"

var abbrevMap = map[string]string{
	"old testament": "[ot]", "genesis": "gen", "exodous": "ex", "leviticus": "lev", "deuteronomy": "deut",
	"numbers": "num", "joshua": "josh", "judges": "judg", "1 samuel": "1-sam",
	"2 samuel": "2-sam", "1 kings": "1-kin", "2 kings": "2-kin", "1 chronicles": "1-chr",
	"2 chronicles": "2-chr", "nehemiah": "neh", "esther": "esth",
	"psalms": "ps", "proverbs": "prov", "ecclesiastes": "eccl", "song of solomon": "song", "ss": "song", "isaiah": "isa",
	"jeremiah": "jer", "lamentations": "lam", "ezekiel": "ezek", "daniel": "dan", "obadiah": "obad",
	"habakkuk": "hab", "zephaniah": "zeph", "haggai": "hag", "zechariah": "zech", "malachi": "mal",

	"new testament": "[nt]", "matthew": "matt",
	"romans": "rom", "1 corinthians": "1-cor", "2 corinthians": "2-cor", "galations": "gal",
	"ephesians": "eph", "philippians": "philip", "colossians": "col", "1 thessalonians": "1-thes",
	"2 thessalonians": "2-thes", "1 timothy": "1-tim", "2 timothy": "2-tim",
	"philemon": "philem", "hebrews": "heb", "1 peter": "1-pet", "2 peter": "2-pet",
	"1 john": "1-jn", "2 john": "2-jn", "3 john": "3-jn", "revelation": "rev",

	"book of mormon": "[bom]", "1 nephi": "1-ne", "2 nephi": "2-ne", "words of mormon": "w-of-m", "msh": "mosiah",
	"helaman": "hel", "3 nephi": "3-ne", "4 nephi": "4-ne", "mormon": "morm", "eth": "ether",
	"moroni": "moro", "mni": "moro",

	"doctrine and covenants": "[dc]", "sections": "dc",

	"pearl of great price": "[pgp]", "abraham": "abr", "js-matthew": "js-m", "js-history": "js-h",
	"articles of faith": "a-of-f", "epistle dedicatory": "dedication", "bofm title page": "bofm-title",
	"title page of the book of mormon": "title-page", "introduction": "introduction",
	"testimony of the three witnesses": "three", "testimony of the eight witnesses": "eight",
	"testimony of the prophet joseph smith": "js", "a brief explanation of the bofm": "explanation",
	"official declarations": "od",
}

// PutAbbrevs replaces instances full book names with the supported abbreviated counterparts
func PutAbbrevs(in string) string {
	for long, abbrev := range abbrevMap {
		in = strings.ToLower(in)
		in = strings.Replace(in, long, abbrev, -1)
	}
	return in
}
