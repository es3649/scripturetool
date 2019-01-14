package scriptures

// Tome is an alias for []string
// it holds a list of books withing a tome
type Tome []string

// OldTestament is a tome with the books of the Old Testament
var OldTestament = Tome([]string{"gen", "ex", "lev", "deut", "num", "josh", "judg", "ruth",
	"1-sam", "2-sam", "1-kgs", "2-kgs", "1-chr", "2-chr", "ezra", "neh", "esth", "job", "ps",
	"prov", "eccl", "song", "isa", "jer", "lam", "ezek", "dan", "hosea", "joel", "amos", "obad",
	"jonah", "micah", "nahum", "hab", "zeph", "hag", "zech", "mal"})

// NewTestament is a tome with the books of the New Testament
var NewTestament = Tome([]string{"matt", "mark", "luke", "john", "acts", "rom", "1-cor", "2-cor",
	"gal", "eph", "philip", "col", "1-thes", "2-thes", "1-tim", "2-tim",
	"titus", "philem", "heb", "james", "1-pet", "2-pet", "1-jn", "2-jn",
	"3-jn", "jude", "rev"})

// BookOfMormon is a tome with the books of the Book of Mormon
var BookOfMormon = Tome([]string{"1-ne", "2-ne", "jacob", "enos", "jarom", "omni", "w-of-m",
	"mosiah", "alma", "hel", "3-ne", "4-ne", "morm", "ether", "moro"})

// DoctrineAndCovenants is a tome with the book of the Doctrin and Covenants
var DoctrineAndCovenants = Tome([]string{"dc"})

// PearlOfGreatPrice is a tome with the books of the Pearl Of Great Price
var PearlOfGreatPrice = Tome([]string{"mos", "abr", "js-h", "js-m", "a-of-f"})
