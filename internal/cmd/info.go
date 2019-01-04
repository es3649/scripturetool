// ScriptureTool
// es3649
// help.go - handling the help command

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays help for argument formatting and abbreviations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Verse Referencing:
reference verses within quotes, one reference per string
a reference should be formatted in any of the following methods
  (or some combination of them):
Reference:
  "Book chap:vs"         or 
  "Book chap"            (prints whole chapter) or
  "Book chap-range"      (doesn't accept verse listing)
Book Name:
  "1 Corinthians"        (by name)
  "1-cor"                (abbreviated: see below)
  "{Matt,Mark,Luke}"     (listed: abbreviated or not)--Not yet implemented
  "[OT]"                 (by class)
Chapter Number:
  "2"                    (number)
  "3-5"                  (range: don't list verses after)
  "20,84,100"            (listed)
  "*"                    (all chapters, same as not listing verse)--Not yet implemented
Verse Numbers:
  "2"                    (number)
  "3-5"                  (range)
  "77,79"                (listed)
  Using * as a chapter will indicate all chapters in the book: John *:12 will indicate
    the 12th verse of every chapter of John (having at least 12 verses)
  Using * as a verse is equivalent to listing the chapter number alone: John 17:* is
    equivalent to John 17

All Scripture (non-verse content excluded)  *  
OLD TESTAMENT   [ot]    NEW TESTAMENT   [nt]   
Genesis         gen     Matthew         matt   
Exodous         ex      Mark            mark   
Leviticus       lev     Luke            luke   
Deuteronomy     deut    John            john   
Numbers         num     Acts            acts   
Joshua          jos     Romans          rom    
Judges          judg    1 Corinthians   1-cor  
Ruth            ruth    2 Corinthians   2-cor  
1 Samuel        1-sam   Galatians       gal    
2 Samuel        2-sam   Ephesians       eph    
1 Kings         1-kin   Philippians     philip 
2 Kings         2-kin   Colossians      col    
1 Chronicles    1-chr   1 Thessalonians 1-thes 
2 Chronicles    2-chr   2 Thessalonians 2-thes 
Ezra            ezra    1 Timothy       1-tim  
Nehemiah        neh     2 Timothy       2-tim  
Esther          esth    Titus           titus  
Job             job     Philemon        philem 
Psalms          ps      Hebrews         heb    
Proverbs        prov    James           jas    
Ecclesiastes    eccl    1 Peter         1-pet  
Song of Solomon song    2 Peter         2-pet  
Isaiah          isa     1 John          1-jn   
Jeremiah        jer     2 John          2-jn   
Lamentations    lam     3 John          3-jn   
Ezekiel         ezek    Jude            jude   
Daniel          dan     Revelation      rev    
Hosea           hosea                          
Joel            joel    BOOK OF MORMON  [bom]  
Amos            amos    1 Nephi         1-ne   
Obadiah         obad    2 Nephi         2-ne   
Jonah           jonah   Jacob           jacob  
Micah           micah   Enos            enos   
Nahum           nahum   Jarom           jarom  
Habakkuk        hab     Omni            omni   
Zephaniah       zeph    Words of Mormon w-of-m 
Haggai          hag     Mosiah          mosiah 
Zechariah       zech    Alma            alma   
Malachi         mal     Helaman         hel    
                        3 Nephi         3-ne   
PEARL OF        [pgp]   4 Nephi         4-ne   
  GREAT PRICE           Mormon          morm   
Moses           moses   Ether           ether  
Abraham         abr     Moroni          moro   
JS-Matthew      js-m                           
JS-History      js-h    DOCTRINE AND    [dc]   
Articles        a-of-f    COVENANTS            
  of Faith              Sections        dc     

NON-VERSE CONTENT:--Not Yet Implemented
                  BIBLE                        
Epistle Dedicatory                      dedication
              BOOK OF MORMON                   
BofM Title Page                         bofm-title
Title Page of the Book of Mormon        title-page
Introduction                            introduction
Testimony of the Three Witnesses        three
Testimony of the Eight Witnesses        eight
Testimony of the ProphetJoseph Smith    js
A Brief Explanation of the BofM         explanation
         DOCTRINE AND COVENANTS                
Official Declarations                   od

LANGUAGE OPTIONS
Libraries must be downloaded for each language,
and downloading is currently only supported for English

Language    Code
English     eng
Spanish     spa
Portuguese  por`)
	},
}
