/**
 *
 */

package parse

import "strings"

var abbrevMap = map[string]string{"OLD TESTAMENT": "[OT]", "Genesis": "gen", "Exodous": "ex", "Leviticus": "lev", "Deuteronomy": "deut",
	"Numbers": "num", "Joshua": "jos", "Judges": "judg", "Ruth": "ruth", "1 Samuel": "1-sam",
	"2 Samuel": "2-sam", "1 Kings": "1-kin", "2 Kings": "2-kin", "1 Chronicles": "1-chr",
	"2 Chronicles": "2-chr", "Ezra": "ezra", "Nehemiah": "neh", "Esther": "esth", "Job": "job",
	"Psalms": "ps", "Proverbs": "prov", "Ecclesiastes": "eccl", "Song of Solomon": "song", "Isaiah": "isa",
	"Jeremiah": "jer", "Lamentations": "lam", "Ezekiel": "ezek", "Daniel": "dan", "Hosea": "hosea",
	"Joel": "joel", "Amos": "amos", "Obadiah": "obad", "Jonah": "jonah", "Micah": "micah", "Nahum": "nahum",
	"Habakkuk": "hab", "Zephaniah": "zeph", "Haggai": "hag", "Zechariah": "zech", "Malachi": "mal",
	"NEW TESTAMENT": "[NT]", "Matthew": "matt", "Mark": "mark", "Luke": "luke", "John": "john",
	"Acts": "acts", "Romans": "rom", "1 Corinthians": "1-cor", "2 Corinthians": "2-cor", "Galations": "gal",
	"Ephesians": "eph", "Philippians": "philip", "Colossians": "col", "1 Thessalonians": "1-thes",
	"2 Thessalonians": "2-thes", "1 Timothy": "1-tim", "2 Timothy": "2-tim", "Titus": "titus",
	"Philemon": "philem", "Hebrews": "heb", "James": "jas", "1 Peter": "1-pet", "2 Peter": "2-pet",
	"1 John": "1-jn", "2 John": "2-jn", "3 John": "3-jn", "Jude": "jude", "Revelation": "rev",
	"BOOK OF MORMON": "[BOM]", "1 Nephi": "1-ne", "2 Nephi": "2-ne", "Jacob": "jacob", "Enos": "enos",
	"Jarom": "jarom", "Omni": "omni", "Words of Mormon": "w-of-m", "Mosiah": "mosiah", "Alma": "alma",
	"Helaman": "hel", "3 Nephi": "3-ne", "4 Nephi": "4-ne", "Mormon": "morm", "Ether": "ether",
	"Moroni": "moro", "DOCTRINE AND COVENANTS": "[DC]", "Sections": "dc", "PEARL OF GREAT PRICE": "[PGP]",
	"Moses": "moses", "Abraham": "abr", "JS-Matthew": "js-m", "JS-History": "js-h",
	"Articles of Faith": "a-of-f", "Epistle Dedicatory": "dedication", "BofM Title Page": "bofm-title",
	"Title Page of the Book of Mormon": "title-page", "Introduction": "introduction",
	"Testimony of the Three Witnesses": "three", "Testimony of the Eight Witnesses": "eight",
	"Testimony of the ProphetJoseph Smith": "js", "A Brief Explanation of the BofM": "explanation",
	"Official Declarations": "od"}

// PutAbbrevs replaces instances full book names with the supported abbreviated counterparts
func PutAbbrevs(in string) string {
	for long, abbrev := range abbrevMap {
		in = strings.Replace(in, long, abbrev, -1)
	}
	return in
}
