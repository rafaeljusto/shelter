package model

// File generated using the language.py script, that is responsable for parsing the IANA
// Language Subtag Registry file obtained from
// http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry
import (
	"strings"
)

// List of possible language types
const (
	LanguageTypeAA     LanguageType = "aa"       // Afar
	LanguageTypeAB     LanguageType = "ab"       // Abkhazian
	LanguageTypeAE     LanguageType = "ae"       // Avestan
	LanguageTypeAF     LanguageType = "af"       // Afrikaans
	LanguageTypeAK     LanguageType = "ak"       // Akan
	LanguageTypeAM     LanguageType = "am"       // Amharic
	LanguageTypeAN     LanguageType = "an"       // Aragonese
	LanguageTypeAR     LanguageType = "ar"       // Arabic
	LanguageTypeAS     LanguageType = "as"       // Assamese
	LanguageTypeAV     LanguageType = "av"       // Avaric
	LanguageTypeAY     LanguageType = "ay"       // Aymara
	LanguageTypeAZ     LanguageType = "az"       // Azerbaijani
	LanguageTypeBA     LanguageType = "ba"       // Bashkir
	LanguageTypeBE     LanguageType = "be"       // Belarusian
	LanguageTypeBG     LanguageType = "bg"       // Bulgarian
	LanguageTypeBH     LanguageType = "bh"       // Bihari languages
	LanguageTypeBI     LanguageType = "bi"       // Bislama
	LanguageTypeBM     LanguageType = "bm"       // Bambara
	LanguageTypeBN     LanguageType = "bn"       // Bengali
	LanguageTypeBO     LanguageType = "bo"       // Tibetan
	LanguageTypeBR     LanguageType = "br"       // Breton
	LanguageTypeBS     LanguageType = "bs"       // Bosnian
	LanguageTypeCA     LanguageType = "ca"       // Catalan and Valencian
	LanguageTypeCE     LanguageType = "ce"       // Chechen
	LanguageTypeCH     LanguageType = "ch"       // Chamorro
	LanguageTypeCO     LanguageType = "co"       // Corsican
	LanguageTypeCR     LanguageType = "cr"       // Cree
	LanguageTypeCS     LanguageType = "cs"       // Czech
	LanguageTypeCU     LanguageType = "cu"       // Church Slavic and Church Slavonic and Old Bulgarian and Old Church Slavonic and Old Slavonic
	LanguageTypeCV     LanguageType = "cv"       // Chuvash
	LanguageTypeCY     LanguageType = "cy"       // Welsh
	LanguageTypeDA     LanguageType = "da"       // Danish
	LanguageTypeDE     LanguageType = "de"       // German
	LanguageTypeDV     LanguageType = "dv"       // Dhivehi and Divehi and Maldivian
	LanguageTypeDZ     LanguageType = "dz"       // Dzongkha
	LanguageTypeEE     LanguageType = "ee"       // Ewe
	LanguageTypeEL     LanguageType = "el"       // Modern Greek (1453-)
	LanguageTypeEN     LanguageType = "en"       // English
	LanguageTypeEO     LanguageType = "eo"       // Esperanto
	LanguageTypeES     LanguageType = "es"       // Spanish and Castilian
	LanguageTypeET     LanguageType = "et"       // Estonian
	LanguageTypeEU     LanguageType = "eu"       // Basque
	LanguageTypeFA     LanguageType = "fa"       // Persian
	LanguageTypeFF     LanguageType = "ff"       // Fulah
	LanguageTypeFI     LanguageType = "fi"       // Finnish
	LanguageTypeFJ     LanguageType = "fj"       // Fijian
	LanguageTypeFO     LanguageType = "fo"       // Faroese
	LanguageTypeFR     LanguageType = "fr"       // French
	LanguageTypeFY     LanguageType = "fy"       // Western Frisian
	LanguageTypeGA     LanguageType = "ga"       // Irish
	LanguageTypeGD     LanguageType = "gd"       // Scottish Gaelic and Gaelic
	LanguageTypeGL     LanguageType = "gl"       // Galician
	LanguageTypeGN     LanguageType = "gn"       // Guarani
	LanguageTypeGU     LanguageType = "gu"       // Gujarati
	LanguageTypeGV     LanguageType = "gv"       // Manx
	LanguageTypeHA     LanguageType = "ha"       // Hausa
	LanguageTypeHE     LanguageType = "he"       // Hebrew
	LanguageTypeHI     LanguageType = "hi"       // Hindi
	LanguageTypeHO     LanguageType = "ho"       // Hiri Motu
	LanguageTypeHR     LanguageType = "hr"       // Croatian
	LanguageTypeHT     LanguageType = "ht"       // Haitian and Haitian Creole
	LanguageTypeHU     LanguageType = "hu"       // Hungarian
	LanguageTypeHY     LanguageType = "hy"       // Armenian
	LanguageTypeHZ     LanguageType = "hz"       // Herero
	LanguageTypeIA     LanguageType = "ia"       // Interlingua (International Auxiliary Language
	LanguageTypeID     LanguageType = "id"       // Indonesian
	LanguageTypeIE     LanguageType = "ie"       // Interlingue and Occidental
	LanguageTypeIG     LanguageType = "ig"       // Igbo
	LanguageTypeII     LanguageType = "ii"       // Sichuan Yi and Nuosu
	LanguageTypeIK     LanguageType = "ik"       // Inupiaq
	LanguageTypeIN     LanguageType = "in"       // Indonesian
	LanguageTypeIO     LanguageType = "io"       // Ido
	LanguageTypeIS     LanguageType = "is"       // Icelandic
	LanguageTypeIT     LanguageType = "it"       // Italian
	LanguageTypeIU     LanguageType = "iu"       // Inuktitut
	LanguageTypeIW     LanguageType = "iw"       // Hebrew
	LanguageTypeJA     LanguageType = "ja"       // Japanese
	LanguageTypeJI     LanguageType = "ji"       // Yiddish
	LanguageTypeJV     LanguageType = "jv"       // Javanese
	LanguageTypeJW     LanguageType = "jw"       // Javanese
	LanguageTypeKA     LanguageType = "ka"       // Georgian
	LanguageTypeKG     LanguageType = "kg"       // Kongo
	LanguageTypeKI     LanguageType = "ki"       // Kikuyu and Gikuyu
	LanguageTypeKJ     LanguageType = "kj"       // Kuanyama and Kwanyama
	LanguageTypeKK     LanguageType = "kk"       // Kazakh
	LanguageTypeKL     LanguageType = "kl"       // Kalaallisut and Greenlandic
	LanguageTypeKM     LanguageType = "km"       // Central Khmer
	LanguageTypeKN     LanguageType = "kn"       // Kannada
	LanguageTypeKO     LanguageType = "ko"       // Korean
	LanguageTypeKR     LanguageType = "kr"       // Kanuri
	LanguageTypeKS     LanguageType = "ks"       // Kashmiri
	LanguageTypeKU     LanguageType = "ku"       // Kurdish
	LanguageTypeKV     LanguageType = "kv"       // Komi
	LanguageTypeKW     LanguageType = "kw"       // Cornish
	LanguageTypeKY     LanguageType = "ky"       // Kirghiz and Kyrgyz
	LanguageTypeLA     LanguageType = "la"       // Latin
	LanguageTypeLB     LanguageType = "lb"       // Luxembourgish and Letzeburgesch
	LanguageTypeLG     LanguageType = "lg"       // Ganda
	LanguageTypeLI     LanguageType = "li"       // Limburgan and Limburger and Limburgish
	LanguageTypeLN     LanguageType = "ln"       // Lingala
	LanguageTypeLO     LanguageType = "lo"       // Lao
	LanguageTypeLT     LanguageType = "lt"       // Lithuanian
	LanguageTypeLU     LanguageType = "lu"       // Luba-Katanga
	LanguageTypeLV     LanguageType = "lv"       // Latvian
	LanguageTypeMG     LanguageType = "mg"       // Malagasy
	LanguageTypeMH     LanguageType = "mh"       // Marshallese
	LanguageTypeMI     LanguageType = "mi"       // Maori
	LanguageTypeMK     LanguageType = "mk"       // Macedonian
	LanguageTypeML     LanguageType = "ml"       // Malayalam
	LanguageTypeMN     LanguageType = "mn"       // Mongolian
	LanguageTypeMO     LanguageType = "mo"       // Moldavian and Moldovan
	LanguageTypeMR     LanguageType = "mr"       // Marathi
	LanguageTypeMS     LanguageType = "ms"       // Malay (macrolanguage)
	LanguageTypeMT     LanguageType = "mt"       // Maltese
	LanguageTypeMY     LanguageType = "my"       // Burmese
	LanguageTypeNA     LanguageType = "na"       // Nauru
	LanguageTypeNB     LanguageType = "nb"       // Norwegian Bokmål
	LanguageTypeND     LanguageType = "nd"       // North Ndebele
	LanguageTypeNE     LanguageType = "ne"       // Nepali (macrolanguage)
	LanguageTypeNG     LanguageType = "ng"       // Ndonga
	LanguageTypeNL     LanguageType = "nl"       // Dutch and Flemish
	LanguageTypeNN     LanguageType = "nn"       // Norwegian Nynorsk
	LanguageTypeNO     LanguageType = "no"       // Norwegian
	LanguageTypeNR     LanguageType = "nr"       // South Ndebele
	LanguageTypeNV     LanguageType = "nv"       // Navajo and Navaho
	LanguageTypeNY     LanguageType = "ny"       // Nyanja and Chewa and Chichewa
	LanguageTypeOC     LanguageType = "oc"       // Occitan (post 1500)
	LanguageTypeOJ     LanguageType = "oj"       // Ojibwa
	LanguageTypeOM     LanguageType = "om"       // Oromo
	LanguageTypeOR     LanguageType = "or"       // Oriya (macrolanguage)
	LanguageTypeOS     LanguageType = "os"       // Ossetian and Ossetic
	LanguageTypePA     LanguageType = "pa"       // Panjabi and Punjabi
	LanguageTypePI     LanguageType = "pi"       // Pali
	LanguageTypePL     LanguageType = "pl"       // Polish
	LanguageTypePS     LanguageType = "ps"       // Pushto and Pashto
	LanguageTypePT     LanguageType = "pt"       // Portuguese
	LanguageTypeQU     LanguageType = "qu"       // Quechua
	LanguageTypeRM     LanguageType = "rm"       // Romansh
	LanguageTypeRN     LanguageType = "rn"       // Rundi
	LanguageTypeRO     LanguageType = "ro"       // Romanian and Moldavian and Moldovan
	LanguageTypeRU     LanguageType = "ru"       // Russian
	LanguageTypeRW     LanguageType = "rw"       // Kinyarwanda
	LanguageTypeSA     LanguageType = "sa"       // Sanskrit
	LanguageTypeSC     LanguageType = "sc"       // Sardinian
	LanguageTypeSD     LanguageType = "sd"       // Sindhi
	LanguageTypeSE     LanguageType = "se"       // Northern Sami
	LanguageTypeSG     LanguageType = "sg"       // Sango
	LanguageTypeSH     LanguageType = "sh"       // Serbo-Croatian
	LanguageTypeSI     LanguageType = "si"       // Sinhala and Sinhalese
	LanguageTypeSK     LanguageType = "sk"       // Slovak
	LanguageTypeSL     LanguageType = "sl"       // Slovenian
	LanguageTypeSM     LanguageType = "sm"       // Samoan
	LanguageTypeSN     LanguageType = "sn"       // Shona
	LanguageTypeSO     LanguageType = "so"       // Somali
	LanguageTypeSQ     LanguageType = "sq"       // Albanian
	LanguageTypeSR     LanguageType = "sr"       // Serbian
	LanguageTypeSS     LanguageType = "ss"       // Swati
	LanguageTypeST     LanguageType = "st"       // Southern Sotho
	LanguageTypeSU     LanguageType = "su"       // Sundanese
	LanguageTypeSV     LanguageType = "sv"       // Swedish
	LanguageTypeSW     LanguageType = "sw"       // Swahili (macrolanguage)
	LanguageTypeTA     LanguageType = "ta"       // Tamil
	LanguageTypeTE     LanguageType = "te"       // Telugu
	LanguageTypeTG     LanguageType = "tg"       // Tajik
	LanguageTypeTH     LanguageType = "th"       // Thai
	LanguageTypeTI     LanguageType = "ti"       // Tigrinya
	LanguageTypeTK     LanguageType = "tk"       // Turkmen
	LanguageTypeTL     LanguageType = "tl"       // Tagalog
	LanguageTypeTN     LanguageType = "tn"       // Tswana
	LanguageTypeTO     LanguageType = "to"       // Tonga (Tonga Islands)
	LanguageTypeTR     LanguageType = "tr"       // Turkish
	LanguageTypeTS     LanguageType = "ts"       // Tsonga
	LanguageTypeTT     LanguageType = "tt"       // Tatar
	LanguageTypeTW     LanguageType = "tw"       // Twi
	LanguageTypeTY     LanguageType = "ty"       // Tahitian
	LanguageTypeUG     LanguageType = "ug"       // Uighur and Uyghur
	LanguageTypeUK     LanguageType = "uk"       // Ukrainian
	LanguageTypeUR     LanguageType = "ur"       // Urdu
	LanguageTypeUZ     LanguageType = "uz"       // Uzbek
	LanguageTypeVE     LanguageType = "ve"       // Venda
	LanguageTypeVI     LanguageType = "vi"       // Vietnamese
	LanguageTypeVO     LanguageType = "vo"       // Volapük
	LanguageTypeWA     LanguageType = "wa"       // Walloon
	LanguageTypeWO     LanguageType = "wo"       // Wolof
	LanguageTypeXH     LanguageType = "xh"       // Xhosa
	LanguageTypeYI     LanguageType = "yi"       // Yiddish
	LanguageTypeYO     LanguageType = "yo"       // Yoruba
	LanguageTypeZA     LanguageType = "za"       // Zhuang and Chuang
	LanguageTypeZH     LanguageType = "zh"       // Chinese
	LanguageTypeZU     LanguageType = "zu"       // Zulu
	LanguageTypeAAA    LanguageType = "aaa"      // Ghotuo
	LanguageTypeAAB    LanguageType = "aab"      // Alumu-Tesu
	LanguageTypeAAC    LanguageType = "aac"      // Ari
	LanguageTypeAAD    LanguageType = "aad"      // Amal
	LanguageTypeAAE    LanguageType = "aae"      // Arbëreshë Albanian
	LanguageTypeAAF    LanguageType = "aaf"      // Aranadan
	LanguageTypeAAG    LanguageType = "aag"      // Ambrak
	LanguageTypeAAH    LanguageType = "aah"      // Abu' Arapesh
	LanguageTypeAAI    LanguageType = "aai"      // Arifama-Miniafia
	LanguageTypeAAK    LanguageType = "aak"      // Ankave
	LanguageTypeAAL    LanguageType = "aal"      // Afade
	LanguageTypeAAM    LanguageType = "aam"      // Aramanik
	LanguageTypeAAN    LanguageType = "aan"      // Anambé
	LanguageTypeAAO    LanguageType = "aao"      // Algerian Saharan Arabic
	LanguageTypeAAP    LanguageType = "aap"      // Pará Arára
	LanguageTypeAAQ    LanguageType = "aaq"      // Eastern Abnaki
	LanguageTypeAAS    LanguageType = "aas"      // Aasáx
	LanguageTypeAAT    LanguageType = "aat"      // Arvanitika Albanian
	LanguageTypeAAU    LanguageType = "aau"      // Abau
	LanguageTypeAAV    LanguageType = "aav"      // Austro-Asiatic languages
	LanguageTypeAAW    LanguageType = "aaw"      // Solong
	LanguageTypeAAX    LanguageType = "aax"      // Mandobo Atas
	LanguageTypeAAZ    LanguageType = "aaz"      // Amarasi
	LanguageTypeABA    LanguageType = "aba"      // Abé
	LanguageTypeABB    LanguageType = "abb"      // Bankon
	LanguageTypeABC    LanguageType = "abc"      // Ambala Ayta
	LanguageTypeABD    LanguageType = "abd"      // Manide
	LanguageTypeABE    LanguageType = "abe"      // Western Abnaki
	LanguageTypeABF    LanguageType = "abf"      // Abai Sungai
	LanguageTypeABG    LanguageType = "abg"      // Abaga
	LanguageTypeABH    LanguageType = "abh"      // Tajiki Arabic
	LanguageTypeABI    LanguageType = "abi"      // Abidji
	LanguageTypeABJ    LanguageType = "abj"      // Aka-Bea
	LanguageTypeABL    LanguageType = "abl"      // Lampung Nyo
	LanguageTypeABM    LanguageType = "abm"      // Abanyom
	LanguageTypeABN    LanguageType = "abn"      // Abua
	LanguageTypeABO    LanguageType = "abo"      // Abon
	LanguageTypeABP    LanguageType = "abp"      // Abellen Ayta
	LanguageTypeABQ    LanguageType = "abq"      // Abaza
	LanguageTypeABR    LanguageType = "abr"      // Abron
	LanguageTypeABS    LanguageType = "abs"      // Ambonese Malay
	LanguageTypeABT    LanguageType = "abt"      // Ambulas
	LanguageTypeABU    LanguageType = "abu"      // Abure
	LanguageTypeABV    LanguageType = "abv"      // Baharna Arabic
	LanguageTypeABW    LanguageType = "abw"      // Pal
	LanguageTypeABX    LanguageType = "abx"      // Inabaknon
	LanguageTypeABY    LanguageType = "aby"      // Aneme Wake
	LanguageTypeABZ    LanguageType = "abz"      // Abui
	LanguageTypeACA    LanguageType = "aca"      // Achagua
	LanguageTypeACB    LanguageType = "acb"      // Áncá
	LanguageTypeACD    LanguageType = "acd"      // Gikyode
	LanguageTypeACE    LanguageType = "ace"      // Achinese
	LanguageTypeACF    LanguageType = "acf"      // Saint Lucian Creole French
	LanguageTypeACH    LanguageType = "ach"      // Acoli
	LanguageTypeACI    LanguageType = "aci"      // Aka-Cari
	LanguageTypeACK    LanguageType = "ack"      // Aka-Kora
	LanguageTypeACL    LanguageType = "acl"      // Akar-Bale
	LanguageTypeACM    LanguageType = "acm"      // Mesopotamian Arabic
	LanguageTypeACN    LanguageType = "acn"      // Achang
	LanguageTypeACP    LanguageType = "acp"      // Eastern Acipa
	LanguageTypeACQ    LanguageType = "acq"      // Ta'izzi-Adeni Arabic
	LanguageTypeACR    LanguageType = "acr"      // Achi
	LanguageTypeACS    LanguageType = "acs"      // Acroá
	LanguageTypeACT    LanguageType = "act"      // Achterhoeks
	LanguageTypeACU    LanguageType = "acu"      // Achuar-Shiwiar
	LanguageTypeACV    LanguageType = "acv"      // Achumawi
	LanguageTypeACW    LanguageType = "acw"      // Hijazi Arabic
	LanguageTypeACX    LanguageType = "acx"      // Omani Arabic
	LanguageTypeACY    LanguageType = "acy"      // Cypriot Arabic
	LanguageTypeACZ    LanguageType = "acz"      // Acheron
	LanguageTypeADA    LanguageType = "ada"      // Adangme
	LanguageTypeADB    LanguageType = "adb"      // Adabe
	LanguageTypeADD    LanguageType = "add"      // Dzodinka
	LanguageTypeADE    LanguageType = "ade"      // Adele
	LanguageTypeADF    LanguageType = "adf"      // Dhofari Arabic
	LanguageTypeADG    LanguageType = "adg"      // Andegerebinha
	LanguageTypeADH    LanguageType = "adh"      // Adhola
	LanguageTypeADI    LanguageType = "adi"      // Adi
	LanguageTypeADJ    LanguageType = "adj"      // Adioukrou
	LanguageTypeADL    LanguageType = "adl"      // Galo
	LanguageTypeADN    LanguageType = "adn"      // Adang
	LanguageTypeADO    LanguageType = "ado"      // Abu
	LanguageTypeADP    LanguageType = "adp"      // Adap
	LanguageTypeADQ    LanguageType = "adq"      // Adangbe
	LanguageTypeADR    LanguageType = "adr"      // Adonara
	LanguageTypeADS    LanguageType = "ads"      // Adamorobe Sign Language
	LanguageTypeADT    LanguageType = "adt"      // Adnyamathanha
	LanguageTypeADU    LanguageType = "adu"      // Aduge
	LanguageTypeADW    LanguageType = "adw"      // Amundava
	LanguageTypeADX    LanguageType = "adx"      // Amdo Tibetan
	LanguageTypeADY    LanguageType = "ady"      // Adyghe and Adygei
	LanguageTypeADZ    LanguageType = "adz"      // Adzera
	LanguageTypeAEA    LanguageType = "aea"      // Areba
	LanguageTypeAEB    LanguageType = "aeb"      // Tunisian Arabic
	LanguageTypeAEC    LanguageType = "aec"      // Saidi Arabic
	LanguageTypeAED    LanguageType = "aed"      // Argentine Sign Language
	LanguageTypeAEE    LanguageType = "aee"      // Northeast Pashayi
	LanguageTypeAEK    LanguageType = "aek"      // Haeke
	LanguageTypeAEL    LanguageType = "ael"      // Ambele
	LanguageTypeAEM    LanguageType = "aem"      // Arem
	LanguageTypeAEN    LanguageType = "aen"      // Armenian Sign Language
	LanguageTypeAEQ    LanguageType = "aeq"      // Aer
	LanguageTypeAER    LanguageType = "aer"      // Eastern Arrernte
	LanguageTypeAES    LanguageType = "aes"      // Alsea
	LanguageTypeAEU    LanguageType = "aeu"      // Akeu
	LanguageTypeAEW    LanguageType = "aew"      // Ambakich
	LanguageTypeAEY    LanguageType = "aey"      // Amele
	LanguageTypeAEZ    LanguageType = "aez"      // Aeka
	LanguageTypeAFA    LanguageType = "afa"      // Afro-Asiatic languages
	LanguageTypeAFB    LanguageType = "afb"      // Gulf Arabic
	LanguageTypeAFD    LanguageType = "afd"      // Andai
	LanguageTypeAFE    LanguageType = "afe"      // Putukwam
	LanguageTypeAFG    LanguageType = "afg"      // Afghan Sign Language
	LanguageTypeAFH    LanguageType = "afh"      // Afrihili
	LanguageTypeAFI    LanguageType = "afi"      // Akrukay
	LanguageTypeAFK    LanguageType = "afk"      // Nanubae
	LanguageTypeAFN    LanguageType = "afn"      // Defaka
	LanguageTypeAFO    LanguageType = "afo"      // Eloyi
	LanguageTypeAFP    LanguageType = "afp"      // Tapei
	LanguageTypeAFS    LanguageType = "afs"      // Afro-Seminole Creole
	LanguageTypeAFT    LanguageType = "aft"      // Afitti
	LanguageTypeAFU    LanguageType = "afu"      // Awutu
	LanguageTypeAFZ    LanguageType = "afz"      // Obokuitai
	LanguageTypeAGA    LanguageType = "aga"      // Aguano
	LanguageTypeAGB    LanguageType = "agb"      // Legbo
	LanguageTypeAGC    LanguageType = "agc"      // Agatu
	LanguageTypeAGD    LanguageType = "agd"      // Agarabi
	LanguageTypeAGE    LanguageType = "age"      // Angal
	LanguageTypeAGF    LanguageType = "agf"      // Arguni
	LanguageTypeAGG    LanguageType = "agg"      // Angor
	LanguageTypeAGH    LanguageType = "agh"      // Ngelima
	LanguageTypeAGI    LanguageType = "agi"      // Agariya
	LanguageTypeAGJ    LanguageType = "agj"      // Argobba
	LanguageTypeAGK    LanguageType = "agk"      // Isarog Agta
	LanguageTypeAGL    LanguageType = "agl"      // Fembe
	LanguageTypeAGM    LanguageType = "agm"      // Angaataha
	LanguageTypeAGN    LanguageType = "agn"      // Agutaynen
	LanguageTypeAGO    LanguageType = "ago"      // Tainae
	LanguageTypeAGP    LanguageType = "agp"      // Paranan
	LanguageTypeAGQ    LanguageType = "agq"      // Aghem
	LanguageTypeAGR    LanguageType = "agr"      // Aguaruna
	LanguageTypeAGS    LanguageType = "ags"      // Esimbi
	LanguageTypeAGT    LanguageType = "agt"      // Central Cagayan Agta
	LanguageTypeAGU    LanguageType = "agu"      // Aguacateco
	LanguageTypeAGV    LanguageType = "agv"      // Remontado Dumagat
	LanguageTypeAGW    LanguageType = "agw"      // Kahua
	LanguageTypeAGX    LanguageType = "agx"      // Aghul
	LanguageTypeAGY    LanguageType = "agy"      // Southern Alta
	LanguageTypeAGZ    LanguageType = "agz"      // Mt. Iriga Agta
	LanguageTypeAHA    LanguageType = "aha"      // Ahanta
	LanguageTypeAHB    LanguageType = "ahb"      // Axamb
	LanguageTypeAHG    LanguageType = "ahg"      // Qimant
	LanguageTypeAHH    LanguageType = "ahh"      // Aghu
	LanguageTypeAHI    LanguageType = "ahi"      // Tiagbamrin Aizi
	LanguageTypeAHK    LanguageType = "ahk"      // Akha
	LanguageTypeAHL    LanguageType = "ahl"      // Igo
	LanguageTypeAHM    LanguageType = "ahm"      // Mobumrin Aizi
	LanguageTypeAHN    LanguageType = "ahn"      // Àhàn
	LanguageTypeAHO    LanguageType = "aho"      // Ahom
	LanguageTypeAHP    LanguageType = "ahp"      // Aproumu Aizi
	LanguageTypeAHR    LanguageType = "ahr"      // Ahirani
	LanguageTypeAHS    LanguageType = "ahs"      // Ashe
	LanguageTypeAHT    LanguageType = "aht"      // Ahtena
	LanguageTypeAIA    LanguageType = "aia"      // Arosi
	LanguageTypeAIB    LanguageType = "aib"      // Ainu (China)
	LanguageTypeAIC    LanguageType = "aic"      // Ainbai
	LanguageTypeAID    LanguageType = "aid"      // Alngith
	LanguageTypeAIE    LanguageType = "aie"      // Amara
	LanguageTypeAIF    LanguageType = "aif"      // Agi
	LanguageTypeAIG    LanguageType = "aig"      // Antigua and Barbuda Creole English
	LanguageTypeAIH    LanguageType = "aih"      // Ai-Cham
	LanguageTypeAII    LanguageType = "aii"      // Assyrian Neo-Aramaic
	LanguageTypeAIJ    LanguageType = "aij"      // Lishanid Noshan
	LanguageTypeAIK    LanguageType = "aik"      // Ake
	LanguageTypeAIL    LanguageType = "ail"      // Aimele
	LanguageTypeAIM    LanguageType = "aim"      // Aimol
	LanguageTypeAIN    LanguageType = "ain"      // Ainu (Japan)
	LanguageTypeAIO    LanguageType = "aio"      // Aiton
	LanguageTypeAIP    LanguageType = "aip"      // Burumakok
	LanguageTypeAIQ    LanguageType = "aiq"      // Aimaq
	LanguageTypeAIR    LanguageType = "air"      // Airoran
	LanguageTypeAIS    LanguageType = "ais"      // Nataoran Amis
	LanguageTypeAIT    LanguageType = "ait"      // Arikem
	LanguageTypeAIW    LanguageType = "aiw"      // Aari
	LanguageTypeAIX    LanguageType = "aix"      // Aighon
	LanguageTypeAIY    LanguageType = "aiy"      // Ali
	LanguageTypeAJA    LanguageType = "aja"      // Aja (Sudan)
	LanguageTypeAJG    LanguageType = "ajg"      // Aja (Benin)
	LanguageTypeAJI    LanguageType = "aji"      // Ajië
	LanguageTypeAJN    LanguageType = "ajn"      // Andajin
	LanguageTypeAJP    LanguageType = "ajp"      // South Levantine Arabic
	LanguageTypeAJT    LanguageType = "ajt"      // Judeo-Tunisian Arabic
	LanguageTypeAJU    LanguageType = "aju"      // Judeo-Moroccan Arabic
	LanguageTypeAJW    LanguageType = "ajw"      // Ajawa
	LanguageTypeAJZ    LanguageType = "ajz"      // Amri Karbi
	LanguageTypeAKB    LanguageType = "akb"      // Batak Angkola
	LanguageTypeAKC    LanguageType = "akc"      // Mpur
	LanguageTypeAKD    LanguageType = "akd"      // Ukpet-Ehom
	LanguageTypeAKE    LanguageType = "ake"      // Akawaio
	LanguageTypeAKF    LanguageType = "akf"      // Akpa
	LanguageTypeAKG    LanguageType = "akg"      // Anakalangu
	LanguageTypeAKH    LanguageType = "akh"      // Angal Heneng
	LanguageTypeAKI    LanguageType = "aki"      // Aiome
	LanguageTypeAKJ    LanguageType = "akj"      // Aka-Jeru
	LanguageTypeAKK    LanguageType = "akk"      // Akkadian
	LanguageTypeAKL    LanguageType = "akl"      // Aklanon
	LanguageTypeAKM    LanguageType = "akm"      // Aka-Bo
	LanguageTypeAKO    LanguageType = "ako"      // Akurio
	LanguageTypeAKP    LanguageType = "akp"      // Siwu
	LanguageTypeAKQ    LanguageType = "akq"      // Ak
	LanguageTypeAKR    LanguageType = "akr"      // Araki
	LanguageTypeAKS    LanguageType = "aks"      // Akaselem
	LanguageTypeAKT    LanguageType = "akt"      // Akolet
	LanguageTypeAKU    LanguageType = "aku"      // Akum
	LanguageTypeAKV    LanguageType = "akv"      // Akhvakh
	LanguageTypeAKW    LanguageType = "akw"      // Akwa
	LanguageTypeAKX    LanguageType = "akx"      // Aka-Kede
	LanguageTypeAKY    LanguageType = "aky"      // Aka-Kol
	LanguageTypeAKZ    LanguageType = "akz"      // Alabama
	LanguageTypeALA    LanguageType = "ala"      // Alago
	LanguageTypeALC    LanguageType = "alc"      // Qawasqar
	LanguageTypeALD    LanguageType = "ald"      // Alladian
	LanguageTypeALE    LanguageType = "ale"      // Aleut
	LanguageTypeALF    LanguageType = "alf"      // Alege
	LanguageTypeALG    LanguageType = "alg"      // Algonquian languages
	LanguageTypeALH    LanguageType = "alh"      // Alawa
	LanguageTypeALI    LanguageType = "ali"      // Amaimon
	LanguageTypeALJ    LanguageType = "alj"      // Alangan
	LanguageTypeALK    LanguageType = "alk"      // Alak
	LanguageTypeALL    LanguageType = "all"      // Allar
	LanguageTypeALM    LanguageType = "alm"      // Amblong
	LanguageTypeALN    LanguageType = "aln"      // Gheg Albanian
	LanguageTypeALO    LanguageType = "alo"      // Larike-Wakasihu
	LanguageTypeALP    LanguageType = "alp"      // Alune
	LanguageTypeALQ    LanguageType = "alq"      // Algonquin
	LanguageTypeALR    LanguageType = "alr"      // Alutor
	LanguageTypeALS    LanguageType = "als"      // Tosk Albanian
	LanguageTypeALT    LanguageType = "alt"      // Southern Altai
	LanguageTypeALU    LanguageType = "alu"      // 'Are'are
	LanguageTypeALV    LanguageType = "alv"      // Atlantic-Congo languages
	LanguageTypeALW    LanguageType = "alw"      // Alaba-K’abeena and Wanbasana
	LanguageTypeALX    LanguageType = "alx"      // Amol
	LanguageTypeALY    LanguageType = "aly"      // Alyawarr
	LanguageTypeALZ    LanguageType = "alz"      // Alur
	LanguageTypeAMA    LanguageType = "ama"      // Amanayé
	LanguageTypeAMB    LanguageType = "amb"      // Ambo
	LanguageTypeAMC    LanguageType = "amc"      // Amahuaca
	LanguageTypeAME    LanguageType = "ame"      // Yanesha'
	LanguageTypeAMF    LanguageType = "amf"      // Hamer-Banna
	LanguageTypeAMG    LanguageType = "amg"      // Amurdak
	LanguageTypeAMI    LanguageType = "ami"      // Amis
	LanguageTypeAMJ    LanguageType = "amj"      // Amdang
	LanguageTypeAMK    LanguageType = "amk"      // Ambai
	LanguageTypeAML    LanguageType = "aml"      // War-Jaintia
	LanguageTypeAMM    LanguageType = "amm"      // Ama (Papua New Guinea)
	LanguageTypeAMN    LanguageType = "amn"      // Amanab
	LanguageTypeAMO    LanguageType = "amo"      // Amo
	LanguageTypeAMP    LanguageType = "amp"      // Alamblak
	LanguageTypeAMQ    LanguageType = "amq"      // Amahai
	LanguageTypeAMR    LanguageType = "amr"      // Amarakaeri
	LanguageTypeAMS    LanguageType = "ams"      // Southern Amami-Oshima
	LanguageTypeAMT    LanguageType = "amt"      // Amto
	LanguageTypeAMU    LanguageType = "amu"      // Guerrero Amuzgo
	LanguageTypeAMV    LanguageType = "amv"      // Ambelau
	LanguageTypeAMW    LanguageType = "amw"      // Western Neo-Aramaic
	LanguageTypeAMX    LanguageType = "amx"      // Anmatyerre
	LanguageTypeAMY    LanguageType = "amy"      // Ami
	LanguageTypeAMZ    LanguageType = "amz"      // Atampaya
	LanguageTypeANA    LanguageType = "ana"      // Andaqui
	LanguageTypeANB    LanguageType = "anb"      // Andoa
	LanguageTypeANC    LanguageType = "anc"      // Ngas
	LanguageTypeAND    LanguageType = "and"      // Ansus
	LanguageTypeANE    LanguageType = "ane"      // Xârâcùù
	LanguageTypeANF    LanguageType = "anf"      // Animere
	LanguageTypeANG    LanguageType = "ang"      // Old English (ca. 450-1100)
	LanguageTypeANH    LanguageType = "anh"      // Nend
	LanguageTypeANI    LanguageType = "ani"      // Andi
	LanguageTypeANJ    LanguageType = "anj"      // Anor
	LanguageTypeANK    LanguageType = "ank"      // Goemai
	LanguageTypeANL    LanguageType = "anl"      // Anu-Hkongso Chin
	LanguageTypeANM    LanguageType = "anm"      // Anal
	LanguageTypeANN    LanguageType = "ann"      // Obolo
	LanguageTypeANO    LanguageType = "ano"      // Andoque
	LanguageTypeANP    LanguageType = "anp"      // Angika
	LanguageTypeANQ    LanguageType = "anq"      // Jarawa (India)
	LanguageTypeANR    LanguageType = "anr"      // Andh
	LanguageTypeANS    LanguageType = "ans"      // Anserma
	LanguageTypeANT    LanguageType = "ant"      // Antakarinya
	LanguageTypeANU    LanguageType = "anu"      // Anuak
	LanguageTypeANV    LanguageType = "anv"      // Denya
	LanguageTypeANW    LanguageType = "anw"      // Anaang
	LanguageTypeANX    LanguageType = "anx"      // Andra-Hus
	LanguageTypeANY    LanguageType = "any"      // Anyin
	LanguageTypeANZ    LanguageType = "anz"      // Anem
	LanguageTypeAOA    LanguageType = "aoa"      // Angolar
	LanguageTypeAOB    LanguageType = "aob"      // Abom
	LanguageTypeAOC    LanguageType = "aoc"      // Pemon
	LanguageTypeAOD    LanguageType = "aod"      // Andarum
	LanguageTypeAOE    LanguageType = "aoe"      // Angal Enen
	LanguageTypeAOF    LanguageType = "aof"      // Bragat
	LanguageTypeAOG    LanguageType = "aog"      // Angoram
	LanguageTypeAOH    LanguageType = "aoh"      // Arma
	LanguageTypeAOI    LanguageType = "aoi"      // Anindilyakwa
	LanguageTypeAOJ    LanguageType = "aoj"      // Mufian
	LanguageTypeAOK    LanguageType = "aok"      // Arhö
	LanguageTypeAOL    LanguageType = "aol"      // Alor
	LanguageTypeAOM    LanguageType = "aom"      // Ömie
	LanguageTypeAON    LanguageType = "aon"      // Bumbita Arapesh
	LanguageTypeAOR    LanguageType = "aor"      // Aore
	LanguageTypeAOS    LanguageType = "aos"      // Taikat
	LanguageTypeAOT    LanguageType = "aot"      // A'tong
	LanguageTypeAOU    LanguageType = "aou"      // A'ou
	LanguageTypeAOX    LanguageType = "aox"      // Atorada
	LanguageTypeAOZ    LanguageType = "aoz"      // Uab Meto
	LanguageTypeAPA    LanguageType = "apa"      // Apache languages
	LanguageTypeAPB    LanguageType = "apb"      // Sa'a
	LanguageTypeAPC    LanguageType = "apc"      // North Levantine Arabic
	LanguageTypeAPD    LanguageType = "apd"      // Sudanese Arabic
	LanguageTypeAPE    LanguageType = "ape"      // Bukiyip
	LanguageTypeAPF    LanguageType = "apf"      // Pahanan Agta
	LanguageTypeAPG    LanguageType = "apg"      // Ampanang
	LanguageTypeAPH    LanguageType = "aph"      // Athpariya
	LanguageTypeAPI    LanguageType = "api"      // Apiaká
	LanguageTypeAPJ    LanguageType = "apj"      // Jicarilla Apache
	LanguageTypeAPK    LanguageType = "apk"      // Kiowa Apache
	LanguageTypeAPL    LanguageType = "apl"      // Lipan Apache
	LanguageTypeAPM    LanguageType = "apm"      // Mescalero-Chiricahua Apache
	LanguageTypeAPN    LanguageType = "apn"      // Apinayé
	LanguageTypeAPO    LanguageType = "apo"      // Ambul
	LanguageTypeAPP    LanguageType = "app"      // Apma
	LanguageTypeAPQ    LanguageType = "apq"      // A-Pucikwar
	LanguageTypeAPR    LanguageType = "apr"      // Arop-Lokep
	LanguageTypeAPS    LanguageType = "aps"      // Arop-Sissano
	LanguageTypeAPT    LanguageType = "apt"      // Apatani
	LanguageTypeAPU    LanguageType = "apu"      // Apurinã
	LanguageTypeAPV    LanguageType = "apv"      // Alapmunte
	LanguageTypeAPW    LanguageType = "apw"      // Western Apache
	LanguageTypeAPX    LanguageType = "apx"      // Aputai
	LanguageTypeAPY    LanguageType = "apy"      // Apalaí
	LanguageTypeAPZ    LanguageType = "apz"      // Safeyoka
	LanguageTypeAQA    LanguageType = "aqa"      // Alacalufan languages
	LanguageTypeAQC    LanguageType = "aqc"      // Archi
	LanguageTypeAQD    LanguageType = "aqd"      // Ampari Dogon
	LanguageTypeAQG    LanguageType = "aqg"      // Arigidi
	LanguageTypeAQL    LanguageType = "aql"      // Algic languages
	LanguageTypeAQM    LanguageType = "aqm"      // Atohwaim
	LanguageTypeAQN    LanguageType = "aqn"      // Northern Alta
	LanguageTypeAQP    LanguageType = "aqp"      // Atakapa
	LanguageTypeAQR    LanguageType = "aqr"      // Arhâ
	LanguageTypeAQZ    LanguageType = "aqz"      // Akuntsu
	LanguageTypeARB    LanguageType = "arb"      // Standard Arabic
	LanguageTypeARC    LanguageType = "arc"      // Official Aramaic (700-300 BCE) and Imperial Aramaic (700-300 BCE)
	LanguageTypeARD    LanguageType = "ard"      // Arabana
	LanguageTypeARE    LanguageType = "are"      // Western Arrarnta
	LanguageTypeARH    LanguageType = "arh"      // Arhuaco
	LanguageTypeARI    LanguageType = "ari"      // Arikara
	LanguageTypeARJ    LanguageType = "arj"      // Arapaso
	LanguageTypeARK    LanguageType = "ark"      // Arikapú
	LanguageTypeARL    LanguageType = "arl"      // Arabela
	LanguageTypeARN    LanguageType = "arn"      // Mapudungun and Mapuche
	LanguageTypeARO    LanguageType = "aro"      // Araona
	LanguageTypeARP    LanguageType = "arp"      // Arapaho
	LanguageTypeARQ    LanguageType = "arq"      // Algerian Arabic
	LanguageTypeARR    LanguageType = "arr"      // Karo (Brazil)
	LanguageTypeARS    LanguageType = "ars"      // Najdi Arabic
	LanguageTypeART    LanguageType = "art"      // Artificial languages
	LanguageTypeARU    LanguageType = "aru"      // Aruá (Amazonas State) and Arawá
	LanguageTypeARV    LanguageType = "arv"      // Arbore
	LanguageTypeARW    LanguageType = "arw"      // Arawak
	LanguageTypeARX    LanguageType = "arx"      // Aruá (Rodonia State)
	LanguageTypeARY    LanguageType = "ary"      // Moroccan Arabic
	LanguageTypeARZ    LanguageType = "arz"      // Egyptian Arabic
	LanguageTypeASA    LanguageType = "asa"      // Asu (Tanzania)
	LanguageTypeASB    LanguageType = "asb"      // Assiniboine
	LanguageTypeASC    LanguageType = "asc"      // Casuarina Coast Asmat
	LanguageTypeASD    LanguageType = "asd"      // Asas
	LanguageTypeASE    LanguageType = "ase"      // American Sign Language
	LanguageTypeASF    LanguageType = "asf"      // Australian Sign Language
	LanguageTypeASG    LanguageType = "asg"      // Cishingini
	LanguageTypeASH    LanguageType = "ash"      // Abishira
	LanguageTypeASI    LanguageType = "asi"      // Buruwai
	LanguageTypeASJ    LanguageType = "asj"      // Sari
	LanguageTypeASK    LanguageType = "ask"      // Ashkun
	LanguageTypeASL    LanguageType = "asl"      // Asilulu
	LanguageTypeASN    LanguageType = "asn"      // Xingú Asuriní
	LanguageTypeASO    LanguageType = "aso"      // Dano
	LanguageTypeASP    LanguageType = "asp"      // Algerian Sign Language
	LanguageTypeASQ    LanguageType = "asq"      // Austrian Sign Language
	LanguageTypeASR    LanguageType = "asr"      // Asuri
	LanguageTypeASS    LanguageType = "ass"      // Ipulo
	LanguageTypeAST    LanguageType = "ast"      // Asturian and Asturleonese and Bable and Leonese
	LanguageTypeASU    LanguageType = "asu"      // Tocantins Asurini
	LanguageTypeASV    LanguageType = "asv"      // Asoa
	LanguageTypeASW    LanguageType = "asw"      // Australian Aborigines Sign Language
	LanguageTypeASX    LanguageType = "asx"      // Muratayak
	LanguageTypeASY    LanguageType = "asy"      // Yaosakor Asmat
	LanguageTypeASZ    LanguageType = "asz"      // As
	LanguageTypeATA    LanguageType = "ata"      // Pele-Ata
	LanguageTypeATB    LanguageType = "atb"      // Zaiwa
	LanguageTypeATC    LanguageType = "atc"      // Atsahuaca
	LanguageTypeATD    LanguageType = "atd"      // Ata Manobo
	LanguageTypeATE    LanguageType = "ate"      // Atemble
	LanguageTypeATG    LanguageType = "atg"      // Ivbie North-Okpela-Arhe
	LanguageTypeATH    LanguageType = "ath"      // Athapascan languages
	LanguageTypeATI    LanguageType = "ati"      // Attié
	LanguageTypeATJ    LanguageType = "atj"      // Atikamekw
	LanguageTypeATK    LanguageType = "atk"      // Ati
	LanguageTypeATL    LanguageType = "atl"      // Mt. Iraya Agta
	LanguageTypeATM    LanguageType = "atm"      // Ata
	LanguageTypeATN    LanguageType = "atn"      // Ashtiani
	LanguageTypeATO    LanguageType = "ato"      // Atong
	LanguageTypeATP    LanguageType = "atp"      // Pudtol Atta
	LanguageTypeATQ    LanguageType = "atq"      // Aralle-Tabulahan
	LanguageTypeATR    LanguageType = "atr"      // Waimiri-Atroari
	LanguageTypeATS    LanguageType = "ats"      // Gros Ventre
	LanguageTypeATT    LanguageType = "att"      // Pamplona Atta
	LanguageTypeATU    LanguageType = "atu"      // Reel
	LanguageTypeATV    LanguageType = "atv"      // Northern Altai
	LanguageTypeATW    LanguageType = "atw"      // Atsugewi
	LanguageTypeATX    LanguageType = "atx"      // Arutani
	LanguageTypeATY    LanguageType = "aty"      // Aneityum
	LanguageTypeATZ    LanguageType = "atz"      // Arta
	LanguageTypeAUA    LanguageType = "aua"      // Asumboa
	LanguageTypeAUB    LanguageType = "aub"      // Alugu
	LanguageTypeAUC    LanguageType = "auc"      // Waorani
	LanguageTypeAUD    LanguageType = "aud"      // Anuta
	LanguageTypeAUE    LanguageType = "aue"      // =/Kx'au//'ein
	LanguageTypeAUF    LanguageType = "auf"      // Arauan languages
	LanguageTypeAUG    LanguageType = "aug"      // Aguna
	LanguageTypeAUH    LanguageType = "auh"      // Aushi
	LanguageTypeAUI    LanguageType = "aui"      // Anuki
	LanguageTypeAUJ    LanguageType = "auj"      // Awjilah
	LanguageTypeAUK    LanguageType = "auk"      // Heyo
	LanguageTypeAUL    LanguageType = "aul"      // Aulua
	LanguageTypeAUM    LanguageType = "aum"      // Asu (Nigeria)
	LanguageTypeAUN    LanguageType = "aun"      // Molmo One
	LanguageTypeAUO    LanguageType = "auo"      // Auyokawa
	LanguageTypeAUP    LanguageType = "aup"      // Makayam
	LanguageTypeAUQ    LanguageType = "auq"      // Anus and Korur
	LanguageTypeAUR    LanguageType = "aur"      // Aruek
	LanguageTypeAUS    LanguageType = "aus"      // Australian languages
	LanguageTypeAUT    LanguageType = "aut"      // Austral
	LanguageTypeAUU    LanguageType = "auu"      // Auye
	LanguageTypeAUW    LanguageType = "auw"      // Awyi
	LanguageTypeAUX    LanguageType = "aux"      // Aurá
	LanguageTypeAUY    LanguageType = "auy"      // Awiyaana
	LanguageTypeAUZ    LanguageType = "auz"      // Uzbeki Arabic
	LanguageTypeAVB    LanguageType = "avb"      // Avau
	LanguageTypeAVD    LanguageType = "avd"      // Alviri-Vidari
	LanguageTypeAVI    LanguageType = "avi"      // Avikam
	LanguageTypeAVK    LanguageType = "avk"      // Kotava
	LanguageTypeAVL    LanguageType = "avl"      // Eastern Egyptian Bedawi Arabic
	LanguageTypeAVM    LanguageType = "avm"      // Angkamuthi
	LanguageTypeAVN    LanguageType = "avn"      // Avatime
	LanguageTypeAVO    LanguageType = "avo"      // Agavotaguerra
	LanguageTypeAVS    LanguageType = "avs"      // Aushiri
	LanguageTypeAVT    LanguageType = "avt"      // Au
	LanguageTypeAVU    LanguageType = "avu"      // Avokaya
	LanguageTypeAVV    LanguageType = "avv"      // Avá-Canoeiro
	LanguageTypeAWA    LanguageType = "awa"      // Awadhi
	LanguageTypeAWB    LanguageType = "awb"      // Awa (Papua New Guinea)
	LanguageTypeAWC    LanguageType = "awc"      // Cicipu
	LanguageTypeAWD    LanguageType = "awd"      // Arawakan languages
	LanguageTypeAWE    LanguageType = "awe"      // Awetí
	LanguageTypeAWG    LanguageType = "awg"      // Anguthimri
	LanguageTypeAWH    LanguageType = "awh"      // Awbono
	LanguageTypeAWI    LanguageType = "awi"      // Aekyom
	LanguageTypeAWK    LanguageType = "awk"      // Awabakal
	LanguageTypeAWM    LanguageType = "awm"      // Arawum
	LanguageTypeAWN    LanguageType = "awn"      // Awngi
	LanguageTypeAWO    LanguageType = "awo"      // Awak
	LanguageTypeAWR    LanguageType = "awr"      // Awera
	LanguageTypeAWS    LanguageType = "aws"      // South Awyu
	LanguageTypeAWT    LanguageType = "awt"      // Araweté
	LanguageTypeAWU    LanguageType = "awu"      // Central Awyu
	LanguageTypeAWV    LanguageType = "awv"      // Jair Awyu
	LanguageTypeAWW    LanguageType = "aww"      // Awun
	LanguageTypeAWX    LanguageType = "awx"      // Awara
	LanguageTypeAWY    LanguageType = "awy"      // Edera Awyu
	LanguageTypeAXB    LanguageType = "axb"      // Abipon
	LanguageTypeAXE    LanguageType = "axe"      // Ayerrerenge
	LanguageTypeAXG    LanguageType = "axg"      // Mato Grosso Arára
	LanguageTypeAXK    LanguageType = "axk"      // Yaka (Central African Republic)
	LanguageTypeAXL    LanguageType = "axl"      // Lower Southern Aranda
	LanguageTypeAXM    LanguageType = "axm"      // Middle Armenian
	LanguageTypeAXX    LanguageType = "axx"      // Xârâgurè
	LanguageTypeAYA    LanguageType = "aya"      // Awar
	LanguageTypeAYB    LanguageType = "ayb"      // Ayizo Gbe
	LanguageTypeAYC    LanguageType = "ayc"      // Southern Aymara
	LanguageTypeAYD    LanguageType = "ayd"      // Ayabadhu
	LanguageTypeAYE    LanguageType = "aye"      // Ayere
	LanguageTypeAYG    LanguageType = "ayg"      // Ginyanga
	LanguageTypeAYH    LanguageType = "ayh"      // Hadrami Arabic
	LanguageTypeAYI    LanguageType = "ayi"      // Leyigha
	LanguageTypeAYK    LanguageType = "ayk"      // Akuku
	LanguageTypeAYL    LanguageType = "ayl"      // Libyan Arabic
	LanguageTypeAYN    LanguageType = "ayn"      // Sanaani Arabic
	LanguageTypeAYO    LanguageType = "ayo"      // Ayoreo
	LanguageTypeAYP    LanguageType = "ayp"      // North Mesopotamian Arabic
	LanguageTypeAYQ    LanguageType = "ayq"      // Ayi (Papua New Guinea)
	LanguageTypeAYR    LanguageType = "ayr"      // Central Aymara
	LanguageTypeAYS    LanguageType = "ays"      // Sorsogon Ayta
	LanguageTypeAYT    LanguageType = "ayt"      // Magbukun Ayta
	LanguageTypeAYU    LanguageType = "ayu"      // Ayu
	LanguageTypeAYX    LanguageType = "ayx"      // Ayi (China)
	LanguageTypeAYY    LanguageType = "ayy"      // Tayabas Ayta
	LanguageTypeAYZ    LanguageType = "ayz"      // Mai Brat
	LanguageTypeAZA    LanguageType = "aza"      // Azha
	LanguageTypeAZB    LanguageType = "azb"      // South Azerbaijani
	LanguageTypeAZC    LanguageType = "azc"      // Uto-Aztecan languages
	LanguageTypeAZD    LanguageType = "azd"      // Eastern Durango Nahuatl
	LanguageTypeAZG    LanguageType = "azg"      // San Pedro Amuzgos Amuzgo
	LanguageTypeAZJ    LanguageType = "azj"      // North Azerbaijani
	LanguageTypeAZM    LanguageType = "azm"      // Ipalapa Amuzgo
	LanguageTypeAZN    LanguageType = "azn"      // Western Durango Nahuatl
	LanguageTypeAZO    LanguageType = "azo"      // Awing
	LanguageTypeAZT    LanguageType = "azt"      // Faire Atta
	LanguageTypeAZZ    LanguageType = "azz"      // Highland Puebla Nahuatl
	LanguageTypeBAA    LanguageType = "baa"      // Babatana
	LanguageTypeBAB    LanguageType = "bab"      // Bainouk-Gunyuño
	LanguageTypeBAC    LanguageType = "bac"      // Badui
	LanguageTypeBAD    LanguageType = "bad"      // Banda languages
	LanguageTypeBAE    LanguageType = "bae"      // Baré
	LanguageTypeBAF    LanguageType = "baf"      // Nubaca
	LanguageTypeBAG    LanguageType = "bag"      // Tuki
	LanguageTypeBAH    LanguageType = "bah"      // Bahamas Creole English
	LanguageTypeBAI    LanguageType = "bai"      // Bamileke languages
	LanguageTypeBAJ    LanguageType = "baj"      // Barakai
	LanguageTypeBAL    LanguageType = "bal"      // Baluchi
	LanguageTypeBAN    LanguageType = "ban"      // Balinese
	LanguageTypeBAO    LanguageType = "bao"      // Waimaha
	LanguageTypeBAP    LanguageType = "bap"      // Bantawa
	LanguageTypeBAR    LanguageType = "bar"      // Bavarian
	LanguageTypeBAS    LanguageType = "bas"      // Basa (Cameroon)
	LanguageTypeBAT    LanguageType = "bat"      // Baltic languages
	LanguageTypeBAU    LanguageType = "bau"      // Bada (Nigeria)
	LanguageTypeBAV    LanguageType = "bav"      // Vengo
	LanguageTypeBAW    LanguageType = "baw"      // Bambili-Bambui
	LanguageTypeBAX    LanguageType = "bax"      // Bamun
	LanguageTypeBAY    LanguageType = "bay"      // Batuley
	LanguageTypeBAZ    LanguageType = "baz"      // Tunen
	LanguageTypeBBA    LanguageType = "bba"      // Baatonum
	LanguageTypeBBB    LanguageType = "bbb"      // Barai
	LanguageTypeBBC    LanguageType = "bbc"      // Batak Toba
	LanguageTypeBBD    LanguageType = "bbd"      // Bau
	LanguageTypeBBE    LanguageType = "bbe"      // Bangba
	LanguageTypeBBF    LanguageType = "bbf"      // Baibai
	LanguageTypeBBG    LanguageType = "bbg"      // Barama
	LanguageTypeBBH    LanguageType = "bbh"      // Bugan
	LanguageTypeBBI    LanguageType = "bbi"      // Barombi
	LanguageTypeBBJ    LanguageType = "bbj"      // Ghomálá'
	LanguageTypeBBK    LanguageType = "bbk"      // Babanki
	LanguageTypeBBL    LanguageType = "bbl"      // Bats
	LanguageTypeBBM    LanguageType = "bbm"      // Babango
	LanguageTypeBBN    LanguageType = "bbn"      // Uneapa
	LanguageTypeBBO    LanguageType = "bbo"      // Northern Bobo Madaré and Konabéré
	LanguageTypeBBP    LanguageType = "bbp"      // West Central Banda
	LanguageTypeBBQ    LanguageType = "bbq"      // Bamali
	LanguageTypeBBR    LanguageType = "bbr"      // Girawa
	LanguageTypeBBS    LanguageType = "bbs"      // Bakpinka
	LanguageTypeBBT    LanguageType = "bbt"      // Mburku
	LanguageTypeBBU    LanguageType = "bbu"      // Kulung (Nigeria)
	LanguageTypeBBV    LanguageType = "bbv"      // Karnai
	LanguageTypeBBW    LanguageType = "bbw"      // Baba
	LanguageTypeBBX    LanguageType = "bbx"      // Bubia
	LanguageTypeBBY    LanguageType = "bby"      // Befang
	LanguageTypeBBZ    LanguageType = "bbz"      // Babalia Creole Arabic
	LanguageTypeBCA    LanguageType = "bca"      // Central Bai
	LanguageTypeBCB    LanguageType = "bcb"      // Bainouk-Samik
	LanguageTypeBCC    LanguageType = "bcc"      // Southern Balochi
	LanguageTypeBCD    LanguageType = "bcd"      // North Babar
	LanguageTypeBCE    LanguageType = "bce"      // Bamenyam
	LanguageTypeBCF    LanguageType = "bcf"      // Bamu
	LanguageTypeBCG    LanguageType = "bcg"      // Baga Binari
	LanguageTypeBCH    LanguageType = "bch"      // Bariai
	LanguageTypeBCI    LanguageType = "bci"      // Baoulé
	LanguageTypeBCJ    LanguageType = "bcj"      // Bardi
	LanguageTypeBCK    LanguageType = "bck"      // Bunaba
	LanguageTypeBCL    LanguageType = "bcl"      // Central Bikol
	LanguageTypeBCM    LanguageType = "bcm"      // Bannoni
	LanguageTypeBCN    LanguageType = "bcn"      // Bali (Nigeria)
	LanguageTypeBCO    LanguageType = "bco"      // Kaluli
	LanguageTypeBCP    LanguageType = "bcp"      // Bali (Democratic Republic of Congo)
	LanguageTypeBCQ    LanguageType = "bcq"      // Bench
	LanguageTypeBCR    LanguageType = "bcr"      // Babine
	LanguageTypeBCS    LanguageType = "bcs"      // Kohumono
	LanguageTypeBCT    LanguageType = "bct"      // Bendi
	LanguageTypeBCU    LanguageType = "bcu"      // Awad Bing
	LanguageTypeBCV    LanguageType = "bcv"      // Shoo-Minda-Nye
	LanguageTypeBCW    LanguageType = "bcw"      // Bana
	LanguageTypeBCY    LanguageType = "bcy"      // Bacama
	LanguageTypeBCZ    LanguageType = "bcz"      // Bainouk-Gunyaamolo
	LanguageTypeBDA    LanguageType = "bda"      // Bayot
	LanguageTypeBDB    LanguageType = "bdb"      // Basap
	LanguageTypeBDC    LanguageType = "bdc"      // Emberá-Baudó
	LanguageTypeBDD    LanguageType = "bdd"      // Bunama
	LanguageTypeBDE    LanguageType = "bde"      // Bade
	LanguageTypeBDF    LanguageType = "bdf"      // Biage
	LanguageTypeBDG    LanguageType = "bdg"      // Bonggi
	LanguageTypeBDH    LanguageType = "bdh"      // Baka (Sudan)
	LanguageTypeBDI    LanguageType = "bdi"      // Burun
	LanguageTypeBDJ    LanguageType = "bdj"      // Bai
	LanguageTypeBDK    LanguageType = "bdk"      // Budukh
	LanguageTypeBDL    LanguageType = "bdl"      // Indonesian Bajau
	LanguageTypeBDM    LanguageType = "bdm"      // Buduma
	LanguageTypeBDN    LanguageType = "bdn"      // Baldemu
	LanguageTypeBDO    LanguageType = "bdo"      // Morom
	LanguageTypeBDP    LanguageType = "bdp"      // Bende
	LanguageTypeBDQ    LanguageType = "bdq"      // Bahnar
	LanguageTypeBDR    LanguageType = "bdr"      // West Coast Bajau
	LanguageTypeBDS    LanguageType = "bds"      // Burunge
	LanguageTypeBDT    LanguageType = "bdt"      // Bokoto
	LanguageTypeBDU    LanguageType = "bdu"      // Oroko
	LanguageTypeBDV    LanguageType = "bdv"      // Bodo Parja
	LanguageTypeBDW    LanguageType = "bdw"      // Baham
	LanguageTypeBDX    LanguageType = "bdx"      // Budong-Budong
	LanguageTypeBDY    LanguageType = "bdy"      // Bandjalang
	LanguageTypeBDZ    LanguageType = "bdz"      // Badeshi
	LanguageTypeBEA    LanguageType = "bea"      // Beaver
	LanguageTypeBEB    LanguageType = "beb"      // Bebele
	LanguageTypeBEC    LanguageType = "bec"      // Iceve-Maci
	LanguageTypeBED    LanguageType = "bed"      // Bedoanas
	LanguageTypeBEE    LanguageType = "bee"      // Byangsi
	LanguageTypeBEF    LanguageType = "bef"      // Benabena
	LanguageTypeBEG    LanguageType = "beg"      // Belait
	LanguageTypeBEH    LanguageType = "beh"      // Biali
	LanguageTypeBEI    LanguageType = "bei"      // Bekati'
	LanguageTypeBEJ    LanguageType = "bej"      // Beja and Bedawiyet
	LanguageTypeBEK    LanguageType = "bek"      // Bebeli
	LanguageTypeBEM    LanguageType = "bem"      // Bemba (Zambia)
	LanguageTypeBEO    LanguageType = "beo"      // Beami
	LanguageTypeBEP    LanguageType = "bep"      // Besoa
	LanguageTypeBEQ    LanguageType = "beq"      // Beembe
	LanguageTypeBER    LanguageType = "ber"      // Berber languages
	LanguageTypeBES    LanguageType = "bes"      // Besme
	LanguageTypeBET    LanguageType = "bet"      // Guiberoua Béte
	LanguageTypeBEU    LanguageType = "beu"      // Blagar
	LanguageTypeBEV    LanguageType = "bev"      // Daloa Bété
	LanguageTypeBEW    LanguageType = "bew"      // Betawi
	LanguageTypeBEX    LanguageType = "bex"      // Jur Modo
	LanguageTypeBEY    LanguageType = "bey"      // Beli (Papua New Guinea)
	LanguageTypeBEZ    LanguageType = "bez"      // Bena (Tanzania)
	LanguageTypeBFA    LanguageType = "bfa"      // Bari
	LanguageTypeBFB    LanguageType = "bfb"      // Pauri Bareli
	LanguageTypeBFC    LanguageType = "bfc"      // Northern Bai
	LanguageTypeBFD    LanguageType = "bfd"      // Bafut
	LanguageTypeBFE    LanguageType = "bfe"      // Betaf and Tena
	LanguageTypeBFF    LanguageType = "bff"      // Bofi
	LanguageTypeBFG    LanguageType = "bfg"      // Busang Kayan
	LanguageTypeBFH    LanguageType = "bfh"      // Blafe
	LanguageTypeBFI    LanguageType = "bfi"      // British Sign Language
	LanguageTypeBFJ    LanguageType = "bfj"      // Bafanji
	LanguageTypeBFK    LanguageType = "bfk"      // Ban Khor Sign Language
	LanguageTypeBFL    LanguageType = "bfl"      // Banda-Ndélé
	LanguageTypeBFM    LanguageType = "bfm"      // Mmen
	LanguageTypeBFN    LanguageType = "bfn"      // Bunak
	LanguageTypeBFO    LanguageType = "bfo"      // Malba Birifor
	LanguageTypeBFP    LanguageType = "bfp"      // Beba
	LanguageTypeBFQ    LanguageType = "bfq"      // Badaga
	LanguageTypeBFR    LanguageType = "bfr"      // Bazigar
	LanguageTypeBFS    LanguageType = "bfs"      // Southern Bai
	LanguageTypeBFT    LanguageType = "bft"      // Balti
	LanguageTypeBFU    LanguageType = "bfu"      // Gahri
	LanguageTypeBFW    LanguageType = "bfw"      // Bondo
	LanguageTypeBFX    LanguageType = "bfx"      // Bantayanon
	LanguageTypeBFY    LanguageType = "bfy"      // Bagheli
	LanguageTypeBFZ    LanguageType = "bfz"      // Mahasu Pahari
	LanguageTypeBGA    LanguageType = "bga"      // Gwamhi-Wuri
	LanguageTypeBGB    LanguageType = "bgb"      // Bobongko
	LanguageTypeBGC    LanguageType = "bgc"      // Haryanvi
	LanguageTypeBGD    LanguageType = "bgd"      // Rathwi Bareli
	LanguageTypeBGE    LanguageType = "bge"      // Bauria
	LanguageTypeBGF    LanguageType = "bgf"      // Bangandu
	LanguageTypeBGG    LanguageType = "bgg"      // Bugun
	LanguageTypeBGI    LanguageType = "bgi"      // Giangan
	LanguageTypeBGJ    LanguageType = "bgj"      // Bangolan
	LanguageTypeBGK    LanguageType = "bgk"      // Bit and Buxinhua
	LanguageTypeBGL    LanguageType = "bgl"      // Bo (Laos)
	LanguageTypeBGM    LanguageType = "bgm"      // Baga Mboteni
	LanguageTypeBGN    LanguageType = "bgn"      // Western Balochi
	LanguageTypeBGO    LanguageType = "bgo"      // Baga Koga
	LanguageTypeBGP    LanguageType = "bgp"      // Eastern Balochi
	LanguageTypeBGQ    LanguageType = "bgq"      // Bagri
	LanguageTypeBGR    LanguageType = "bgr"      // Bawm Chin
	LanguageTypeBGS    LanguageType = "bgs"      // Tagabawa
	LanguageTypeBGT    LanguageType = "bgt"      // Bughotu
	LanguageTypeBGU    LanguageType = "bgu"      // Mbongno
	LanguageTypeBGV    LanguageType = "bgv"      // Warkay-Bipim
	LanguageTypeBGW    LanguageType = "bgw"      // Bhatri
	LanguageTypeBGX    LanguageType = "bgx"      // Balkan Gagauz Turkish
	LanguageTypeBGY    LanguageType = "bgy"      // Benggoi
	LanguageTypeBGZ    LanguageType = "bgz"      // Banggai
	LanguageTypeBHA    LanguageType = "bha"      // Bharia
	LanguageTypeBHB    LanguageType = "bhb"      // Bhili
	LanguageTypeBHC    LanguageType = "bhc"      // Biga
	LanguageTypeBHD    LanguageType = "bhd"      // Bhadrawahi
	LanguageTypeBHE    LanguageType = "bhe"      // Bhaya
	LanguageTypeBHF    LanguageType = "bhf"      // Odiai
	LanguageTypeBHG    LanguageType = "bhg"      // Binandere
	LanguageTypeBHH    LanguageType = "bhh"      // Bukharic
	LanguageTypeBHI    LanguageType = "bhi"      // Bhilali
	LanguageTypeBHJ    LanguageType = "bhj"      // Bahing
	LanguageTypeBHK    LanguageType = "bhk"      // Albay Bicolano
	LanguageTypeBHL    LanguageType = "bhl"      // Bimin
	LanguageTypeBHM    LanguageType = "bhm"      // Bathari
	LanguageTypeBHN    LanguageType = "bhn"      // Bohtan Neo-Aramaic
	LanguageTypeBHO    LanguageType = "bho"      // Bhojpuri
	LanguageTypeBHP    LanguageType = "bhp"      // Bima
	LanguageTypeBHQ    LanguageType = "bhq"      // Tukang Besi South
	LanguageTypeBHR    LanguageType = "bhr"      // Bara Malagasy
	LanguageTypeBHS    LanguageType = "bhs"      // Buwal
	LanguageTypeBHT    LanguageType = "bht"      // Bhattiyali
	LanguageTypeBHU    LanguageType = "bhu"      // Bhunjia
	LanguageTypeBHV    LanguageType = "bhv"      // Bahau
	LanguageTypeBHW    LanguageType = "bhw"      // Biak
	LanguageTypeBHX    LanguageType = "bhx"      // Bhalay
	LanguageTypeBHY    LanguageType = "bhy"      // Bhele
	LanguageTypeBHZ    LanguageType = "bhz"      // Bada (Indonesia)
	LanguageTypeBIA    LanguageType = "bia"      // Badimaya
	LanguageTypeBIB    LanguageType = "bib"      // Bissa and Bisa
	LanguageTypeBIC    LanguageType = "bic"      // Bikaru
	LanguageTypeBID    LanguageType = "bid"      // Bidiyo
	LanguageTypeBIE    LanguageType = "bie"      // Bepour
	LanguageTypeBIF    LanguageType = "bif"      // Biafada
	LanguageTypeBIG    LanguageType = "big"      // Biangai
	LanguageTypeBIJ    LanguageType = "bij"      // Vaghat-Ya-Bijim-Legeri
	LanguageTypeBIK    LanguageType = "bik"      // Bikol
	LanguageTypeBIL    LanguageType = "bil"      // Bile
	LanguageTypeBIM    LanguageType = "bim"      // Bimoba
	LanguageTypeBIN    LanguageType = "bin"      // Bini and Edo
	LanguageTypeBIO    LanguageType = "bio"      // Nai
	LanguageTypeBIP    LanguageType = "bip"      // Bila
	LanguageTypeBIQ    LanguageType = "biq"      // Bipi
	LanguageTypeBIR    LanguageType = "bir"      // Bisorio
	LanguageTypeBIT    LanguageType = "bit"      // Berinomo
	LanguageTypeBIU    LanguageType = "biu"      // Biete
	LanguageTypeBIV    LanguageType = "biv"      // Southern Birifor
	LanguageTypeBIW    LanguageType = "biw"      // Kol (Cameroon)
	LanguageTypeBIX    LanguageType = "bix"      // Bijori
	LanguageTypeBIY    LanguageType = "biy"      // Birhor
	LanguageTypeBIZ    LanguageType = "biz"      // Baloi
	LanguageTypeBJA    LanguageType = "bja"      // Budza
	LanguageTypeBJB    LanguageType = "bjb"      // Banggarla
	LanguageTypeBJC    LanguageType = "bjc"      // Bariji
	LanguageTypeBJD    LanguageType = "bjd"      // Bandjigali
	LanguageTypeBJE    LanguageType = "bje"      // Biao-Jiao Mien
	LanguageTypeBJF    LanguageType = "bjf"      // Barzani Jewish Neo-Aramaic
	LanguageTypeBJG    LanguageType = "bjg"      // Bidyogo
	LanguageTypeBJH    LanguageType = "bjh"      // Bahinemo
	LanguageTypeBJI    LanguageType = "bji"      // Burji
	LanguageTypeBJJ    LanguageType = "bjj"      // Kanauji
	LanguageTypeBJK    LanguageType = "bjk"      // Barok
	LanguageTypeBJL    LanguageType = "bjl"      // Bulu (Papua New Guinea)
	LanguageTypeBJM    LanguageType = "bjm"      // Bajelani
	LanguageTypeBJN    LanguageType = "bjn"      // Banjar
	LanguageTypeBJO    LanguageType = "bjo"      // Mid-Southern Banda
	LanguageTypeBJP    LanguageType = "bjp"      // Fanamaket
	LanguageTypeBJQ    LanguageType = "bjq"      // Southern Betsimisaraka Malagasy
	LanguageTypeBJR    LanguageType = "bjr"      // Binumarien
	LanguageTypeBJS    LanguageType = "bjs"      // Bajan
	LanguageTypeBJT    LanguageType = "bjt"      // Balanta-Ganja
	LanguageTypeBJU    LanguageType = "bju"      // Busuu
	LanguageTypeBJV    LanguageType = "bjv"      // Bedjond
	LanguageTypeBJW    LanguageType = "bjw"      // Bakwé
	LanguageTypeBJX    LanguageType = "bjx"      // Banao Itneg
	LanguageTypeBJY    LanguageType = "bjy"      // Bayali
	LanguageTypeBJZ    LanguageType = "bjz"      // Baruga
	LanguageTypeBKA    LanguageType = "bka"      // Kyak
	LanguageTypeBKB    LanguageType = "bkb"      // Finallig
	LanguageTypeBKC    LanguageType = "bkc"      // Baka (Cameroon)
	LanguageTypeBKD    LanguageType = "bkd"      // Binukid and Talaandig
	LanguageTypeBKF    LanguageType = "bkf"      // Beeke
	LanguageTypeBKG    LanguageType = "bkg"      // Buraka
	LanguageTypeBKH    LanguageType = "bkh"      // Bakoko
	LanguageTypeBKI    LanguageType = "bki"      // Baki
	LanguageTypeBKJ    LanguageType = "bkj"      // Pande
	LanguageTypeBKK    LanguageType = "bkk"      // Brokskat
	LanguageTypeBKL    LanguageType = "bkl"      // Berik
	LanguageTypeBKM    LanguageType = "bkm"      // Kom (Cameroon)
	LanguageTypeBKN    LanguageType = "bkn"      // Bukitan
	LanguageTypeBKO    LanguageType = "bko"      // Kwa'
	LanguageTypeBKP    LanguageType = "bkp"      // Boko (Democratic Republic of Congo)
	LanguageTypeBKQ    LanguageType = "bkq"      // Bakairí
	LanguageTypeBKR    LanguageType = "bkr"      // Bakumpai
	LanguageTypeBKS    LanguageType = "bks"      // Northern Sorsoganon
	LanguageTypeBKT    LanguageType = "bkt"      // Boloki
	LanguageTypeBKU    LanguageType = "bku"      // Buhid
	LanguageTypeBKV    LanguageType = "bkv"      // Bekwarra
	LanguageTypeBKW    LanguageType = "bkw"      // Bekwel
	LanguageTypeBKX    LanguageType = "bkx"      // Baikeno
	LanguageTypeBKY    LanguageType = "bky"      // Bokyi
	LanguageTypeBKZ    LanguageType = "bkz"      // Bungku
	LanguageTypeBLA    LanguageType = "bla"      // Siksika
	LanguageTypeBLB    LanguageType = "blb"      // Bilua
	LanguageTypeBLC    LanguageType = "blc"      // Bella Coola
	LanguageTypeBLD    LanguageType = "bld"      // Bolango
	LanguageTypeBLE    LanguageType = "ble"      // Balanta-Kentohe
	LanguageTypeBLF    LanguageType = "blf"      // Buol
	LanguageTypeBLG    LanguageType = "blg"      // Balau
	LanguageTypeBLH    LanguageType = "blh"      // Kuwaa
	LanguageTypeBLI    LanguageType = "bli"      // Bolia
	LanguageTypeBLJ    LanguageType = "blj"      // Bolongan
	LanguageTypeBLK    LanguageType = "blk"      // Pa'o Karen and Pa'O
	LanguageTypeBLL    LanguageType = "bll"      // Biloxi
	LanguageTypeBLM    LanguageType = "blm"      // Beli (Sudan)
	LanguageTypeBLN    LanguageType = "bln"      // Southern Catanduanes Bikol
	LanguageTypeBLO    LanguageType = "blo"      // Anii
	LanguageTypeBLP    LanguageType = "blp"      // Blablanga
	LanguageTypeBLQ    LanguageType = "blq"      // Baluan-Pam
	LanguageTypeBLR    LanguageType = "blr"      // Blang
	LanguageTypeBLS    LanguageType = "bls"      // Balaesang
	LanguageTypeBLT    LanguageType = "blt"      // Tai Dam
	LanguageTypeBLV    LanguageType = "blv"      // Bolo
	LanguageTypeBLW    LanguageType = "blw"      // Balangao
	LanguageTypeBLX    LanguageType = "blx"      // Mag-Indi Ayta
	LanguageTypeBLY    LanguageType = "bly"      // Notre
	LanguageTypeBLZ    LanguageType = "blz"      // Balantak
	LanguageTypeBMA    LanguageType = "bma"      // Lame
	LanguageTypeBMB    LanguageType = "bmb"      // Bembe
	LanguageTypeBMC    LanguageType = "bmc"      // Biem
	LanguageTypeBMD    LanguageType = "bmd"      // Baga Manduri
	LanguageTypeBME    LanguageType = "bme"      // Limassa
	LanguageTypeBMF    LanguageType = "bmf"      // Bom
	LanguageTypeBMG    LanguageType = "bmg"      // Bamwe
	LanguageTypeBMH    LanguageType = "bmh"      // Kein
	LanguageTypeBMI    LanguageType = "bmi"      // Bagirmi
	LanguageTypeBMJ    LanguageType = "bmj"      // Bote-Majhi
	LanguageTypeBMK    LanguageType = "bmk"      // Ghayavi
	LanguageTypeBML    LanguageType = "bml"      // Bomboli
	LanguageTypeBMM    LanguageType = "bmm"      // Northern Betsimisaraka Malagasy
	LanguageTypeBMN    LanguageType = "bmn"      // Bina (Papua New Guinea)
	LanguageTypeBMO    LanguageType = "bmo"      // Bambalang
	LanguageTypeBMP    LanguageType = "bmp"      // Bulgebi
	LanguageTypeBMQ    LanguageType = "bmq"      // Bomu
	LanguageTypeBMR    LanguageType = "bmr"      // Muinane
	LanguageTypeBMS    LanguageType = "bms"      // Bilma Kanuri
	LanguageTypeBMT    LanguageType = "bmt"      // Biao Mon
	LanguageTypeBMU    LanguageType = "bmu"      // Somba-Siawari
	LanguageTypeBMV    LanguageType = "bmv"      // Bum
	LanguageTypeBMW    LanguageType = "bmw"      // Bomwali
	LanguageTypeBMX    LanguageType = "bmx"      // Baimak
	LanguageTypeBMY    LanguageType = "bmy"      // Bemba (Democratic Republic of Congo)
	LanguageTypeBMZ    LanguageType = "bmz"      // Baramu
	LanguageTypeBNA    LanguageType = "bna"      // Bonerate
	LanguageTypeBNB    LanguageType = "bnb"      // Bookan
	LanguageTypeBNC    LanguageType = "bnc"      // Bontok
	LanguageTypeBND    LanguageType = "bnd"      // Banda (Indonesia)
	LanguageTypeBNE    LanguageType = "bne"      // Bintauna
	LanguageTypeBNF    LanguageType = "bnf"      // Masiwang
	LanguageTypeBNG    LanguageType = "bng"      // Benga
	LanguageTypeBNI    LanguageType = "bni"      // Bangi
	LanguageTypeBNJ    LanguageType = "bnj"      // Eastern Tawbuid
	LanguageTypeBNK    LanguageType = "bnk"      // Bierebo
	LanguageTypeBNL    LanguageType = "bnl"      // Boon
	LanguageTypeBNM    LanguageType = "bnm"      // Batanga
	LanguageTypeBNN    LanguageType = "bnn"      // Bunun
	LanguageTypeBNO    LanguageType = "bno"      // Bantoanon
	LanguageTypeBNP    LanguageType = "bnp"      // Bola
	LanguageTypeBNQ    LanguageType = "bnq"      // Bantik
	LanguageTypeBNR    LanguageType = "bnr"      // Butmas-Tur
	LanguageTypeBNS    LanguageType = "bns"      // Bundeli
	LanguageTypeBNT    LanguageType = "bnt"      // Bantu languages
	LanguageTypeBNU    LanguageType = "bnu"      // Bentong
	LanguageTypeBNV    LanguageType = "bnv"      // Bonerif and Beneraf and Edwas
	LanguageTypeBNW    LanguageType = "bnw"      // Bisis
	LanguageTypeBNX    LanguageType = "bnx"      // Bangubangu
	LanguageTypeBNY    LanguageType = "bny"      // Bintulu
	LanguageTypeBNZ    LanguageType = "bnz"      // Beezen
	LanguageTypeBOA    LanguageType = "boa"      // Bora
	LanguageTypeBOB    LanguageType = "bob"      // Aweer
	LanguageTypeBOE    LanguageType = "boe"      // Mundabli
	LanguageTypeBOF    LanguageType = "bof"      // Bolon
	LanguageTypeBOG    LanguageType = "bog"      // Bamako Sign Language
	LanguageTypeBOH    LanguageType = "boh"      // Boma
	LanguageTypeBOI    LanguageType = "boi"      // Barbareño
	LanguageTypeBOJ    LanguageType = "boj"      // Anjam
	LanguageTypeBOK    LanguageType = "bok"      // Bonjo
	LanguageTypeBOL    LanguageType = "bol"      // Bole
	LanguageTypeBOM    LanguageType = "bom"      // Berom
	LanguageTypeBON    LanguageType = "bon"      // Bine
	LanguageTypeBOO    LanguageType = "boo"      // Tiemacèwè Bozo
	LanguageTypeBOP    LanguageType = "bop"      // Bonkiman
	LanguageTypeBOQ    LanguageType = "boq"      // Bogaya
	LanguageTypeBOR    LanguageType = "bor"      // Borôro
	LanguageTypeBOT    LanguageType = "bot"      // Bongo
	LanguageTypeBOU    LanguageType = "bou"      // Bondei
	LanguageTypeBOV    LanguageType = "bov"      // Tuwuli
	LanguageTypeBOW    LanguageType = "bow"      // Rema
	LanguageTypeBOX    LanguageType = "box"      // Buamu
	LanguageTypeBOY    LanguageType = "boy"      // Bodo (Central African Republic)
	LanguageTypeBOZ    LanguageType = "boz"      // Tiéyaxo Bozo
	LanguageTypeBPA    LanguageType = "bpa"      // Daakaka
	LanguageTypeBPB    LanguageType = "bpb"      // Barbacoas
	LanguageTypeBPD    LanguageType = "bpd"      // Banda-Banda
	LanguageTypeBPG    LanguageType = "bpg"      // Bonggo
	LanguageTypeBPH    LanguageType = "bph"      // Botlikh
	LanguageTypeBPI    LanguageType = "bpi"      // Bagupi
	LanguageTypeBPJ    LanguageType = "bpj"      // Binji
	LanguageTypeBPK    LanguageType = "bpk"      // Orowe and 'Ôrôê
	LanguageTypeBPL    LanguageType = "bpl"      // Broome Pearling Lugger Pidgin
	LanguageTypeBPM    LanguageType = "bpm"      // Biyom
	LanguageTypeBPN    LanguageType = "bpn"      // Dzao Min
	LanguageTypeBPO    LanguageType = "bpo"      // Anasi
	LanguageTypeBPP    LanguageType = "bpp"      // Kaure
	LanguageTypeBPQ    LanguageType = "bpq"      // Banda Malay
	LanguageTypeBPR    LanguageType = "bpr"      // Koronadal Blaan
	LanguageTypeBPS    LanguageType = "bps"      // Sarangani Blaan
	LanguageTypeBPT    LanguageType = "bpt"      // Barrow Point
	LanguageTypeBPU    LanguageType = "bpu"      // Bongu
	LanguageTypeBPV    LanguageType = "bpv"      // Bian Marind
	LanguageTypeBPW    LanguageType = "bpw"      // Bo (Papua New Guinea)
	LanguageTypeBPX    LanguageType = "bpx"      // Palya Bareli
	LanguageTypeBPY    LanguageType = "bpy"      // Bishnupriya
	LanguageTypeBPZ    LanguageType = "bpz"      // Bilba
	LanguageTypeBQA    LanguageType = "bqa"      // Tchumbuli
	LanguageTypeBQB    LanguageType = "bqb"      // Bagusa
	LanguageTypeBQC    LanguageType = "bqc"      // Boko (Benin) and Boo
	LanguageTypeBQD    LanguageType = "bqd"      // Bung
	LanguageTypeBQF    LanguageType = "bqf"      // Baga Kaloum
	LanguageTypeBQG    LanguageType = "bqg"      // Bago-Kusuntu
	LanguageTypeBQH    LanguageType = "bqh"      // Baima
	LanguageTypeBQI    LanguageType = "bqi"      // Bakhtiari
	LanguageTypeBQJ    LanguageType = "bqj"      // Bandial
	LanguageTypeBQK    LanguageType = "bqk"      // Banda-Mbrès
	LanguageTypeBQL    LanguageType = "bql"      // Bilakura
	LanguageTypeBQM    LanguageType = "bqm"      // Wumboko
	LanguageTypeBQN    LanguageType = "bqn"      // Bulgarian Sign Language
	LanguageTypeBQO    LanguageType = "bqo"      // Balo
	LanguageTypeBQP    LanguageType = "bqp"      // Busa
	LanguageTypeBQQ    LanguageType = "bqq"      // Biritai
	LanguageTypeBQR    LanguageType = "bqr"      // Burusu
	LanguageTypeBQS    LanguageType = "bqs"      // Bosngun
	LanguageTypeBQT    LanguageType = "bqt"      // Bamukumbit
	LanguageTypeBQU    LanguageType = "bqu"      // Boguru
	LanguageTypeBQV    LanguageType = "bqv"      // Koro Wachi and Begbere-Ejar
	LanguageTypeBQW    LanguageType = "bqw"      // Buru (Nigeria)
	LanguageTypeBQX    LanguageType = "bqx"      // Baangi
	LanguageTypeBQY    LanguageType = "bqy"      // Bengkala Sign Language
	LanguageTypeBQZ    LanguageType = "bqz"      // Bakaka
	LanguageTypeBRA    LanguageType = "bra"      // Braj
	LanguageTypeBRB    LanguageType = "brb"      // Lave
	LanguageTypeBRC    LanguageType = "brc"      // Berbice Creole Dutch
	LanguageTypeBRD    LanguageType = "brd"      // Baraamu
	LanguageTypeBRF    LanguageType = "brf"      // Bera
	LanguageTypeBRG    LanguageType = "brg"      // Baure
	LanguageTypeBRH    LanguageType = "brh"      // Brahui
	LanguageTypeBRI    LanguageType = "bri"      // Mokpwe
	LanguageTypeBRJ    LanguageType = "brj"      // Bieria
	LanguageTypeBRK    LanguageType = "brk"      // Birked
	LanguageTypeBRL    LanguageType = "brl"      // Birwa
	LanguageTypeBRM    LanguageType = "brm"      // Barambu
	LanguageTypeBRN    LanguageType = "brn"      // Boruca
	LanguageTypeBRO    LanguageType = "bro"      // Brokkat
	LanguageTypeBRP    LanguageType = "brp"      // Barapasi
	LanguageTypeBRQ    LanguageType = "brq"      // Breri
	LanguageTypeBRR    LanguageType = "brr"      // Birao
	LanguageTypeBRS    LanguageType = "brs"      // Baras
	LanguageTypeBRT    LanguageType = "brt"      // Bitare
	LanguageTypeBRU    LanguageType = "bru"      // Eastern Bru
	LanguageTypeBRV    LanguageType = "brv"      // Western Bru
	LanguageTypeBRW    LanguageType = "brw"      // Bellari
	LanguageTypeBRX    LanguageType = "brx"      // Bodo (India)
	LanguageTypeBRY    LanguageType = "bry"      // Burui
	LanguageTypeBRZ    LanguageType = "brz"      // Bilbil
	LanguageTypeBSA    LanguageType = "bsa"      // Abinomn
	LanguageTypeBSB    LanguageType = "bsb"      // Brunei Bisaya
	LanguageTypeBSC    LanguageType = "bsc"      // Bassari and Oniyan
	LanguageTypeBSE    LanguageType = "bse"      // Wushi
	LanguageTypeBSF    LanguageType = "bsf"      // Bauchi
	LanguageTypeBSG    LanguageType = "bsg"      // Bashkardi
	LanguageTypeBSH    LanguageType = "bsh"      // Kati
	LanguageTypeBSI    LanguageType = "bsi"      // Bassossi
	LanguageTypeBSJ    LanguageType = "bsj"      // Bangwinji
	LanguageTypeBSK    LanguageType = "bsk"      // Burushaski
	LanguageTypeBSL    LanguageType = "bsl"      // Basa-Gumna
	LanguageTypeBSM    LanguageType = "bsm"      // Busami
	LanguageTypeBSN    LanguageType = "bsn"      // Barasana-Eduria
	LanguageTypeBSO    LanguageType = "bso"      // Buso
	LanguageTypeBSP    LanguageType = "bsp"      // Baga Sitemu
	LanguageTypeBSQ    LanguageType = "bsq"      // Bassa
	LanguageTypeBSR    LanguageType = "bsr"      // Bassa-Kontagora
	LanguageTypeBSS    LanguageType = "bss"      // Akoose
	LanguageTypeBST    LanguageType = "bst"      // Basketo
	LanguageTypeBSU    LanguageType = "bsu"      // Bahonsuai
	LanguageTypeBSV    LanguageType = "bsv"      // Baga Sobané
	LanguageTypeBSW    LanguageType = "bsw"      // Baiso
	LanguageTypeBSX    LanguageType = "bsx"      // Yangkam
	LanguageTypeBSY    LanguageType = "bsy"      // Sabah Bisaya
	LanguageTypeBTA    LanguageType = "bta"      // Bata
	LanguageTypeBTB    LanguageType = "btb"      // Beti (Cameroon)
	LanguageTypeBTC    LanguageType = "btc"      // Bati (Cameroon)
	LanguageTypeBTD    LanguageType = "btd"      // Batak Dairi
	LanguageTypeBTE    LanguageType = "bte"      // Gamo-Ningi
	LanguageTypeBTF    LanguageType = "btf"      // Birgit
	LanguageTypeBTG    LanguageType = "btg"      // Gagnoa Bété
	LanguageTypeBTH    LanguageType = "bth"      // Biatah Bidayuh
	LanguageTypeBTI    LanguageType = "bti"      // Burate
	LanguageTypeBTJ    LanguageType = "btj"      // Bacanese Malay
	LanguageTypeBTK    LanguageType = "btk"      // Batak languages
	LanguageTypeBTL    LanguageType = "btl"      // Bhatola
	LanguageTypeBTM    LanguageType = "btm"      // Batak Mandailing
	LanguageTypeBTN    LanguageType = "btn"      // Ratagnon
	LanguageTypeBTO    LanguageType = "bto"      // Rinconada Bikol
	LanguageTypeBTP    LanguageType = "btp"      // Budibud
	LanguageTypeBTQ    LanguageType = "btq"      // Batek
	LanguageTypeBTR    LanguageType = "btr"      // Baetora
	LanguageTypeBTS    LanguageType = "bts"      // Batak Simalungun
	LanguageTypeBTT    LanguageType = "btt"      // Bete-Bendi
	LanguageTypeBTU    LanguageType = "btu"      // Batu
	LanguageTypeBTV    LanguageType = "btv"      // Bateri
	LanguageTypeBTW    LanguageType = "btw"      // Butuanon
	LanguageTypeBTX    LanguageType = "btx"      // Batak Karo
	LanguageTypeBTY    LanguageType = "bty"      // Bobot
	LanguageTypeBTZ    LanguageType = "btz"      // Batak Alas-Kluet
	LanguageTypeBUA    LanguageType = "bua"      // Buriat
	LanguageTypeBUB    LanguageType = "bub"      // Bua
	LanguageTypeBUC    LanguageType = "buc"      // Bushi
	LanguageTypeBUD    LanguageType = "bud"      // Ntcham
	LanguageTypeBUE    LanguageType = "bue"      // Beothuk
	LanguageTypeBUF    LanguageType = "buf"      // Bushoong
	LanguageTypeBUG    LanguageType = "bug"      // Buginese
	LanguageTypeBUH    LanguageType = "buh"      // Younuo Bunu
	LanguageTypeBUI    LanguageType = "bui"      // Bongili
	LanguageTypeBUJ    LanguageType = "buj"      // Basa-Gurmana
	LanguageTypeBUK    LanguageType = "buk"      // Bugawac
	LanguageTypeBUM    LanguageType = "bum"      // Bulu (Cameroon)
	LanguageTypeBUN    LanguageType = "bun"      // Sherbro
	LanguageTypeBUO    LanguageType = "buo"      // Terei
	LanguageTypeBUP    LanguageType = "bup"      // Busoa
	LanguageTypeBUQ    LanguageType = "buq"      // Brem
	LanguageTypeBUS    LanguageType = "bus"      // Bokobaru
	LanguageTypeBUT    LanguageType = "but"      // Bungain
	LanguageTypeBUU    LanguageType = "buu"      // Budu
	LanguageTypeBUV    LanguageType = "buv"      // Bun
	LanguageTypeBUW    LanguageType = "buw"      // Bubi
	LanguageTypeBUX    LanguageType = "bux"      // Boghom
	LanguageTypeBUY    LanguageType = "buy"      // Bullom So
	LanguageTypeBUZ    LanguageType = "buz"      // Bukwen
	LanguageTypeBVA    LanguageType = "bva"      // Barein
	LanguageTypeBVB    LanguageType = "bvb"      // Bube
	LanguageTypeBVC    LanguageType = "bvc"      // Baelelea
	LanguageTypeBVD    LanguageType = "bvd"      // Baeggu
	LanguageTypeBVE    LanguageType = "bve"      // Berau Malay
	LanguageTypeBVF    LanguageType = "bvf"      // Boor
	LanguageTypeBVG    LanguageType = "bvg"      // Bonkeng
	LanguageTypeBVH    LanguageType = "bvh"      // Bure
	LanguageTypeBVI    LanguageType = "bvi"      // Belanda Viri
	LanguageTypeBVJ    LanguageType = "bvj"      // Baan
	LanguageTypeBVK    LanguageType = "bvk"      // Bukat
	LanguageTypeBVL    LanguageType = "bvl"      // Bolivian Sign Language
	LanguageTypeBVM    LanguageType = "bvm"      // Bamunka
	LanguageTypeBVN    LanguageType = "bvn"      // Buna
	LanguageTypeBVO    LanguageType = "bvo"      // Bolgo
	LanguageTypeBVP    LanguageType = "bvp"      // Bumang
	LanguageTypeBVQ    LanguageType = "bvq"      // Birri
	LanguageTypeBVR    LanguageType = "bvr"      // Burarra
	LanguageTypeBVT    LanguageType = "bvt"      // Bati (Indonesia)
	LanguageTypeBVU    LanguageType = "bvu"      // Bukit Malay
	LanguageTypeBVV    LanguageType = "bvv"      // Baniva
	LanguageTypeBVW    LanguageType = "bvw"      // Boga
	LanguageTypeBVX    LanguageType = "bvx"      // Dibole
	LanguageTypeBVY    LanguageType = "bvy"      // Baybayanon
	LanguageTypeBVZ    LanguageType = "bvz"      // Bauzi
	LanguageTypeBWA    LanguageType = "bwa"      // Bwatoo
	LanguageTypeBWB    LanguageType = "bwb"      // Namosi-Naitasiri-Serua
	LanguageTypeBWC    LanguageType = "bwc"      // Bwile
	LanguageTypeBWD    LanguageType = "bwd"      // Bwaidoka
	LanguageTypeBWE    LanguageType = "bwe"      // Bwe Karen
	LanguageTypeBWF    LanguageType = "bwf"      // Boselewa
	LanguageTypeBWG    LanguageType = "bwg"      // Barwe
	LanguageTypeBWH    LanguageType = "bwh"      // Bishuo
	LanguageTypeBWI    LanguageType = "bwi"      // Baniwa
	LanguageTypeBWJ    LanguageType = "bwj"      // Láá Láá Bwamu
	LanguageTypeBWK    LanguageType = "bwk"      // Bauwaki
	LanguageTypeBWL    LanguageType = "bwl"      // Bwela
	LanguageTypeBWM    LanguageType = "bwm"      // Biwat
	LanguageTypeBWN    LanguageType = "bwn"      // Wunai Bunu
	LanguageTypeBWO    LanguageType = "bwo"      // Boro (Ethiopia) and Borna (Ethiopia)
	LanguageTypeBWP    LanguageType = "bwp"      // Mandobo Bawah
	LanguageTypeBWQ    LanguageType = "bwq"      // Southern Bobo Madaré
	LanguageTypeBWR    LanguageType = "bwr"      // Bura-Pabir
	LanguageTypeBWS    LanguageType = "bws"      // Bomboma
	LanguageTypeBWT    LanguageType = "bwt"      // Bafaw-Balong
	LanguageTypeBWU    LanguageType = "bwu"      // Buli (Ghana)
	LanguageTypeBWW    LanguageType = "bww"      // Bwa
	LanguageTypeBWX    LanguageType = "bwx"      // Bu-Nao Bunu
	LanguageTypeBWY    LanguageType = "bwy"      // Cwi Bwamu
	LanguageTypeBWZ    LanguageType = "bwz"      // Bwisi
	LanguageTypeBXA    LanguageType = "bxa"      // Tairaha
	LanguageTypeBXB    LanguageType = "bxb"      // Belanda Bor
	LanguageTypeBXC    LanguageType = "bxc"      // Molengue
	LanguageTypeBXD    LanguageType = "bxd"      // Pela
	LanguageTypeBXE    LanguageType = "bxe"      // Birale
	LanguageTypeBXF    LanguageType = "bxf"      // Bilur and Minigir
	LanguageTypeBXG    LanguageType = "bxg"      // Bangala
	LanguageTypeBXH    LanguageType = "bxh"      // Buhutu
	LanguageTypeBXI    LanguageType = "bxi"      // Pirlatapa
	LanguageTypeBXJ    LanguageType = "bxj"      // Bayungu
	LanguageTypeBXK    LanguageType = "bxk"      // Bukusu and Lubukusu
	LanguageTypeBXL    LanguageType = "bxl"      // Jalkunan
	LanguageTypeBXM    LanguageType = "bxm"      // Mongolia Buriat
	LanguageTypeBXN    LanguageType = "bxn"      // Burduna
	LanguageTypeBXO    LanguageType = "bxo"      // Barikanchi
	LanguageTypeBXP    LanguageType = "bxp"      // Bebil
	LanguageTypeBXQ    LanguageType = "bxq"      // Beele
	LanguageTypeBXR    LanguageType = "bxr"      // Russia Buriat
	LanguageTypeBXS    LanguageType = "bxs"      // Busam
	LanguageTypeBXU    LanguageType = "bxu"      // China Buriat
	LanguageTypeBXV    LanguageType = "bxv"      // Berakou
	LanguageTypeBXW    LanguageType = "bxw"      // Bankagooma
	LanguageTypeBXX    LanguageType = "bxx"      // Borna (Democratic Republic of Congo)
	LanguageTypeBXZ    LanguageType = "bxz"      // Binahari
	LanguageTypeBYA    LanguageType = "bya"      // Batak
	LanguageTypeBYB    LanguageType = "byb"      // Bikya
	LanguageTypeBYC    LanguageType = "byc"      // Ubaghara
	LanguageTypeBYD    LanguageType = "byd"      // Benyadu'
	LanguageTypeBYE    LanguageType = "bye"      // Pouye
	LanguageTypeBYF    LanguageType = "byf"      // Bete
	LanguageTypeBYG    LanguageType = "byg"      // Baygo
	LanguageTypeBYH    LanguageType = "byh"      // Bhujel
	LanguageTypeBYI    LanguageType = "byi"      // Buyu
	LanguageTypeBYJ    LanguageType = "byj"      // Bina (Nigeria)
	LanguageTypeBYK    LanguageType = "byk"      // Biao
	LanguageTypeBYL    LanguageType = "byl"      // Bayono
	LanguageTypeBYM    LanguageType = "bym"      // Bidyara
	LanguageTypeBYN    LanguageType = "byn"      // Bilin and Blin
	LanguageTypeBYO    LanguageType = "byo"      // Biyo
	LanguageTypeBYP    LanguageType = "byp"      // Bumaji
	LanguageTypeBYQ    LanguageType = "byq"      // Basay
	LanguageTypeBYR    LanguageType = "byr"      // Baruya and Yipma
	LanguageTypeBYS    LanguageType = "bys"      // Burak
	LanguageTypeBYT    LanguageType = "byt"      // Berti
	LanguageTypeBYV    LanguageType = "byv"      // Medumba
	LanguageTypeBYW    LanguageType = "byw"      // Belhariya
	LanguageTypeBYX    LanguageType = "byx"      // Qaqet
	LanguageTypeBYY    LanguageType = "byy"      // Buya
	LanguageTypeBYZ    LanguageType = "byz"      // Banaro
	LanguageTypeBZA    LanguageType = "bza"      // Bandi
	LanguageTypeBZB    LanguageType = "bzb"      // Andio
	LanguageTypeBZC    LanguageType = "bzc"      // Southern Betsimisaraka Malagasy
	LanguageTypeBZD    LanguageType = "bzd"      // Bribri
	LanguageTypeBZE    LanguageType = "bze"      // Jenaama Bozo
	LanguageTypeBZF    LanguageType = "bzf"      // Boikin
	LanguageTypeBZG    LanguageType = "bzg"      // Babuza
	LanguageTypeBZH    LanguageType = "bzh"      // Mapos Buang
	LanguageTypeBZI    LanguageType = "bzi"      // Bisu
	LanguageTypeBZJ    LanguageType = "bzj"      // Belize Kriol English
	LanguageTypeBZK    LanguageType = "bzk"      // Nicaragua Creole English
	LanguageTypeBZL    LanguageType = "bzl"      // Boano (Sulawesi)
	LanguageTypeBZM    LanguageType = "bzm"      // Bolondo
	LanguageTypeBZN    LanguageType = "bzn"      // Boano (Maluku)
	LanguageTypeBZO    LanguageType = "bzo"      // Bozaba
	LanguageTypeBZP    LanguageType = "bzp"      // Kemberano
	LanguageTypeBZQ    LanguageType = "bzq"      // Buli (Indonesia)
	LanguageTypeBZR    LanguageType = "bzr"      // Biri
	LanguageTypeBZS    LanguageType = "bzs"      // Brazilian Sign Language
	LanguageTypeBZT    LanguageType = "bzt"      // Brithenig
	LanguageTypeBZU    LanguageType = "bzu"      // Burmeso
	LanguageTypeBZV    LanguageType = "bzv"      // Naami
	LanguageTypeBZW    LanguageType = "bzw"      // Basa (Nigeria)
	LanguageTypeBZX    LanguageType = "bzx"      // Kɛlɛngaxo Bozo
	LanguageTypeBZY    LanguageType = "bzy"      // Obanliku
	LanguageTypeBZZ    LanguageType = "bzz"      // Evant
	LanguageTypeCAA    LanguageType = "caa"      // Chortí
	LanguageTypeCAB    LanguageType = "cab"      // Garifuna
	LanguageTypeCAC    LanguageType = "cac"      // Chuj
	LanguageTypeCAD    LanguageType = "cad"      // Caddo
	LanguageTypeCAE    LanguageType = "cae"      // Lehar and Laalaa
	LanguageTypeCAF    LanguageType = "caf"      // Southern Carrier
	LanguageTypeCAG    LanguageType = "cag"      // Nivaclé
	LanguageTypeCAH    LanguageType = "cah"      // Cahuarano
	LanguageTypeCAI    LanguageType = "cai"      // Central American Indian languages
	LanguageTypeCAJ    LanguageType = "caj"      // Chané
	LanguageTypeCAK    LanguageType = "cak"      // Kaqchikel and Cakchiquel
	LanguageTypeCAL    LanguageType = "cal"      // Carolinian
	LanguageTypeCAM    LanguageType = "cam"      // Cemuhî
	LanguageTypeCAN    LanguageType = "can"      // Chambri
	LanguageTypeCAO    LanguageType = "cao"      // Chácobo
	LanguageTypeCAP    LanguageType = "cap"      // Chipaya
	LanguageTypeCAQ    LanguageType = "caq"      // Car Nicobarese
	LanguageTypeCAR    LanguageType = "car"      // Galibi Carib
	LanguageTypeCAS    LanguageType = "cas"      // Tsimané
	LanguageTypeCAU    LanguageType = "cau"      // Caucasian languages
	LanguageTypeCAV    LanguageType = "cav"      // Cavineña
	LanguageTypeCAW    LanguageType = "caw"      // Callawalla
	LanguageTypeCAX    LanguageType = "cax"      // Chiquitano
	LanguageTypeCAY    LanguageType = "cay"      // Cayuga
	LanguageTypeCAZ    LanguageType = "caz"      // Canichana
	LanguageTypeCBA    LanguageType = "cba"      // Chibchan languages
	LanguageTypeCBB    LanguageType = "cbb"      // Cabiyarí
	LanguageTypeCBC    LanguageType = "cbc"      // Carapana
	LanguageTypeCBD    LanguageType = "cbd"      // Carijona
	LanguageTypeCBE    LanguageType = "cbe"      // Chipiajes
	LanguageTypeCBG    LanguageType = "cbg"      // Chimila
	LanguageTypeCBH    LanguageType = "cbh"      // Cagua
	LanguageTypeCBI    LanguageType = "cbi"      // Chachi
	LanguageTypeCBJ    LanguageType = "cbj"      // Ede Cabe
	LanguageTypeCBK    LanguageType = "cbk"      // Chavacano
	LanguageTypeCBL    LanguageType = "cbl"      // Bualkhaw Chin
	LanguageTypeCBN    LanguageType = "cbn"      // Nyahkur
	LanguageTypeCBO    LanguageType = "cbo"      // Izora
	LanguageTypeCBR    LanguageType = "cbr"      // Cashibo-Cacataibo
	LanguageTypeCBS    LanguageType = "cbs"      // Cashinahua
	LanguageTypeCBT    LanguageType = "cbt"      // Chayahuita
	LanguageTypeCBU    LanguageType = "cbu"      // Candoshi-Shapra
	LanguageTypeCBV    LanguageType = "cbv"      // Cacua
	LanguageTypeCBW    LanguageType = "cbw"      // Kinabalian
	LanguageTypeCBY    LanguageType = "cby"      // Carabayo
	LanguageTypeCCA    LanguageType = "cca"      // Cauca
	LanguageTypeCCC    LanguageType = "ccc"      // Chamicuro
	LanguageTypeCCD    LanguageType = "ccd"      // Cafundo Creole
	LanguageTypeCCE    LanguageType = "cce"      // Chopi
	LanguageTypeCCG    LanguageType = "ccg"      // Samba Daka
	LanguageTypeCCH    LanguageType = "cch"      // Atsam
	LanguageTypeCCJ    LanguageType = "ccj"      // Kasanga
	LanguageTypeCCL    LanguageType = "ccl"      // Cutchi-Swahili
	LanguageTypeCCM    LanguageType = "ccm"      // Malaccan Creole Malay
	LanguageTypeCCN    LanguageType = "ccn"      // North Caucasian languages
	LanguageTypeCCO    LanguageType = "cco"      // Comaltepec Chinantec
	LanguageTypeCCP    LanguageType = "ccp"      // Chakma
	LanguageTypeCCQ    LanguageType = "ccq"      // Chaungtha
	LanguageTypeCCR    LanguageType = "ccr"      // Cacaopera
	LanguageTypeCCS    LanguageType = "ccs"      // South Caucasian languages
	LanguageTypeCDA    LanguageType = "cda"      // Choni
	LanguageTypeCDC    LanguageType = "cdc"      // Chadic languages
	LanguageTypeCDD    LanguageType = "cdd"      // Caddoan languages
	LanguageTypeCDE    LanguageType = "cde"      // Chenchu
	LanguageTypeCDF    LanguageType = "cdf"      // Chiru
	LanguageTypeCDG    LanguageType = "cdg"      // Chamari
	LanguageTypeCDH    LanguageType = "cdh"      // Chambeali
	LanguageTypeCDI    LanguageType = "cdi"      // Chodri
	LanguageTypeCDJ    LanguageType = "cdj"      // Churahi
	LanguageTypeCDM    LanguageType = "cdm"      // Chepang
	LanguageTypeCDN    LanguageType = "cdn"      // Chaudangsi
	LanguageTypeCDO    LanguageType = "cdo"      // Min Dong Chinese
	LanguageTypeCDR    LanguageType = "cdr"      // Cinda-Regi-Tiyal
	LanguageTypeCDS    LanguageType = "cds"      // Chadian Sign Language
	LanguageTypeCDY    LanguageType = "cdy"      // Chadong
	LanguageTypeCDZ    LanguageType = "cdz"      // Koda
	LanguageTypeCEA    LanguageType = "cea"      // Lower Chehalis
	LanguageTypeCEB    LanguageType = "ceb"      // Cebuano
	LanguageTypeCEG    LanguageType = "ceg"      // Chamacoco
	LanguageTypeCEK    LanguageType = "cek"      // Eastern Khumi Chin
	LanguageTypeCEL    LanguageType = "cel"      // Celtic languages
	LanguageTypeCEN    LanguageType = "cen"      // Cen
	LanguageTypeCET    LanguageType = "cet"      // Centúúm
	LanguageTypeCFA    LanguageType = "cfa"      // Dijim-Bwilim
	LanguageTypeCFD    LanguageType = "cfd"      // Cara
	LanguageTypeCFG    LanguageType = "cfg"      // Como Karim
	LanguageTypeCFM    LanguageType = "cfm"      // Falam Chin
	LanguageTypeCGA    LanguageType = "cga"      // Changriwa
	LanguageTypeCGC    LanguageType = "cgc"      // Kagayanen
	LanguageTypeCGG    LanguageType = "cgg"      // Chiga
	LanguageTypeCGK    LanguageType = "cgk"      // Chocangacakha
	LanguageTypeCHB    LanguageType = "chb"      // Chibcha
	LanguageTypeCHC    LanguageType = "chc"      // Catawba
	LanguageTypeCHD    LanguageType = "chd"      // Highland Oaxaca Chontal
	LanguageTypeCHF    LanguageType = "chf"      // Tabasco Chontal
	LanguageTypeCHG    LanguageType = "chg"      // Chagatai
	LanguageTypeCHH    LanguageType = "chh"      // Chinook
	LanguageTypeCHJ    LanguageType = "chj"      // Ojitlán Chinantec
	LanguageTypeCHK    LanguageType = "chk"      // Chuukese
	LanguageTypeCHL    LanguageType = "chl"      // Cahuilla
	LanguageTypeCHM    LanguageType = "chm"      // Mari (Russia)
	LanguageTypeCHN    LanguageType = "chn"      // Chinook jargon
	LanguageTypeCHO    LanguageType = "cho"      // Choctaw
	LanguageTypeCHP    LanguageType = "chp"      // Chipewyan and Dene Suline
	LanguageTypeCHQ    LanguageType = "chq"      // Quiotepec Chinantec
	LanguageTypeCHR    LanguageType = "chr"      // Cherokee
	LanguageTypeCHT    LanguageType = "cht"      // Cholón
	LanguageTypeCHW    LanguageType = "chw"      // Chuwabu
	LanguageTypeCHX    LanguageType = "chx"      // Chantyal
	LanguageTypeCHY    LanguageType = "chy"      // Cheyenne
	LanguageTypeCHZ    LanguageType = "chz"      // Ozumacín Chinantec
	LanguageTypeCIA    LanguageType = "cia"      // Cia-Cia
	LanguageTypeCIB    LanguageType = "cib"      // Ci Gbe
	LanguageTypeCIC    LanguageType = "cic"      // Chickasaw
	LanguageTypeCID    LanguageType = "cid"      // Chimariko
	LanguageTypeCIE    LanguageType = "cie"      // Cineni
	LanguageTypeCIH    LanguageType = "cih"      // Chinali
	LanguageTypeCIK    LanguageType = "cik"      // Chitkuli Kinnauri
	LanguageTypeCIM    LanguageType = "cim"      // Cimbrian
	LanguageTypeCIN    LanguageType = "cin"      // Cinta Larga
	LanguageTypeCIP    LanguageType = "cip"      // Chiapanec
	LanguageTypeCIR    LanguageType = "cir"      // Tiri and Haméa and Méa
	LanguageTypeCIW    LanguageType = "ciw"      // Chippewa
	LanguageTypeCIY    LanguageType = "ciy"      // Chaima
	LanguageTypeCJA    LanguageType = "cja"      // Western Cham
	LanguageTypeCJE    LanguageType = "cje"      // Chru
	LanguageTypeCJH    LanguageType = "cjh"      // Upper Chehalis
	LanguageTypeCJI    LanguageType = "cji"      // Chamalal
	LanguageTypeCJK    LanguageType = "cjk"      // Chokwe
	LanguageTypeCJM    LanguageType = "cjm"      // Eastern Cham
	LanguageTypeCJN    LanguageType = "cjn"      // Chenapian
	LanguageTypeCJO    LanguageType = "cjo"      // Ashéninka Pajonal
	LanguageTypeCJP    LanguageType = "cjp"      // Cabécar
	LanguageTypeCJR    LanguageType = "cjr"      // Chorotega
	LanguageTypeCJS    LanguageType = "cjs"      // Shor
	LanguageTypeCJV    LanguageType = "cjv"      // Chuave
	LanguageTypeCJY    LanguageType = "cjy"      // Jinyu Chinese
	LanguageTypeCKA    LanguageType = "cka"      // Khumi Awa Chin
	LanguageTypeCKB    LanguageType = "ckb"      // Central Kurdish
	LanguageTypeCKH    LanguageType = "ckh"      // Chak
	LanguageTypeCKL    LanguageType = "ckl"      // Cibak
	LanguageTypeCKN    LanguageType = "ckn"      // Kaang Chin
	LanguageTypeCKO    LanguageType = "cko"      // Anufo
	LanguageTypeCKQ    LanguageType = "ckq"      // Kajakse
	LanguageTypeCKR    LanguageType = "ckr"      // Kairak
	LanguageTypeCKS    LanguageType = "cks"      // Tayo
	LanguageTypeCKT    LanguageType = "ckt"      // Chukot
	LanguageTypeCKU    LanguageType = "cku"      // Koasati
	LanguageTypeCKV    LanguageType = "ckv"      // Kavalan
	LanguageTypeCKX    LanguageType = "ckx"      // Caka
	LanguageTypeCKY    LanguageType = "cky"      // Cakfem-Mushere
	LanguageTypeCKZ    LanguageType = "ckz"      // Cakchiquel-Quiché Mixed Language
	LanguageTypeCLA    LanguageType = "cla"      // Ron
	LanguageTypeCLC    LanguageType = "clc"      // Chilcotin
	LanguageTypeCLD    LanguageType = "cld"      // Chaldean Neo-Aramaic
	LanguageTypeCLE    LanguageType = "cle"      // Lealao Chinantec
	LanguageTypeCLH    LanguageType = "clh"      // Chilisso
	LanguageTypeCLI    LanguageType = "cli"      // Chakali
	LanguageTypeCLJ    LanguageType = "clj"      // Laitu Chin
	LanguageTypeCLK    LanguageType = "clk"      // Idu-Mishmi
	LanguageTypeCLL    LanguageType = "cll"      // Chala
	LanguageTypeCLM    LanguageType = "clm"      // Clallam
	LanguageTypeCLO    LanguageType = "clo"      // Lowland Oaxaca Chontal
	LanguageTypeCLT    LanguageType = "clt"      // Lautu Chin
	LanguageTypeCLU    LanguageType = "clu"      // Caluyanun
	LanguageTypeCLW    LanguageType = "clw"      // Chulym
	LanguageTypeCLY    LanguageType = "cly"      // Eastern Highland Chatino
	LanguageTypeCMA    LanguageType = "cma"      // Maa
	LanguageTypeCMC    LanguageType = "cmc"      // Chamic languages
	LanguageTypeCME    LanguageType = "cme"      // Cerma
	LanguageTypeCMG    LanguageType = "cmg"      // Classical Mongolian
	LanguageTypeCMI    LanguageType = "cmi"      // Emberá-Chamí
	LanguageTypeCMK    LanguageType = "cmk"      // Chimakum
	LanguageTypeCML    LanguageType = "cml"      // Campalagian
	LanguageTypeCMM    LanguageType = "cmm"      // Michigamea
	LanguageTypeCMN    LanguageType = "cmn"      // Mandarin Chinese
	LanguageTypeCMO    LanguageType = "cmo"      // Central Mnong
	LanguageTypeCMR    LanguageType = "cmr"      // Mro-Khimi Chin
	LanguageTypeCMS    LanguageType = "cms"      // Messapic
	LanguageTypeCMT    LanguageType = "cmt"      // Camtho
	LanguageTypeCNA    LanguageType = "cna"      // Changthang
	LanguageTypeCNB    LanguageType = "cnb"      // Chinbon Chin
	LanguageTypeCNC    LanguageType = "cnc"      // Côông
	LanguageTypeCNG    LanguageType = "cng"      // Northern Qiang
	LanguageTypeCNH    LanguageType = "cnh"      // Haka Chin
	LanguageTypeCNI    LanguageType = "cni"      // Asháninka
	LanguageTypeCNK    LanguageType = "cnk"      // Khumi Chin
	LanguageTypeCNL    LanguageType = "cnl"      // Lalana Chinantec
	LanguageTypeCNO    LanguageType = "cno"      // Con
	LanguageTypeCNS    LanguageType = "cns"      // Central Asmat
	LanguageTypeCNT    LanguageType = "cnt"      // Tepetotutla Chinantec
	LanguageTypeCNU    LanguageType = "cnu"      // Chenoua
	LanguageTypeCNW    LanguageType = "cnw"      // Ngawn Chin
	LanguageTypeCNX    LanguageType = "cnx"      // Middle Cornish
	LanguageTypeCOA    LanguageType = "coa"      // Cocos Islands Malay
	LanguageTypeCOB    LanguageType = "cob"      // Chicomuceltec
	LanguageTypeCOC    LanguageType = "coc"      // Cocopa
	LanguageTypeCOD    LanguageType = "cod"      // Cocama-Cocamilla
	LanguageTypeCOE    LanguageType = "coe"      // Koreguaje
	LanguageTypeCOF    LanguageType = "cof"      // Colorado
	LanguageTypeCOG    LanguageType = "cog"      // Chong
	LanguageTypeCOH    LanguageType = "coh"      // Chonyi-Dzihana-Kauma and Chichonyi-Chidzihana-Chikauma
	LanguageTypeCOJ    LanguageType = "coj"      // Cochimi
	LanguageTypeCOK    LanguageType = "cok"      // Santa Teresa Cora
	LanguageTypeCOL    LanguageType = "col"      // Columbia-Wenatchi
	LanguageTypeCOM    LanguageType = "com"      // Comanche
	LanguageTypeCON    LanguageType = "con"      // Cofán
	LanguageTypeCOO    LanguageType = "coo"      // Comox
	LanguageTypeCOP    LanguageType = "cop"      // Coptic
	LanguageTypeCOQ    LanguageType = "coq"      // Coquille
	LanguageTypeCOT    LanguageType = "cot"      // Caquinte
	LanguageTypeCOU    LanguageType = "cou"      // Wamey
	LanguageTypeCOV    LanguageType = "cov"      // Cao Miao
	LanguageTypeCOW    LanguageType = "cow"      // Cowlitz
	LanguageTypeCOX    LanguageType = "cox"      // Nanti
	LanguageTypeCOY    LanguageType = "coy"      // Coyaima
	LanguageTypeCOZ    LanguageType = "coz"      // Chochotec
	LanguageTypeCPA    LanguageType = "cpa"      // Palantla Chinantec
	LanguageTypeCPB    LanguageType = "cpb"      // Ucayali-Yurúa Ashéninka
	LanguageTypeCPC    LanguageType = "cpc"      // Ajyíninka Apurucayali
	LanguageTypeCPE    LanguageType = "cpe"      // English-based creoles and pidgins
	LanguageTypeCPF    LanguageType = "cpf"      // French-based creoles and pidgins
	LanguageTypeCPG    LanguageType = "cpg"      // Cappadocian Greek
	LanguageTypeCPI    LanguageType = "cpi"      // Chinese Pidgin English
	LanguageTypeCPN    LanguageType = "cpn"      // Cherepon
	LanguageTypeCPO    LanguageType = "cpo"      // Kpeego
	LanguageTypeCPP    LanguageType = "cpp"      // Portuguese-based creoles and pidgins
	LanguageTypeCPS    LanguageType = "cps"      // Capiznon
	LanguageTypeCPU    LanguageType = "cpu"      // Pichis Ashéninka
	LanguageTypeCPX    LanguageType = "cpx"      // Pu-Xian Chinese
	LanguageTypeCPY    LanguageType = "cpy"      // South Ucayali Ashéninka
	LanguageTypeCQD    LanguageType = "cqd"      // Chuanqiandian Cluster Miao
	LanguageTypeCQU    LanguageType = "cqu"      // Chilean Quechua
	LanguageTypeCRA    LanguageType = "cra"      // Chara
	LanguageTypeCRB    LanguageType = "crb"      // Island Carib
	LanguageTypeCRC    LanguageType = "crc"      // Lonwolwol
	LanguageTypeCRD    LanguageType = "crd"      // Coeur d'Alene
	LanguageTypeCRF    LanguageType = "crf"      // Caramanta
	LanguageTypeCRG    LanguageType = "crg"      // Michif
	LanguageTypeCRH    LanguageType = "crh"      // Crimean Tatar and Crimean Turkish
	LanguageTypeCRI    LanguageType = "cri"      // Sãotomense
	LanguageTypeCRJ    LanguageType = "crj"      // Southern East Cree
	LanguageTypeCRK    LanguageType = "crk"      // Plains Cree
	LanguageTypeCRL    LanguageType = "crl"      // Northern East Cree
	LanguageTypeCRM    LanguageType = "crm"      // Moose Cree
	LanguageTypeCRN    LanguageType = "crn"      // El Nayar Cora
	LanguageTypeCRO    LanguageType = "cro"      // Crow
	LanguageTypeCRP    LanguageType = "crp"      // Creoles and pidgins
	LanguageTypeCRQ    LanguageType = "crq"      // Iyo'wujwa Chorote
	LanguageTypeCRR    LanguageType = "crr"      // Carolina Algonquian
	LanguageTypeCRS    LanguageType = "crs"      // Seselwa Creole French
	LanguageTypeCRT    LanguageType = "crt"      // Iyojwa'ja Chorote
	LanguageTypeCRV    LanguageType = "crv"      // Chaura
	LanguageTypeCRW    LanguageType = "crw"      // Chrau
	LanguageTypeCRX    LanguageType = "crx"      // Carrier
	LanguageTypeCRY    LanguageType = "cry"      // Cori
	LanguageTypeCRZ    LanguageType = "crz"      // Cruzeño
	LanguageTypeCSA    LanguageType = "csa"      // Chiltepec Chinantec
	LanguageTypeCSB    LanguageType = "csb"      // Kashubian
	LanguageTypeCSC    LanguageType = "csc"      // Catalan Sign Language and Lengua de señas catalana and Llengua de Signes Catalana
	LanguageTypeCSD    LanguageType = "csd"      // Chiangmai Sign Language
	LanguageTypeCSE    LanguageType = "cse"      // Czech Sign Language
	LanguageTypeCSF    LanguageType = "csf"      // Cuba Sign Language
	LanguageTypeCSG    LanguageType = "csg"      // Chilean Sign Language
	LanguageTypeCSH    LanguageType = "csh"      // Asho Chin
	LanguageTypeCSI    LanguageType = "csi"      // Coast Miwok
	LanguageTypeCSJ    LanguageType = "csj"      // Songlai Chin
	LanguageTypeCSK    LanguageType = "csk"      // Jola-Kasa
	LanguageTypeCSL    LanguageType = "csl"      // Chinese Sign Language
	LanguageTypeCSM    LanguageType = "csm"      // Central Sierra Miwok
	LanguageTypeCSN    LanguageType = "csn"      // Colombian Sign Language
	LanguageTypeCSO    LanguageType = "cso"      // Sochiapam Chinantec and Sochiapan Chinantec
	LanguageTypeCSQ    LanguageType = "csq"      // Croatia Sign Language
	LanguageTypeCSR    LanguageType = "csr"      // Costa Rican Sign Language
	LanguageTypeCSS    LanguageType = "css"      // Southern Ohlone
	LanguageTypeCST    LanguageType = "cst"      // Northern Ohlone
	LanguageTypeCSU    LanguageType = "csu"      // Central Sudanic languages
	LanguageTypeCSV    LanguageType = "csv"      // Sumtu Chin
	LanguageTypeCSW    LanguageType = "csw"      // Swampy Cree
	LanguageTypeCSY    LanguageType = "csy"      // Siyin Chin
	LanguageTypeCSZ    LanguageType = "csz"      // Coos
	LanguageTypeCTA    LanguageType = "cta"      // Tataltepec Chatino
	LanguageTypeCTC    LanguageType = "ctc"      // Chetco
	LanguageTypeCTD    LanguageType = "ctd"      // Tedim Chin
	LanguageTypeCTE    LanguageType = "cte"      // Tepinapa Chinantec
	LanguageTypeCTG    LanguageType = "ctg"      // Chittagonian
	LanguageTypeCTH    LanguageType = "cth"      // Thaiphum Chin
	LanguageTypeCTL    LanguageType = "ctl"      // Tlacoatzintepec Chinantec
	LanguageTypeCTM    LanguageType = "ctm"      // Chitimacha
	LanguageTypeCTN    LanguageType = "ctn"      // Chhintange
	LanguageTypeCTO    LanguageType = "cto"      // Emberá-Catío
	LanguageTypeCTP    LanguageType = "ctp"      // Western Highland Chatino
	LanguageTypeCTS    LanguageType = "cts"      // Northern Catanduanes Bikol
	LanguageTypeCTT    LanguageType = "ctt"      // Wayanad Chetti
	LanguageTypeCTU    LanguageType = "ctu"      // Chol
	LanguageTypeCTZ    LanguageType = "ctz"      // Zacatepec Chatino
	LanguageTypeCUA    LanguageType = "cua"      // Cua
	LanguageTypeCUB    LanguageType = "cub"      // Cubeo
	LanguageTypeCUC    LanguageType = "cuc"      // Usila Chinantec
	LanguageTypeCUG    LanguageType = "cug"      // Cung
	LanguageTypeCUH    LanguageType = "cuh"      // Chuka and Gichuka
	LanguageTypeCUI    LanguageType = "cui"      // Cuiba
	LanguageTypeCUJ    LanguageType = "cuj"      // Mashco Piro
	LanguageTypeCUK    LanguageType = "cuk"      // San Blas Kuna
	LanguageTypeCUL    LanguageType = "cul"      // Culina and Kulina
	LanguageTypeCUM    LanguageType = "cum"      // Cumeral
	LanguageTypeCUO    LanguageType = "cuo"      // Cumanagoto
	LanguageTypeCUP    LanguageType = "cup"      // Cupeño
	LanguageTypeCUQ    LanguageType = "cuq"      // Cun
	LanguageTypeCUR    LanguageType = "cur"      // Chhulung
	LanguageTypeCUS    LanguageType = "cus"      // Cushitic languages
	LanguageTypeCUT    LanguageType = "cut"      // Teutila Cuicatec
	LanguageTypeCUU    LanguageType = "cuu"      // Tai Ya
	LanguageTypeCUV    LanguageType = "cuv"      // Cuvok
	LanguageTypeCUW    LanguageType = "cuw"      // Chukwa
	LanguageTypeCUX    LanguageType = "cux"      // Tepeuxila Cuicatec
	LanguageTypeCVG    LanguageType = "cvg"      // Chug
	LanguageTypeCVN    LanguageType = "cvn"      // Valle Nacional Chinantec
	LanguageTypeCWA    LanguageType = "cwa"      // Kabwa
	LanguageTypeCWB    LanguageType = "cwb"      // Maindo
	LanguageTypeCWD    LanguageType = "cwd"      // Woods Cree
	LanguageTypeCWE    LanguageType = "cwe"      // Kwere
	LanguageTypeCWG    LanguageType = "cwg"      // Chewong and Cheq Wong
	LanguageTypeCWT    LanguageType = "cwt"      // Kuwaataay
	LanguageTypeCYA    LanguageType = "cya"      // Nopala Chatino
	LanguageTypeCYB    LanguageType = "cyb"      // Cayubaba
	LanguageTypeCYO    LanguageType = "cyo"      // Cuyonon
	LanguageTypeCZH    LanguageType = "czh"      // Huizhou Chinese
	LanguageTypeCZK    LanguageType = "czk"      // Knaanic
	LanguageTypeCZN    LanguageType = "czn"      // Zenzontepec Chatino
	LanguageTypeCZO    LanguageType = "czo"      // Min Zhong Chinese
	LanguageTypeCZT    LanguageType = "czt"      // Zotung Chin
	LanguageTypeDAA    LanguageType = "daa"      // Dangaléat
	LanguageTypeDAC    LanguageType = "dac"      // Dambi
	LanguageTypeDAD    LanguageType = "dad"      // Marik
	LanguageTypeDAE    LanguageType = "dae"      // Duupa
	LanguageTypeDAF    LanguageType = "daf"      // Dan
	LanguageTypeDAG    LanguageType = "dag"      // Dagbani
	LanguageTypeDAH    LanguageType = "dah"      // Gwahatike
	LanguageTypeDAI    LanguageType = "dai"      // Day
	LanguageTypeDAJ    LanguageType = "daj"      // Dar Fur Daju
	LanguageTypeDAK    LanguageType = "dak"      // Dakota
	LanguageTypeDAL    LanguageType = "dal"      // Dahalo
	LanguageTypeDAM    LanguageType = "dam"      // Damakawa
	LanguageTypeDAO    LanguageType = "dao"      // Daai Chin
	LanguageTypeDAP    LanguageType = "dap"      // Nisi (India)
	LanguageTypeDAQ    LanguageType = "daq"      // Dandami Maria
	LanguageTypeDAR    LanguageType = "dar"      // Dargwa
	LanguageTypeDAS    LanguageType = "das"      // Daho-Doo
	LanguageTypeDAU    LanguageType = "dau"      // Dar Sila Daju
	LanguageTypeDAV    LanguageType = "dav"      // Taita and Dawida
	LanguageTypeDAW    LanguageType = "daw"      // Davawenyo
	LanguageTypeDAX    LanguageType = "dax"      // Dayi
	LanguageTypeDAY    LanguageType = "day"      // Land Dayak languages
	LanguageTypeDAZ    LanguageType = "daz"      // Dao
	LanguageTypeDBA    LanguageType = "dba"      // Bangime
	LanguageTypeDBB    LanguageType = "dbb"      // Deno
	LanguageTypeDBD    LanguageType = "dbd"      // Dadiya
	LanguageTypeDBE    LanguageType = "dbe"      // Dabe
	LanguageTypeDBF    LanguageType = "dbf"      // Edopi
	LanguageTypeDBG    LanguageType = "dbg"      // Dogul Dom Dogon
	LanguageTypeDBI    LanguageType = "dbi"      // Doka
	LanguageTypeDBJ    LanguageType = "dbj"      // Ida'an
	LanguageTypeDBL    LanguageType = "dbl"      // Dyirbal
	LanguageTypeDBM    LanguageType = "dbm"      // Duguri
	LanguageTypeDBN    LanguageType = "dbn"      // Duriankere
	LanguageTypeDBO    LanguageType = "dbo"      // Dulbu
	LanguageTypeDBP    LanguageType = "dbp"      // Duwai
	LanguageTypeDBQ    LanguageType = "dbq"      // Daba
	LanguageTypeDBR    LanguageType = "dbr"      // Dabarre
	LanguageTypeDBT    LanguageType = "dbt"      // Ben Tey Dogon
	LanguageTypeDBU    LanguageType = "dbu"      // Bondum Dom Dogon
	LanguageTypeDBV    LanguageType = "dbv"      // Dungu
	LanguageTypeDBW    LanguageType = "dbw"      // Bankan Tey Dogon
	LanguageTypeDBY    LanguageType = "dby"      // Dibiyaso
	LanguageTypeDCC    LanguageType = "dcc"      // Deccan
	LanguageTypeDCR    LanguageType = "dcr"      // Negerhollands
	LanguageTypeDDA    LanguageType = "dda"      // Dadi Dadi
	LanguageTypeDDD    LanguageType = "ddd"      // Dongotono
	LanguageTypeDDE    LanguageType = "dde"      // Doondo
	LanguageTypeDDG    LanguageType = "ddg"      // Fataluku
	LanguageTypeDDI    LanguageType = "ddi"      // West Goodenough
	LanguageTypeDDJ    LanguageType = "ddj"      // Jaru
	LanguageTypeDDN    LanguageType = "ddn"      // Dendi (Benin)
	LanguageTypeDDO    LanguageType = "ddo"      // Dido
	LanguageTypeDDR    LanguageType = "ddr"      // Dhudhuroa
	LanguageTypeDDS    LanguageType = "dds"      // Donno So Dogon
	LanguageTypeDDW    LanguageType = "ddw"      // Dawera-Daweloor
	LanguageTypeDEC    LanguageType = "dec"      // Dagik
	LanguageTypeDED    LanguageType = "ded"      // Dedua
	LanguageTypeDEE    LanguageType = "dee"      // Dewoin
	LanguageTypeDEF    LanguageType = "def"      // Dezfuli
	LanguageTypeDEG    LanguageType = "deg"      // Degema
	LanguageTypeDEH    LanguageType = "deh"      // Dehwari
	LanguageTypeDEI    LanguageType = "dei"      // Demisa
	LanguageTypeDEK    LanguageType = "dek"      // Dek
	LanguageTypeDEL    LanguageType = "del"      // Delaware
	LanguageTypeDEM    LanguageType = "dem"      // Dem
	LanguageTypeDEN    LanguageType = "den"      // Slave (Athapascan)
	LanguageTypeDEP    LanguageType = "dep"      // Pidgin Delaware
	LanguageTypeDEQ    LanguageType = "deq"      // Dendi (Central African Republic)
	LanguageTypeDER    LanguageType = "der"      // Deori
	LanguageTypeDES    LanguageType = "des"      // Desano
	LanguageTypeDEV    LanguageType = "dev"      // Domung
	LanguageTypeDEZ    LanguageType = "dez"      // Dengese
	LanguageTypeDGA    LanguageType = "dga"      // Southern Dagaare
	LanguageTypeDGB    LanguageType = "dgb"      // Bunoge Dogon
	LanguageTypeDGC    LanguageType = "dgc"      // Casiguran Dumagat Agta
	LanguageTypeDGD    LanguageType = "dgd"      // Dagaari Dioula
	LanguageTypeDGE    LanguageType = "dge"      // Degenan
	LanguageTypeDGG    LanguageType = "dgg"      // Doga
	LanguageTypeDGH    LanguageType = "dgh"      // Dghwede
	LanguageTypeDGI    LanguageType = "dgi"      // Northern Dagara
	LanguageTypeDGK    LanguageType = "dgk"      // Dagba
	LanguageTypeDGL    LanguageType = "dgl"      // Andaandi and Dongolawi
	LanguageTypeDGN    LanguageType = "dgn"      // Dagoman
	LanguageTypeDGO    LanguageType = "dgo"      // Dogri (individual language)
	LanguageTypeDGR    LanguageType = "dgr"      // Dogrib
	LanguageTypeDGS    LanguageType = "dgs"      // Dogoso
	LanguageTypeDGT    LanguageType = "dgt"      // Ndra'ngith
	LanguageTypeDGU    LanguageType = "dgu"      // Degaru
	LanguageTypeDGW    LanguageType = "dgw"      // Daungwurrung
	LanguageTypeDGX    LanguageType = "dgx"      // Doghoro
	LanguageTypeDGZ    LanguageType = "dgz"      // Daga
	LanguageTypeDHA    LanguageType = "dha"      // Dhanwar (India)
	LanguageTypeDHD    LanguageType = "dhd"      // Dhundari
	LanguageTypeDHG    LanguageType = "dhg"      // Djangu and Dhangu
	LanguageTypeDHI    LanguageType = "dhi"      // Dhimal
	LanguageTypeDHL    LanguageType = "dhl"      // Dhalandji
	LanguageTypeDHM    LanguageType = "dhm"      // Zemba
	LanguageTypeDHN    LanguageType = "dhn"      // Dhanki
	LanguageTypeDHO    LanguageType = "dho"      // Dhodia
	LanguageTypeDHR    LanguageType = "dhr"      // Dhargari
	LanguageTypeDHS    LanguageType = "dhs"      // Dhaiso
	LanguageTypeDHU    LanguageType = "dhu"      // Dhurga
	LanguageTypeDHV    LanguageType = "dhv"      // Dehu and Drehu
	LanguageTypeDHW    LanguageType = "dhw"      // Dhanwar (Nepal)
	LanguageTypeDHX    LanguageType = "dhx"      // Dhungaloo
	LanguageTypeDIA    LanguageType = "dia"      // Dia
	LanguageTypeDIB    LanguageType = "dib"      // South Central Dinka
	LanguageTypeDIC    LanguageType = "dic"      // Lakota Dida
	LanguageTypeDID    LanguageType = "did"      // Didinga
	LanguageTypeDIF    LanguageType = "dif"      // Dieri
	LanguageTypeDIG    LanguageType = "dig"      // Digo and Chidigo
	LanguageTypeDIH    LanguageType = "dih"      // Kumiai
	LanguageTypeDII    LanguageType = "dii"      // Dimbong
	LanguageTypeDIJ    LanguageType = "dij"      // Dai
	LanguageTypeDIK    LanguageType = "dik"      // Southwestern Dinka
	LanguageTypeDIL    LanguageType = "dil"      // Dilling
	LanguageTypeDIM    LanguageType = "dim"      // Dime
	LanguageTypeDIN    LanguageType = "din"      // Dinka
	LanguageTypeDIO    LanguageType = "dio"      // Dibo
	LanguageTypeDIP    LanguageType = "dip"      // Northeastern Dinka
	LanguageTypeDIQ    LanguageType = "diq"      // Dimli (individual language)
	LanguageTypeDIR    LanguageType = "dir"      // Dirim
	LanguageTypeDIS    LanguageType = "dis"      // Dimasa
	LanguageTypeDIT    LanguageType = "dit"      // Dirari
	LanguageTypeDIU    LanguageType = "diu"      // Diriku
	LanguageTypeDIW    LanguageType = "diw"      // Northwestern Dinka
	LanguageTypeDIX    LanguageType = "dix"      // Dixon Reef
	LanguageTypeDIY    LanguageType = "diy"      // Diuwe
	LanguageTypeDIZ    LanguageType = "diz"      // Ding
	LanguageTypeDJA    LanguageType = "dja"      // Djadjawurrung
	LanguageTypeDJB    LanguageType = "djb"      // Djinba
	LanguageTypeDJC    LanguageType = "djc"      // Dar Daju Daju
	LanguageTypeDJD    LanguageType = "djd"      // Djamindjung
	LanguageTypeDJE    LanguageType = "dje"      // Zarma
	LanguageTypeDJF    LanguageType = "djf"      // Djangun
	LanguageTypeDJI    LanguageType = "dji"      // Djinang
	LanguageTypeDJJ    LanguageType = "djj"      // Djeebbana
	LanguageTypeDJK    LanguageType = "djk"      // Eastern Maroon Creole and Businenge Tongo and Nenge
	LanguageTypeDJL    LanguageType = "djl"      // Djiwarli
	LanguageTypeDJM    LanguageType = "djm"      // Jamsay Dogon
	LanguageTypeDJN    LanguageType = "djn"      // Djauan
	LanguageTypeDJO    LanguageType = "djo"      // Jangkang
	LanguageTypeDJR    LanguageType = "djr"      // Djambarrpuyngu
	LanguageTypeDJU    LanguageType = "dju"      // Kapriman
	LanguageTypeDJW    LanguageType = "djw"      // Djawi
	LanguageTypeDKA    LanguageType = "dka"      // Dakpakha
	LanguageTypeDKK    LanguageType = "dkk"      // Dakka
	LanguageTypeDKL    LanguageType = "dkl"      // Kolum So Dogon
	LanguageTypeDKR    LanguageType = "dkr"      // Kuijau
	LanguageTypeDKS    LanguageType = "dks"      // Southeastern Dinka
	LanguageTypeDKX    LanguageType = "dkx"      // Mazagway
	LanguageTypeDLG    LanguageType = "dlg"      // Dolgan
	LanguageTypeDLK    LanguageType = "dlk"      // Dahalik
	LanguageTypeDLM    LanguageType = "dlm"      // Dalmatian
	LanguageTypeDLN    LanguageType = "dln"      // Darlong
	LanguageTypeDMA    LanguageType = "dma"      // Duma
	LanguageTypeDMB    LanguageType = "dmb"      // Mombo Dogon
	LanguageTypeDMC    LanguageType = "dmc"      // Gavak
	LanguageTypeDMD    LanguageType = "dmd"      // Madhi Madhi
	LanguageTypeDME    LanguageType = "dme"      // Dugwor
	LanguageTypeDMG    LanguageType = "dmg"      // Upper Kinabatangan
	LanguageTypeDMK    LanguageType = "dmk"      // Domaaki
	LanguageTypeDML    LanguageType = "dml"      // Dameli
	LanguageTypeDMM    LanguageType = "dmm"      // Dama
	LanguageTypeDMN    LanguageType = "dmn"      // Mande languages
	LanguageTypeDMO    LanguageType = "dmo"      // Kemedzung
	LanguageTypeDMR    LanguageType = "dmr"      // East Damar
	LanguageTypeDMS    LanguageType = "dms"      // Dampelas
	LanguageTypeDMU    LanguageType = "dmu"      // Dubu and Tebi
	LanguageTypeDMV    LanguageType = "dmv"      // Dumpas
	LanguageTypeDMW    LanguageType = "dmw"      // Mudburra
	LanguageTypeDMX    LanguageType = "dmx"      // Dema
	LanguageTypeDMY    LanguageType = "dmy"      // Demta and Sowari
	LanguageTypeDNA    LanguageType = "dna"      // Upper Grand Valley Dani
	LanguageTypeDND    LanguageType = "dnd"      // Daonda
	LanguageTypeDNE    LanguageType = "dne"      // Ndendeule
	LanguageTypeDNG    LanguageType = "dng"      // Dungan
	LanguageTypeDNI    LanguageType = "dni"      // Lower Grand Valley Dani
	LanguageTypeDNJ    LanguageType = "dnj"      // Dan
	LanguageTypeDNK    LanguageType = "dnk"      // Dengka
	LanguageTypeDNN    LanguageType = "dnn"      // Dzùùngoo
	LanguageTypeDNR    LanguageType = "dnr"      // Danaru
	LanguageTypeDNT    LanguageType = "dnt"      // Mid Grand Valley Dani
	LanguageTypeDNU    LanguageType = "dnu"      // Danau
	LanguageTypeDNV    LanguageType = "dnv"      // Danu
	LanguageTypeDNW    LanguageType = "dnw"      // Western Dani
	LanguageTypeDNY    LanguageType = "dny"      // Dení
	LanguageTypeDOA    LanguageType = "doa"      // Dom
	LanguageTypeDOB    LanguageType = "dob"      // Dobu
	LanguageTypeDOC    LanguageType = "doc"      // Northern Dong
	LanguageTypeDOE    LanguageType = "doe"      // Doe
	LanguageTypeDOF    LanguageType = "dof"      // Domu
	LanguageTypeDOH    LanguageType = "doh"      // Dong
	LanguageTypeDOI    LanguageType = "doi"      // Dogri (macrolanguage)
	LanguageTypeDOK    LanguageType = "dok"      // Dondo
	LanguageTypeDOL    LanguageType = "dol"      // Doso
	LanguageTypeDON    LanguageType = "don"      // Toura (Papua New Guinea)
	LanguageTypeDOO    LanguageType = "doo"      // Dongo
	LanguageTypeDOP    LanguageType = "dop"      // Lukpa
	LanguageTypeDOQ    LanguageType = "doq"      // Dominican Sign Language
	LanguageTypeDOR    LanguageType = "dor"      // Dori'o
	LanguageTypeDOS    LanguageType = "dos"      // Dogosé
	LanguageTypeDOT    LanguageType = "dot"      // Dass
	LanguageTypeDOV    LanguageType = "dov"      // Dombe
	LanguageTypeDOW    LanguageType = "dow"      // Doyayo
	LanguageTypeDOX    LanguageType = "dox"      // Bussa
	LanguageTypeDOY    LanguageType = "doy"      // Dompo
	LanguageTypeDOZ    LanguageType = "doz"      // Dorze
	LanguageTypeDPP    LanguageType = "dpp"      // Papar
	LanguageTypeDRA    LanguageType = "dra"      // Dravidian languages
	LanguageTypeDRB    LanguageType = "drb"      // Dair
	LanguageTypeDRC    LanguageType = "drc"      // Minderico
	LanguageTypeDRD    LanguageType = "drd"      // Darmiya
	LanguageTypeDRE    LanguageType = "dre"      // Dolpo
	LanguageTypeDRG    LanguageType = "drg"      // Rungus
	LanguageTypeDRH    LanguageType = "drh"      // Darkhat
	LanguageTypeDRI    LanguageType = "dri"      // C'lela
	LanguageTypeDRL    LanguageType = "drl"      // Paakantyi
	LanguageTypeDRN    LanguageType = "drn"      // West Damar
	LanguageTypeDRO    LanguageType = "dro"      // Daro-Matu Melanau
	LanguageTypeDRQ    LanguageType = "drq"      // Dura
	LanguageTypeDRR    LanguageType = "drr"      // Dororo
	LanguageTypeDRS    LanguageType = "drs"      // Gedeo
	LanguageTypeDRT    LanguageType = "drt"      // Drents
	LanguageTypeDRU    LanguageType = "dru"      // Rukai
	LanguageTypeDRW    LanguageType = "drw"      // Darwazi
	LanguageTypeDRY    LanguageType = "dry"      // Darai
	LanguageTypeDSB    LanguageType = "dsb"      // Lower Sorbian
	LanguageTypeDSE    LanguageType = "dse"      // Dutch Sign Language
	LanguageTypeDSH    LanguageType = "dsh"      // Daasanach
	LanguageTypeDSI    LanguageType = "dsi"      // Disa
	LanguageTypeDSL    LanguageType = "dsl"      // Danish Sign Language
	LanguageTypeDSN    LanguageType = "dsn"      // Dusner
	LanguageTypeDSO    LanguageType = "dso"      // Desiya
	LanguageTypeDSQ    LanguageType = "dsq"      // Tadaksahak
	LanguageTypeDTA    LanguageType = "dta"      // Daur
	LanguageTypeDTB    LanguageType = "dtb"      // Labuk-Kinabatangan Kadazan
	LanguageTypeDTD    LanguageType = "dtd"      // Ditidaht
	LanguageTypeDTH    LanguageType = "dth"      // Adithinngithigh
	LanguageTypeDTI    LanguageType = "dti"      // Ana Tinga Dogon
	LanguageTypeDTK    LanguageType = "dtk"      // Tene Kan Dogon
	LanguageTypeDTM    LanguageType = "dtm"      // Tomo Kan Dogon
	LanguageTypeDTO    LanguageType = "dto"      // Tommo So Dogon
	LanguageTypeDTP    LanguageType = "dtp"      // Central Dusun
	LanguageTypeDTR    LanguageType = "dtr"      // Lotud
	LanguageTypeDTS    LanguageType = "dts"      // Toro So Dogon
	LanguageTypeDTT    LanguageType = "dtt"      // Toro Tegu Dogon
	LanguageTypeDTU    LanguageType = "dtu"      // Tebul Ure Dogon
	LanguageTypeDTY    LanguageType = "dty"      // Dotyali
	LanguageTypeDUA    LanguageType = "dua"      // Duala
	LanguageTypeDUB    LanguageType = "dub"      // Dubli
	LanguageTypeDUC    LanguageType = "duc"      // Duna
	LanguageTypeDUD    LanguageType = "dud"      // Hun-Saare
	LanguageTypeDUE    LanguageType = "due"      // Umiray Dumaget Agta
	LanguageTypeDUF    LanguageType = "duf"      // Dumbea and Drubea
	LanguageTypeDUG    LanguageType = "dug"      // Duruma and Chiduruma
	LanguageTypeDUH    LanguageType = "duh"      // Dungra Bhil
	LanguageTypeDUI    LanguageType = "dui"      // Dumun
	LanguageTypeDUJ    LanguageType = "duj"      // Dhuwal
	LanguageTypeDUK    LanguageType = "duk"      // Uyajitaya
	LanguageTypeDUL    LanguageType = "dul"      // Alabat Island Agta
	LanguageTypeDUM    LanguageType = "dum"      // Middle Dutch (ca. 1050-1350)
	LanguageTypeDUN    LanguageType = "dun"      // Dusun Deyah
	LanguageTypeDUO    LanguageType = "duo"      // Dupaninan Agta
	LanguageTypeDUP    LanguageType = "dup"      // Duano
	LanguageTypeDUQ    LanguageType = "duq"      // Dusun Malang
	LanguageTypeDUR    LanguageType = "dur"      // Dii
	LanguageTypeDUS    LanguageType = "dus"      // Dumi
	LanguageTypeDUU    LanguageType = "duu"      // Drung
	LanguageTypeDUV    LanguageType = "duv"      // Duvle
	LanguageTypeDUW    LanguageType = "duw"      // Dusun Witu
	LanguageTypeDUX    LanguageType = "dux"      // Duungooma
	LanguageTypeDUY    LanguageType = "duy"      // Dicamay Agta
	LanguageTypeDUZ    LanguageType = "duz"      // Duli
	LanguageTypeDVA    LanguageType = "dva"      // Duau
	LanguageTypeDWA    LanguageType = "dwa"      // Diri
	LanguageTypeDWL    LanguageType = "dwl"      // Walo Kumbe Dogon
	LanguageTypeDWR    LanguageType = "dwr"      // Dawro
	LanguageTypeDWS    LanguageType = "dws"      // Dutton World Speedwords
	LanguageTypeDWW    LanguageType = "dww"      // Dawawa
	LanguageTypeDYA    LanguageType = "dya"      // Dyan
	LanguageTypeDYB    LanguageType = "dyb"      // Dyaberdyaber
	LanguageTypeDYD    LanguageType = "dyd"      // Dyugun
	LanguageTypeDYG    LanguageType = "dyg"      // Villa Viciosa Agta
	LanguageTypeDYI    LanguageType = "dyi"      // Djimini Senoufo
	LanguageTypeDYM    LanguageType = "dym"      // Yanda Dom Dogon
	LanguageTypeDYN    LanguageType = "dyn"      // Dyangadi
	LanguageTypeDYO    LanguageType = "dyo"      // Jola-Fonyi
	LanguageTypeDYU    LanguageType = "dyu"      // Dyula
	LanguageTypeDYY    LanguageType = "dyy"      // Dyaabugay
	LanguageTypeDZA    LanguageType = "dza"      // Tunzu
	LanguageTypeDZD    LanguageType = "dzd"      // Daza
	LanguageTypeDZE    LanguageType = "dze"      // Djiwarli
	LanguageTypeDZG    LanguageType = "dzg"      // Dazaga
	LanguageTypeDZL    LanguageType = "dzl"      // Dzalakha
	LanguageTypeDZN    LanguageType = "dzn"      // Dzando
	LanguageTypeEAA    LanguageType = "eaa"      // Karenggapa
	LanguageTypeEBG    LanguageType = "ebg"      // Ebughu
	LanguageTypeEBK    LanguageType = "ebk"      // Eastern Bontok
	LanguageTypeEBO    LanguageType = "ebo"      // Teke-Ebo
	LanguageTypeEBR    LanguageType = "ebr"      // Ebrié
	LanguageTypeEBU    LanguageType = "ebu"      // Embu and Kiembu
	LanguageTypeECR    LanguageType = "ecr"      // Eteocretan
	LanguageTypeECS    LanguageType = "ecs"      // Ecuadorian Sign Language
	LanguageTypeECY    LanguageType = "ecy"      // Eteocypriot
	LanguageTypeEEE    LanguageType = "eee"      // E
	LanguageTypeEFA    LanguageType = "efa"      // Efai
	LanguageTypeEFE    LanguageType = "efe"      // Efe
	LanguageTypeEFI    LanguageType = "efi"      // Efik
	LanguageTypeEGA    LanguageType = "ega"      // Ega
	LanguageTypeEGL    LanguageType = "egl"      // Emilian
	LanguageTypeEGO    LanguageType = "ego"      // Eggon
	LanguageTypeEGX    LanguageType = "egx"      // Egyptian languages
	LanguageTypeEGY    LanguageType = "egy"      // Egyptian (Ancient)
	LanguageTypeEHU    LanguageType = "ehu"      // Ehueun
	LanguageTypeEIP    LanguageType = "eip"      // Eipomek
	LanguageTypeEIT    LanguageType = "eit"      // Eitiep
	LanguageTypeEIV    LanguageType = "eiv"      // Askopan
	LanguageTypeEJA    LanguageType = "eja"      // Ejamat
	LanguageTypeEKA    LanguageType = "eka"      // Ekajuk
	LanguageTypeEKC    LanguageType = "ekc"      // Eastern Karnic
	LanguageTypeEKE    LanguageType = "eke"      // Ekit
	LanguageTypeEKG    LanguageType = "ekg"      // Ekari
	LanguageTypeEKI    LanguageType = "eki"      // Eki
	LanguageTypeEKK    LanguageType = "ekk"      // Standard Estonian
	LanguageTypeEKL    LanguageType = "ekl"      // Kol (Bangladesh) and Kol
	LanguageTypeEKM    LanguageType = "ekm"      // Elip
	LanguageTypeEKO    LanguageType = "eko"      // Koti
	LanguageTypeEKP    LanguageType = "ekp"      // Ekpeye
	LanguageTypeEKR    LanguageType = "ekr"      // Yace
	LanguageTypeEKY    LanguageType = "eky"      // Eastern Kayah
	LanguageTypeELE    LanguageType = "ele"      // Elepi
	LanguageTypeELH    LanguageType = "elh"      // El Hugeirat
	LanguageTypeELI    LanguageType = "eli"      // Nding
	LanguageTypeELK    LanguageType = "elk"      // Elkei
	LanguageTypeELM    LanguageType = "elm"      // Eleme
	LanguageTypeELO    LanguageType = "elo"      // El Molo
	LanguageTypeELP    LanguageType = "elp"      // Elpaputih
	LanguageTypeELU    LanguageType = "elu"      // Elu
	LanguageTypeELX    LanguageType = "elx"      // Elamite
	LanguageTypeEMA    LanguageType = "ema"      // Emai-Iuleha-Ora
	LanguageTypeEMB    LanguageType = "emb"      // Embaloh
	LanguageTypeEME    LanguageType = "eme"      // Emerillon
	LanguageTypeEMG    LanguageType = "emg"      // Eastern Meohang
	LanguageTypeEMI    LanguageType = "emi"      // Mussau-Emira
	LanguageTypeEMK    LanguageType = "emk"      // Eastern Maninkakan
	LanguageTypeEMM    LanguageType = "emm"      // Mamulique
	LanguageTypeEMN    LanguageType = "emn"      // Eman
	LanguageTypeEMO    LanguageType = "emo"      // Emok
	LanguageTypeEMP    LanguageType = "emp"      // Northern Emberá
	LanguageTypeEMS    LanguageType = "ems"      // Pacific Gulf Yupik
	LanguageTypeEMU    LanguageType = "emu"      // Eastern Muria
	LanguageTypeEMW    LanguageType = "emw"      // Emplawas
	LanguageTypeEMX    LanguageType = "emx"      // Erromintxela
	LanguageTypeEMY    LanguageType = "emy"      // Epigraphic Mayan
	LanguageTypeENA    LanguageType = "ena"      // Apali
	LanguageTypeENB    LanguageType = "enb"      // Markweeta
	LanguageTypeENC    LanguageType = "enc"      // En
	LanguageTypeEND    LanguageType = "end"      // Ende
	LanguageTypeENF    LanguageType = "enf"      // Forest Enets
	LanguageTypeENH    LanguageType = "enh"      // Tundra Enets
	LanguageTypeENM    LanguageType = "enm"      // Middle English (1100-1500)
	LanguageTypeENN    LanguageType = "enn"      // Engenni
	LanguageTypeENO    LanguageType = "eno"      // Enggano
	LanguageTypeENQ    LanguageType = "enq"      // Enga
	LanguageTypeENR    LanguageType = "enr"      // Emumu and Emem
	LanguageTypeENU    LanguageType = "enu"      // Enu
	LanguageTypeENV    LanguageType = "env"      // Enwan (Edu State)
	LanguageTypeENW    LanguageType = "enw"      // Enwan (Akwa Ibom State)
	LanguageTypeEOT    LanguageType = "eot"      // Beti (Côte d'Ivoire)
	LanguageTypeEPI    LanguageType = "epi"      // Epie
	LanguageTypeERA    LanguageType = "era"      // Eravallan
	LanguageTypeERG    LanguageType = "erg"      // Sie
	LanguageTypeERH    LanguageType = "erh"      // Eruwa
	LanguageTypeERI    LanguageType = "eri"      // Ogea
	LanguageTypeERK    LanguageType = "erk"      // South Efate
	LanguageTypeERO    LanguageType = "ero"      // Horpa
	LanguageTypeERR    LanguageType = "err"      // Erre
	LanguageTypeERS    LanguageType = "ers"      // Ersu
	LanguageTypeERT    LanguageType = "ert"      // Eritai
	LanguageTypeERW    LanguageType = "erw"      // Erokwanas
	LanguageTypeESE    LanguageType = "ese"      // Ese Ejja
	LanguageTypeESH    LanguageType = "esh"      // Eshtehardi
	LanguageTypeESI    LanguageType = "esi"      // North Alaskan Inupiatun
	LanguageTypeESK    LanguageType = "esk"      // Northwest Alaska Inupiatun
	LanguageTypeESL    LanguageType = "esl"      // Egypt Sign Language
	LanguageTypeESM    LanguageType = "esm"      // Esuma
	LanguageTypeESN    LanguageType = "esn"      // Salvadoran Sign Language
	LanguageTypeESO    LanguageType = "eso"      // Estonian Sign Language
	LanguageTypeESQ    LanguageType = "esq"      // Esselen
	LanguageTypeESS    LanguageType = "ess"      // Central Siberian Yupik
	LanguageTypeESU    LanguageType = "esu"      // Central Yupik
	LanguageTypeESX    LanguageType = "esx"      // Eskimo-Aleut languages
	LanguageTypeETB    LanguageType = "etb"      // Etebi
	LanguageTypeETC    LanguageType = "etc"      // Etchemin
	LanguageTypeETH    LanguageType = "eth"      // Ethiopian Sign Language
	LanguageTypeETN    LanguageType = "etn"      // Eton (Vanuatu)
	LanguageTypeETO    LanguageType = "eto"      // Eton (Cameroon)
	LanguageTypeETR    LanguageType = "etr"      // Edolo
	LanguageTypeETS    LanguageType = "ets"      // Yekhee
	LanguageTypeETT    LanguageType = "ett"      // Etruscan
	LanguageTypeETU    LanguageType = "etu"      // Ejagham
	LanguageTypeETX    LanguageType = "etx"      // Eten
	LanguageTypeETZ    LanguageType = "etz"      // Semimi
	LanguageTypeEUQ    LanguageType = "euq"      // Basque (family)
	LanguageTypeEVE    LanguageType = "eve"      // Even
	LanguageTypeEVH    LanguageType = "evh"      // Uvbie
	LanguageTypeEVN    LanguageType = "evn"      // Evenki
	LanguageTypeEWO    LanguageType = "ewo"      // Ewondo
	LanguageTypeEXT    LanguageType = "ext"      // Extremaduran
	LanguageTypeEYA    LanguageType = "eya"      // Eyak
	LanguageTypeEYO    LanguageType = "eyo"      // Keiyo
	LanguageTypeEZA    LanguageType = "eza"      // Ezaa
	LanguageTypeEZE    LanguageType = "eze"      // Uzekwe
	LanguageTypeFAA    LanguageType = "faa"      // Fasu
	LanguageTypeFAB    LanguageType = "fab"      // Fa d'Ambu
	LanguageTypeFAD    LanguageType = "fad"      // Wagi
	LanguageTypeFAF    LanguageType = "faf"      // Fagani
	LanguageTypeFAG    LanguageType = "fag"      // Finongan
	LanguageTypeFAH    LanguageType = "fah"      // Baissa Fali
	LanguageTypeFAI    LanguageType = "fai"      // Faiwol
	LanguageTypeFAJ    LanguageType = "faj"      // Faita
	LanguageTypeFAK    LanguageType = "fak"      // Fang (Cameroon)
	LanguageTypeFAL    LanguageType = "fal"      // South Fali
	LanguageTypeFAM    LanguageType = "fam"      // Fam
	LanguageTypeFAN    LanguageType = "fan"      // Fang (Equatorial Guinea)
	LanguageTypeFAP    LanguageType = "fap"      // Palor
	LanguageTypeFAR    LanguageType = "far"      // Fataleka
	LanguageTypeFAT    LanguageType = "fat"      // Fanti
	LanguageTypeFAU    LanguageType = "fau"      // Fayu
	LanguageTypeFAX    LanguageType = "fax"      // Fala
	LanguageTypeFAY    LanguageType = "fay"      // Southwestern Fars
	LanguageTypeFAZ    LanguageType = "faz"      // Northwestern Fars
	LanguageTypeFBL    LanguageType = "fbl"      // West Albay Bikol
	LanguageTypeFCS    LanguageType = "fcs"      // Quebec Sign Language
	LanguageTypeFER    LanguageType = "fer"      // Feroge
	LanguageTypeFFI    LanguageType = "ffi"      // Foia Foia
	LanguageTypeFFM    LanguageType = "ffm"      // Maasina Fulfulde
	LanguageTypeFGR    LanguageType = "fgr"      // Fongoro
	LanguageTypeFIA    LanguageType = "fia"      // Nobiin
	LanguageTypeFIE    LanguageType = "fie"      // Fyer
	LanguageTypeFIL    LanguageType = "fil"      // Filipino and Pilipino
	LanguageTypeFIP    LanguageType = "fip"      // Fipa
	LanguageTypeFIR    LanguageType = "fir"      // Firan
	LanguageTypeFIT    LanguageType = "fit"      // Tornedalen Finnish
	LanguageTypeFIU    LanguageType = "fiu"      // Finno-Ugrian languages
	LanguageTypeFIW    LanguageType = "fiw"      // Fiwaga
	LanguageTypeFKK    LanguageType = "fkk"      // Kirya-Konzəl
	LanguageTypeFKV    LanguageType = "fkv"      // Kven Finnish
	LanguageTypeFLA    LanguageType = "fla"      // Kalispel-Pend d'Oreille
	LanguageTypeFLH    LanguageType = "flh"      // Foau
	LanguageTypeFLI    LanguageType = "fli"      // Fali
	LanguageTypeFLL    LanguageType = "fll"      // North Fali
	LanguageTypeFLN    LanguageType = "fln"      // Flinders Island
	LanguageTypeFLR    LanguageType = "flr"      // Fuliiru
	LanguageTypeFLY    LanguageType = "fly"      // Tsotsitaal
	LanguageTypeFMP    LanguageType = "fmp"      // Fe'fe'
	LanguageTypeFMU    LanguageType = "fmu"      // Far Western Muria
	LanguageTypeFNG    LanguageType = "fng"      // Fanagalo
	LanguageTypeFNI    LanguageType = "fni"      // Fania
	LanguageTypeFOD    LanguageType = "fod"      // Foodo
	LanguageTypeFOI    LanguageType = "foi"      // Foi
	LanguageTypeFOM    LanguageType = "fom"      // Foma
	LanguageTypeFON    LanguageType = "fon"      // Fon
	LanguageTypeFOR    LanguageType = "for"      // Fore
	LanguageTypeFOS    LanguageType = "fos"      // Siraya
	LanguageTypeFOX    LanguageType = "fox"      // Formosan languages
	LanguageTypeFPE    LanguageType = "fpe"      // Fernando Po Creole English
	LanguageTypeFQS    LanguageType = "fqs"      // Fas
	LanguageTypeFRC    LanguageType = "frc"      // Cajun French
	LanguageTypeFRD    LanguageType = "frd"      // Fordata
	LanguageTypeFRK    LanguageType = "frk"      // Frankish
	LanguageTypeFRM    LanguageType = "frm"      // Middle French (ca. 1400-1600)
	LanguageTypeFRO    LanguageType = "fro"      // Old French (842-ca. 1400)
	LanguageTypeFRP    LanguageType = "frp"      // Arpitan and Francoprovençal
	LanguageTypeFRQ    LanguageType = "frq"      // Forak
	LanguageTypeFRR    LanguageType = "frr"      // Northern Frisian
	LanguageTypeFRS    LanguageType = "frs"      // Eastern Frisian
	LanguageTypeFRT    LanguageType = "frt"      // Fortsenal
	LanguageTypeFSE    LanguageType = "fse"      // Finnish Sign Language
	LanguageTypeFSL    LanguageType = "fsl"      // French Sign Language
	LanguageTypeFSS    LanguageType = "fss"      // Finland-Swedish Sign Language and finlandssvenskt teckenspråk and suomenruotsalainen viittomakieli
	LanguageTypeFUB    LanguageType = "fub"      // Adamawa Fulfulde
	LanguageTypeFUC    LanguageType = "fuc"      // Pulaar
	LanguageTypeFUD    LanguageType = "fud"      // East Futuna
	LanguageTypeFUE    LanguageType = "fue"      // Borgu Fulfulde
	LanguageTypeFUF    LanguageType = "fuf"      // Pular
	LanguageTypeFUH    LanguageType = "fuh"      // Western Niger Fulfulde
	LanguageTypeFUI    LanguageType = "fui"      // Bagirmi Fulfulde
	LanguageTypeFUJ    LanguageType = "fuj"      // Ko
	LanguageTypeFUM    LanguageType = "fum"      // Fum
	LanguageTypeFUN    LanguageType = "fun"      // Fulniô
	LanguageTypeFUQ    LanguageType = "fuq"      // Central-Eastern Niger Fulfulde
	LanguageTypeFUR    LanguageType = "fur"      // Friulian
	LanguageTypeFUT    LanguageType = "fut"      // Futuna-Aniwa
	LanguageTypeFUU    LanguageType = "fuu"      // Furu
	LanguageTypeFUV    LanguageType = "fuv"      // Nigerian Fulfulde
	LanguageTypeFUY    LanguageType = "fuy"      // Fuyug
	LanguageTypeFVR    LanguageType = "fvr"      // Fur
	LanguageTypeFWA    LanguageType = "fwa"      // Fwâi
	LanguageTypeFWE    LanguageType = "fwe"      // Fwe
	LanguageTypeGAA    LanguageType = "gaa"      // Ga
	LanguageTypeGAB    LanguageType = "gab"      // Gabri
	LanguageTypeGAC    LanguageType = "gac"      // Mixed Great Andamanese
	LanguageTypeGAD    LanguageType = "gad"      // Gaddang
	LanguageTypeGAE    LanguageType = "gae"      // Guarequena
	LanguageTypeGAF    LanguageType = "gaf"      // Gende
	LanguageTypeGAG    LanguageType = "gag"      // Gagauz
	LanguageTypeGAH    LanguageType = "gah"      // Alekano
	LanguageTypeGAI    LanguageType = "gai"      // Borei
	LanguageTypeGAJ    LanguageType = "gaj"      // Gadsup
	LanguageTypeGAK    LanguageType = "gak"      // Gamkonora
	LanguageTypeGAL    LanguageType = "gal"      // Galolen
	LanguageTypeGAM    LanguageType = "gam"      // Kandawo
	LanguageTypeGAN    LanguageType = "gan"      // Gan Chinese
	LanguageTypeGAO    LanguageType = "gao"      // Gants
	LanguageTypeGAP    LanguageType = "gap"      // Gal
	LanguageTypeGAQ    LanguageType = "gaq"      // Gata'
	LanguageTypeGAR    LanguageType = "gar"      // Galeya
	LanguageTypeGAS    LanguageType = "gas"      // Adiwasi Garasia
	LanguageTypeGAT    LanguageType = "gat"      // Kenati
	LanguageTypeGAU    LanguageType = "gau"      // Mudhili Gadaba
	LanguageTypeGAV    LanguageType = "gav"      // Gabutamon
	LanguageTypeGAW    LanguageType = "gaw"      // Nobonob
	LanguageTypeGAX    LanguageType = "gax"      // Borana-Arsi-Guji Oromo
	LanguageTypeGAY    LanguageType = "gay"      // Gayo
	LanguageTypeGAZ    LanguageType = "gaz"      // West Central Oromo
	LanguageTypeGBA    LanguageType = "gba"      // Gbaya (Central African Republic)
	LanguageTypeGBB    LanguageType = "gbb"      // Kaytetye
	LanguageTypeGBC    LanguageType = "gbc"      // Garawa
	LanguageTypeGBD    LanguageType = "gbd"      // Karadjeri
	LanguageTypeGBE    LanguageType = "gbe"      // Niksek
	LanguageTypeGBF    LanguageType = "gbf"      // Gaikundi
	LanguageTypeGBG    LanguageType = "gbg"      // Gbanziri
	LanguageTypeGBH    LanguageType = "gbh"      // Defi Gbe
	LanguageTypeGBI    LanguageType = "gbi"      // Galela
	LanguageTypeGBJ    LanguageType = "gbj"      // Bodo Gadaba
	LanguageTypeGBK    LanguageType = "gbk"      // Gaddi
	LanguageTypeGBL    LanguageType = "gbl"      // Gamit
	LanguageTypeGBM    LanguageType = "gbm"      // Garhwali
	LanguageTypeGBN    LanguageType = "gbn"      // Mo'da
	LanguageTypeGBO    LanguageType = "gbo"      // Northern Grebo
	LanguageTypeGBP    LanguageType = "gbp"      // Gbaya-Bossangoa
	LanguageTypeGBQ    LanguageType = "gbq"      // Gbaya-Bozoum
	LanguageTypeGBR    LanguageType = "gbr"      // Gbagyi
	LanguageTypeGBS    LanguageType = "gbs"      // Gbesi Gbe
	LanguageTypeGBU    LanguageType = "gbu"      // Gagadu
	LanguageTypeGBV    LanguageType = "gbv"      // Gbanu
	LanguageTypeGBW    LanguageType = "gbw"      // Gabi-Gabi
	LanguageTypeGBX    LanguageType = "gbx"      // Eastern Xwla Gbe
	LanguageTypeGBY    LanguageType = "gby"      // Gbari
	LanguageTypeGBZ    LanguageType = "gbz"      // Zoroastrian Dari
	LanguageTypeGCC    LanguageType = "gcc"      // Mali
	LanguageTypeGCD    LanguageType = "gcd"      // Ganggalida
	LanguageTypeGCE    LanguageType = "gce"      // Galice
	LanguageTypeGCF    LanguageType = "gcf"      // Guadeloupean Creole French
	LanguageTypeGCL    LanguageType = "gcl"      // Grenadian Creole English
	LanguageTypeGCN    LanguageType = "gcn"      // Gaina
	LanguageTypeGCR    LanguageType = "gcr"      // Guianese Creole French
	LanguageTypeGCT    LanguageType = "gct"      // Colonia Tovar German
	LanguageTypeGDA    LanguageType = "gda"      // Gade Lohar
	LanguageTypeGDB    LanguageType = "gdb"      // Pottangi Ollar Gadaba
	LanguageTypeGDC    LanguageType = "gdc"      // Gugu Badhun
	LanguageTypeGDD    LanguageType = "gdd"      // Gedaged
	LanguageTypeGDE    LanguageType = "gde"      // Gude
	LanguageTypeGDF    LanguageType = "gdf"      // Guduf-Gava
	LanguageTypeGDG    LanguageType = "gdg"      // Ga'dang
	LanguageTypeGDH    LanguageType = "gdh"      // Gadjerawang
	LanguageTypeGDI    LanguageType = "gdi"      // Gundi
	LanguageTypeGDJ    LanguageType = "gdj"      // Gurdjar
	LanguageTypeGDK    LanguageType = "gdk"      // Gadang
	LanguageTypeGDL    LanguageType = "gdl"      // Dirasha
	LanguageTypeGDM    LanguageType = "gdm"      // Laal
	LanguageTypeGDN    LanguageType = "gdn"      // Umanakaina
	LanguageTypeGDO    LanguageType = "gdo"      // Ghodoberi
	LanguageTypeGDQ    LanguageType = "gdq"      // Mehri
	LanguageTypeGDR    LanguageType = "gdr"      // Wipi
	LanguageTypeGDS    LanguageType = "gds"      // Ghandruk Sign Language
	LanguageTypeGDT    LanguageType = "gdt"      // Kungardutyi
	LanguageTypeGDU    LanguageType = "gdu"      // Gudu
	LanguageTypeGDX    LanguageType = "gdx"      // Godwari
	LanguageTypeGEA    LanguageType = "gea"      // Geruma
	LanguageTypeGEB    LanguageType = "geb"      // Kire
	LanguageTypeGEC    LanguageType = "gec"      // Gboloo Grebo
	LanguageTypeGED    LanguageType = "ged"      // Gade
	LanguageTypeGEG    LanguageType = "geg"      // Gengle
	LanguageTypeGEH    LanguageType = "geh"      // Hutterite German and Hutterisch
	LanguageTypeGEI    LanguageType = "gei"      // Gebe
	LanguageTypeGEJ    LanguageType = "gej"      // Gen
	LanguageTypeGEK    LanguageType = "gek"      // Yiwom
	LanguageTypeGEL    LanguageType = "gel"      // ut-Ma'in
	LanguageTypeGEM    LanguageType = "gem"      // Germanic languages
	LanguageTypeGEQ    LanguageType = "geq"      // Geme
	LanguageTypeGES    LanguageType = "ges"      // Geser-Gorom
	LanguageTypeGEW    LanguageType = "gew"      // Gera
	LanguageTypeGEX    LanguageType = "gex"      // Garre
	LanguageTypeGEY    LanguageType = "gey"      // Enya
	LanguageTypeGEZ    LanguageType = "gez"      // Geez
	LanguageTypeGFK    LanguageType = "gfk"      // Patpatar
	LanguageTypeGFT    LanguageType = "gft"      // Gafat
	LanguageTypeGFX    LanguageType = "gfx"      // Mangetti Dune !Xung
	LanguageTypeGGA    LanguageType = "gga"      // Gao
	LanguageTypeGGB    LanguageType = "ggb"      // Gbii
	LanguageTypeGGD    LanguageType = "ggd"      // Gugadj
	LanguageTypeGGE    LanguageType = "gge"      // Guragone
	LanguageTypeGGG    LanguageType = "ggg"      // Gurgula
	LanguageTypeGGK    LanguageType = "ggk"      // Kungarakany
	LanguageTypeGGL    LanguageType = "ggl"      // Ganglau
	LanguageTypeGGN    LanguageType = "ggn"      // Eastern Gurung
	LanguageTypeGGO    LanguageType = "ggo"      // Southern Gondi
	LanguageTypeGGR    LanguageType = "ggr"      // Aghu Tharnggalu
	LanguageTypeGGT    LanguageType = "ggt"      // Gitua
	LanguageTypeGGU    LanguageType = "ggu"      // Gagu and Gban
	LanguageTypeGGW    LanguageType = "ggw"      // Gogodala
	LanguageTypeGHA    LanguageType = "gha"      // Ghadamès
	LanguageTypeGHC    LanguageType = "ghc"      // Hiberno-Scottish Gaelic
	LanguageTypeGHE    LanguageType = "ghe"      // Southern Ghale
	LanguageTypeGHH    LanguageType = "ghh"      // Northern Ghale
	LanguageTypeGHK    LanguageType = "ghk"      // Geko Karen
	LanguageTypeGHL    LanguageType = "ghl"      // Ghulfan
	LanguageTypeGHN    LanguageType = "ghn"      // Ghanongga
	LanguageTypeGHO    LanguageType = "gho"      // Ghomara
	LanguageTypeGHR    LanguageType = "ghr"      // Ghera
	LanguageTypeGHS    LanguageType = "ghs"      // Guhu-Samane
	LanguageTypeGHT    LanguageType = "ght"      // Kuke and Kutang Ghale
	LanguageTypeGIA    LanguageType = "gia"      // Kitja
	LanguageTypeGIB    LanguageType = "gib"      // Gibanawa
	LanguageTypeGIC    LanguageType = "gic"      // Gail
	LanguageTypeGID    LanguageType = "gid"      // Gidar
	LanguageTypeGIG    LanguageType = "gig"      // Goaria
	LanguageTypeGIH    LanguageType = "gih"      // Githabul
	LanguageTypeGIL    LanguageType = "gil"      // Gilbertese
	LanguageTypeGIM    LanguageType = "gim"      // Gimi (Eastern Highlands)
	LanguageTypeGIN    LanguageType = "gin"      // Hinukh
	LanguageTypeGIO    LanguageType = "gio"      // Gelao
	LanguageTypeGIP    LanguageType = "gip"      // Gimi (West New Britain)
	LanguageTypeGIQ    LanguageType = "giq"      // Green Gelao
	LanguageTypeGIR    LanguageType = "gir"      // Red Gelao
	LanguageTypeGIS    LanguageType = "gis"      // North Giziga
	LanguageTypeGIT    LanguageType = "git"      // Gitxsan
	LanguageTypeGIU    LanguageType = "giu"      // Mulao
	LanguageTypeGIW    LanguageType = "giw"      // White Gelao
	LanguageTypeGIX    LanguageType = "gix"      // Gilima
	LanguageTypeGIY    LanguageType = "giy"      // Giyug
	LanguageTypeGIZ    LanguageType = "giz"      // South Giziga
	LanguageTypeGJI    LanguageType = "gji"      // Geji
	LanguageTypeGJK    LanguageType = "gjk"      // Kachi Koli
	LanguageTypeGJM    LanguageType = "gjm"      // Gunditjmara
	LanguageTypeGJN    LanguageType = "gjn"      // Gonja
	LanguageTypeGJU    LanguageType = "gju"      // Gujari
	LanguageTypeGKA    LanguageType = "gka"      // Guya
	LanguageTypeGKE    LanguageType = "gke"      // Ndai
	LanguageTypeGKN    LanguageType = "gkn"      // Gokana
	LanguageTypeGKO    LanguageType = "gko"      // Kok-Nar
	LanguageTypeGKP    LanguageType = "gkp"      // Guinea Kpelle
	LanguageTypeGLC    LanguageType = "glc"      // Bon Gula
	LanguageTypeGLD    LanguageType = "gld"      // Nanai
	LanguageTypeGLH    LanguageType = "glh"      // Northwest Pashayi
	LanguageTypeGLI    LanguageType = "gli"      // Guliguli
	LanguageTypeGLJ    LanguageType = "glj"      // Gula Iro
	LanguageTypeGLK    LanguageType = "glk"      // Gilaki
	LanguageTypeGLL    LanguageType = "gll"      // Garlali
	LanguageTypeGLO    LanguageType = "glo"      // Galambu
	LanguageTypeGLR    LanguageType = "glr"      // Glaro-Twabo
	LanguageTypeGLU    LanguageType = "glu"      // Gula (Chad)
	LanguageTypeGLW    LanguageType = "glw"      // Glavda
	LanguageTypeGLY    LanguageType = "gly"      // Gule
	LanguageTypeGMA    LanguageType = "gma"      // Gambera
	LanguageTypeGMB    LanguageType = "gmb"      // Gula'alaa
	LanguageTypeGMD    LanguageType = "gmd"      // Mághdì
	LanguageTypeGME    LanguageType = "gme"      // East Germanic languages
	LanguageTypeGMH    LanguageType = "gmh"      // Middle High German (ca. 1050-1500)
	LanguageTypeGML    LanguageType = "gml"      // Middle Low German
	LanguageTypeGMM    LanguageType = "gmm"      // Gbaya-Mbodomo
	LanguageTypeGMN    LanguageType = "gmn"      // Gimnime
	LanguageTypeGMQ    LanguageType = "gmq"      // North Germanic languages
	LanguageTypeGMU    LanguageType = "gmu"      // Gumalu
	LanguageTypeGMV    LanguageType = "gmv"      // Gamo
	LanguageTypeGMW    LanguageType = "gmw"      // West Germanic languages
	LanguageTypeGMX    LanguageType = "gmx"      // Magoma
	LanguageTypeGMY    LanguageType = "gmy"      // Mycenaean Greek
	LanguageTypeGMZ    LanguageType = "gmz"      // Mgbolizhia
	LanguageTypeGNA    LanguageType = "gna"      // Kaansa
	LanguageTypeGNB    LanguageType = "gnb"      // Gangte
	LanguageTypeGNC    LanguageType = "gnc"      // Guanche
	LanguageTypeGND    LanguageType = "gnd"      // Zulgo-Gemzek
	LanguageTypeGNE    LanguageType = "gne"      // Ganang
	LanguageTypeGNG    LanguageType = "gng"      // Ngangam
	LanguageTypeGNH    LanguageType = "gnh"      // Lere
	LanguageTypeGNI    LanguageType = "gni"      // Gooniyandi
	LanguageTypeGNK    LanguageType = "gnk"      // //Gana
	LanguageTypeGNL    LanguageType = "gnl"      // Gangulu
	LanguageTypeGNM    LanguageType = "gnm"      // Ginuman
	LanguageTypeGNN    LanguageType = "gnn"      // Gumatj
	LanguageTypeGNO    LanguageType = "gno"      // Northern Gondi
	LanguageTypeGNQ    LanguageType = "gnq"      // Gana
	LanguageTypeGNR    LanguageType = "gnr"      // Gureng Gureng
	LanguageTypeGNT    LanguageType = "gnt"      // Guntai
	LanguageTypeGNU    LanguageType = "gnu"      // Gnau
	LanguageTypeGNW    LanguageType = "gnw"      // Western Bolivian Guaraní
	LanguageTypeGNZ    LanguageType = "gnz"      // Ganzi
	LanguageTypeGOA    LanguageType = "goa"      // Guro
	LanguageTypeGOB    LanguageType = "gob"      // Playero
	LanguageTypeGOC    LanguageType = "goc"      // Gorakor
	LanguageTypeGOD    LanguageType = "god"      // Godié
	LanguageTypeGOE    LanguageType = "goe"      // Gongduk
	LanguageTypeGOF    LanguageType = "gof"      // Gofa
	LanguageTypeGOG    LanguageType = "gog"      // Gogo
	LanguageTypeGOH    LanguageType = "goh"      // Old High German (ca. 750-1050)
	LanguageTypeGOI    LanguageType = "goi"      // Gobasi
	LanguageTypeGOJ    LanguageType = "goj"      // Gowlan
	LanguageTypeGOK    LanguageType = "gok"      // Gowli
	LanguageTypeGOL    LanguageType = "gol"      // Gola
	LanguageTypeGOM    LanguageType = "gom"      // Goan Konkani
	LanguageTypeGON    LanguageType = "gon"      // Gondi
	LanguageTypeGOO    LanguageType = "goo"      // Gone Dau
	LanguageTypeGOP    LanguageType = "gop"      // Yeretuar
	LanguageTypeGOQ    LanguageType = "goq"      // Gorap
	LanguageTypeGOR    LanguageType = "gor"      // Gorontalo
	LanguageTypeGOS    LanguageType = "gos"      // Gronings
	LanguageTypeGOT    LanguageType = "got"      // Gothic
	LanguageTypeGOU    LanguageType = "gou"      // Gavar
	LanguageTypeGOW    LanguageType = "gow"      // Gorowa
	LanguageTypeGOX    LanguageType = "gox"      // Gobu
	LanguageTypeGOY    LanguageType = "goy"      // Goundo
	LanguageTypeGOZ    LanguageType = "goz"      // Gozarkhani
	LanguageTypeGPA    LanguageType = "gpa"      // Gupa-Abawa
	LanguageTypeGPE    LanguageType = "gpe"      // Ghanaian Pidgin English
	LanguageTypeGPN    LanguageType = "gpn"      // Taiap
	LanguageTypeGQA    LanguageType = "gqa"      // Ga'anda
	LanguageTypeGQI    LanguageType = "gqi"      // Guiqiong
	LanguageTypeGQN    LanguageType = "gqn"      // Guana (Brazil)
	LanguageTypeGQR    LanguageType = "gqr"      // Gor
	LanguageTypeGQU    LanguageType = "gqu"      // Qau
	LanguageTypeGRA    LanguageType = "gra"      // Rajput Garasia
	LanguageTypeGRB    LanguageType = "grb"      // Grebo
	LanguageTypeGRC    LanguageType = "grc"      // Ancient Greek (to 1453)
	LanguageTypeGRD    LanguageType = "grd"      // Guruntum-Mbaaru
	LanguageTypeGRG    LanguageType = "grg"      // Madi
	LanguageTypeGRH    LanguageType = "grh"      // Gbiri-Niragu
	LanguageTypeGRI    LanguageType = "gri"      // Ghari
	LanguageTypeGRJ    LanguageType = "grj"      // Southern Grebo
	LanguageTypeGRK    LanguageType = "grk"      // Greek languages
	LanguageTypeGRM    LanguageType = "grm"      // Kota Marudu Talantang
	LanguageTypeGRO    LanguageType = "gro"      // Groma
	LanguageTypeGRQ    LanguageType = "grq"      // Gorovu
	LanguageTypeGRR    LanguageType = "grr"      // Taznatit
	LanguageTypeGRS    LanguageType = "grs"      // Gresi
	LanguageTypeGRT    LanguageType = "grt"      // Garo
	LanguageTypeGRU    LanguageType = "gru"      // Kistane
	LanguageTypeGRV    LanguageType = "grv"      // Central Grebo
	LanguageTypeGRW    LanguageType = "grw"      // Gweda
	LanguageTypeGRX    LanguageType = "grx"      // Guriaso
	LanguageTypeGRY    LanguageType = "gry"      // Barclayville Grebo
	LanguageTypeGRZ    LanguageType = "grz"      // Guramalum
	LanguageTypeGSE    LanguageType = "gse"      // Ghanaian Sign Language
	LanguageTypeGSG    LanguageType = "gsg"      // German Sign Language
	LanguageTypeGSL    LanguageType = "gsl"      // Gusilay
	LanguageTypeGSM    LanguageType = "gsm"      // Guatemalan Sign Language
	LanguageTypeGSN    LanguageType = "gsn"      // Gusan
	LanguageTypeGSO    LanguageType = "gso"      // Southwest Gbaya
	LanguageTypeGSP    LanguageType = "gsp"      // Wasembo
	LanguageTypeGSS    LanguageType = "gss"      // Greek Sign Language
	LanguageTypeGSW    LanguageType = "gsw"      // Swiss German and Alemannic and Alsatian
	LanguageTypeGTA    LanguageType = "gta"      // Guató
	LanguageTypeGTI    LanguageType = "gti"      // Gbati-ri
	LanguageTypeGTU    LanguageType = "gtu"      // Aghu-Tharnggala
	LanguageTypeGUA    LanguageType = "gua"      // Shiki
	LanguageTypeGUB    LanguageType = "gub"      // Guajajára
	LanguageTypeGUC    LanguageType = "guc"      // Wayuu
	LanguageTypeGUD    LanguageType = "gud"      // Yocoboué Dida
	LanguageTypeGUE    LanguageType = "gue"      // Gurinji
	LanguageTypeGUF    LanguageType = "guf"      // Gupapuyngu
	LanguageTypeGUG    LanguageType = "gug"      // Paraguayan Guaraní
	LanguageTypeGUH    LanguageType = "guh"      // Guahibo
	LanguageTypeGUI    LanguageType = "gui"      // Eastern Bolivian Guaraní
	LanguageTypeGUK    LanguageType = "guk"      // Gumuz
	LanguageTypeGUL    LanguageType = "gul"      // Sea Island Creole English
	LanguageTypeGUM    LanguageType = "gum"      // Guambiano
	LanguageTypeGUN    LanguageType = "gun"      // Mbyá Guaraní
	LanguageTypeGUO    LanguageType = "guo"      // Guayabero
	LanguageTypeGUP    LanguageType = "gup"      // Gunwinggu
	LanguageTypeGUQ    LanguageType = "guq"      // Aché
	LanguageTypeGUR    LanguageType = "gur"      // Farefare
	LanguageTypeGUS    LanguageType = "gus"      // Guinean Sign Language
	LanguageTypeGUT    LanguageType = "gut"      // Maléku Jaíka
	LanguageTypeGUU    LanguageType = "guu"      // Yanomamö
	LanguageTypeGUV    LanguageType = "guv"      // Gey
	LanguageTypeGUW    LanguageType = "guw"      // Gun
	LanguageTypeGUX    LanguageType = "gux"      // Gourmanchéma
	LanguageTypeGUZ    LanguageType = "guz"      // Gusii and Ekegusii
	LanguageTypeGVA    LanguageType = "gva"      // Guana (Paraguay)
	LanguageTypeGVC    LanguageType = "gvc"      // Guanano
	LanguageTypeGVE    LanguageType = "gve"      // Duwet
	LanguageTypeGVF    LanguageType = "gvf"      // Golin
	LanguageTypeGVJ    LanguageType = "gvj"      // Guajá
	LanguageTypeGVL    LanguageType = "gvl"      // Gulay
	LanguageTypeGVM    LanguageType = "gvm"      // Gurmana
	LanguageTypeGVN    LanguageType = "gvn"      // Kuku-Yalanji
	LanguageTypeGVO    LanguageType = "gvo"      // Gavião Do Jiparaná
	LanguageTypeGVP    LanguageType = "gvp"      // Pará Gavião
	LanguageTypeGVR    LanguageType = "gvr"      // Western Gurung
	LanguageTypeGVS    LanguageType = "gvs"      // Gumawana
	LanguageTypeGVY    LanguageType = "gvy"      // Guyani
	LanguageTypeGWA    LanguageType = "gwa"      // Mbato
	LanguageTypeGWB    LanguageType = "gwb"      // Gwa
	LanguageTypeGWC    LanguageType = "gwc"      // Kalami
	LanguageTypeGWD    LanguageType = "gwd"      // Gawwada
	LanguageTypeGWE    LanguageType = "gwe"      // Gweno
	LanguageTypeGWF    LanguageType = "gwf"      // Gowro
	LanguageTypeGWG    LanguageType = "gwg"      // Moo
	LanguageTypeGWI    LanguageType = "gwi"      // Gwichʼin
	LanguageTypeGWJ    LanguageType = "gwj"      // /Gwi
	LanguageTypeGWM    LanguageType = "gwm"      // Awngthim
	LanguageTypeGWN    LanguageType = "gwn"      // Gwandara
	LanguageTypeGWR    LanguageType = "gwr"      // Gwere
	LanguageTypeGWT    LanguageType = "gwt"      // Gawar-Bati
	LanguageTypeGWU    LanguageType = "gwu"      // Guwamu
	LanguageTypeGWW    LanguageType = "gww"      // Kwini
	LanguageTypeGWX    LanguageType = "gwx"      // Gua
	LanguageTypeGXX    LanguageType = "gxx"      // Wè Southern
	LanguageTypeGYA    LanguageType = "gya"      // Northwest Gbaya
	LanguageTypeGYB    LanguageType = "gyb"      // Garus
	LanguageTypeGYD    LanguageType = "gyd"      // Kayardild
	LanguageTypeGYE    LanguageType = "gye"      // Gyem
	LanguageTypeGYF    LanguageType = "gyf"      // Gungabula
	LanguageTypeGYG    LanguageType = "gyg"      // Gbayi
	LanguageTypeGYI    LanguageType = "gyi"      // Gyele
	LanguageTypeGYL    LanguageType = "gyl"      // Gayil
	LanguageTypeGYM    LanguageType = "gym"      // Ngäbere
	LanguageTypeGYN    LanguageType = "gyn"      // Guyanese Creole English
	LanguageTypeGYR    LanguageType = "gyr"      // Guarayu
	LanguageTypeGYY    LanguageType = "gyy"      // Gunya
	LanguageTypeGZA    LanguageType = "gza"      // Ganza
	LanguageTypeGZI    LanguageType = "gzi"      // Gazi
	LanguageTypeGZN    LanguageType = "gzn"      // Gane
	LanguageTypeHAA    LanguageType = "haa"      // Han
	LanguageTypeHAB    LanguageType = "hab"      // Hanoi Sign Language
	LanguageTypeHAC    LanguageType = "hac"      // Gurani
	LanguageTypeHAD    LanguageType = "had"      // Hatam
	LanguageTypeHAE    LanguageType = "hae"      // Eastern Oromo
	LanguageTypeHAF    LanguageType = "haf"      // Haiphong Sign Language
	LanguageTypeHAG    LanguageType = "hag"      // Hanga
	LanguageTypeHAH    LanguageType = "hah"      // Hahon
	LanguageTypeHAI    LanguageType = "hai"      // Haida
	LanguageTypeHAJ    LanguageType = "haj"      // Hajong
	LanguageTypeHAK    LanguageType = "hak"      // Hakka Chinese
	LanguageTypeHAL    LanguageType = "hal"      // Halang
	LanguageTypeHAM    LanguageType = "ham"      // Hewa
	LanguageTypeHAN    LanguageType = "han"      // Hangaza
	LanguageTypeHAO    LanguageType = "hao"      // Hakö
	LanguageTypeHAP    LanguageType = "hap"      // Hupla
	LanguageTypeHAQ    LanguageType = "haq"      // Ha
	LanguageTypeHAR    LanguageType = "har"      // Harari
	LanguageTypeHAS    LanguageType = "has"      // Haisla
	LanguageTypeHAV    LanguageType = "hav"      // Havu
	LanguageTypeHAW    LanguageType = "haw"      // Hawaiian
	LanguageTypeHAX    LanguageType = "hax"      // Southern Haida
	LanguageTypeHAY    LanguageType = "hay"      // Haya
	LanguageTypeHAZ    LanguageType = "haz"      // Hazaragi
	LanguageTypeHBA    LanguageType = "hba"      // Hamba
	LanguageTypeHBB    LanguageType = "hbb"      // Huba
	LanguageTypeHBN    LanguageType = "hbn"      // Heiban
	LanguageTypeHBO    LanguageType = "hbo"      // Ancient Hebrew
	LanguageTypeHBU    LanguageType = "hbu"      // Habu
	LanguageTypeHCA    LanguageType = "hca"      // Andaman Creole Hindi
	LanguageTypeHCH    LanguageType = "hch"      // Huichol
	LanguageTypeHDN    LanguageType = "hdn"      // Northern Haida
	LanguageTypeHDS    LanguageType = "hds"      // Honduras Sign Language
	LanguageTypeHDY    LanguageType = "hdy"      // Hadiyya
	LanguageTypeHEA    LanguageType = "hea"      // Northern Qiandong Miao
	LanguageTypeHED    LanguageType = "hed"      // Herdé
	LanguageTypeHEG    LanguageType = "heg"      // Helong
	LanguageTypeHEH    LanguageType = "heh"      // Hehe
	LanguageTypeHEI    LanguageType = "hei"      // Heiltsuk
	LanguageTypeHEM    LanguageType = "hem"      // Hemba
	LanguageTypeHGM    LanguageType = "hgm"      // Hai//om
	LanguageTypeHGW    LanguageType = "hgw"      // Haigwai
	LanguageTypeHHI    LanguageType = "hhi"      // Hoia Hoia
	LanguageTypeHHR    LanguageType = "hhr"      // Kerak
	LanguageTypeHHY    LanguageType = "hhy"      // Hoyahoya
	LanguageTypeHIA    LanguageType = "hia"      // Lamang
	LanguageTypeHIB    LanguageType = "hib"      // Hibito
	LanguageTypeHID    LanguageType = "hid"      // Hidatsa
	LanguageTypeHIF    LanguageType = "hif"      // Fiji Hindi
	LanguageTypeHIG    LanguageType = "hig"      // Kamwe
	LanguageTypeHIH    LanguageType = "hih"      // Pamosu
	LanguageTypeHII    LanguageType = "hii"      // Hinduri
	LanguageTypeHIJ    LanguageType = "hij"      // Hijuk
	LanguageTypeHIK    LanguageType = "hik"      // Seit-Kaitetu
	LanguageTypeHIL    LanguageType = "hil"      // Hiligaynon
	LanguageTypeHIM    LanguageType = "him"      // Himachali languages and Western Pahari languages
	LanguageTypeHIO    LanguageType = "hio"      // Tsoa
	LanguageTypeHIR    LanguageType = "hir"      // Himarimã
	LanguageTypeHIT    LanguageType = "hit"      // Hittite
	LanguageTypeHIW    LanguageType = "hiw"      // Hiw
	LanguageTypeHIX    LanguageType = "hix"      // Hixkaryána
	LanguageTypeHJI    LanguageType = "hji"      // Haji
	LanguageTypeHKA    LanguageType = "hka"      // Kahe
	LanguageTypeHKE    LanguageType = "hke"      // Hunde
	LanguageTypeHKK    LanguageType = "hkk"      // Hunjara-Kaina Ke
	LanguageTypeHKS    LanguageType = "hks"      // Hong Kong Sign Language and Heung Kong Sau Yue
	LanguageTypeHLA    LanguageType = "hla"      // Halia
	LanguageTypeHLB    LanguageType = "hlb"      // Halbi
	LanguageTypeHLD    LanguageType = "hld"      // Halang Doan
	LanguageTypeHLE    LanguageType = "hle"      // Hlersu
	LanguageTypeHLT    LanguageType = "hlt"      // Matu Chin
	LanguageTypeHLU    LanguageType = "hlu"      // Hieroglyphic Luwian
	LanguageTypeHMA    LanguageType = "hma"      // Southern Mashan Hmong and Southern Mashan Miao
	LanguageTypeHMB    LanguageType = "hmb"      // Humburi Senni Songhay
	LanguageTypeHMC    LanguageType = "hmc"      // Central Huishui Hmong and Central Huishui Miao
	LanguageTypeHMD    LanguageType = "hmd"      // Large Flowery Miao and A-hmaos and Da-Hua Miao
	LanguageTypeHME    LanguageType = "hme"      // Eastern Huishui Hmong and Eastern Huishui Miao
	LanguageTypeHMF    LanguageType = "hmf"      // Hmong Don
	LanguageTypeHMG    LanguageType = "hmg"      // Southwestern Guiyang Hmong
	LanguageTypeHMH    LanguageType = "hmh"      // Southwestern Huishui Hmong and Southwestern Huishui Miao
	LanguageTypeHMI    LanguageType = "hmi"      // Northern Huishui Hmong and Northern Huishui Miao
	LanguageTypeHMJ    LanguageType = "hmj"      // Ge and Gejia
	LanguageTypeHMK    LanguageType = "hmk"      // Maek
	LanguageTypeHML    LanguageType = "hml"      // Luopohe Hmong and Luopohe Miao
	LanguageTypeHMM    LanguageType = "hmm"      // Central Mashan Hmong and Central Mashan Miao
	LanguageTypeHMN    LanguageType = "hmn"      // Hmong and Mong
	LanguageTypeHMP    LanguageType = "hmp"      // Northern Mashan Hmong and Northern Mashan Miao
	LanguageTypeHMQ    LanguageType = "hmq"      // Eastern Qiandong Miao
	LanguageTypeHMR    LanguageType = "hmr"      // Hmar
	LanguageTypeHMS    LanguageType = "hms"      // Southern Qiandong Miao
	LanguageTypeHMT    LanguageType = "hmt"      // Hamtai
	LanguageTypeHMU    LanguageType = "hmu"      // Hamap
	LanguageTypeHMV    LanguageType = "hmv"      // Hmong Dô
	LanguageTypeHMW    LanguageType = "hmw"      // Western Mashan Hmong and Western Mashan Miao
	LanguageTypeHMX    LanguageType = "hmx"      // Hmong-Mien languages
	LanguageTypeHMY    LanguageType = "hmy"      // Southern Guiyang Hmong and Southern Guiyang Miao
	LanguageTypeHMZ    LanguageType = "hmz"      // Hmong Shua and Sinicized Miao
	LanguageTypeHNA    LanguageType = "hna"      // Mina (Cameroon)
	LanguageTypeHND    LanguageType = "hnd"      // Southern Hindko
	LanguageTypeHNE    LanguageType = "hne"      // Chhattisgarhi
	LanguageTypeHNH    LanguageType = "hnh"      // //Ani
	LanguageTypeHNI    LanguageType = "hni"      // Hani
	LanguageTypeHNJ    LanguageType = "hnj"      // Hmong Njua and Mong Leng and Mong Njua
	LanguageTypeHNN    LanguageType = "hnn"      // Hanunoo
	LanguageTypeHNO    LanguageType = "hno"      // Northern Hindko
	LanguageTypeHNS    LanguageType = "hns"      // Caribbean Hindustani
	LanguageTypeHNU    LanguageType = "hnu"      // Hung
	LanguageTypeHOA    LanguageType = "hoa"      // Hoava
	LanguageTypeHOB    LanguageType = "hob"      // Mari (Madang Province)
	LanguageTypeHOC    LanguageType = "hoc"      // Ho
	LanguageTypeHOD    LanguageType = "hod"      // Holma
	LanguageTypeHOE    LanguageType = "hoe"      // Horom
	LanguageTypeHOH    LanguageType = "hoh"      // Hobyót
	LanguageTypeHOI    LanguageType = "hoi"      // Holikachuk
	LanguageTypeHOJ    LanguageType = "hoj"      // Hadothi and Haroti
	LanguageTypeHOK    LanguageType = "hok"      // Hokan languages
	LanguageTypeHOL    LanguageType = "hol"      // Holu
	LanguageTypeHOM    LanguageType = "hom"      // Homa
	LanguageTypeHOO    LanguageType = "hoo"      // Holoholo
	LanguageTypeHOP    LanguageType = "hop"      // Hopi
	LanguageTypeHOR    LanguageType = "hor"      // Horo
	LanguageTypeHOS    LanguageType = "hos"      // Ho Chi Minh City Sign Language
	LanguageTypeHOT    LanguageType = "hot"      // Hote and Malê
	LanguageTypeHOV    LanguageType = "hov"      // Hovongan
	LanguageTypeHOW    LanguageType = "how"      // Honi
	LanguageTypeHOY    LanguageType = "hoy"      // Holiya
	LanguageTypeHOZ    LanguageType = "hoz"      // Hozo
	LanguageTypeHPO    LanguageType = "hpo"      // Hpon
	LanguageTypeHPS    LanguageType = "hps"      // Hawai'i Pidgin Sign Language
	LanguageTypeHRA    LanguageType = "hra"      // Hrangkhol
	LanguageTypeHRC    LanguageType = "hrc"      // Niwer Mil
	LanguageTypeHRE    LanguageType = "hre"      // Hre
	LanguageTypeHRK    LanguageType = "hrk"      // Haruku
	LanguageTypeHRM    LanguageType = "hrm"      // Horned Miao
	LanguageTypeHRO    LanguageType = "hro"      // Haroi
	LanguageTypeHRP    LanguageType = "hrp"      // Nhirrpi
	LanguageTypeHRR    LanguageType = "hrr"      // Horuru
	LanguageTypeHRT    LanguageType = "hrt"      // Hértevin
	LanguageTypeHRU    LanguageType = "hru"      // Hruso
	LanguageTypeHRW    LanguageType = "hrw"      // Warwar Feni
	LanguageTypeHRX    LanguageType = "hrx"      // Hunsrik
	LanguageTypeHRZ    LanguageType = "hrz"      // Harzani
	LanguageTypeHSB    LanguageType = "hsb"      // Upper Sorbian
	LanguageTypeHSH    LanguageType = "hsh"      // Hungarian Sign Language
	LanguageTypeHSL    LanguageType = "hsl"      // Hausa Sign Language
	LanguageTypeHSN    LanguageType = "hsn"      // Xiang Chinese
	LanguageTypeHSS    LanguageType = "hss"      // Harsusi
	LanguageTypeHTI    LanguageType = "hti"      // Hoti
	LanguageTypeHTO    LanguageType = "hto"      // Minica Huitoto
	LanguageTypeHTS    LanguageType = "hts"      // Hadza
	LanguageTypeHTU    LanguageType = "htu"      // Hitu
	LanguageTypeHTX    LanguageType = "htx"      // Middle Hittite
	LanguageTypeHUB    LanguageType = "hub"      // Huambisa
	LanguageTypeHUC    LanguageType = "huc"      // =/Hua
	LanguageTypeHUD    LanguageType = "hud"      // Huaulu
	LanguageTypeHUE    LanguageType = "hue"      // San Francisco Del Mar Huave
	LanguageTypeHUF    LanguageType = "huf"      // Humene
	LanguageTypeHUG    LanguageType = "hug"      // Huachipaeri
	LanguageTypeHUH    LanguageType = "huh"      // Huilliche
	LanguageTypeHUI    LanguageType = "hui"      // Huli
	LanguageTypeHUJ    LanguageType = "huj"      // Northern Guiyang Hmong and Northern Guiyang Miao
	LanguageTypeHUK    LanguageType = "huk"      // Hulung
	LanguageTypeHUL    LanguageType = "hul"      // Hula
	LanguageTypeHUM    LanguageType = "hum"      // Hungana
	LanguageTypeHUO    LanguageType = "huo"      // Hu
	LanguageTypeHUP    LanguageType = "hup"      // Hupa
	LanguageTypeHUQ    LanguageType = "huq"      // Tsat
	LanguageTypeHUR    LanguageType = "hur"      // Halkomelem
	LanguageTypeHUS    LanguageType = "hus"      // Huastec
	LanguageTypeHUT    LanguageType = "hut"      // Humla
	LanguageTypeHUU    LanguageType = "huu"      // Murui Huitoto
	LanguageTypeHUV    LanguageType = "huv"      // San Mateo Del Mar Huave
	LanguageTypeHUW    LanguageType = "huw"      // Hukumina
	LanguageTypeHUX    LanguageType = "hux"      // Nüpode Huitoto
	LanguageTypeHUY    LanguageType = "huy"      // Hulaulá
	LanguageTypeHUZ    LanguageType = "huz"      // Hunzib
	LanguageTypeHVC    LanguageType = "hvc"      // Haitian Vodoun Culture Language
	LanguageTypeHVE    LanguageType = "hve"      // San Dionisio Del Mar Huave
	LanguageTypeHVK    LanguageType = "hvk"      // Haveke
	LanguageTypeHVN    LanguageType = "hvn"      // Sabu
	LanguageTypeHVV    LanguageType = "hvv"      // Santa María Del Mar Huave
	LanguageTypeHWA    LanguageType = "hwa"      // Wané
	LanguageTypeHWC    LanguageType = "hwc"      // Hawai'i Creole English and Hawai'i Pidgin
	LanguageTypeHWO    LanguageType = "hwo"      // Hwana
	LanguageTypeHYA    LanguageType = "hya"      // Hya
	LanguageTypeHYX    LanguageType = "hyx"      // Armenian (family)
	LanguageTypeIAI    LanguageType = "iai"      // Iaai
	LanguageTypeIAN    LanguageType = "ian"      // Iatmul
	LanguageTypeIAP    LanguageType = "iap"      // Iapama
	LanguageTypeIAR    LanguageType = "iar"      // Purari
	LanguageTypeIBA    LanguageType = "iba"      // Iban
	LanguageTypeIBB    LanguageType = "ibb"      // Ibibio
	LanguageTypeIBD    LanguageType = "ibd"      // Iwaidja
	LanguageTypeIBE    LanguageType = "ibe"      // Akpes
	LanguageTypeIBG    LanguageType = "ibg"      // Ibanag
	LanguageTypeIBI    LanguageType = "ibi"      // Ibilo
	LanguageTypeIBL    LanguageType = "ibl"      // Ibaloi
	LanguageTypeIBM    LanguageType = "ibm"      // Agoi
	LanguageTypeIBN    LanguageType = "ibn"      // Ibino
	LanguageTypeIBR    LanguageType = "ibr"      // Ibuoro
	LanguageTypeIBU    LanguageType = "ibu"      // Ibu
	LanguageTypeIBY    LanguageType = "iby"      // Ibani
	LanguageTypeICA    LanguageType = "ica"      // Ede Ica
	LanguageTypeICH    LanguageType = "ich"      // Etkywan
	LanguageTypeICL    LanguageType = "icl"      // Icelandic Sign Language
	LanguageTypeICR    LanguageType = "icr"      // Islander Creole English
	LanguageTypeIDA    LanguageType = "ida"      // Idakho-Isukha-Tiriki and Luidakho-Luisukha-Lutirichi
	LanguageTypeIDB    LanguageType = "idb"      // Indo-Portuguese
	LanguageTypeIDC    LanguageType = "idc"      // Idon and Ajiya
	LanguageTypeIDD    LanguageType = "idd"      // Ede Idaca
	LanguageTypeIDE    LanguageType = "ide"      // Idere
	LanguageTypeIDI    LanguageType = "idi"      // Idi
	LanguageTypeIDR    LanguageType = "idr"      // Indri
	LanguageTypeIDS    LanguageType = "ids"      // Idesa
	LanguageTypeIDT    LanguageType = "idt"      // Idaté
	LanguageTypeIDU    LanguageType = "idu"      // Idoma
	LanguageTypeIFA    LanguageType = "ifa"      // Amganad Ifugao
	LanguageTypeIFB    LanguageType = "ifb"      // Batad Ifugao and Ayangan Ifugao
	LanguageTypeIFE    LanguageType = "ife"      // Ifè
	LanguageTypeIFF    LanguageType = "iff"      // Ifo
	LanguageTypeIFK    LanguageType = "ifk"      // Tuwali Ifugao
	LanguageTypeIFM    LanguageType = "ifm"      // Teke-Fuumu
	LanguageTypeIFU    LanguageType = "ifu"      // Mayoyao Ifugao
	LanguageTypeIFY    LanguageType = "ify"      // Keley-I Kallahan
	LanguageTypeIGB    LanguageType = "igb"      // Ebira
	LanguageTypeIGE    LanguageType = "ige"      // Igede
	LanguageTypeIGG    LanguageType = "igg"      // Igana
	LanguageTypeIGL    LanguageType = "igl"      // Igala
	LanguageTypeIGM    LanguageType = "igm"      // Kanggape
	LanguageTypeIGN    LanguageType = "ign"      // Ignaciano
	LanguageTypeIGO    LanguageType = "igo"      // Isebe
	LanguageTypeIGS    LanguageType = "igs"      // Interglossa
	LanguageTypeIGW    LanguageType = "igw"      // Igwe
	LanguageTypeIHB    LanguageType = "ihb"      // Iha Based Pidgin
	LanguageTypeIHI    LanguageType = "ihi"      // Ihievbe
	LanguageTypeIHP    LanguageType = "ihp"      // Iha
	LanguageTypeIHW    LanguageType = "ihw"      // Bidhawal
	LanguageTypeIIN    LanguageType = "iin"      // Thiin
	LanguageTypeIIR    LanguageType = "iir"      // Indo-Iranian languages
	LanguageTypeIJC    LanguageType = "ijc"      // Izon
	LanguageTypeIJE    LanguageType = "ije"      // Biseni
	LanguageTypeIJJ    LanguageType = "ijj"      // Ede Ije
	LanguageTypeIJN    LanguageType = "ijn"      // Kalabari
	LanguageTypeIJO    LanguageType = "ijo"      // Ijo languages
	LanguageTypeIJS    LanguageType = "ijs"      // Southeast Ijo
	LanguageTypeIKE    LanguageType = "ike"      // Eastern Canadian Inuktitut
	LanguageTypeIKI    LanguageType = "iki"      // Iko
	LanguageTypeIKK    LanguageType = "ikk"      // Ika
	LanguageTypeIKL    LanguageType = "ikl"      // Ikulu
	LanguageTypeIKO    LanguageType = "iko"      // Olulumo-Ikom
	LanguageTypeIKP    LanguageType = "ikp"      // Ikpeshi
	LanguageTypeIKR    LanguageType = "ikr"      // Ikaranggal
	LanguageTypeIKT    LanguageType = "ikt"      // Inuinnaqtun and Western Canadian Inuktitut
	LanguageTypeIKV    LanguageType = "ikv"      // Iku-Gora-Ankwa
	LanguageTypeIKW    LanguageType = "ikw"      // Ikwere
	LanguageTypeIKX    LanguageType = "ikx"      // Ik
	LanguageTypeIKZ    LanguageType = "ikz"      // Ikizu
	LanguageTypeILA    LanguageType = "ila"      // Ile Ape
	LanguageTypeILB    LanguageType = "ilb"      // Ila
	LanguageTypeILG    LanguageType = "ilg"      // Garig-Ilgar
	LanguageTypeILI    LanguageType = "ili"      // Ili Turki
	LanguageTypeILK    LanguageType = "ilk"      // Ilongot
	LanguageTypeILL    LanguageType = "ill"      // Iranun
	LanguageTypeILO    LanguageType = "ilo"      // Iloko
	LanguageTypeILS    LanguageType = "ils"      // International Sign
	LanguageTypeILU    LanguageType = "ilu"      // Ili'uun
	LanguageTypeILV    LanguageType = "ilv"      // Ilue
	LanguageTypeILW    LanguageType = "ilw"      // Talur
	LanguageTypeIMA    LanguageType = "ima"      // Mala Malasar
	LanguageTypeIME    LanguageType = "ime"      // Imeraguen
	LanguageTypeIMI    LanguageType = "imi"      // Anamgura
	LanguageTypeIML    LanguageType = "iml"      // Miluk
	LanguageTypeIMN    LanguageType = "imn"      // Imonda
	LanguageTypeIMO    LanguageType = "imo"      // Imbongu
	LanguageTypeIMR    LanguageType = "imr"      // Imroing
	LanguageTypeIMS    LanguageType = "ims"      // Marsian
	LanguageTypeIMY    LanguageType = "imy"      // Milyan
	LanguageTypeINB    LanguageType = "inb"      // Inga
	LanguageTypeINC    LanguageType = "inc"      // Indic languages
	LanguageTypeINE    LanguageType = "ine"      // Indo-European languages
	LanguageTypeING    LanguageType = "ing"      // Degexit'an
	LanguageTypeINH    LanguageType = "inh"      // Ingush
	LanguageTypeINJ    LanguageType = "inj"      // Jungle Inga
	LanguageTypeINL    LanguageType = "inl"      // Indonesian Sign Language
	LanguageTypeINM    LanguageType = "inm"      // Minaean
	LanguageTypeINN    LanguageType = "inn"      // Isinai
	LanguageTypeINO    LanguageType = "ino"      // Inoke-Yate
	LanguageTypeINP    LanguageType = "inp"      // Iñapari
	LanguageTypeINS    LanguageType = "ins"      // Indian Sign Language
	LanguageTypeINT    LanguageType = "int"      // Intha
	LanguageTypeINZ    LanguageType = "inz"      // Ineseño
	LanguageTypeIOR    LanguageType = "ior"      // Inor
	LanguageTypeIOU    LanguageType = "iou"      // Tuma-Irumu
	LanguageTypeIOW    LanguageType = "iow"      // Iowa-Oto
	LanguageTypeIPI    LanguageType = "ipi"      // Ipili
	LanguageTypeIPO    LanguageType = "ipo"      // Ipiko
	LanguageTypeIQU    LanguageType = "iqu"      // Iquito
	LanguageTypeIQW    LanguageType = "iqw"      // Ikwo
	LanguageTypeIRA    LanguageType = "ira"      // Iranian languages
	LanguageTypeIRE    LanguageType = "ire"      // Iresim
	LanguageTypeIRH    LanguageType = "irh"      // Irarutu
	LanguageTypeIRI    LanguageType = "iri"      // Irigwe
	LanguageTypeIRK    LanguageType = "irk"      // Iraqw
	LanguageTypeIRN    LanguageType = "irn"      // Irántxe
	LanguageTypeIRO    LanguageType = "iro"      // Iroquoian languages
	LanguageTypeIRR    LanguageType = "irr"      // Ir
	LanguageTypeIRU    LanguageType = "iru"      // Irula
	LanguageTypeIRX    LanguageType = "irx"      // Kamberau
	LanguageTypeIRY    LanguageType = "iry"      // Iraya
	LanguageTypeISA    LanguageType = "isa"      // Isabi
	LanguageTypeISC    LanguageType = "isc"      // Isconahua
	LanguageTypeISD    LanguageType = "isd"      // Isnag
	LanguageTypeISE    LanguageType = "ise"      // Italian Sign Language
	LanguageTypeISG    LanguageType = "isg"      // Irish Sign Language
	LanguageTypeISH    LanguageType = "ish"      // Esan
	LanguageTypeISI    LanguageType = "isi"      // Nkem-Nkum
	LanguageTypeISK    LanguageType = "isk"      // Ishkashimi
	LanguageTypeISM    LanguageType = "ism"      // Masimasi
	LanguageTypeISN    LanguageType = "isn"      // Isanzu
	LanguageTypeISO    LanguageType = "iso"      // Isoko
	LanguageTypeISR    LanguageType = "isr"      // Israeli Sign Language
	LanguageTypeIST    LanguageType = "ist"      // Istriot
	LanguageTypeISU    LanguageType = "isu"      // Isu (Menchum Division)
	LanguageTypeITB    LanguageType = "itb"      // Binongan Itneg
	LanguageTypeITC    LanguageType = "itc"      // Italic languages
	LanguageTypeITE    LanguageType = "ite"      // Itene
	LanguageTypeITI    LanguageType = "iti"      // Inlaod Itneg
	LanguageTypeITK    LanguageType = "itk"      // Judeo-Italian
	LanguageTypeITL    LanguageType = "itl"      // Itelmen
	LanguageTypeITM    LanguageType = "itm"      // Itu Mbon Uzo
	LanguageTypeITO    LanguageType = "ito"      // Itonama
	LanguageTypeITR    LanguageType = "itr"      // Iteri
	LanguageTypeITS    LanguageType = "its"      // Isekiri
	LanguageTypeITT    LanguageType = "itt"      // Maeng Itneg
	LanguageTypeITV    LanguageType = "itv"      // Itawit
	LanguageTypeITW    LanguageType = "itw"      // Ito
	LanguageTypeITX    LanguageType = "itx"      // Itik
	LanguageTypeITY    LanguageType = "ity"      // Moyadan Itneg
	LanguageTypeITZ    LanguageType = "itz"      // Itzá
	LanguageTypeIUM    LanguageType = "ium"      // Iu Mien
	LanguageTypeIVB    LanguageType = "ivb"      // Ibatan
	LanguageTypeIVV    LanguageType = "ivv"      // Ivatan
	LanguageTypeIWK    LanguageType = "iwk"      // I-Wak
	LanguageTypeIWM    LanguageType = "iwm"      // Iwam
	LanguageTypeIWO    LanguageType = "iwo"      // Iwur
	LanguageTypeIWS    LanguageType = "iws"      // Sepik Iwam
	LanguageTypeIXC    LanguageType = "ixc"      // Ixcatec
	LanguageTypeIXL    LanguageType = "ixl"      // Ixil
	LanguageTypeIYA    LanguageType = "iya"      // Iyayu
	LanguageTypeIYO    LanguageType = "iyo"      // Mesaka
	LanguageTypeIYX    LanguageType = "iyx"      // Yaka (Congo)
	LanguageTypeIZH    LanguageType = "izh"      // Ingrian
	LanguageTypeIZI    LanguageType = "izi"      // Izi-Ezaa-Ikwo-Mgbo
	LanguageTypeIZR    LanguageType = "izr"      // Izere
	LanguageTypeIZZ    LanguageType = "izz"      // Izii
	LanguageTypeJAA    LanguageType = "jaa"      // Jamamadí
	LanguageTypeJAB    LanguageType = "jab"      // Hyam
	LanguageTypeJAC    LanguageType = "jac"      // Popti' and Jakalteko
	LanguageTypeJAD    LanguageType = "jad"      // Jahanka
	LanguageTypeJAE    LanguageType = "jae"      // Yabem
	LanguageTypeJAF    LanguageType = "jaf"      // Jara
	LanguageTypeJAH    LanguageType = "jah"      // Jah Hut
	LanguageTypeJAJ    LanguageType = "jaj"      // Zazao
	LanguageTypeJAK    LanguageType = "jak"      // Jakun
	LanguageTypeJAL    LanguageType = "jal"      // Yalahatan
	LanguageTypeJAM    LanguageType = "jam"      // Jamaican Creole English
	LanguageTypeJAN    LanguageType = "jan"      // Jandai
	LanguageTypeJAO    LanguageType = "jao"      // Yanyuwa
	LanguageTypeJAQ    LanguageType = "jaq"      // Yaqay
	LanguageTypeJAR    LanguageType = "jar"      // Jarawa (Nigeria)
	LanguageTypeJAS    LanguageType = "jas"      // New Caledonian Javanese
	LanguageTypeJAT    LanguageType = "jat"      // Jakati
	LanguageTypeJAU    LanguageType = "jau"      // Yaur
	LanguageTypeJAX    LanguageType = "jax"      // Jambi Malay
	LanguageTypeJAY    LanguageType = "jay"      // Yan-nhangu
	LanguageTypeJAZ    LanguageType = "jaz"      // Jawe
	LanguageTypeJBE    LanguageType = "jbe"      // Judeo-Berber
	LanguageTypeJBI    LanguageType = "jbi"      // Badjiri
	LanguageTypeJBJ    LanguageType = "jbj"      // Arandai
	LanguageTypeJBK    LanguageType = "jbk"      // Barikewa
	LanguageTypeJBN    LanguageType = "jbn"      // Nafusi
	LanguageTypeJBO    LanguageType = "jbo"      // Lojban
	LanguageTypeJBR    LanguageType = "jbr"      // Jofotek-Bromnya
	LanguageTypeJBT    LanguageType = "jbt"      // Jabutí
	LanguageTypeJBU    LanguageType = "jbu"      // Jukun Takum
	LanguageTypeJBW    LanguageType = "jbw"      // Yawijibaya
	LanguageTypeJCS    LanguageType = "jcs"      // Jamaican Country Sign Language
	LanguageTypeJCT    LanguageType = "jct"      // Krymchak
	LanguageTypeJDA    LanguageType = "jda"      // Jad
	LanguageTypeJDG    LanguageType = "jdg"      // Jadgali
	LanguageTypeJDT    LanguageType = "jdt"      // Judeo-Tat
	LanguageTypeJEB    LanguageType = "jeb"      // Jebero
	LanguageTypeJEE    LanguageType = "jee"      // Jerung
	LanguageTypeJEG    LanguageType = "jeg"      // Jeng
	LanguageTypeJEH    LanguageType = "jeh"      // Jeh
	LanguageTypeJEI    LanguageType = "jei"      // Yei
	LanguageTypeJEK    LanguageType = "jek"      // Jeri Kuo
	LanguageTypeJEL    LanguageType = "jel"      // Yelmek
	LanguageTypeJEN    LanguageType = "jen"      // Dza
	LanguageTypeJER    LanguageType = "jer"      // Jere
	LanguageTypeJET    LanguageType = "jet"      // Manem
	LanguageTypeJEU    LanguageType = "jeu"      // Jonkor Bourmataguil
	LanguageTypeJGB    LanguageType = "jgb"      // Ngbee
	LanguageTypeJGE    LanguageType = "jge"      // Judeo-Georgian
	LanguageTypeJGK    LanguageType = "jgk"      // Gwak
	LanguageTypeJGO    LanguageType = "jgo"      // Ngomba
	LanguageTypeJHI    LanguageType = "jhi"      // Jehai
	LanguageTypeJHS    LanguageType = "jhs"      // Jhankot Sign Language
	LanguageTypeJIA    LanguageType = "jia"      // Jina
	LanguageTypeJIB    LanguageType = "jib"      // Jibu
	LanguageTypeJIC    LanguageType = "jic"      // Tol
	LanguageTypeJID    LanguageType = "jid"      // Bu
	LanguageTypeJIE    LanguageType = "jie"      // Jilbe
	LanguageTypeJIG    LanguageType = "jig"      // Djingili
	LanguageTypeJIH    LanguageType = "jih"      // sTodsde and Shangzhai
	LanguageTypeJII    LanguageType = "jii"      // Jiiddu
	LanguageTypeJIL    LanguageType = "jil"      // Jilim
	LanguageTypeJIM    LanguageType = "jim"      // Jimi (Cameroon)
	LanguageTypeJIO    LanguageType = "jio"      // Jiamao
	LanguageTypeJIQ    LanguageType = "jiq"      // Guanyinqiao and Lavrung
	LanguageTypeJIT    LanguageType = "jit"      // Jita
	LanguageTypeJIU    LanguageType = "jiu"      // Youle Jinuo
	LanguageTypeJIV    LanguageType = "jiv"      // Shuar
	LanguageTypeJIY    LanguageType = "jiy"      // Buyuan Jinuo
	LanguageTypeJJR    LanguageType = "jjr"      // Bankal
	LanguageTypeJKM    LanguageType = "jkm"      // Mobwa Karen
	LanguageTypeJKO    LanguageType = "jko"      // Kubo
	LanguageTypeJKP    LanguageType = "jkp"      // Paku Karen
	LanguageTypeJKR    LanguageType = "jkr"      // Koro (India)
	LanguageTypeJKU    LanguageType = "jku"      // Labir
	LanguageTypeJLE    LanguageType = "jle"      // Ngile
	LanguageTypeJLS    LanguageType = "jls"      // Jamaican Sign Language
	LanguageTypeJMA    LanguageType = "jma"      // Dima
	LanguageTypeJMB    LanguageType = "jmb"      // Zumbun
	LanguageTypeJMC    LanguageType = "jmc"      // Machame
	LanguageTypeJMD    LanguageType = "jmd"      // Yamdena
	LanguageTypeJMI    LanguageType = "jmi"      // Jimi (Nigeria)
	LanguageTypeJML    LanguageType = "jml"      // Jumli
	LanguageTypeJMN    LanguageType = "jmn"      // Makuri Naga
	LanguageTypeJMR    LanguageType = "jmr"      // Kamara
	LanguageTypeJMS    LanguageType = "jms"      // Mashi (Nigeria)
	LanguageTypeJMW    LanguageType = "jmw"      // Mouwase
	LanguageTypeJMX    LanguageType = "jmx"      // Western Juxtlahuaca Mixtec
	LanguageTypeJNA    LanguageType = "jna"      // Jangshung
	LanguageTypeJND    LanguageType = "jnd"      // Jandavra
	LanguageTypeJNG    LanguageType = "jng"      // Yangman
	LanguageTypeJNI    LanguageType = "jni"      // Janji
	LanguageTypeJNJ    LanguageType = "jnj"      // Yemsa
	LanguageTypeJNL    LanguageType = "jnl"      // Rawat
	LanguageTypeJNS    LanguageType = "jns"      // Jaunsari
	LanguageTypeJOB    LanguageType = "job"      // Joba
	LanguageTypeJOD    LanguageType = "jod"      // Wojenaka
	LanguageTypeJOR    LanguageType = "jor"      // Jorá
	LanguageTypeJOS    LanguageType = "jos"      // Jordanian Sign Language
	LanguageTypeJOW    LanguageType = "jow"      // Jowulu
	LanguageTypeJPA    LanguageType = "jpa"      // Jewish Palestinian Aramaic
	LanguageTypeJPR    LanguageType = "jpr"      // Judeo-Persian
	LanguageTypeJPX    LanguageType = "jpx"      // Japanese (family)
	LanguageTypeJQR    LanguageType = "jqr"      // Jaqaru
	LanguageTypeJRA    LanguageType = "jra"      // Jarai
	LanguageTypeJRB    LanguageType = "jrb"      // Judeo-Arabic
	LanguageTypeJRR    LanguageType = "jrr"      // Jiru
	LanguageTypeJRT    LanguageType = "jrt"      // Jorto
	LanguageTypeJRU    LanguageType = "jru"      // Japrería
	LanguageTypeJSL    LanguageType = "jsl"      // Japanese Sign Language
	LanguageTypeJUA    LanguageType = "jua"      // Júma
	LanguageTypeJUB    LanguageType = "jub"      // Wannu
	LanguageTypeJUC    LanguageType = "juc"      // Jurchen
	LanguageTypeJUD    LanguageType = "jud"      // Worodougou
	LanguageTypeJUH    LanguageType = "juh"      // Hõne
	LanguageTypeJUI    LanguageType = "jui"      // Ngadjuri
	LanguageTypeJUK    LanguageType = "juk"      // Wapan
	LanguageTypeJUL    LanguageType = "jul"      // Jirel
	LanguageTypeJUM    LanguageType = "jum"      // Jumjum
	LanguageTypeJUN    LanguageType = "jun"      // Juang
	LanguageTypeJUO    LanguageType = "juo"      // Jiba
	LanguageTypeJUP    LanguageType = "jup"      // Hupdë
	LanguageTypeJUR    LanguageType = "jur"      // Jurúna
	LanguageTypeJUS    LanguageType = "jus"      // Jumla Sign Language
	LanguageTypeJUT    LanguageType = "jut"      // Jutish
	LanguageTypeJUU    LanguageType = "juu"      // Ju
	LanguageTypeJUW    LanguageType = "juw"      // Wãpha
	LanguageTypeJUY    LanguageType = "juy"      // Juray
	LanguageTypeJVD    LanguageType = "jvd"      // Javindo
	LanguageTypeJVN    LanguageType = "jvn"      // Caribbean Javanese
	LanguageTypeJWI    LanguageType = "jwi"      // Jwira-Pepesa
	LanguageTypeJYA    LanguageType = "jya"      // Jiarong
	LanguageTypeJYE    LanguageType = "jye"      // Judeo-Yemeni Arabic
	LanguageTypeJYY    LanguageType = "jyy"      // Jaya
	LanguageTypeKAA    LanguageType = "kaa"      // Kara-Kalpak
	LanguageTypeKAB    LanguageType = "kab"      // Kabyle
	LanguageTypeKAC    LanguageType = "kac"      // Kachin and Jingpho
	LanguageTypeKAD    LanguageType = "kad"      // Adara
	LanguageTypeKAE    LanguageType = "kae"      // Ketangalan
	LanguageTypeKAF    LanguageType = "kaf"      // Katso
	LanguageTypeKAG    LanguageType = "kag"      // Kajaman
	LanguageTypeKAH    LanguageType = "kah"      // Kara (Central African Republic)
	LanguageTypeKAI    LanguageType = "kai"      // Karekare
	LanguageTypeKAJ    LanguageType = "kaj"      // Jju
	LanguageTypeKAK    LanguageType = "kak"      // Kayapa Kallahan
	LanguageTypeKAM    LanguageType = "kam"      // Kamba (Kenya)
	LanguageTypeKAO    LanguageType = "kao"      // Xaasongaxango
	LanguageTypeKAP    LanguageType = "kap"      // Bezhta
	LanguageTypeKAQ    LanguageType = "kaq"      // Capanahua
	LanguageTypeKAR    LanguageType = "kar"      // Karen languages
	LanguageTypeKAV    LanguageType = "kav"      // Katukína
	LanguageTypeKAW    LanguageType = "kaw"      // Kawi
	LanguageTypeKAX    LanguageType = "kax"      // Kao
	LanguageTypeKAY    LanguageType = "kay"      // Kamayurá
	LanguageTypeKBA    LanguageType = "kba"      // Kalarko
	LanguageTypeKBB    LanguageType = "kbb"      // Kaxuiâna
	LanguageTypeKBC    LanguageType = "kbc"      // Kadiwéu
	LanguageTypeKBD    LanguageType = "kbd"      // Kabardian
	LanguageTypeKBE    LanguageType = "kbe"      // Kanju
	LanguageTypeKBF    LanguageType = "kbf"      // Kakauhua
	LanguageTypeKBG    LanguageType = "kbg"      // Khamba
	LanguageTypeKBH    LanguageType = "kbh"      // Camsá
	LanguageTypeKBI    LanguageType = "kbi"      // Kaptiau
	LanguageTypeKBJ    LanguageType = "kbj"      // Kari
	LanguageTypeKBK    LanguageType = "kbk"      // Grass Koiari
	LanguageTypeKBL    LanguageType = "kbl"      // Kanembu
	LanguageTypeKBM    LanguageType = "kbm"      // Iwal
	LanguageTypeKBN    LanguageType = "kbn"      // Kare (Central African Republic)
	LanguageTypeKBO    LanguageType = "kbo"      // Keliko
	LanguageTypeKBP    LanguageType = "kbp"      // Kabiyè
	LanguageTypeKBQ    LanguageType = "kbq"      // Kamano
	LanguageTypeKBR    LanguageType = "kbr"      // Kafa
	LanguageTypeKBS    LanguageType = "kbs"      // Kande
	LanguageTypeKBT    LanguageType = "kbt"      // Abadi
	LanguageTypeKBU    LanguageType = "kbu"      // Kabutra
	LanguageTypeKBV    LanguageType = "kbv"      // Dera (Indonesia)
	LanguageTypeKBW    LanguageType = "kbw"      // Kaiep
	LanguageTypeKBX    LanguageType = "kbx"      // Ap Ma
	LanguageTypeKBY    LanguageType = "kby"      // Manga Kanuri
	LanguageTypeKBZ    LanguageType = "kbz"      // Duhwa
	LanguageTypeKCA    LanguageType = "kca"      // Khanty
	LanguageTypeKCB    LanguageType = "kcb"      // Kawacha
	LanguageTypeKCC    LanguageType = "kcc"      // Lubila
	LanguageTypeKCD    LanguageType = "kcd"      // Ngkâlmpw Kanum
	LanguageTypeKCE    LanguageType = "kce"      // Kaivi
	LanguageTypeKCF    LanguageType = "kcf"      // Ukaan
	LanguageTypeKCG    LanguageType = "kcg"      // Tyap
	LanguageTypeKCH    LanguageType = "kch"      // Vono
	LanguageTypeKCI    LanguageType = "kci"      // Kamantan
	LanguageTypeKCJ    LanguageType = "kcj"      // Kobiana
	LanguageTypeKCK    LanguageType = "kck"      // Kalanga
	LanguageTypeKCL    LanguageType = "kcl"      // Kela (Papua New Guinea) and Kala
	LanguageTypeKCM    LanguageType = "kcm"      // Gula (Central African Republic)
	LanguageTypeKCN    LanguageType = "kcn"      // Nubi
	LanguageTypeKCO    LanguageType = "kco"      // Kinalakna
	LanguageTypeKCP    LanguageType = "kcp"      // Kanga
	LanguageTypeKCQ    LanguageType = "kcq"      // Kamo
	LanguageTypeKCR    LanguageType = "kcr"      // Katla
	LanguageTypeKCS    LanguageType = "kcs"      // Koenoem
	LanguageTypeKCT    LanguageType = "kct"      // Kaian
	LanguageTypeKCU    LanguageType = "kcu"      // Kami (Tanzania)
	LanguageTypeKCV    LanguageType = "kcv"      // Kete
	LanguageTypeKCW    LanguageType = "kcw"      // Kabwari
	LanguageTypeKCX    LanguageType = "kcx"      // Kachama-Ganjule
	LanguageTypeKCY    LanguageType = "kcy"      // Korandje
	LanguageTypeKCZ    LanguageType = "kcz"      // Konongo
	LanguageTypeKDA    LanguageType = "kda"      // Worimi
	LanguageTypeKDC    LanguageType = "kdc"      // Kutu
	LanguageTypeKDD    LanguageType = "kdd"      // Yankunytjatjara
	LanguageTypeKDE    LanguageType = "kde"      // Makonde
	LanguageTypeKDF    LanguageType = "kdf"      // Mamusi
	LanguageTypeKDG    LanguageType = "kdg"      // Seba
	LanguageTypeKDH    LanguageType = "kdh"      // Tem
	LanguageTypeKDI    LanguageType = "kdi"      // Kumam
	LanguageTypeKDJ    LanguageType = "kdj"      // Karamojong
	LanguageTypeKDK    LanguageType = "kdk"      // Numèè and Kwényi
	LanguageTypeKDL    LanguageType = "kdl"      // Tsikimba
	LanguageTypeKDM    LanguageType = "kdm"      // Kagoma
	LanguageTypeKDN    LanguageType = "kdn"      // Kunda
	LanguageTypeKDO    LanguageType = "kdo"      // Kordofanian languages
	LanguageTypeKDP    LanguageType = "kdp"      // Kaningdon-Nindem
	LanguageTypeKDQ    LanguageType = "kdq"      // Koch
	LanguageTypeKDR    LanguageType = "kdr"      // Karaim
	LanguageTypeKDT    LanguageType = "kdt"      // Kuy
	LanguageTypeKDU    LanguageType = "kdu"      // Kadaru
	LanguageTypeKDV    LanguageType = "kdv"      // Kado
	LanguageTypeKDW    LanguageType = "kdw"      // Koneraw
	LanguageTypeKDX    LanguageType = "kdx"      // Kam
	LanguageTypeKDY    LanguageType = "kdy"      // Keder and Keijar
	LanguageTypeKDZ    LanguageType = "kdz"      // Kwaja
	LanguageTypeKEA    LanguageType = "kea"      // Kabuverdianu
	LanguageTypeKEB    LanguageType = "keb"      // Kélé
	LanguageTypeKEC    LanguageType = "kec"      // Keiga
	LanguageTypeKED    LanguageType = "ked"      // Kerewe
	LanguageTypeKEE    LanguageType = "kee"      // Eastern Keres
	LanguageTypeKEF    LanguageType = "kef"      // Kpessi
	LanguageTypeKEG    LanguageType = "keg"      // Tese
	LanguageTypeKEH    LanguageType = "keh"      // Keak
	LanguageTypeKEI    LanguageType = "kei"      // Kei
	LanguageTypeKEJ    LanguageType = "kej"      // Kadar
	LanguageTypeKEK    LanguageType = "kek"      // Kekchí
	LanguageTypeKEL    LanguageType = "kel"      // Kela (Democratic Republic of Congo)
	LanguageTypeKEM    LanguageType = "kem"      // Kemak
	LanguageTypeKEN    LanguageType = "ken"      // Kenyang
	LanguageTypeKEO    LanguageType = "keo"      // Kakwa
	LanguageTypeKEP    LanguageType = "kep"      // Kaikadi
	LanguageTypeKEQ    LanguageType = "keq"      // Kamar
	LanguageTypeKER    LanguageType = "ker"      // Kera
	LanguageTypeKES    LanguageType = "kes"      // Kugbo
	LanguageTypeKET    LanguageType = "ket"      // Ket
	LanguageTypeKEU    LanguageType = "keu"      // Akebu
	LanguageTypeKEV    LanguageType = "kev"      // Kanikkaran
	LanguageTypeKEW    LanguageType = "kew"      // West Kewa
	LanguageTypeKEX    LanguageType = "kex"      // Kukna
	LanguageTypeKEY    LanguageType = "key"      // Kupia
	LanguageTypeKEZ    LanguageType = "kez"      // Kukele
	LanguageTypeKFA    LanguageType = "kfa"      // Kodava
	LanguageTypeKFB    LanguageType = "kfb"      // Northwestern Kolami
	LanguageTypeKFC    LanguageType = "kfc"      // Konda-Dora
	LanguageTypeKFD    LanguageType = "kfd"      // Korra Koraga
	LanguageTypeKFE    LanguageType = "kfe"      // Kota (India)
	LanguageTypeKFF    LanguageType = "kff"      // Koya
	LanguageTypeKFG    LanguageType = "kfg"      // Kudiya
	LanguageTypeKFH    LanguageType = "kfh"      // Kurichiya
	LanguageTypeKFI    LanguageType = "kfi"      // Kannada Kurumba
	LanguageTypeKFJ    LanguageType = "kfj"      // Kemiehua
	LanguageTypeKFK    LanguageType = "kfk"      // Kinnauri
	LanguageTypeKFL    LanguageType = "kfl"      // Kung
	LanguageTypeKFM    LanguageType = "kfm"      // Khunsari
	LanguageTypeKFN    LanguageType = "kfn"      // Kuk
	LanguageTypeKFO    LanguageType = "kfo"      // Koro (Côte d'Ivoire)
	LanguageTypeKFP    LanguageType = "kfp"      // Korwa
	LanguageTypeKFQ    LanguageType = "kfq"      // Korku
	LanguageTypeKFR    LanguageType = "kfr"      // Kachchi
	LanguageTypeKFS    LanguageType = "kfs"      // Bilaspuri
	LanguageTypeKFT    LanguageType = "kft"      // Kanjari
	LanguageTypeKFU    LanguageType = "kfu"      // Katkari
	LanguageTypeKFV    LanguageType = "kfv"      // Kurmukar
	LanguageTypeKFW    LanguageType = "kfw"      // Kharam Naga
	LanguageTypeKFX    LanguageType = "kfx"      // Kullu Pahari
	LanguageTypeKFY    LanguageType = "kfy"      // Kumaoni
	LanguageTypeKFZ    LanguageType = "kfz"      // Koromfé
	LanguageTypeKGA    LanguageType = "kga"      // Koyaga
	LanguageTypeKGB    LanguageType = "kgb"      // Kawe
	LanguageTypeKGC    LanguageType = "kgc"      // Kasseng
	LanguageTypeKGD    LanguageType = "kgd"      // Kataang
	LanguageTypeKGE    LanguageType = "kge"      // Komering
	LanguageTypeKGF    LanguageType = "kgf"      // Kube
	LanguageTypeKGG    LanguageType = "kgg"      // Kusunda
	LanguageTypeKGH    LanguageType = "kgh"      // Upper Tanudan Kalinga
	LanguageTypeKGI    LanguageType = "kgi"      // Selangor Sign Language
	LanguageTypeKGJ    LanguageType = "kgj"      // Gamale Kham
	LanguageTypeKGK    LanguageType = "kgk"      // Kaiwá
	LanguageTypeKGL    LanguageType = "kgl"      // Kunggari
	LanguageTypeKGM    LanguageType = "kgm"      // Karipúna
	LanguageTypeKGN    LanguageType = "kgn"      // Karingani
	LanguageTypeKGO    LanguageType = "kgo"      // Krongo
	LanguageTypeKGP    LanguageType = "kgp"      // Kaingang
	LanguageTypeKGQ    LanguageType = "kgq"      // Kamoro
	LanguageTypeKGR    LanguageType = "kgr"      // Abun
	LanguageTypeKGS    LanguageType = "kgs"      // Kumbainggar
	LanguageTypeKGT    LanguageType = "kgt"      // Somyev
	LanguageTypeKGU    LanguageType = "kgu"      // Kobol
	LanguageTypeKGV    LanguageType = "kgv"      // Karas
	LanguageTypeKGW    LanguageType = "kgw"      // Karon Dori
	LanguageTypeKGX    LanguageType = "kgx"      // Kamaru
	LanguageTypeKGY    LanguageType = "kgy"      // Kyerung
	LanguageTypeKHA    LanguageType = "kha"      // Khasi
	LanguageTypeKHB    LanguageType = "khb"      // Lü
	LanguageTypeKHC    LanguageType = "khc"      // Tukang Besi North
	LanguageTypeKHD    LanguageType = "khd"      // Bädi Kanum
	LanguageTypeKHE    LanguageType = "khe"      // Korowai
	LanguageTypeKHF    LanguageType = "khf"      // Khuen
	LanguageTypeKHG    LanguageType = "khg"      // Khams Tibetan
	LanguageTypeKHH    LanguageType = "khh"      // Kehu
	LanguageTypeKHI    LanguageType = "khi"      // Khoisan languages
	LanguageTypeKHJ    LanguageType = "khj"      // Kuturmi
	LanguageTypeKHK    LanguageType = "khk"      // Halh Mongolian
	LanguageTypeKHL    LanguageType = "khl"      // Lusi
	LanguageTypeKHN    LanguageType = "khn"      // Khandesi
	LanguageTypeKHO    LanguageType = "kho"      // Khotanese and Sakan
	LanguageTypeKHP    LanguageType = "khp"      // Kapori and Kapauri
	LanguageTypeKHQ    LanguageType = "khq"      // Koyra Chiini Songhay
	LanguageTypeKHR    LanguageType = "khr"      // Kharia
	LanguageTypeKHS    LanguageType = "khs"      // Kasua
	LanguageTypeKHT    LanguageType = "kht"      // Khamti
	LanguageTypeKHU    LanguageType = "khu"      // Nkhumbi
	LanguageTypeKHV    LanguageType = "khv"      // Khvarshi
	LanguageTypeKHW    LanguageType = "khw"      // Khowar
	LanguageTypeKHX    LanguageType = "khx"      // Kanu
	LanguageTypeKHY    LanguageType = "khy"      // Kele (Democratic Republic of Congo)
	LanguageTypeKHZ    LanguageType = "khz"      // Keapara
	LanguageTypeKIA    LanguageType = "kia"      // Kim
	LanguageTypeKIB    LanguageType = "kib"      // Koalib
	LanguageTypeKIC    LanguageType = "kic"      // Kickapoo
	LanguageTypeKID    LanguageType = "kid"      // Koshin
	LanguageTypeKIE    LanguageType = "kie"      // Kibet
	LanguageTypeKIF    LanguageType = "kif"      // Eastern Parbate Kham
	LanguageTypeKIG    LanguageType = "kig"      // Kimaama and Kimaghima
	LanguageTypeKIH    LanguageType = "kih"      // Kilmeri
	LanguageTypeKII    LanguageType = "kii"      // Kitsai
	LanguageTypeKIJ    LanguageType = "kij"      // Kilivila
	LanguageTypeKIL    LanguageType = "kil"      // Kariya
	LanguageTypeKIM    LanguageType = "kim"      // Karagas
	LanguageTypeKIO    LanguageType = "kio"      // Kiowa
	LanguageTypeKIP    LanguageType = "kip"      // Sheshi Kham
	LanguageTypeKIQ    LanguageType = "kiq"      // Kosadle and Kosare
	LanguageTypeKIS    LanguageType = "kis"      // Kis
	LanguageTypeKIT    LanguageType = "kit"      // Agob
	LanguageTypeKIU    LanguageType = "kiu"      // Kirmanjki (individual language)
	LanguageTypeKIV    LanguageType = "kiv"      // Kimbu
	LanguageTypeKIW    LanguageType = "kiw"      // Northeast Kiwai
	LanguageTypeKIX    LanguageType = "kix"      // Khiamniungan Naga
	LanguageTypeKIY    LanguageType = "kiy"      // Kirikiri
	LanguageTypeKIZ    LanguageType = "kiz"      // Kisi
	LanguageTypeKJA    LanguageType = "kja"      // Mlap
	LanguageTypeKJB    LanguageType = "kjb"      // Q'anjob'al and Kanjobal
	LanguageTypeKJC    LanguageType = "kjc"      // Coastal Konjo
	LanguageTypeKJD    LanguageType = "kjd"      // Southern Kiwai
	LanguageTypeKJE    LanguageType = "kje"      // Kisar
	LanguageTypeKJF    LanguageType = "kjf"      // Khalaj
	LanguageTypeKJG    LanguageType = "kjg"      // Khmu
	LanguageTypeKJH    LanguageType = "kjh"      // Khakas
	LanguageTypeKJI    LanguageType = "kji"      // Zabana
	LanguageTypeKJJ    LanguageType = "kjj"      // Khinalugh
	LanguageTypeKJK    LanguageType = "kjk"      // Highland Konjo
	LanguageTypeKJL    LanguageType = "kjl"      // Western Parbate Kham
	LanguageTypeKJM    LanguageType = "kjm"      // Kháng
	LanguageTypeKJN    LanguageType = "kjn"      // Kunjen
	LanguageTypeKJO    LanguageType = "kjo"      // Harijan Kinnauri
	LanguageTypeKJP    LanguageType = "kjp"      // Pwo Eastern Karen
	LanguageTypeKJQ    LanguageType = "kjq"      // Western Keres
	LanguageTypeKJR    LanguageType = "kjr"      // Kurudu
	LanguageTypeKJS    LanguageType = "kjs"      // East Kewa
	LanguageTypeKJT    LanguageType = "kjt"      // Phrae Pwo Karen
	LanguageTypeKJU    LanguageType = "kju"      // Kashaya
	LanguageTypeKJX    LanguageType = "kjx"      // Ramopa
	LanguageTypeKJY    LanguageType = "kjy"      // Erave
	LanguageTypeKJZ    LanguageType = "kjz"      // Bumthangkha
	LanguageTypeKKA    LanguageType = "kka"      // Kakanda
	LanguageTypeKKB    LanguageType = "kkb"      // Kwerisa
	LanguageTypeKKC    LanguageType = "kkc"      // Odoodee
	LanguageTypeKKD    LanguageType = "kkd"      // Kinuku
	LanguageTypeKKE    LanguageType = "kke"      // Kakabe
	LanguageTypeKKF    LanguageType = "kkf"      // Kalaktang Monpa
	LanguageTypeKKG    LanguageType = "kkg"      // Mabaka Valley Kalinga
	LanguageTypeKKH    LanguageType = "kkh"      // Khün
	LanguageTypeKKI    LanguageType = "kki"      // Kagulu
	LanguageTypeKKJ    LanguageType = "kkj"      // Kako
	LanguageTypeKKK    LanguageType = "kkk"      // Kokota
	LanguageTypeKKL    LanguageType = "kkl"      // Kosarek Yale
	LanguageTypeKKM    LanguageType = "kkm"      // Kiong
	LanguageTypeKKN    LanguageType = "kkn"      // Kon Keu
	LanguageTypeKKO    LanguageType = "kko"      // Karko
	LanguageTypeKKP    LanguageType = "kkp"      // Gugubera
	LanguageTypeKKQ    LanguageType = "kkq"      // Kaiku
	LanguageTypeKKR    LanguageType = "kkr"      // Kir-Balar
	LanguageTypeKKS    LanguageType = "kks"      // Giiwo
	LanguageTypeKKT    LanguageType = "kkt"      // Koi
	LanguageTypeKKU    LanguageType = "kku"      // Tumi
	LanguageTypeKKV    LanguageType = "kkv"      // Kangean
	LanguageTypeKKW    LanguageType = "kkw"      // Teke-Kukuya
	LanguageTypeKKX    LanguageType = "kkx"      // Kohin
	LanguageTypeKKY    LanguageType = "kky"      // Guguyimidjir
	LanguageTypeKKZ    LanguageType = "kkz"      // Kaska
	LanguageTypeKLA    LanguageType = "kla"      // Klamath-Modoc
	LanguageTypeKLB    LanguageType = "klb"      // Kiliwa
	LanguageTypeKLC    LanguageType = "klc"      // Kolbila
	LanguageTypeKLD    LanguageType = "kld"      // Gamilaraay
	LanguageTypeKLE    LanguageType = "kle"      // Kulung (Nepal)
	LanguageTypeKLF    LanguageType = "klf"      // Kendeje
	LanguageTypeKLG    LanguageType = "klg"      // Tagakaulo
	LanguageTypeKLH    LanguageType = "klh"      // Weliki
	LanguageTypeKLI    LanguageType = "kli"      // Kalumpang
	LanguageTypeKLJ    LanguageType = "klj"      // Turkic Khalaj
	LanguageTypeKLK    LanguageType = "klk"      // Kono (Nigeria)
	LanguageTypeKLL    LanguageType = "kll"      // Kagan Kalagan
	LanguageTypeKLM    LanguageType = "klm"      // Migum
	LanguageTypeKLN    LanguageType = "kln"      // Kalenjin
	LanguageTypeKLO    LanguageType = "klo"      // Kapya
	LanguageTypeKLP    LanguageType = "klp"      // Kamasa
	LanguageTypeKLQ    LanguageType = "klq"      // Rumu
	LanguageTypeKLR    LanguageType = "klr"      // Khaling
	LanguageTypeKLS    LanguageType = "kls"      // Kalasha
	LanguageTypeKLT    LanguageType = "klt"      // Nukna
	LanguageTypeKLU    LanguageType = "klu"      // Klao
	LanguageTypeKLV    LanguageType = "klv"      // Maskelynes
	LanguageTypeKLW    LanguageType = "klw"      // Lindu
	LanguageTypeKLX    LanguageType = "klx"      // Koluwawa
	LanguageTypeKLY    LanguageType = "kly"      // Kalao
	LanguageTypeKLZ    LanguageType = "klz"      // Kabola
	LanguageTypeKMA    LanguageType = "kma"      // Konni
	LanguageTypeKMB    LanguageType = "kmb"      // Kimbundu
	LanguageTypeKMC    LanguageType = "kmc"      // Southern Dong
	LanguageTypeKMD    LanguageType = "kmd"      // Majukayang Kalinga
	LanguageTypeKME    LanguageType = "kme"      // Bakole
	LanguageTypeKMF    LanguageType = "kmf"      // Kare (Papua New Guinea)
	LanguageTypeKMG    LanguageType = "kmg"      // Kâte
	LanguageTypeKMH    LanguageType = "kmh"      // Kalam
	LanguageTypeKMI    LanguageType = "kmi"      // Kami (Nigeria)
	LanguageTypeKMJ    LanguageType = "kmj"      // Kumarbhag Paharia
	LanguageTypeKMK    LanguageType = "kmk"      // Limos Kalinga
	LanguageTypeKML    LanguageType = "kml"      // Tanudan Kalinga
	LanguageTypeKMM    LanguageType = "kmm"      // Kom (India)
	LanguageTypeKMN    LanguageType = "kmn"      // Awtuw
	LanguageTypeKMO    LanguageType = "kmo"      // Kwoma
	LanguageTypeKMP    LanguageType = "kmp"      // Gimme
	LanguageTypeKMQ    LanguageType = "kmq"      // Kwama
	LanguageTypeKMR    LanguageType = "kmr"      // Northern Kurdish
	LanguageTypeKMS    LanguageType = "kms"      // Kamasau
	LanguageTypeKMT    LanguageType = "kmt"      // Kemtuik
	LanguageTypeKMU    LanguageType = "kmu"      // Kanite
	LanguageTypeKMV    LanguageType = "kmv"      // Karipúna Creole French
	LanguageTypeKMW    LanguageType = "kmw"      // Komo (Democratic Republic of Congo)
	LanguageTypeKMX    LanguageType = "kmx"      // Waboda
	LanguageTypeKMY    LanguageType = "kmy"      // Koma
	LanguageTypeKMZ    LanguageType = "kmz"      // Khorasani Turkish
	LanguageTypeKNA    LanguageType = "kna"      // Dera (Nigeria)
	LanguageTypeKNB    LanguageType = "knb"      // Lubuagan Kalinga
	LanguageTypeKNC    LanguageType = "knc"      // Central Kanuri
	LanguageTypeKND    LanguageType = "knd"      // Konda
	LanguageTypeKNE    LanguageType = "kne"      // Kankanaey
	LanguageTypeKNF    LanguageType = "knf"      // Mankanya
	LanguageTypeKNG    LanguageType = "kng"      // Koongo
	LanguageTypeKNI    LanguageType = "kni"      // Kanufi
	LanguageTypeKNJ    LanguageType = "knj"      // Western Kanjobal
	LanguageTypeKNK    LanguageType = "knk"      // Kuranko
	LanguageTypeKNL    LanguageType = "knl"      // Keninjal
	LanguageTypeKNM    LanguageType = "knm"      // Kanamarí
	LanguageTypeKNN    LanguageType = "knn"      // Konkani (individual language)
	LanguageTypeKNO    LanguageType = "kno"      // Kono (Sierra Leone)
	LanguageTypeKNP    LanguageType = "knp"      // Kwanja
	LanguageTypeKNQ    LanguageType = "knq"      // Kintaq
	LanguageTypeKNR    LanguageType = "knr"      // Kaningra
	LanguageTypeKNS    LanguageType = "kns"      // Kensiu
	LanguageTypeKNT    LanguageType = "knt"      // Panoan Katukína
	LanguageTypeKNU    LanguageType = "knu"      // Kono (Guinea)
	LanguageTypeKNV    LanguageType = "knv"      // Tabo
	LanguageTypeKNW    LanguageType = "knw"      // Kung-Ekoka
	LanguageTypeKNX    LanguageType = "knx"      // Kendayan and Salako
	LanguageTypeKNY    LanguageType = "kny"      // Kanyok
	LanguageTypeKNZ    LanguageType = "knz"      // Kalamsé
	LanguageTypeKOA    LanguageType = "koa"      // Konomala
	LanguageTypeKOC    LanguageType = "koc"      // Kpati
	LanguageTypeKOD    LanguageType = "kod"      // Kodi
	LanguageTypeKOE    LanguageType = "koe"      // Kacipo-Balesi
	LanguageTypeKOF    LanguageType = "kof"      // Kubi
	LanguageTypeKOG    LanguageType = "kog"      // Cogui and Kogi
	LanguageTypeKOH    LanguageType = "koh"      // Koyo
	LanguageTypeKOI    LanguageType = "koi"      // Komi-Permyak
	LanguageTypeKOJ    LanguageType = "koj"      // Sara Dunjo
	LanguageTypeKOK    LanguageType = "kok"      // Konkani (macrolanguage)
	LanguageTypeKOL    LanguageType = "kol"      // Kol (Papua New Guinea)
	LanguageTypeKOO    LanguageType = "koo"      // Konzo
	LanguageTypeKOP    LanguageType = "kop"      // Waube
	LanguageTypeKOQ    LanguageType = "koq"      // Kota (Gabon)
	LanguageTypeKOS    LanguageType = "kos"      // Kosraean
	LanguageTypeKOT    LanguageType = "kot"      // Lagwan
	LanguageTypeKOU    LanguageType = "kou"      // Koke
	LanguageTypeKOV    LanguageType = "kov"      // Kudu-Camo
	LanguageTypeKOW    LanguageType = "kow"      // Kugama
	LanguageTypeKOX    LanguageType = "kox"      // Coxima
	LanguageTypeKOY    LanguageType = "koy"      // Koyukon
	LanguageTypeKOZ    LanguageType = "koz"      // Korak
	LanguageTypeKPA    LanguageType = "kpa"      // Kutto
	LanguageTypeKPB    LanguageType = "kpb"      // Mullu Kurumba
	LanguageTypeKPC    LanguageType = "kpc"      // Curripaco
	LanguageTypeKPD    LanguageType = "kpd"      // Koba
	LanguageTypeKPE    LanguageType = "kpe"      // Kpelle
	LanguageTypeKPF    LanguageType = "kpf"      // Komba
	LanguageTypeKPG    LanguageType = "kpg"      // Kapingamarangi
	LanguageTypeKPH    LanguageType = "kph"      // Kplang
	LanguageTypeKPI    LanguageType = "kpi"      // Kofei
	LanguageTypeKPJ    LanguageType = "kpj"      // Karajá
	LanguageTypeKPK    LanguageType = "kpk"      // Kpan
	LanguageTypeKPL    LanguageType = "kpl"      // Kpala
	LanguageTypeKPM    LanguageType = "kpm"      // Koho
	LanguageTypeKPN    LanguageType = "kpn"      // Kepkiriwát
	LanguageTypeKPO    LanguageType = "kpo"      // Ikposo
	LanguageTypeKPP    LanguageType = "kpp"      // Paku Karen
	LanguageTypeKPQ    LanguageType = "kpq"      // Korupun-Sela
	LanguageTypeKPR    LanguageType = "kpr"      // Korafe-Yegha
	LanguageTypeKPS    LanguageType = "kps"      // Tehit
	LanguageTypeKPT    LanguageType = "kpt"      // Karata
	LanguageTypeKPU    LanguageType = "kpu"      // Kafoa
	LanguageTypeKPV    LanguageType = "kpv"      // Komi-Zyrian
	LanguageTypeKPW    LanguageType = "kpw"      // Kobon
	LanguageTypeKPX    LanguageType = "kpx"      // Mountain Koiali
	LanguageTypeKPY    LanguageType = "kpy"      // Koryak
	LanguageTypeKPZ    LanguageType = "kpz"      // Kupsabiny
	LanguageTypeKQA    LanguageType = "kqa"      // Mum
	LanguageTypeKQB    LanguageType = "kqb"      // Kovai
	LanguageTypeKQC    LanguageType = "kqc"      // Doromu-Koki
	LanguageTypeKQD    LanguageType = "kqd"      // Koy Sanjaq Surat
	LanguageTypeKQE    LanguageType = "kqe"      // Kalagan
	LanguageTypeKQF    LanguageType = "kqf"      // Kakabai
	LanguageTypeKQG    LanguageType = "kqg"      // Khe
	LanguageTypeKQH    LanguageType = "kqh"      // Kisankasa
	LanguageTypeKQI    LanguageType = "kqi"      // Koitabu
	LanguageTypeKQJ    LanguageType = "kqj"      // Koromira
	LanguageTypeKQK    LanguageType = "kqk"      // Kotafon Gbe
	LanguageTypeKQL    LanguageType = "kql"      // Kyenele
	LanguageTypeKQM    LanguageType = "kqm"      // Khisa
	LanguageTypeKQN    LanguageType = "kqn"      // Kaonde
	LanguageTypeKQO    LanguageType = "kqo"      // Eastern Krahn
	LanguageTypeKQP    LanguageType = "kqp"      // Kimré
	LanguageTypeKQQ    LanguageType = "kqq"      // Krenak
	LanguageTypeKQR    LanguageType = "kqr"      // Kimaragang
	LanguageTypeKQS    LanguageType = "kqs"      // Northern Kissi
	LanguageTypeKQT    LanguageType = "kqt"      // Klias River Kadazan
	LanguageTypeKQU    LanguageType = "kqu"      // Seroa
	LanguageTypeKQV    LanguageType = "kqv"      // Okolod
	LanguageTypeKQW    LanguageType = "kqw"      // Kandas
	LanguageTypeKQX    LanguageType = "kqx"      // Mser
	LanguageTypeKQY    LanguageType = "kqy"      // Koorete
	LanguageTypeKQZ    LanguageType = "kqz"      // Korana
	LanguageTypeKRA    LanguageType = "kra"      // Kumhali
	LanguageTypeKRB    LanguageType = "krb"      // Karkin
	LanguageTypeKRC    LanguageType = "krc"      // Karachay-Balkar
	LanguageTypeKRD    LanguageType = "krd"      // Kairui-Midiki
	LanguageTypeKRE    LanguageType = "kre"      // Panará
	LanguageTypeKRF    LanguageType = "krf"      // Koro (Vanuatu)
	LanguageTypeKRH    LanguageType = "krh"      // Kurama
	LanguageTypeKRI    LanguageType = "kri"      // Krio
	LanguageTypeKRJ    LanguageType = "krj"      // Kinaray-A
	LanguageTypeKRK    LanguageType = "krk"      // Kerek
	LanguageTypeKRL    LanguageType = "krl"      // Karelian
	LanguageTypeKRM    LanguageType = "krm"      // Krim
	LanguageTypeKRN    LanguageType = "krn"      // Sapo
	LanguageTypeKRO    LanguageType = "kro"      // Kru languages
	LanguageTypeKRP    LanguageType = "krp"      // Korop
	LanguageTypeKRR    LanguageType = "krr"      // Kru'ng 2
	LanguageTypeKRS    LanguageType = "krs"      // Gbaya (Sudan)
	LanguageTypeKRT    LanguageType = "krt"      // Tumari Kanuri
	LanguageTypeKRU    LanguageType = "kru"      // Kurukh
	LanguageTypeKRV    LanguageType = "krv"      // Kavet
	LanguageTypeKRW    LanguageType = "krw"      // Western Krahn
	LanguageTypeKRX    LanguageType = "krx"      // Karon
	LanguageTypeKRY    LanguageType = "kry"      // Kryts
	LanguageTypeKRZ    LanguageType = "krz"      // Sota Kanum
	LanguageTypeKSA    LanguageType = "ksa"      // Shuwa-Zamani
	LanguageTypeKSB    LanguageType = "ksb"      // Shambala
	LanguageTypeKSC    LanguageType = "ksc"      // Southern Kalinga
	LanguageTypeKSD    LanguageType = "ksd"      // Kuanua
	LanguageTypeKSE    LanguageType = "kse"      // Kuni
	LanguageTypeKSF    LanguageType = "ksf"      // Bafia
	LanguageTypeKSG    LanguageType = "ksg"      // Kusaghe
	LanguageTypeKSH    LanguageType = "ksh"      // Kölsch
	LanguageTypeKSI    LanguageType = "ksi"      // Krisa and I'saka
	LanguageTypeKSJ    LanguageType = "ksj"      // Uare
	LanguageTypeKSK    LanguageType = "ksk"      // Kansa
	LanguageTypeKSL    LanguageType = "ksl"      // Kumalu
	LanguageTypeKSM    LanguageType = "ksm"      // Kumba
	LanguageTypeKSN    LanguageType = "ksn"      // Kasiguranin
	LanguageTypeKSO    LanguageType = "kso"      // Kofa
	LanguageTypeKSP    LanguageType = "ksp"      // Kaba
	LanguageTypeKSQ    LanguageType = "ksq"      // Kwaami
	LanguageTypeKSR    LanguageType = "ksr"      // Borong
	LanguageTypeKSS    LanguageType = "kss"      // Southern Kisi
	LanguageTypeKST    LanguageType = "kst"      // Winyé
	LanguageTypeKSU    LanguageType = "ksu"      // Khamyang
	LanguageTypeKSV    LanguageType = "ksv"      // Kusu
	LanguageTypeKSW    LanguageType = "ksw"      // S'gaw Karen
	LanguageTypeKSX    LanguageType = "ksx"      // Kedang
	LanguageTypeKSY    LanguageType = "ksy"      // Kharia Thar
	LanguageTypeKSZ    LanguageType = "ksz"      // Kodaku
	LanguageTypeKTA    LanguageType = "kta"      // Katua
	LanguageTypeKTB    LanguageType = "ktb"      // Kambaata
	LanguageTypeKTC    LanguageType = "ktc"      // Kholok
	LanguageTypeKTD    LanguageType = "ktd"      // Kokata
	LanguageTypeKTE    LanguageType = "kte"      // Nubri
	LanguageTypeKTF    LanguageType = "ktf"      // Kwami
	LanguageTypeKTG    LanguageType = "ktg"      // Kalkutung
	LanguageTypeKTH    LanguageType = "kth"      // Karanga
	LanguageTypeKTI    LanguageType = "kti"      // North Muyu
	LanguageTypeKTJ    LanguageType = "ktj"      // Plapo Krumen
	LanguageTypeKTK    LanguageType = "ktk"      // Kaniet
	LanguageTypeKTL    LanguageType = "ktl"      // Koroshi
	LanguageTypeKTM    LanguageType = "ktm"      // Kurti
	LanguageTypeKTN    LanguageType = "ktn"      // Karitiâna
	LanguageTypeKTO    LanguageType = "kto"      // Kuot
	LanguageTypeKTP    LanguageType = "ktp"      // Kaduo
	LanguageTypeKTQ    LanguageType = "ktq"      // Katabaga
	LanguageTypeKTR    LanguageType = "ktr"      // Kota Marudu Tinagas
	LanguageTypeKTS    LanguageType = "kts"      // South Muyu
	LanguageTypeKTT    LanguageType = "ktt"      // Ketum
	LanguageTypeKTU    LanguageType = "ktu"      // Kituba (Democratic Republic of Congo)
	LanguageTypeKTV    LanguageType = "ktv"      // Eastern Katu
	LanguageTypeKTW    LanguageType = "ktw"      // Kato
	LanguageTypeKTX    LanguageType = "ktx"      // Kaxararí
	LanguageTypeKTY    LanguageType = "kty"      // Kango (Bas-Uélé District)
	LanguageTypeKTZ    LanguageType = "ktz"      // Ju/'hoan
	LanguageTypeKUB    LanguageType = "kub"      // Kutep
	LanguageTypeKUC    LanguageType = "kuc"      // Kwinsu
	LanguageTypeKUD    LanguageType = "kud"      // 'Auhelawa
	LanguageTypeKUE    LanguageType = "kue"      // Kuman
	LanguageTypeKUF    LanguageType = "kuf"      // Western Katu
	LanguageTypeKUG    LanguageType = "kug"      // Kupa
	LanguageTypeKUH    LanguageType = "kuh"      // Kushi
	LanguageTypeKUI    LanguageType = "kui"      // Kuikúro-Kalapálo
	LanguageTypeKUJ    LanguageType = "kuj"      // Kuria
	LanguageTypeKUK    LanguageType = "kuk"      // Kepo'
	LanguageTypeKUL    LanguageType = "kul"      // Kulere
	LanguageTypeKUM    LanguageType = "kum"      // Kumyk
	LanguageTypeKUN    LanguageType = "kun"      // Kunama
	LanguageTypeKUO    LanguageType = "kuo"      // Kumukio
	LanguageTypeKUP    LanguageType = "kup"      // Kunimaipa
	LanguageTypeKUQ    LanguageType = "kuq"      // Karipuna
	LanguageTypeKUS    LanguageType = "kus"      // Kusaal
	LanguageTypeKUT    LanguageType = "kut"      // Kutenai
	LanguageTypeKUU    LanguageType = "kuu"      // Upper Kuskokwim
	LanguageTypeKUV    LanguageType = "kuv"      // Kur
	LanguageTypeKUW    LanguageType = "kuw"      // Kpagua
	LanguageTypeKUX    LanguageType = "kux"      // Kukatja
	LanguageTypeKUY    LanguageType = "kuy"      // Kuuku-Ya'u
	LanguageTypeKUZ    LanguageType = "kuz"      // Kunza
	LanguageTypeKVA    LanguageType = "kva"      // Bagvalal
	LanguageTypeKVB    LanguageType = "kvb"      // Kubu
	LanguageTypeKVC    LanguageType = "kvc"      // Kove
	LanguageTypeKVD    LanguageType = "kvd"      // Kui (Indonesia)
	LanguageTypeKVE    LanguageType = "kve"      // Kalabakan
	LanguageTypeKVF    LanguageType = "kvf"      // Kabalai
	LanguageTypeKVG    LanguageType = "kvg"      // Kuni-Boazi
	LanguageTypeKVH    LanguageType = "kvh"      // Komodo
	LanguageTypeKVI    LanguageType = "kvi"      // Kwang
	LanguageTypeKVJ    LanguageType = "kvj"      // Psikye
	LanguageTypeKVK    LanguageType = "kvk"      // Korean Sign Language
	LanguageTypeKVL    LanguageType = "kvl"      // Kayaw
	LanguageTypeKVM    LanguageType = "kvm"      // Kendem
	LanguageTypeKVN    LanguageType = "kvn"      // Border Kuna
	LanguageTypeKVO    LanguageType = "kvo"      // Dobel
	LanguageTypeKVP    LanguageType = "kvp"      // Kompane
	LanguageTypeKVQ    LanguageType = "kvq"      // Geba Karen
	LanguageTypeKVR    LanguageType = "kvr"      // Kerinci
	LanguageTypeKVS    LanguageType = "kvs"      // Kunggara
	LanguageTypeKVT    LanguageType = "kvt"      // Lahta Karen and Lahta
	LanguageTypeKVU    LanguageType = "kvu"      // Yinbaw Karen
	LanguageTypeKVV    LanguageType = "kvv"      // Kola
	LanguageTypeKVW    LanguageType = "kvw"      // Wersing
	LanguageTypeKVX    LanguageType = "kvx"      // Parkari Koli
	LanguageTypeKVY    LanguageType = "kvy"      // Yintale Karen and Yintale
	LanguageTypeKVZ    LanguageType = "kvz"      // Tsakwambo and Tsaukambo
	LanguageTypeKWA    LanguageType = "kwa"      // Dâw
	LanguageTypeKWB    LanguageType = "kwb"      // Kwa
	LanguageTypeKWC    LanguageType = "kwc"      // Likwala
	LanguageTypeKWD    LanguageType = "kwd"      // Kwaio
	LanguageTypeKWE    LanguageType = "kwe"      // Kwerba
	LanguageTypeKWF    LanguageType = "kwf"      // Kwara'ae
	LanguageTypeKWG    LanguageType = "kwg"      // Sara Kaba Deme
	LanguageTypeKWH    LanguageType = "kwh"      // Kowiai
	LanguageTypeKWI    LanguageType = "kwi"      // Awa-Cuaiquer
	LanguageTypeKWJ    LanguageType = "kwj"      // Kwanga
	LanguageTypeKWK    LanguageType = "kwk"      // Kwakiutl
	LanguageTypeKWL    LanguageType = "kwl"      // Kofyar
	LanguageTypeKWM    LanguageType = "kwm"      // Kwambi
	LanguageTypeKWN    LanguageType = "kwn"      // Kwangali
	LanguageTypeKWO    LanguageType = "kwo"      // Kwomtari
	LanguageTypeKWP    LanguageType = "kwp"      // Kodia
	LanguageTypeKWQ    LanguageType = "kwq"      // Kwak
	LanguageTypeKWR    LanguageType = "kwr"      // Kwer
	LanguageTypeKWS    LanguageType = "kws"      // Kwese
	LanguageTypeKWT    LanguageType = "kwt"      // Kwesten
	LanguageTypeKWU    LanguageType = "kwu"      // Kwakum
	LanguageTypeKWV    LanguageType = "kwv"      // Sara Kaba Náà
	LanguageTypeKWW    LanguageType = "kww"      // Kwinti
	LanguageTypeKWX    LanguageType = "kwx"      // Khirwar
	LanguageTypeKWY    LanguageType = "kwy"      // San Salvador Kongo
	LanguageTypeKWZ    LanguageType = "kwz"      // Kwadi
	LanguageTypeKXA    LanguageType = "kxa"      // Kairiru
	LanguageTypeKXB    LanguageType = "kxb"      // Krobu
	LanguageTypeKXC    LanguageType = "kxc"      // Konso and Khonso
	LanguageTypeKXD    LanguageType = "kxd"      // Brunei
	LanguageTypeKXE    LanguageType = "kxe"      // Kakihum
	LanguageTypeKXF    LanguageType = "kxf"      // Manumanaw Karen and Manumanaw
	LanguageTypeKXH    LanguageType = "kxh"      // Karo (Ethiopia)
	LanguageTypeKXI    LanguageType = "kxi"      // Keningau Murut
	LanguageTypeKXJ    LanguageType = "kxj"      // Kulfa
	LanguageTypeKXK    LanguageType = "kxk"      // Zayein Karen
	LanguageTypeKXL    LanguageType = "kxl"      // Nepali Kurux
	LanguageTypeKXM    LanguageType = "kxm"      // Northern Khmer
	LanguageTypeKXN    LanguageType = "kxn"      // Kanowit-Tanjong Melanau
	LanguageTypeKXO    LanguageType = "kxo"      // Kanoé
	LanguageTypeKXP    LanguageType = "kxp"      // Wadiyara Koli
	LanguageTypeKXQ    LanguageType = "kxq"      // Smärky Kanum
	LanguageTypeKXR    LanguageType = "kxr"      // Koro (Papua New Guinea)
	LanguageTypeKXS    LanguageType = "kxs"      // Kangjia
	LanguageTypeKXT    LanguageType = "kxt"      // Koiwat
	LanguageTypeKXU    LanguageType = "kxu"      // Kui (India)
	LanguageTypeKXV    LanguageType = "kxv"      // Kuvi
	LanguageTypeKXW    LanguageType = "kxw"      // Konai
	LanguageTypeKXX    LanguageType = "kxx"      // Likuba
	LanguageTypeKXY    LanguageType = "kxy"      // Kayong
	LanguageTypeKXZ    LanguageType = "kxz"      // Kerewo
	LanguageTypeKYA    LanguageType = "kya"      // Kwaya
	LanguageTypeKYB    LanguageType = "kyb"      // Butbut Kalinga
	LanguageTypeKYC    LanguageType = "kyc"      // Kyaka
	LanguageTypeKYD    LanguageType = "kyd"      // Karey
	LanguageTypeKYE    LanguageType = "kye"      // Krache
	LanguageTypeKYF    LanguageType = "kyf"      // Kouya
	LanguageTypeKYG    LanguageType = "kyg"      // Keyagana
	LanguageTypeKYH    LanguageType = "kyh"      // Karok
	LanguageTypeKYI    LanguageType = "kyi"      // Kiput
	LanguageTypeKYJ    LanguageType = "kyj"      // Karao
	LanguageTypeKYK    LanguageType = "kyk"      // Kamayo
	LanguageTypeKYL    LanguageType = "kyl"      // Kalapuya
	LanguageTypeKYM    LanguageType = "kym"      // Kpatili
	LanguageTypeKYN    LanguageType = "kyn"      // Northern Binukidnon
	LanguageTypeKYO    LanguageType = "kyo"      // Kelon
	LanguageTypeKYP    LanguageType = "kyp"      // Kang
	LanguageTypeKYQ    LanguageType = "kyq"      // Kenga
	LanguageTypeKYR    LanguageType = "kyr"      // Kuruáya
	LanguageTypeKYS    LanguageType = "kys"      // Baram Kayan
	LanguageTypeKYT    LanguageType = "kyt"      // Kayagar
	LanguageTypeKYU    LanguageType = "kyu"      // Western Kayah
	LanguageTypeKYV    LanguageType = "kyv"      // Kayort
	LanguageTypeKYW    LanguageType = "kyw"      // Kudmali
	LanguageTypeKYX    LanguageType = "kyx"      // Rapoisi
	LanguageTypeKYY    LanguageType = "kyy"      // Kambaira
	LanguageTypeKYZ    LanguageType = "kyz"      // Kayabí
	LanguageTypeKZA    LanguageType = "kza"      // Western Karaboro
	LanguageTypeKZB    LanguageType = "kzb"      // Kaibobo
	LanguageTypeKZC    LanguageType = "kzc"      // Bondoukou Kulango
	LanguageTypeKZD    LanguageType = "kzd"      // Kadai
	LanguageTypeKZE    LanguageType = "kze"      // Kosena
	LanguageTypeKZF    LanguageType = "kzf"      // Da'a Kaili
	LanguageTypeKZG    LanguageType = "kzg"      // Kikai
	LanguageTypeKZH    LanguageType = "kzh"      // Kenuzi-Dongola
	LanguageTypeKZI    LanguageType = "kzi"      // Kelabit
	LanguageTypeKZJ    LanguageType = "kzj"      // Coastal Kadazan
	LanguageTypeKZK    LanguageType = "kzk"      // Kazukuru
	LanguageTypeKZL    LanguageType = "kzl"      // Kayeli
	LanguageTypeKZM    LanguageType = "kzm"      // Kais
	LanguageTypeKZN    LanguageType = "kzn"      // Kokola
	LanguageTypeKZO    LanguageType = "kzo"      // Kaningi
	LanguageTypeKZP    LanguageType = "kzp"      // Kaidipang
	LanguageTypeKZQ    LanguageType = "kzq"      // Kaike
	LanguageTypeKZR    LanguageType = "kzr"      // Karang
	LanguageTypeKZS    LanguageType = "kzs"      // Sugut Dusun
	LanguageTypeKZT    LanguageType = "kzt"      // Tambunan Dusun
	LanguageTypeKZU    LanguageType = "kzu"      // Kayupulau
	LanguageTypeKZV    LanguageType = "kzv"      // Komyandaret
	LanguageTypeKZW    LanguageType = "kzw"      // Karirí-Xocó
	LanguageTypeKZX    LanguageType = "kzx"      // Kamarian
	LanguageTypeKZY    LanguageType = "kzy"      // Kango (Tshopo District)
	LanguageTypeKZZ    LanguageType = "kzz"      // Kalabra
	LanguageTypeLAA    LanguageType = "laa"      // Southern Subanen
	LanguageTypeLAB    LanguageType = "lab"      // Linear A
	LanguageTypeLAC    LanguageType = "lac"      // Lacandon
	LanguageTypeLAD    LanguageType = "lad"      // Ladino
	LanguageTypeLAE    LanguageType = "lae"      // Pattani
	LanguageTypeLAF    LanguageType = "laf"      // Lafofa
	LanguageTypeLAG    LanguageType = "lag"      // Langi
	LanguageTypeLAH    LanguageType = "lah"      // Lahnda
	LanguageTypeLAI    LanguageType = "lai"      // Lambya
	LanguageTypeLAJ    LanguageType = "laj"      // Lango (Uganda)
	LanguageTypeLAK    LanguageType = "lak"      // Laka (Nigeria)
	LanguageTypeLAL    LanguageType = "lal"      // Lalia
	LanguageTypeLAM    LanguageType = "lam"      // Lamba
	LanguageTypeLAN    LanguageType = "lan"      // Laru
	LanguageTypeLAP    LanguageType = "lap"      // Laka (Chad)
	LanguageTypeLAQ    LanguageType = "laq"      // Qabiao
	LanguageTypeLAR    LanguageType = "lar"      // Larteh
	LanguageTypeLAS    LanguageType = "las"      // Lama (Togo)
	LanguageTypeLAU    LanguageType = "lau"      // Laba
	LanguageTypeLAW    LanguageType = "law"      // Lauje
	LanguageTypeLAX    LanguageType = "lax"      // Tiwa
	LanguageTypeLAY    LanguageType = "lay"      // Lama (Myanmar)
	LanguageTypeLAZ    LanguageType = "laz"      // Aribwatsa
	LanguageTypeLBA    LanguageType = "lba"      // Lui
	LanguageTypeLBB    LanguageType = "lbb"      // Label
	LanguageTypeLBC    LanguageType = "lbc"      // Lakkia
	LanguageTypeLBE    LanguageType = "lbe"      // Lak
	LanguageTypeLBF    LanguageType = "lbf"      // Tinani
	LanguageTypeLBG    LanguageType = "lbg"      // Laopang
	LanguageTypeLBI    LanguageType = "lbi"      // La'bi
	LanguageTypeLBJ    LanguageType = "lbj"      // Ladakhi
	LanguageTypeLBK    LanguageType = "lbk"      // Central Bontok
	LanguageTypeLBL    LanguageType = "lbl"      // Libon Bikol
	LanguageTypeLBM    LanguageType = "lbm"      // Lodhi
	LanguageTypeLBN    LanguageType = "lbn"      // Lamet
	LanguageTypeLBO    LanguageType = "lbo"      // Laven
	LanguageTypeLBQ    LanguageType = "lbq"      // Wampar
	LanguageTypeLBR    LanguageType = "lbr"      // Lohorung
	LanguageTypeLBS    LanguageType = "lbs"      // Libyan Sign Language
	LanguageTypeLBT    LanguageType = "lbt"      // Lachi
	LanguageTypeLBU    LanguageType = "lbu"      // Labu
	LanguageTypeLBV    LanguageType = "lbv"      // Lavatbura-Lamusong
	LanguageTypeLBW    LanguageType = "lbw"      // Tolaki
	LanguageTypeLBX    LanguageType = "lbx"      // Lawangan
	LanguageTypeLBY    LanguageType = "lby"      // Lamu-Lamu
	LanguageTypeLBZ    LanguageType = "lbz"      // Lardil
	LanguageTypeLCC    LanguageType = "lcc"      // Legenyem
	LanguageTypeLCD    LanguageType = "lcd"      // Lola
	LanguageTypeLCE    LanguageType = "lce"      // Loncong
	LanguageTypeLCF    LanguageType = "lcf"      // Lubu
	LanguageTypeLCH    LanguageType = "lch"      // Luchazi
	LanguageTypeLCL    LanguageType = "lcl"      // Lisela
	LanguageTypeLCM    LanguageType = "lcm"      // Tungag
	LanguageTypeLCP    LanguageType = "lcp"      // Western Lawa
	LanguageTypeLCQ    LanguageType = "lcq"      // Luhu
	LanguageTypeLCS    LanguageType = "lcs"      // Lisabata-Nuniali
	LanguageTypeLDA    LanguageType = "lda"      // Kla-Dan
	LanguageTypeLDB    LanguageType = "ldb"      // Dũya
	LanguageTypeLDD    LanguageType = "ldd"      // Luri
	LanguageTypeLDG    LanguageType = "ldg"      // Lenyima
	LanguageTypeLDH    LanguageType = "ldh"      // Lamja-Dengsa-Tola
	LanguageTypeLDI    LanguageType = "ldi"      // Laari
	LanguageTypeLDJ    LanguageType = "ldj"      // Lemoro
	LanguageTypeLDK    LanguageType = "ldk"      // Leelau
	LanguageTypeLDL    LanguageType = "ldl"      // Kaan
	LanguageTypeLDM    LanguageType = "ldm"      // Landoma
	LanguageTypeLDN    LanguageType = "ldn"      // Láadan
	LanguageTypeLDO    LanguageType = "ldo"      // Loo
	LanguageTypeLDP    LanguageType = "ldp"      // Tso
	LanguageTypeLDQ    LanguageType = "ldq"      // Lufu
	LanguageTypeLEA    LanguageType = "lea"      // Lega-Shabunda
	LanguageTypeLEB    LanguageType = "leb"      // Lala-Bisa
	LanguageTypeLEC    LanguageType = "lec"      // Leco
	LanguageTypeLED    LanguageType = "led"      // Lendu
	LanguageTypeLEE    LanguageType = "lee"      // Lyélé
	LanguageTypeLEF    LanguageType = "lef"      // Lelemi
	LanguageTypeLEG    LanguageType = "leg"      // Lengua
	LanguageTypeLEH    LanguageType = "leh"      // Lenje
	LanguageTypeLEI    LanguageType = "lei"      // Lemio
	LanguageTypeLEJ    LanguageType = "lej"      // Lengola
	LanguageTypeLEK    LanguageType = "lek"      // Leipon
	LanguageTypeLEL    LanguageType = "lel"      // Lele (Democratic Republic of Congo)
	LanguageTypeLEM    LanguageType = "lem"      // Nomaande
	LanguageTypeLEN    LanguageType = "len"      // Lenca
	LanguageTypeLEO    LanguageType = "leo"      // Leti (Cameroon)
	LanguageTypeLEP    LanguageType = "lep"      // Lepcha
	LanguageTypeLEQ    LanguageType = "leq"      // Lembena
	LanguageTypeLER    LanguageType = "ler"      // Lenkau
	LanguageTypeLES    LanguageType = "les"      // Lese
	LanguageTypeLET    LanguageType = "let"      // Lesing-Gelimi and Amio-Gelimi
	LanguageTypeLEU    LanguageType = "leu"      // Kara (Papua New Guinea)
	LanguageTypeLEV    LanguageType = "lev"      // Lamma
	LanguageTypeLEW    LanguageType = "lew"      // Ledo Kaili
	LanguageTypeLEX    LanguageType = "lex"      // Luang
	LanguageTypeLEY    LanguageType = "ley"      // Lemolang
	LanguageTypeLEZ    LanguageType = "lez"      // Lezghian
	LanguageTypeLFA    LanguageType = "lfa"      // Lefa
	LanguageTypeLFN    LanguageType = "lfn"      // Lingua Franca Nova
	LanguageTypeLGA    LanguageType = "lga"      // Lungga
	LanguageTypeLGB    LanguageType = "lgb"      // Laghu
	LanguageTypeLGG    LanguageType = "lgg"      // Lugbara
	LanguageTypeLGH    LanguageType = "lgh"      // Laghuu
	LanguageTypeLGI    LanguageType = "lgi"      // Lengilu
	LanguageTypeLGK    LanguageType = "lgk"      // Lingarak and Neverver
	LanguageTypeLGL    LanguageType = "lgl"      // Wala
	LanguageTypeLGM    LanguageType = "lgm"      // Lega-Mwenga
	LanguageTypeLGN    LanguageType = "lgn"      // Opuuo
	LanguageTypeLGQ    LanguageType = "lgq"      // Logba
	LanguageTypeLGR    LanguageType = "lgr"      // Lengo
	LanguageTypeLGT    LanguageType = "lgt"      // Pahi
	LanguageTypeLGU    LanguageType = "lgu"      // Longgu
	LanguageTypeLGZ    LanguageType = "lgz"      // Ligenza
	LanguageTypeLHA    LanguageType = "lha"      // Laha (Viet Nam)
	LanguageTypeLHH    LanguageType = "lhh"      // Laha (Indonesia)
	LanguageTypeLHI    LanguageType = "lhi"      // Lahu Shi
	LanguageTypeLHL    LanguageType = "lhl"      // Lahul Lohar
	LanguageTypeLHM    LanguageType = "lhm"      // Lhomi
	LanguageTypeLHN    LanguageType = "lhn"      // Lahanan
	LanguageTypeLHP    LanguageType = "lhp"      // Lhokpu
	LanguageTypeLHS    LanguageType = "lhs"      // Mlahsö
	LanguageTypeLHT    LanguageType = "lht"      // Lo-Toga
	LanguageTypeLHU    LanguageType = "lhu"      // Lahu
	LanguageTypeLIA    LanguageType = "lia"      // West-Central Limba
	LanguageTypeLIB    LanguageType = "lib"      // Likum
	LanguageTypeLIC    LanguageType = "lic"      // Hlai
	LanguageTypeLID    LanguageType = "lid"      // Nyindrou
	LanguageTypeLIE    LanguageType = "lie"      // Likila
	LanguageTypeLIF    LanguageType = "lif"      // Limbu
	LanguageTypeLIG    LanguageType = "lig"      // Ligbi
	LanguageTypeLIH    LanguageType = "lih"      // Lihir
	LanguageTypeLII    LanguageType = "lii"      // Lingkhim
	LanguageTypeLIJ    LanguageType = "lij"      // Ligurian
	LanguageTypeLIK    LanguageType = "lik"      // Lika
	LanguageTypeLIL    LanguageType = "lil"      // Lillooet
	LanguageTypeLIO    LanguageType = "lio"      // Liki
	LanguageTypeLIP    LanguageType = "lip"      // Sekpele
	LanguageTypeLIQ    LanguageType = "liq"      // Libido
	LanguageTypeLIR    LanguageType = "lir"      // Liberian English
	LanguageTypeLIS    LanguageType = "lis"      // Lisu
	LanguageTypeLIU    LanguageType = "liu"      // Logorik
	LanguageTypeLIV    LanguageType = "liv"      // Liv
	LanguageTypeLIW    LanguageType = "liw"      // Col
	LanguageTypeLIX    LanguageType = "lix"      // Liabuku
	LanguageTypeLIY    LanguageType = "liy"      // Banda-Bambari
	LanguageTypeLIZ    LanguageType = "liz"      // Libinza
	LanguageTypeLJA    LanguageType = "lja"      // Golpa
	LanguageTypeLJE    LanguageType = "lje"      // Rampi
	LanguageTypeLJI    LanguageType = "lji"      // Laiyolo
	LanguageTypeLJL    LanguageType = "ljl"      // Li'o
	LanguageTypeLJP    LanguageType = "ljp"      // Lampung Api
	LanguageTypeLJW    LanguageType = "ljw"      // Yirandali
	LanguageTypeLJX    LanguageType = "ljx"      // Yuru
	LanguageTypeLKA    LanguageType = "lka"      // Lakalei
	LanguageTypeLKB    LanguageType = "lkb"      // Kabras and Lukabaras
	LanguageTypeLKC    LanguageType = "lkc"      // Kucong
	LanguageTypeLKD    LanguageType = "lkd"      // Lakondê
	LanguageTypeLKE    LanguageType = "lke"      // Kenyi
	LanguageTypeLKH    LanguageType = "lkh"      // Lakha
	LanguageTypeLKI    LanguageType = "lki"      // Laki
	LanguageTypeLKJ    LanguageType = "lkj"      // Remun
	LanguageTypeLKL    LanguageType = "lkl"      // Laeko-Libuat
	LanguageTypeLKM    LanguageType = "lkm"      // Kalaamaya
	LanguageTypeLKN    LanguageType = "lkn"      // Lakon and Vure
	LanguageTypeLKO    LanguageType = "lko"      // Khayo and Olukhayo
	LanguageTypeLKR    LanguageType = "lkr"      // Päri
	LanguageTypeLKS    LanguageType = "lks"      // Kisa and Olushisa
	LanguageTypeLKT    LanguageType = "lkt"      // Lakota
	LanguageTypeLKU    LanguageType = "lku"      // Kungkari
	LanguageTypeLKY    LanguageType = "lky"      // Lokoya
	LanguageTypeLLA    LanguageType = "lla"      // Lala-Roba
	LanguageTypeLLB    LanguageType = "llb"      // Lolo
	LanguageTypeLLC    LanguageType = "llc"      // Lele (Guinea)
	LanguageTypeLLD    LanguageType = "lld"      // Ladin
	LanguageTypeLLE    LanguageType = "lle"      // Lele (Papua New Guinea)
	LanguageTypeLLF    LanguageType = "llf"      // Hermit
	LanguageTypeLLG    LanguageType = "llg"      // Lole
	LanguageTypeLLH    LanguageType = "llh"      // Lamu
	LanguageTypeLLI    LanguageType = "lli"      // Teke-Laali
	LanguageTypeLLJ    LanguageType = "llj"      // Ladji Ladji
	LanguageTypeLLK    LanguageType = "llk"      // Lelak
	LanguageTypeLLL    LanguageType = "lll"      // Lilau
	LanguageTypeLLM    LanguageType = "llm"      // Lasalimu
	LanguageTypeLLN    LanguageType = "lln"      // Lele (Chad)
	LanguageTypeLLO    LanguageType = "llo"      // Khlor
	LanguageTypeLLP    LanguageType = "llp"      // North Efate
	LanguageTypeLLQ    LanguageType = "llq"      // Lolak
	LanguageTypeLLS    LanguageType = "lls"      // Lithuanian Sign Language
	LanguageTypeLLU    LanguageType = "llu"      // Lau
	LanguageTypeLLX    LanguageType = "llx"      // Lauan
	LanguageTypeLMA    LanguageType = "lma"      // East Limba
	LanguageTypeLMB    LanguageType = "lmb"      // Merei
	LanguageTypeLMC    LanguageType = "lmc"      // Limilngan
	LanguageTypeLMD    LanguageType = "lmd"      // Lumun
	LanguageTypeLME    LanguageType = "lme"      // Pévé
	LanguageTypeLMF    LanguageType = "lmf"      // South Lembata
	LanguageTypeLMG    LanguageType = "lmg"      // Lamogai
	LanguageTypeLMH    LanguageType = "lmh"      // Lambichhong
	LanguageTypeLMI    LanguageType = "lmi"      // Lombi
	LanguageTypeLMJ    LanguageType = "lmj"      // West Lembata
	LanguageTypeLMK    LanguageType = "lmk"      // Lamkang
	LanguageTypeLML    LanguageType = "lml"      // Hano
	LanguageTypeLMM    LanguageType = "lmm"      // Lamam
	LanguageTypeLMN    LanguageType = "lmn"      // Lambadi
	LanguageTypeLMO    LanguageType = "lmo"      // Lombard
	LanguageTypeLMP    LanguageType = "lmp"      // Limbum
	LanguageTypeLMQ    LanguageType = "lmq"      // Lamatuka
	LanguageTypeLMR    LanguageType = "lmr"      // Lamalera
	LanguageTypeLMU    LanguageType = "lmu"      // Lamenu
	LanguageTypeLMV    LanguageType = "lmv"      // Lomaiviti
	LanguageTypeLMW    LanguageType = "lmw"      // Lake Miwok
	LanguageTypeLMX    LanguageType = "lmx"      // Laimbue
	LanguageTypeLMY    LanguageType = "lmy"      // Lamboya
	LanguageTypeLMZ    LanguageType = "lmz"      // Lumbee
	LanguageTypeLNA    LanguageType = "lna"      // Langbashe
	LanguageTypeLNB    LanguageType = "lnb"      // Mbalanhu
	LanguageTypeLND    LanguageType = "lnd"      // Lundayeh and Lun Bawang
	LanguageTypeLNG    LanguageType = "lng"      // Langobardic
	LanguageTypeLNH    LanguageType = "lnh"      // Lanoh
	LanguageTypeLNI    LanguageType = "lni"      // Daantanai'
	LanguageTypeLNJ    LanguageType = "lnj"      // Leningitij
	LanguageTypeLNL    LanguageType = "lnl"      // South Central Banda
	LanguageTypeLNM    LanguageType = "lnm"      // Langam
	LanguageTypeLNN    LanguageType = "lnn"      // Lorediakarkar
	LanguageTypeLNO    LanguageType = "lno"      // Lango (Sudan)
	LanguageTypeLNS    LanguageType = "lns"      // Lamnso'
	LanguageTypeLNU    LanguageType = "lnu"      // Longuda
	LanguageTypeLNW    LanguageType = "lnw"      // Lanima
	LanguageTypeLNZ    LanguageType = "lnz"      // Lonzo
	LanguageTypeLOA    LanguageType = "loa"      // Loloda
	LanguageTypeLOB    LanguageType = "lob"      // Lobi
	LanguageTypeLOC    LanguageType = "loc"      // Inonhan
	LanguageTypeLOE    LanguageType = "loe"      // Saluan
	LanguageTypeLOF    LanguageType = "lof"      // Logol
	LanguageTypeLOG    LanguageType = "log"      // Logo
	LanguageTypeLOH    LanguageType = "loh"      // Narim
	LanguageTypeLOI    LanguageType = "loi"      // Loma (Côte d'Ivoire)
	LanguageTypeLOJ    LanguageType = "loj"      // Lou
	LanguageTypeLOK    LanguageType = "lok"      // Loko
	LanguageTypeLOL    LanguageType = "lol"      // Mongo
	LanguageTypeLOM    LanguageType = "lom"      // Loma (Liberia)
	LanguageTypeLON    LanguageType = "lon"      // Malawi Lomwe
	LanguageTypeLOO    LanguageType = "loo"      // Lombo
	LanguageTypeLOP    LanguageType = "lop"      // Lopa
	LanguageTypeLOQ    LanguageType = "loq"      // Lobala
	LanguageTypeLOR    LanguageType = "lor"      // Téén
	LanguageTypeLOS    LanguageType = "los"      // Loniu
	LanguageTypeLOT    LanguageType = "lot"      // Otuho
	LanguageTypeLOU    LanguageType = "lou"      // Louisiana Creole French
	LanguageTypeLOV    LanguageType = "lov"      // Lopi
	LanguageTypeLOW    LanguageType = "low"      // Tampias Lobu
	LanguageTypeLOX    LanguageType = "lox"      // Loun
	LanguageTypeLOY    LanguageType = "loy"      // Loke
	LanguageTypeLOZ    LanguageType = "loz"      // Lozi
	LanguageTypeLPA    LanguageType = "lpa"      // Lelepa
	LanguageTypeLPE    LanguageType = "lpe"      // Lepki
	LanguageTypeLPN    LanguageType = "lpn"      // Long Phuri Naga
	LanguageTypeLPO    LanguageType = "lpo"      // Lipo
	LanguageTypeLPX    LanguageType = "lpx"      // Lopit
	LanguageTypeLRA    LanguageType = "lra"      // Rara Bakati'
	LanguageTypeLRC    LanguageType = "lrc"      // Northern Luri
	LanguageTypeLRE    LanguageType = "lre"      // Laurentian
	LanguageTypeLRG    LanguageType = "lrg"      // Laragia
	LanguageTypeLRI    LanguageType = "lri"      // Marachi and Olumarachi
	LanguageTypeLRK    LanguageType = "lrk"      // Loarki
	LanguageTypeLRL    LanguageType = "lrl"      // Lari
	LanguageTypeLRM    LanguageType = "lrm"      // Marama and Olumarama
	LanguageTypeLRN    LanguageType = "lrn"      // Lorang
	LanguageTypeLRO    LanguageType = "lro"      // Laro
	LanguageTypeLRR    LanguageType = "lrr"      // Southern Yamphu
	LanguageTypeLRT    LanguageType = "lrt"      // Larantuka Malay
	LanguageTypeLRV    LanguageType = "lrv"      // Larevat
	LanguageTypeLRZ    LanguageType = "lrz"      // Lemerig
	LanguageTypeLSA    LanguageType = "lsa"      // Lasgerdi
	LanguageTypeLSD    LanguageType = "lsd"      // Lishana Deni
	LanguageTypeLSE    LanguageType = "lse"      // Lusengo
	LanguageTypeLSG    LanguageType = "lsg"      // Lyons Sign Language
	LanguageTypeLSH    LanguageType = "lsh"      // Lish
	LanguageTypeLSI    LanguageType = "lsi"      // Lashi
	LanguageTypeLSL    LanguageType = "lsl"      // Latvian Sign Language
	LanguageTypeLSM    LanguageType = "lsm"      // Saamia and Olusamia
	LanguageTypeLSO    LanguageType = "lso"      // Laos Sign Language
	LanguageTypeLSP    LanguageType = "lsp"      // Panamanian Sign Language and Lengua de Señas Panameñas
	LanguageTypeLSR    LanguageType = "lsr"      // Aruop
	LanguageTypeLSS    LanguageType = "lss"      // Lasi
	LanguageTypeLST    LanguageType = "lst"      // Trinidad and Tobago Sign Language
	LanguageTypeLSY    LanguageType = "lsy"      // Mauritian Sign Language
	LanguageTypeLTC    LanguageType = "ltc"      // Late Middle Chinese
	LanguageTypeLTG    LanguageType = "ltg"      // Latgalian
	LanguageTypeLTI    LanguageType = "lti"      // Leti (Indonesia)
	LanguageTypeLTN    LanguageType = "ltn"      // Latundê
	LanguageTypeLTO    LanguageType = "lto"      // Tsotso and Olutsotso
	LanguageTypeLTS    LanguageType = "lts"      // Tachoni and Lutachoni
	LanguageTypeLTU    LanguageType = "ltu"      // Latu
	LanguageTypeLUA    LanguageType = "lua"      // Luba-Lulua
	LanguageTypeLUC    LanguageType = "luc"      // Aringa
	LanguageTypeLUD    LanguageType = "lud"      // Ludian
	LanguageTypeLUE    LanguageType = "lue"      // Luvale
	LanguageTypeLUF    LanguageType = "luf"      // Laua
	LanguageTypeLUI    LanguageType = "lui"      // Luiseno
	LanguageTypeLUJ    LanguageType = "luj"      // Luna
	LanguageTypeLUK    LanguageType = "luk"      // Lunanakha
	LanguageTypeLUL    LanguageType = "lul"      // Olu'bo
	LanguageTypeLUM    LanguageType = "lum"      // Luimbi
	LanguageTypeLUN    LanguageType = "lun"      // Lunda
	LanguageTypeLUO    LanguageType = "luo"      // Luo (Kenya and Tanzania) and Dholuo
	LanguageTypeLUP    LanguageType = "lup"      // Lumbu
	LanguageTypeLUQ    LanguageType = "luq"      // Lucumi
	LanguageTypeLUR    LanguageType = "lur"      // Laura
	LanguageTypeLUS    LanguageType = "lus"      // Lushai
	LanguageTypeLUT    LanguageType = "lut"      // Lushootseed
	LanguageTypeLUU    LanguageType = "luu"      // Lumba-Yakkha
	LanguageTypeLUV    LanguageType = "luv"      // Luwati
	LanguageTypeLUW    LanguageType = "luw"      // Luo (Cameroon)
	LanguageTypeLUY    LanguageType = "luy"      // Luyia and Oluluyia
	LanguageTypeLUZ    LanguageType = "luz"      // Southern Luri
	LanguageTypeLVA    LanguageType = "lva"      // Maku'a
	LanguageTypeLVK    LanguageType = "lvk"      // Lavukaleve
	LanguageTypeLVS    LanguageType = "lvs"      // Standard Latvian
	LanguageTypeLVU    LanguageType = "lvu"      // Levuka
	LanguageTypeLWA    LanguageType = "lwa"      // Lwalu
	LanguageTypeLWE    LanguageType = "lwe"      // Lewo Eleng
	LanguageTypeLWG    LanguageType = "lwg"      // Wanga and Oluwanga
	LanguageTypeLWH    LanguageType = "lwh"      // White Lachi
	LanguageTypeLWL    LanguageType = "lwl"      // Eastern Lawa
	LanguageTypeLWM    LanguageType = "lwm"      // Laomian
	LanguageTypeLWO    LanguageType = "lwo"      // Luwo
	LanguageTypeLWT    LanguageType = "lwt"      // Lewotobi
	LanguageTypeLWU    LanguageType = "lwu"      // Lawu
	LanguageTypeLWW    LanguageType = "lww"      // Lewo
	LanguageTypeLYA    LanguageType = "lya"      // Layakha
	LanguageTypeLYG    LanguageType = "lyg"      // Lyngngam
	LanguageTypeLYN    LanguageType = "lyn"      // Luyana
	LanguageTypeLZH    LanguageType = "lzh"      // Literary Chinese
	LanguageTypeLZL    LanguageType = "lzl"      // Litzlitz
	LanguageTypeLZN    LanguageType = "lzn"      // Leinong Naga
	LanguageTypeLZZ    LanguageType = "lzz"      // Laz
	LanguageTypeMAA    LanguageType = "maa"      // San Jerónimo Tecóatl Mazatec
	LanguageTypeMAB    LanguageType = "mab"      // Yutanduchi Mixtec
	LanguageTypeMAD    LanguageType = "mad"      // Madurese
	LanguageTypeMAE    LanguageType = "mae"      // Bo-Rukul
	LanguageTypeMAF    LanguageType = "maf"      // Mafa
	LanguageTypeMAG    LanguageType = "mag"      // Magahi
	LanguageTypeMAI    LanguageType = "mai"      // Maithili
	LanguageTypeMAJ    LanguageType = "maj"      // Jalapa De Díaz Mazatec
	LanguageTypeMAK    LanguageType = "mak"      // Makasar
	LanguageTypeMAM    LanguageType = "mam"      // Mam
	LanguageTypeMAN    LanguageType = "man"      // Mandingo and Manding
	LanguageTypeMAP    LanguageType = "map"      // Austronesian languages
	LanguageTypeMAQ    LanguageType = "maq"      // Chiquihuitlán Mazatec
	LanguageTypeMAS    LanguageType = "mas"      // Masai
	LanguageTypeMAT    LanguageType = "mat"      // San Francisco Matlatzinca
	LanguageTypeMAU    LanguageType = "mau"      // Huautla Mazatec
	LanguageTypeMAV    LanguageType = "mav"      // Sateré-Mawé
	LanguageTypeMAW    LanguageType = "maw"      // Mampruli
	LanguageTypeMAX    LanguageType = "max"      // North Moluccan Malay
	LanguageTypeMAZ    LanguageType = "maz"      // Central Mazahua
	LanguageTypeMBA    LanguageType = "mba"      // Higaonon
	LanguageTypeMBB    LanguageType = "mbb"      // Western Bukidnon Manobo
	LanguageTypeMBC    LanguageType = "mbc"      // Macushi
	LanguageTypeMBD    LanguageType = "mbd"      // Dibabawon Manobo
	LanguageTypeMBE    LanguageType = "mbe"      // Molale
	LanguageTypeMBF    LanguageType = "mbf"      // Baba Malay
	LanguageTypeMBH    LanguageType = "mbh"      // Mangseng
	LanguageTypeMBI    LanguageType = "mbi"      // Ilianen Manobo
	LanguageTypeMBJ    LanguageType = "mbj"      // Nadëb
	LanguageTypeMBK    LanguageType = "mbk"      // Malol
	LanguageTypeMBL    LanguageType = "mbl"      // Maxakalí
	LanguageTypeMBM    LanguageType = "mbm"      // Ombamba
	LanguageTypeMBN    LanguageType = "mbn"      // Macaguán
	LanguageTypeMBO    LanguageType = "mbo"      // Mbo (Cameroon)
	LanguageTypeMBP    LanguageType = "mbp"      // Malayo
	LanguageTypeMBQ    LanguageType = "mbq"      // Maisin
	LanguageTypeMBR    LanguageType = "mbr"      // Nukak Makú
	LanguageTypeMBS    LanguageType = "mbs"      // Sarangani Manobo
	LanguageTypeMBT    LanguageType = "mbt"      // Matigsalug Manobo
	LanguageTypeMBU    LanguageType = "mbu"      // Mbula-Bwazza
	LanguageTypeMBV    LanguageType = "mbv"      // Mbulungish
	LanguageTypeMBW    LanguageType = "mbw"      // Maring
	LanguageTypeMBX    LanguageType = "mbx"      // Mari (East Sepik Province)
	LanguageTypeMBY    LanguageType = "mby"      // Memoni
	LanguageTypeMBZ    LanguageType = "mbz"      // Amoltepec Mixtec
	LanguageTypeMCA    LanguageType = "mca"      // Maca
	LanguageTypeMCB    LanguageType = "mcb"      // Machiguenga
	LanguageTypeMCC    LanguageType = "mcc"      // Bitur
	LanguageTypeMCD    LanguageType = "mcd"      // Sharanahua
	LanguageTypeMCE    LanguageType = "mce"      // Itundujia Mixtec
	LanguageTypeMCF    LanguageType = "mcf"      // Matsés
	LanguageTypeMCG    LanguageType = "mcg"      // Mapoyo
	LanguageTypeMCH    LanguageType = "mch"      // Maquiritari
	LanguageTypeMCI    LanguageType = "mci"      // Mese
	LanguageTypeMCJ    LanguageType = "mcj"      // Mvanip
	LanguageTypeMCK    LanguageType = "mck"      // Mbunda
	LanguageTypeMCL    LanguageType = "mcl"      // Macaguaje
	LanguageTypeMCM    LanguageType = "mcm"      // Malaccan Creole Portuguese
	LanguageTypeMCN    LanguageType = "mcn"      // Masana
	LanguageTypeMCO    LanguageType = "mco"      // Coatlán Mixe
	LanguageTypeMCP    LanguageType = "mcp"      // Makaa
	LanguageTypeMCQ    LanguageType = "mcq"      // Ese
	LanguageTypeMCR    LanguageType = "mcr"      // Menya
	LanguageTypeMCS    LanguageType = "mcs"      // Mambai
	LanguageTypeMCT    LanguageType = "mct"      // Mengisa
	LanguageTypeMCU    LanguageType = "mcu"      // Cameroon Mambila
	LanguageTypeMCV    LanguageType = "mcv"      // Minanibai
	LanguageTypeMCW    LanguageType = "mcw"      // Mawa (Chad)
	LanguageTypeMCX    LanguageType = "mcx"      // Mpiemo
	LanguageTypeMCY    LanguageType = "mcy"      // South Watut
	LanguageTypeMCZ    LanguageType = "mcz"      // Mawan
	LanguageTypeMDA    LanguageType = "mda"      // Mada (Nigeria)
	LanguageTypeMDB    LanguageType = "mdb"      // Morigi
	LanguageTypeMDC    LanguageType = "mdc"      // Male (Papua New Guinea)
	LanguageTypeMDD    LanguageType = "mdd"      // Mbum
	LanguageTypeMDE    LanguageType = "mde"      // Maba (Chad)
	LanguageTypeMDF    LanguageType = "mdf"      // Moksha
	LanguageTypeMDG    LanguageType = "mdg"      // Massalat
	LanguageTypeMDH    LanguageType = "mdh"      // Maguindanaon
	LanguageTypeMDI    LanguageType = "mdi"      // Mamvu
	LanguageTypeMDJ    LanguageType = "mdj"      // Mangbetu
	LanguageTypeMDK    LanguageType = "mdk"      // Mangbutu
	LanguageTypeMDL    LanguageType = "mdl"      // Maltese Sign Language
	LanguageTypeMDM    LanguageType = "mdm"      // Mayogo
	LanguageTypeMDN    LanguageType = "mdn"      // Mbati
	LanguageTypeMDP    LanguageType = "mdp"      // Mbala
	LanguageTypeMDQ    LanguageType = "mdq"      // Mbole
	LanguageTypeMDR    LanguageType = "mdr"      // Mandar
	LanguageTypeMDS    LanguageType = "mds"      // Maria (Papua New Guinea)
	LanguageTypeMDT    LanguageType = "mdt"      // Mbere
	LanguageTypeMDU    LanguageType = "mdu"      // Mboko
	LanguageTypeMDV    LanguageType = "mdv"      // Santa Lucía Monteverde Mixtec
	LanguageTypeMDW    LanguageType = "mdw"      // Mbosi
	LanguageTypeMDX    LanguageType = "mdx"      // Dizin
	LanguageTypeMDY    LanguageType = "mdy"      // Male (Ethiopia)
	LanguageTypeMDZ    LanguageType = "mdz"      // Suruí Do Pará
	LanguageTypeMEA    LanguageType = "mea"      // Menka
	LanguageTypeMEB    LanguageType = "meb"      // Ikobi
	LanguageTypeMEC    LanguageType = "mec"      // Mara
	LanguageTypeMED    LanguageType = "med"      // Melpa
	LanguageTypeMEE    LanguageType = "mee"      // Mengen
	LanguageTypeMEF    LanguageType = "mef"      // Megam
	LanguageTypeMEG    LanguageType = "meg"      // Mea
	LanguageTypeMEH    LanguageType = "meh"      // Southwestern Tlaxiaco Mixtec
	LanguageTypeMEI    LanguageType = "mei"      // Midob
	LanguageTypeMEJ    LanguageType = "mej"      // Meyah
	LanguageTypeMEK    LanguageType = "mek"      // Mekeo
	LanguageTypeMEL    LanguageType = "mel"      // Central Melanau
	LanguageTypeMEM    LanguageType = "mem"      // Mangala
	LanguageTypeMEN    LanguageType = "men"      // Mende (Sierra Leone)
	LanguageTypeMEO    LanguageType = "meo"      // Kedah Malay
	LanguageTypeMEP    LanguageType = "mep"      // Miriwung
	LanguageTypeMEQ    LanguageType = "meq"      // Merey
	LanguageTypeMER    LanguageType = "mer"      // Meru
	LanguageTypeMES    LanguageType = "mes"      // Masmaje
	LanguageTypeMET    LanguageType = "met"      // Mato
	LanguageTypeMEU    LanguageType = "meu"      // Motu
	LanguageTypeMEV    LanguageType = "mev"      // Mano
	LanguageTypeMEW    LanguageType = "mew"      // Maaka
	LanguageTypeMEY    LanguageType = "mey"      // Hassaniyya
	LanguageTypeMEZ    LanguageType = "mez"      // Menominee
	LanguageTypeMFA    LanguageType = "mfa"      // Pattani Malay
	LanguageTypeMFB    LanguageType = "mfb"      // Bangka
	LanguageTypeMFC    LanguageType = "mfc"      // Mba
	LanguageTypeMFD    LanguageType = "mfd"      // Mendankwe-Nkwen
	LanguageTypeMFE    LanguageType = "mfe"      // Morisyen
	LanguageTypeMFF    LanguageType = "mff"      // Naki
	LanguageTypeMFG    LanguageType = "mfg"      // Mogofin
	LanguageTypeMFH    LanguageType = "mfh"      // Matal
	LanguageTypeMFI    LanguageType = "mfi"      // Wandala
	LanguageTypeMFJ    LanguageType = "mfj"      // Mefele
	LanguageTypeMFK    LanguageType = "mfk"      // North Mofu
	LanguageTypeMFL    LanguageType = "mfl"      // Putai
	LanguageTypeMFM    LanguageType = "mfm"      // Marghi South
	LanguageTypeMFN    LanguageType = "mfn"      // Cross River Mbembe
	LanguageTypeMFO    LanguageType = "mfo"      // Mbe
	LanguageTypeMFP    LanguageType = "mfp"      // Makassar Malay
	LanguageTypeMFQ    LanguageType = "mfq"      // Moba
	LanguageTypeMFR    LanguageType = "mfr"      // Marithiel
	LanguageTypeMFS    LanguageType = "mfs"      // Mexican Sign Language
	LanguageTypeMFT    LanguageType = "mft"      // Mokerang
	LanguageTypeMFU    LanguageType = "mfu"      // Mbwela
	LanguageTypeMFV    LanguageType = "mfv"      // Mandjak
	LanguageTypeMFW    LanguageType = "mfw"      // Mulaha
	LanguageTypeMFX    LanguageType = "mfx"      // Melo
	LanguageTypeMFY    LanguageType = "mfy"      // Mayo
	LanguageTypeMFZ    LanguageType = "mfz"      // Mabaan
	LanguageTypeMGA    LanguageType = "mga"      // Middle Irish (900-1200)
	LanguageTypeMGB    LanguageType = "mgb"      // Mararit
	LanguageTypeMGC    LanguageType = "mgc"      // Morokodo
	LanguageTypeMGD    LanguageType = "mgd"      // Moru
	LanguageTypeMGE    LanguageType = "mge"      // Mango
	LanguageTypeMGF    LanguageType = "mgf"      // Maklew
	LanguageTypeMGG    LanguageType = "mgg"      // Mpumpong
	LanguageTypeMGH    LanguageType = "mgh"      // Makhuwa-Meetto
	LanguageTypeMGI    LanguageType = "mgi"      // Lijili
	LanguageTypeMGJ    LanguageType = "mgj"      // Abureni
	LanguageTypeMGK    LanguageType = "mgk"      // Mawes
	LanguageTypeMGL    LanguageType = "mgl"      // Maleu-Kilenge
	LanguageTypeMGM    LanguageType = "mgm"      // Mambae
	LanguageTypeMGN    LanguageType = "mgn"      // Mbangi
	LanguageTypeMGO    LanguageType = "mgo"      // Meta'
	LanguageTypeMGP    LanguageType = "mgp"      // Eastern Magar
	LanguageTypeMGQ    LanguageType = "mgq"      // Malila
	LanguageTypeMGR    LanguageType = "mgr"      // Mambwe-Lungu
	LanguageTypeMGS    LanguageType = "mgs"      // Manda (Tanzania)
	LanguageTypeMGT    LanguageType = "mgt"      // Mongol
	LanguageTypeMGU    LanguageType = "mgu"      // Mailu
	LanguageTypeMGV    LanguageType = "mgv"      // Matengo
	LanguageTypeMGW    LanguageType = "mgw"      // Matumbi
	LanguageTypeMGX    LanguageType = "mgx"      // Omati
	LanguageTypeMGY    LanguageType = "mgy"      // Mbunga
	LanguageTypeMGZ    LanguageType = "mgz"      // Mbugwe
	LanguageTypeMHA    LanguageType = "mha"      // Manda (India)
	LanguageTypeMHB    LanguageType = "mhb"      // Mahongwe
	LanguageTypeMHC    LanguageType = "mhc"      // Mocho
	LanguageTypeMHD    LanguageType = "mhd"      // Mbugu
	LanguageTypeMHE    LanguageType = "mhe"      // Besisi and Mah Meri
	LanguageTypeMHF    LanguageType = "mhf"      // Mamaa
	LanguageTypeMHG    LanguageType = "mhg"      // Margu
	LanguageTypeMHH    LanguageType = "mhh"      // Maskoy Pidgin
	LanguageTypeMHI    LanguageType = "mhi"      // Ma'di
	LanguageTypeMHJ    LanguageType = "mhj"      // Mogholi
	LanguageTypeMHK    LanguageType = "mhk"      // Mungaka
	LanguageTypeMHL    LanguageType = "mhl"      // Mauwake
	LanguageTypeMHM    LanguageType = "mhm"      // Makhuwa-Moniga
	LanguageTypeMHN    LanguageType = "mhn"      // Mócheno
	LanguageTypeMHO    LanguageType = "mho"      // Mashi (Zambia)
	LanguageTypeMHP    LanguageType = "mhp"      // Balinese Malay
	LanguageTypeMHQ    LanguageType = "mhq"      // Mandan
	LanguageTypeMHR    LanguageType = "mhr"      // Eastern Mari
	LanguageTypeMHS    LanguageType = "mhs"      // Buru (Indonesia)
	LanguageTypeMHT    LanguageType = "mht"      // Mandahuaca
	LanguageTypeMHU    LanguageType = "mhu"      // Digaro-Mishmi and Darang Deng
	LanguageTypeMHW    LanguageType = "mhw"      // Mbukushu
	LanguageTypeMHX    LanguageType = "mhx"      // Maru and Lhaovo
	LanguageTypeMHY    LanguageType = "mhy"      // Ma'anyan
	LanguageTypeMHZ    LanguageType = "mhz"      // Mor (Mor Islands)
	LanguageTypeMIA    LanguageType = "mia"      // Miami
	LanguageTypeMIB    LanguageType = "mib"      // Atatláhuca Mixtec
	LanguageTypeMIC    LanguageType = "mic"      // Mi'kmaq and Micmac
	LanguageTypeMID    LanguageType = "mid"      // Mandaic
	LanguageTypeMIE    LanguageType = "mie"      // Ocotepec Mixtec
	LanguageTypeMIF    LanguageType = "mif"      // Mofu-Gudur
	LanguageTypeMIG    LanguageType = "mig"      // San Miguel El Grande Mixtec
	LanguageTypeMIH    LanguageType = "mih"      // Chayuco Mixtec
	LanguageTypeMII    LanguageType = "mii"      // Chigmecatitlán Mixtec
	LanguageTypeMIJ    LanguageType = "mij"      // Abar and Mungbam
	LanguageTypeMIK    LanguageType = "mik"      // Mikasuki
	LanguageTypeMIL    LanguageType = "mil"      // Peñoles Mixtec
	LanguageTypeMIM    LanguageType = "mim"      // Alacatlatzala Mixtec
	LanguageTypeMIN    LanguageType = "min"      // Minangkabau
	LanguageTypeMIO    LanguageType = "mio"      // Pinotepa Nacional Mixtec
	LanguageTypeMIP    LanguageType = "mip"      // Apasco-Apoala Mixtec
	LanguageTypeMIQ    LanguageType = "miq"      // Mískito
	LanguageTypeMIR    LanguageType = "mir"      // Isthmus Mixe
	LanguageTypeMIS    LanguageType = "mis"      // Uncoded languages
	LanguageTypeMIT    LanguageType = "mit"      // Southern Puebla Mixtec
	LanguageTypeMIU    LanguageType = "miu"      // Cacaloxtepec Mixtec
	LanguageTypeMIW    LanguageType = "miw"      // Akoye
	LanguageTypeMIX    LanguageType = "mix"      // Mixtepec Mixtec
	LanguageTypeMIY    LanguageType = "miy"      // Ayutla Mixtec
	LanguageTypeMIZ    LanguageType = "miz"      // Coatzospan Mixtec
	LanguageTypeMJA    LanguageType = "mja"      // Mahei
	LanguageTypeMJC    LanguageType = "mjc"      // San Juan Colorado Mixtec
	LanguageTypeMJD    LanguageType = "mjd"      // Northwest Maidu
	LanguageTypeMJE    LanguageType = "mje"      // Muskum
	LanguageTypeMJG    LanguageType = "mjg"      // Tu
	LanguageTypeMJH    LanguageType = "mjh"      // Mwera (Nyasa)
	LanguageTypeMJI    LanguageType = "mji"      // Kim Mun
	LanguageTypeMJJ    LanguageType = "mjj"      // Mawak
	LanguageTypeMJK    LanguageType = "mjk"      // Matukar
	LanguageTypeMJL    LanguageType = "mjl"      // Mandeali
	LanguageTypeMJM    LanguageType = "mjm"      // Medebur
	LanguageTypeMJN    LanguageType = "mjn"      // Ma (Papua New Guinea)
	LanguageTypeMJO    LanguageType = "mjo"      // Malankuravan
	LanguageTypeMJP    LanguageType = "mjp"      // Malapandaram
	LanguageTypeMJQ    LanguageType = "mjq"      // Malaryan
	LanguageTypeMJR    LanguageType = "mjr"      // Malavedan
	LanguageTypeMJS    LanguageType = "mjs"      // Miship
	LanguageTypeMJT    LanguageType = "mjt"      // Sauria Paharia
	LanguageTypeMJU    LanguageType = "mju"      // Manna-Dora
	LanguageTypeMJV    LanguageType = "mjv"      // Mannan
	LanguageTypeMJW    LanguageType = "mjw"      // Karbi
	LanguageTypeMJX    LanguageType = "mjx"      // Mahali
	LanguageTypeMJY    LanguageType = "mjy"      // Mahican
	LanguageTypeMJZ    LanguageType = "mjz"      // Majhi
	LanguageTypeMKA    LanguageType = "mka"      // Mbre
	LanguageTypeMKB    LanguageType = "mkb"      // Mal Paharia
	LanguageTypeMKC    LanguageType = "mkc"      // Siliput
	LanguageTypeMKE    LanguageType = "mke"      // Mawchi
	LanguageTypeMKF    LanguageType = "mkf"      // Miya
	LanguageTypeMKG    LanguageType = "mkg"      // Mak (China)
	LanguageTypeMKH    LanguageType = "mkh"      // Mon-Khmer languages
	LanguageTypeMKI    LanguageType = "mki"      // Dhatki
	LanguageTypeMKJ    LanguageType = "mkj"      // Mokilese
	LanguageTypeMKK    LanguageType = "mkk"      // Byep
	LanguageTypeMKL    LanguageType = "mkl"      // Mokole
	LanguageTypeMKM    LanguageType = "mkm"      // Moklen
	LanguageTypeMKN    LanguageType = "mkn"      // Kupang Malay
	LanguageTypeMKO    LanguageType = "mko"      // Mingang Doso
	LanguageTypeMKP    LanguageType = "mkp"      // Moikodi
	LanguageTypeMKQ    LanguageType = "mkq"      // Bay Miwok
	LanguageTypeMKR    LanguageType = "mkr"      // Malas
	LanguageTypeMKS    LanguageType = "mks"      // Silacayoapan Mixtec
	LanguageTypeMKT    LanguageType = "mkt"      // Vamale
	LanguageTypeMKU    LanguageType = "mku"      // Konyanka Maninka
	LanguageTypeMKV    LanguageType = "mkv"      // Mafea
	LanguageTypeMKW    LanguageType = "mkw"      // Kituba (Congo)
	LanguageTypeMKX    LanguageType = "mkx"      // Kinamiging Manobo
	LanguageTypeMKY    LanguageType = "mky"      // East Makian
	LanguageTypeMKZ    LanguageType = "mkz"      // Makasae
	LanguageTypeMLA    LanguageType = "mla"      // Malo
	LanguageTypeMLB    LanguageType = "mlb"      // Mbule
	LanguageTypeMLC    LanguageType = "mlc"      // Cao Lan
	LanguageTypeMLD    LanguageType = "mld"      // Malakhel
	LanguageTypeMLE    LanguageType = "mle"      // Manambu
	LanguageTypeMLF    LanguageType = "mlf"      // Mal
	LanguageTypeMLH    LanguageType = "mlh"      // Mape
	LanguageTypeMLI    LanguageType = "mli"      // Malimpung
	LanguageTypeMLJ    LanguageType = "mlj"      // Miltu
	LanguageTypeMLK    LanguageType = "mlk"      // Ilwana and Kiwilwana
	LanguageTypeMLL    LanguageType = "mll"      // Malua Bay
	LanguageTypeMLM    LanguageType = "mlm"      // Mulam
	LanguageTypeMLN    LanguageType = "mln"      // Malango
	LanguageTypeMLO    LanguageType = "mlo"      // Mlomp
	LanguageTypeMLP    LanguageType = "mlp"      // Bargam
	LanguageTypeMLQ    LanguageType = "mlq"      // Western Maninkakan
	LanguageTypeMLR    LanguageType = "mlr"      // Vame
	LanguageTypeMLS    LanguageType = "mls"      // Masalit
	LanguageTypeMLU    LanguageType = "mlu"      // To'abaita
	LanguageTypeMLV    LanguageType = "mlv"      // Motlav and Mwotlap
	LanguageTypeMLW    LanguageType = "mlw"      // Moloko
	LanguageTypeMLX    LanguageType = "mlx"      // Malfaxal and Naha'ai
	LanguageTypeMLZ    LanguageType = "mlz"      // Malaynon
	LanguageTypeMMA    LanguageType = "mma"      // Mama
	LanguageTypeMMB    LanguageType = "mmb"      // Momina
	LanguageTypeMMC    LanguageType = "mmc"      // Michoacán Mazahua
	LanguageTypeMMD    LanguageType = "mmd"      // Maonan
	LanguageTypeMME    LanguageType = "mme"      // Mae
	LanguageTypeMMF    LanguageType = "mmf"      // Mundat
	LanguageTypeMMG    LanguageType = "mmg"      // North Ambrym
	LanguageTypeMMH    LanguageType = "mmh"      // Mehináku
	LanguageTypeMMI    LanguageType = "mmi"      // Musar
	LanguageTypeMMJ    LanguageType = "mmj"      // Majhwar
	LanguageTypeMMK    LanguageType = "mmk"      // Mukha-Dora
	LanguageTypeMML    LanguageType = "mml"      // Man Met
	LanguageTypeMMM    LanguageType = "mmm"      // Maii
	LanguageTypeMMN    LanguageType = "mmn"      // Mamanwa
	LanguageTypeMMO    LanguageType = "mmo"      // Mangga Buang
	LanguageTypeMMP    LanguageType = "mmp"      // Siawi
	LanguageTypeMMQ    LanguageType = "mmq"      // Musak
	LanguageTypeMMR    LanguageType = "mmr"      // Western Xiangxi Miao
	LanguageTypeMMT    LanguageType = "mmt"      // Malalamai
	LanguageTypeMMU    LanguageType = "mmu"      // Mmaala
	LanguageTypeMMV    LanguageType = "mmv"      // Miriti
	LanguageTypeMMW    LanguageType = "mmw"      // Emae
	LanguageTypeMMX    LanguageType = "mmx"      // Madak
	LanguageTypeMMY    LanguageType = "mmy"      // Migaama
	LanguageTypeMMZ    LanguageType = "mmz"      // Mabaale
	LanguageTypeMNA    LanguageType = "mna"      // Mbula
	LanguageTypeMNB    LanguageType = "mnb"      // Muna
	LanguageTypeMNC    LanguageType = "mnc"      // Manchu
	LanguageTypeMND    LanguageType = "mnd"      // Mondé
	LanguageTypeMNE    LanguageType = "mne"      // Naba
	LanguageTypeMNF    LanguageType = "mnf"      // Mundani
	LanguageTypeMNG    LanguageType = "mng"      // Eastern Mnong
	LanguageTypeMNH    LanguageType = "mnh"      // Mono (Democratic Republic of Congo)
	LanguageTypeMNI    LanguageType = "mni"      // Manipuri
	LanguageTypeMNJ    LanguageType = "mnj"      // Munji
	LanguageTypeMNK    LanguageType = "mnk"      // Mandinka
	LanguageTypeMNL    LanguageType = "mnl"      // Tiale
	LanguageTypeMNM    LanguageType = "mnm"      // Mapena
	LanguageTypeMNN    LanguageType = "mnn"      // Southern Mnong
	LanguageTypeMNO    LanguageType = "mno"      // Manobo languages
	LanguageTypeMNP    LanguageType = "mnp"      // Min Bei Chinese
	LanguageTypeMNQ    LanguageType = "mnq"      // Minriq
	LanguageTypeMNR    LanguageType = "mnr"      // Mono (USA)
	LanguageTypeMNS    LanguageType = "mns"      // Mansi
	LanguageTypeMNT    LanguageType = "mnt"      // Maykulan
	LanguageTypeMNU    LanguageType = "mnu"      // Mer
	LanguageTypeMNV    LanguageType = "mnv"      // Rennell-Bellona
	LanguageTypeMNW    LanguageType = "mnw"      // Mon
	LanguageTypeMNX    LanguageType = "mnx"      // Manikion
	LanguageTypeMNY    LanguageType = "mny"      // Manyawa
	LanguageTypeMNZ    LanguageType = "mnz"      // Moni
	LanguageTypeMOA    LanguageType = "moa"      // Mwan
	LanguageTypeMOC    LanguageType = "moc"      // Mocoví
	LanguageTypeMOD    LanguageType = "mod"      // Mobilian
	LanguageTypeMOE    LanguageType = "moe"      // Montagnais
	LanguageTypeMOF    LanguageType = "mof"      // Mohegan-Montauk-Narragansett
	LanguageTypeMOG    LanguageType = "mog"      // Mongondow
	LanguageTypeMOH    LanguageType = "moh"      // Mohawk
	LanguageTypeMOI    LanguageType = "moi"      // Mboi
	LanguageTypeMOJ    LanguageType = "moj"      // Monzombo
	LanguageTypeMOK    LanguageType = "mok"      // Morori
	LanguageTypeMOM    LanguageType = "mom"      // Mangue
	LanguageTypeMOO    LanguageType = "moo"      // Monom
	LanguageTypeMOP    LanguageType = "mop"      // Mopán Maya
	LanguageTypeMOQ    LanguageType = "moq"      // Mor (Bomberai Peninsula)
	LanguageTypeMOR    LanguageType = "mor"      // Moro
	LanguageTypeMOS    LanguageType = "mos"      // Mossi
	LanguageTypeMOT    LanguageType = "mot"      // Barí
	LanguageTypeMOU    LanguageType = "mou"      // Mogum
	LanguageTypeMOV    LanguageType = "mov"      // Mohave
	LanguageTypeMOW    LanguageType = "mow"      // Moi (Congo)
	LanguageTypeMOX    LanguageType = "mox"      // Molima
	LanguageTypeMOY    LanguageType = "moy"      // Shekkacho
	LanguageTypeMOZ    LanguageType = "moz"      // Mukulu and Gergiko
	LanguageTypeMPA    LanguageType = "mpa"      // Mpoto
	LanguageTypeMPB    LanguageType = "mpb"      // Mullukmulluk
	LanguageTypeMPC    LanguageType = "mpc"      // Mangarayi
	LanguageTypeMPD    LanguageType = "mpd"      // Machinere
	LanguageTypeMPE    LanguageType = "mpe"      // Majang
	LanguageTypeMPG    LanguageType = "mpg"      // Marba
	LanguageTypeMPH    LanguageType = "mph"      // Maung
	LanguageTypeMPI    LanguageType = "mpi"      // Mpade
	LanguageTypeMPJ    LanguageType = "mpj"      // Martu Wangka
	LanguageTypeMPK    LanguageType = "mpk"      // Mbara (Chad)
	LanguageTypeMPL    LanguageType = "mpl"      // Middle Watut
	LanguageTypeMPM    LanguageType = "mpm"      // Yosondúa Mixtec
	LanguageTypeMPN    LanguageType = "mpn"      // Mindiri
	LanguageTypeMPO    LanguageType = "mpo"      // Miu
	LanguageTypeMPP    LanguageType = "mpp"      // Migabac
	LanguageTypeMPQ    LanguageType = "mpq"      // Matís
	LanguageTypeMPR    LanguageType = "mpr"      // Vangunu
	LanguageTypeMPS    LanguageType = "mps"      // Dadibi
	LanguageTypeMPT    LanguageType = "mpt"      // Mian
	LanguageTypeMPU    LanguageType = "mpu"      // Makuráp
	LanguageTypeMPV    LanguageType = "mpv"      // Mungkip
	LanguageTypeMPW    LanguageType = "mpw"      // Mapidian
	LanguageTypeMPX    LanguageType = "mpx"      // Misima-Panaeati
	LanguageTypeMPY    LanguageType = "mpy"      // Mapia
	LanguageTypeMPZ    LanguageType = "mpz"      // Mpi
	LanguageTypeMQA    LanguageType = "mqa"      // Maba (Indonesia)
	LanguageTypeMQB    LanguageType = "mqb"      // Mbuko
	LanguageTypeMQC    LanguageType = "mqc"      // Mangole
	LanguageTypeMQE    LanguageType = "mqe"      // Matepi
	LanguageTypeMQF    LanguageType = "mqf"      // Momuna
	LanguageTypeMQG    LanguageType = "mqg"      // Kota Bangun Kutai Malay
	LanguageTypeMQH    LanguageType = "mqh"      // Tlazoyaltepec Mixtec
	LanguageTypeMQI    LanguageType = "mqi"      // Mariri
	LanguageTypeMQJ    LanguageType = "mqj"      // Mamasa
	LanguageTypeMQK    LanguageType = "mqk"      // Rajah Kabunsuwan Manobo
	LanguageTypeMQL    LanguageType = "mql"      // Mbelime
	LanguageTypeMQM    LanguageType = "mqm"      // South Marquesan
	LanguageTypeMQN    LanguageType = "mqn"      // Moronene
	LanguageTypeMQO    LanguageType = "mqo"      // Modole
	LanguageTypeMQP    LanguageType = "mqp"      // Manipa
	LanguageTypeMQQ    LanguageType = "mqq"      // Minokok
	LanguageTypeMQR    LanguageType = "mqr"      // Mander
	LanguageTypeMQS    LanguageType = "mqs"      // West Makian
	LanguageTypeMQT    LanguageType = "mqt"      // Mok
	LanguageTypeMQU    LanguageType = "mqu"      // Mandari
	LanguageTypeMQV    LanguageType = "mqv"      // Mosimo
	LanguageTypeMQW    LanguageType = "mqw"      // Murupi
	LanguageTypeMQX    LanguageType = "mqx"      // Mamuju
	LanguageTypeMQY    LanguageType = "mqy"      // Manggarai
	LanguageTypeMQZ    LanguageType = "mqz"      // Pano
	LanguageTypeMRA    LanguageType = "mra"      // Mlabri
	LanguageTypeMRB    LanguageType = "mrb"      // Marino
	LanguageTypeMRC    LanguageType = "mrc"      // Maricopa
	LanguageTypeMRD    LanguageType = "mrd"      // Western Magar
	LanguageTypeMRE    LanguageType = "mre"      // Martha's Vineyard Sign Language
	LanguageTypeMRF    LanguageType = "mrf"      // Elseng
	LanguageTypeMRG    LanguageType = "mrg"      // Mising
	LanguageTypeMRH    LanguageType = "mrh"      // Mara Chin
	LanguageTypeMRJ    LanguageType = "mrj"      // Western Mari
	LanguageTypeMRK    LanguageType = "mrk"      // Hmwaveke
	LanguageTypeMRL    LanguageType = "mrl"      // Mortlockese
	LanguageTypeMRM    LanguageType = "mrm"      // Merlav and Mwerlap
	LanguageTypeMRN    LanguageType = "mrn"      // Cheke Holo
	LanguageTypeMRO    LanguageType = "mro"      // Mru
	LanguageTypeMRP    LanguageType = "mrp"      // Morouas
	LanguageTypeMRQ    LanguageType = "mrq"      // North Marquesan
	LanguageTypeMRR    LanguageType = "mrr"      // Maria (India)
	LanguageTypeMRS    LanguageType = "mrs"      // Maragus
	LanguageTypeMRT    LanguageType = "mrt"      // Marghi Central
	LanguageTypeMRU    LanguageType = "mru"      // Mono (Cameroon)
	LanguageTypeMRV    LanguageType = "mrv"      // Mangareva
	LanguageTypeMRW    LanguageType = "mrw"      // Maranao
	LanguageTypeMRX    LanguageType = "mrx"      // Maremgi and Dineor
	LanguageTypeMRY    LanguageType = "mry"      // Mandaya
	LanguageTypeMRZ    LanguageType = "mrz"      // Marind
	LanguageTypeMSB    LanguageType = "msb"      // Masbatenyo
	LanguageTypeMSC    LanguageType = "msc"      // Sankaran Maninka
	LanguageTypeMSD    LanguageType = "msd"      // Yucatec Maya Sign Language
	LanguageTypeMSE    LanguageType = "mse"      // Musey
	LanguageTypeMSF    LanguageType = "msf"      // Mekwei
	LanguageTypeMSG    LanguageType = "msg"      // Moraid
	LanguageTypeMSH    LanguageType = "msh"      // Masikoro Malagasy
	LanguageTypeMSI    LanguageType = "msi"      // Sabah Malay
	LanguageTypeMSJ    LanguageType = "msj"      // Ma (Democratic Republic of Congo)
	LanguageTypeMSK    LanguageType = "msk"      // Mansaka
	LanguageTypeMSL    LanguageType = "msl"      // Molof and Poule
	LanguageTypeMSM    LanguageType = "msm"      // Agusan Manobo
	LanguageTypeMSN    LanguageType = "msn"      // Vurës
	LanguageTypeMSO    LanguageType = "mso"      // Mombum
	LanguageTypeMSP    LanguageType = "msp"      // Maritsauá
	LanguageTypeMSQ    LanguageType = "msq"      // Caac
	LanguageTypeMSR    LanguageType = "msr"      // Mongolian Sign Language
	LanguageTypeMSS    LanguageType = "mss"      // West Masela
	LanguageTypeMST    LanguageType = "mst"      // Cataelano Mandaya
	LanguageTypeMSU    LanguageType = "msu"      // Musom
	LanguageTypeMSV    LanguageType = "msv"      // Maslam
	LanguageTypeMSW    LanguageType = "msw"      // Mansoanka
	LanguageTypeMSX    LanguageType = "msx"      // Moresada
	LanguageTypeMSY    LanguageType = "msy"      // Aruamu
	LanguageTypeMSZ    LanguageType = "msz"      // Momare
	LanguageTypeMTA    LanguageType = "mta"      // Cotabato Manobo
	LanguageTypeMTB    LanguageType = "mtb"      // Anyin Morofo
	LanguageTypeMTC    LanguageType = "mtc"      // Munit
	LanguageTypeMTD    LanguageType = "mtd"      // Mualang
	LanguageTypeMTE    LanguageType = "mte"      // Mono (Solomon Islands)
	LanguageTypeMTF    LanguageType = "mtf"      // Murik (Papua New Guinea)
	LanguageTypeMTG    LanguageType = "mtg"      // Una
	LanguageTypeMTH    LanguageType = "mth"      // Munggui
	LanguageTypeMTI    LanguageType = "mti"      // Maiwa (Papua New Guinea)
	LanguageTypeMTJ    LanguageType = "mtj"      // Moskona
	LanguageTypeMTK    LanguageType = "mtk"      // Mbe'
	LanguageTypeMTL    LanguageType = "mtl"      // Montol
	LanguageTypeMTM    LanguageType = "mtm"      // Mator
	LanguageTypeMTN    LanguageType = "mtn"      // Matagalpa
	LanguageTypeMTO    LanguageType = "mto"      // Totontepec Mixe
	LanguageTypeMTP    LanguageType = "mtp"      // Wichí Lhamtés Nocten
	LanguageTypeMTQ    LanguageType = "mtq"      // Muong
	LanguageTypeMTR    LanguageType = "mtr"      // Mewari
	LanguageTypeMTS    LanguageType = "mts"      // Yora
	LanguageTypeMTT    LanguageType = "mtt"      // Mota
	LanguageTypeMTU    LanguageType = "mtu"      // Tututepec Mixtec
	LanguageTypeMTV    LanguageType = "mtv"      // Asaro'o
	LanguageTypeMTW    LanguageType = "mtw"      // Southern Binukidnon
	LanguageTypeMTX    LanguageType = "mtx"      // Tidaá Mixtec
	LanguageTypeMTY    LanguageType = "mty"      // Nabi
	LanguageTypeMUA    LanguageType = "mua"      // Mundang
	LanguageTypeMUB    LanguageType = "mub"      // Mubi
	LanguageTypeMUC    LanguageType = "muc"      // Ajumbu
	LanguageTypeMUD    LanguageType = "mud"      // Mednyj Aleut
	LanguageTypeMUE    LanguageType = "mue"      // Media Lengua
	LanguageTypeMUG    LanguageType = "mug"      // Musgu
	LanguageTypeMUH    LanguageType = "muh"      // Mündü
	LanguageTypeMUI    LanguageType = "mui"      // Musi
	LanguageTypeMUJ    LanguageType = "muj"      // Mabire
	LanguageTypeMUK    LanguageType = "muk"      // Mugom
	LanguageTypeMUL    LanguageType = "mul"      // Multiple languages
	LanguageTypeMUM    LanguageType = "mum"      // Maiwala
	LanguageTypeMUN    LanguageType = "mun"      // Munda languages
	LanguageTypeMUO    LanguageType = "muo"      // Nyong
	LanguageTypeMUP    LanguageType = "mup"      // Malvi
	LanguageTypeMUQ    LanguageType = "muq"      // Eastern Xiangxi Miao
	LanguageTypeMUR    LanguageType = "mur"      // Murle
	LanguageTypeMUS    LanguageType = "mus"      // Creek
	LanguageTypeMUT    LanguageType = "mut"      // Western Muria
	LanguageTypeMUU    LanguageType = "muu"      // Yaaku
	LanguageTypeMUV    LanguageType = "muv"      // Muthuvan
	LanguageTypeMUX    LanguageType = "mux"      // Bo-Ung
	LanguageTypeMUY    LanguageType = "muy"      // Muyang
	LanguageTypeMUZ    LanguageType = "muz"      // Mursi
	LanguageTypeMVA    LanguageType = "mva"      // Manam
	LanguageTypeMVB    LanguageType = "mvb"      // Mattole
	LanguageTypeMVD    LanguageType = "mvd"      // Mamboru
	LanguageTypeMVE    LanguageType = "mve"      // Marwari (Pakistan)
	LanguageTypeMVF    LanguageType = "mvf"      // Peripheral Mongolian
	LanguageTypeMVG    LanguageType = "mvg"      // Yucuañe Mixtec
	LanguageTypeMVH    LanguageType = "mvh"      // Mulgi
	LanguageTypeMVI    LanguageType = "mvi"      // Miyako
	LanguageTypeMVK    LanguageType = "mvk"      // Mekmek
	LanguageTypeMVL    LanguageType = "mvl"      // Mbara (Australia)
	LanguageTypeMVM    LanguageType = "mvm"      // Muya
	LanguageTypeMVN    LanguageType = "mvn"      // Minaveha
	LanguageTypeMVO    LanguageType = "mvo"      // Marovo
	LanguageTypeMVP    LanguageType = "mvp"      // Duri
	LanguageTypeMVQ    LanguageType = "mvq"      // Moere
	LanguageTypeMVR    LanguageType = "mvr"      // Marau
	LanguageTypeMVS    LanguageType = "mvs"      // Massep
	LanguageTypeMVT    LanguageType = "mvt"      // Mpotovoro
	LanguageTypeMVU    LanguageType = "mvu"      // Marfa
	LanguageTypeMVV    LanguageType = "mvv"      // Tagal Murut
	LanguageTypeMVW    LanguageType = "mvw"      // Machinga
	LanguageTypeMVX    LanguageType = "mvx"      // Meoswar
	LanguageTypeMVY    LanguageType = "mvy"      // Indus Kohistani
	LanguageTypeMVZ    LanguageType = "mvz"      // Mesqan
	LanguageTypeMWA    LanguageType = "mwa"      // Mwatebu
	LanguageTypeMWB    LanguageType = "mwb"      // Juwal
	LanguageTypeMWC    LanguageType = "mwc"      // Are
	LanguageTypeMWD    LanguageType = "mwd"      // Mudbura
	LanguageTypeMWE    LanguageType = "mwe"      // Mwera (Chimwera)
	LanguageTypeMWF    LanguageType = "mwf"      // Murrinh-Patha
	LanguageTypeMWG    LanguageType = "mwg"      // Aiklep
	LanguageTypeMWH    LanguageType = "mwh"      // Mouk-Aria
	LanguageTypeMWI    LanguageType = "mwi"      // Labo and Ninde
	LanguageTypeMWJ    LanguageType = "mwj"      // Maligo
	LanguageTypeMWK    LanguageType = "mwk"      // Kita Maninkakan
	LanguageTypeMWL    LanguageType = "mwl"      // Mirandese
	LanguageTypeMWM    LanguageType = "mwm"      // Sar
	LanguageTypeMWN    LanguageType = "mwn"      // Nyamwanga
	LanguageTypeMWO    LanguageType = "mwo"      // Central Maewo
	LanguageTypeMWP    LanguageType = "mwp"      // Kala Lagaw Ya
	LanguageTypeMWQ    LanguageType = "mwq"      // Mün Chin
	LanguageTypeMWR    LanguageType = "mwr"      // Marwari
	LanguageTypeMWS    LanguageType = "mws"      // Mwimbi-Muthambi
	LanguageTypeMWT    LanguageType = "mwt"      // Moken
	LanguageTypeMWU    LanguageType = "mwu"      // Mittu
	LanguageTypeMWV    LanguageType = "mwv"      // Mentawai
	LanguageTypeMWW    LanguageType = "mww"      // Hmong Daw
	LanguageTypeMWX    LanguageType = "mwx"      // Mediak
	LanguageTypeMWY    LanguageType = "mwy"      // Mosiro
	LanguageTypeMWZ    LanguageType = "mwz"      // Moingi
	LanguageTypeMXA    LanguageType = "mxa"      // Northwest Oaxaca Mixtec
	LanguageTypeMXB    LanguageType = "mxb"      // Tezoatlán Mixtec
	LanguageTypeMXC    LanguageType = "mxc"      // Manyika
	LanguageTypeMXD    LanguageType = "mxd"      // Modang
	LanguageTypeMXE    LanguageType = "mxe"      // Mele-Fila
	LanguageTypeMXF    LanguageType = "mxf"      // Malgbe
	LanguageTypeMXG    LanguageType = "mxg"      // Mbangala
	LanguageTypeMXH    LanguageType = "mxh"      // Mvuba
	LanguageTypeMXI    LanguageType = "mxi"      // Mozarabic
	LanguageTypeMXJ    LanguageType = "mxj"      // Miju-Mishmi and Geman Deng
	LanguageTypeMXK    LanguageType = "mxk"      // Monumbo
	LanguageTypeMXL    LanguageType = "mxl"      // Maxi Gbe
	LanguageTypeMXM    LanguageType = "mxm"      // Meramera
	LanguageTypeMXN    LanguageType = "mxn"      // Moi (Indonesia)
	LanguageTypeMXO    LanguageType = "mxo"      // Mbowe
	LanguageTypeMXP    LanguageType = "mxp"      // Tlahuitoltepec Mixe
	LanguageTypeMXQ    LanguageType = "mxq"      // Juquila Mixe
	LanguageTypeMXR    LanguageType = "mxr"      // Murik (Malaysia)
	LanguageTypeMXS    LanguageType = "mxs"      // Huitepec Mixtec
	LanguageTypeMXT    LanguageType = "mxt"      // Jamiltepec Mixtec
	LanguageTypeMXU    LanguageType = "mxu"      // Mada (Cameroon)
	LanguageTypeMXV    LanguageType = "mxv"      // Metlatónoc Mixtec
	LanguageTypeMXW    LanguageType = "mxw"      // Namo
	LanguageTypeMXX    LanguageType = "mxx"      // Mahou and Mawukakan
	LanguageTypeMXY    LanguageType = "mxy"      // Southeastern Nochixtlán Mixtec
	LanguageTypeMXZ    LanguageType = "mxz"      // Central Masela
	LanguageTypeMYB    LanguageType = "myb"      // Mbay
	LanguageTypeMYC    LanguageType = "myc"      // Mayeka
	LanguageTypeMYD    LanguageType = "myd"      // Maramba
	LanguageTypeMYE    LanguageType = "mye"      // Myene
	LanguageTypeMYF    LanguageType = "myf"      // Bambassi
	LanguageTypeMYG    LanguageType = "myg"      // Manta
	LanguageTypeMYH    LanguageType = "myh"      // Makah
	LanguageTypeMYI    LanguageType = "myi"      // Mina (India)
	LanguageTypeMYJ    LanguageType = "myj"      // Mangayat
	LanguageTypeMYK    LanguageType = "myk"      // Mamara Senoufo
	LanguageTypeMYL    LanguageType = "myl"      // Moma
	LanguageTypeMYM    LanguageType = "mym"      // Me'en
	LanguageTypeMYN    LanguageType = "myn"      // Mayan languages
	LanguageTypeMYO    LanguageType = "myo"      // Anfillo
	LanguageTypeMYP    LanguageType = "myp"      // Pirahã
	LanguageTypeMYQ    LanguageType = "myq"      // Forest Maninka
	LanguageTypeMYR    LanguageType = "myr"      // Muniche
	LanguageTypeMYS    LanguageType = "mys"      // Mesmes
	LanguageTypeMYT    LanguageType = "myt"      // Sangab Mandaya
	LanguageTypeMYU    LanguageType = "myu"      // Mundurukú
	LanguageTypeMYV    LanguageType = "myv"      // Erzya
	LanguageTypeMYW    LanguageType = "myw"      // Muyuw
	LanguageTypeMYX    LanguageType = "myx"      // Masaaba
	LanguageTypeMYY    LanguageType = "myy"      // Macuna
	LanguageTypeMYZ    LanguageType = "myz"      // Classical Mandaic
	LanguageTypeMZA    LanguageType = "mza"      // Santa María Zacatepec Mixtec
	LanguageTypeMZB    LanguageType = "mzb"      // Tumzabt
	LanguageTypeMZC    LanguageType = "mzc"      // Madagascar Sign Language
	LanguageTypeMZD    LanguageType = "mzd"      // Malimba
	LanguageTypeMZE    LanguageType = "mze"      // Morawa
	LanguageTypeMZG    LanguageType = "mzg"      // Monastic Sign Language
	LanguageTypeMZH    LanguageType = "mzh"      // Wichí Lhamtés Güisnay
	LanguageTypeMZI    LanguageType = "mzi"      // Ixcatlán Mazatec
	LanguageTypeMZJ    LanguageType = "mzj"      // Manya
	LanguageTypeMZK    LanguageType = "mzk"      // Nigeria Mambila
	LanguageTypeMZL    LanguageType = "mzl"      // Mazatlán Mixe
	LanguageTypeMZM    LanguageType = "mzm"      // Mumuye
	LanguageTypeMZN    LanguageType = "mzn"      // Mazanderani
	LanguageTypeMZO    LanguageType = "mzo"      // Matipuhy
	LanguageTypeMZP    LanguageType = "mzp"      // Movima
	LanguageTypeMZQ    LanguageType = "mzq"      // Mori Atas
	LanguageTypeMZR    LanguageType = "mzr"      // Marúbo
	LanguageTypeMZS    LanguageType = "mzs"      // Macanese
	LanguageTypeMZT    LanguageType = "mzt"      // Mintil
	LanguageTypeMZU    LanguageType = "mzu"      // Inapang
	LanguageTypeMZV    LanguageType = "mzv"      // Manza
	LanguageTypeMZW    LanguageType = "mzw"      // Deg
	LanguageTypeMZX    LanguageType = "mzx"      // Mawayana
	LanguageTypeMZY    LanguageType = "mzy"      // Mozambican Sign Language
	LanguageTypeMZZ    LanguageType = "mzz"      // Maiadomu
	LanguageTypeNAA    LanguageType = "naa"      // Namla
	LanguageTypeNAB    LanguageType = "nab"      // Southern Nambikuára
	LanguageTypeNAC    LanguageType = "nac"      // Narak
	LanguageTypeNAD    LanguageType = "nad"      // Nijadali
	LanguageTypeNAE    LanguageType = "nae"      // Naka'ela
	LanguageTypeNAF    LanguageType = "naf"      // Nabak
	LanguageTypeNAG    LanguageType = "nag"      // Naga Pidgin
	LanguageTypeNAH    LanguageType = "nah"      // Nahuatl languages
	LanguageTypeNAI    LanguageType = "nai"      // North American Indian languages
	LanguageTypeNAJ    LanguageType = "naj"      // Nalu
	LanguageTypeNAK    LanguageType = "nak"      // Nakanai
	LanguageTypeNAL    LanguageType = "nal"      // Nalik
	LanguageTypeNAM    LanguageType = "nam"      // Ngan'gityemerri
	LanguageTypeNAN    LanguageType = "nan"      // Min Nan Chinese
	LanguageTypeNAO    LanguageType = "nao"      // Naaba
	LanguageTypeNAP    LanguageType = "nap"      // Neapolitan
	LanguageTypeNAQ    LanguageType = "naq"      // Nama (Namibia)
	LanguageTypeNAR    LanguageType = "nar"      // Iguta
	LanguageTypeNAS    LanguageType = "nas"      // Naasioi
	LanguageTypeNAT    LanguageType = "nat"      // Hungworo
	LanguageTypeNAW    LanguageType = "naw"      // Nawuri
	LanguageTypeNAX    LanguageType = "nax"      // Nakwi
	LanguageTypeNAY    LanguageType = "nay"      // Narrinyeri
	LanguageTypeNAZ    LanguageType = "naz"      // Coatepec Nahuatl
	LanguageTypeNBA    LanguageType = "nba"      // Nyemba
	LanguageTypeNBB    LanguageType = "nbb"      // Ndoe
	LanguageTypeNBC    LanguageType = "nbc"      // Chang Naga
	LanguageTypeNBD    LanguageType = "nbd"      // Ngbinda
	LanguageTypeNBE    LanguageType = "nbe"      // Konyak Naga
	LanguageTypeNBF    LanguageType = "nbf"      // Naxi
	LanguageTypeNBG    LanguageType = "nbg"      // Nagarchal
	LanguageTypeNBH    LanguageType = "nbh"      // Ngamo
	LanguageTypeNBI    LanguageType = "nbi"      // Mao Naga
	LanguageTypeNBJ    LanguageType = "nbj"      // Ngarinman
	LanguageTypeNBK    LanguageType = "nbk"      // Nake
	LanguageTypeNBM    LanguageType = "nbm"      // Ngbaka Ma'bo
	LanguageTypeNBN    LanguageType = "nbn"      // Kuri
	LanguageTypeNBO    LanguageType = "nbo"      // Nkukoli
	LanguageTypeNBP    LanguageType = "nbp"      // Nnam
	LanguageTypeNBQ    LanguageType = "nbq"      // Nggem
	LanguageTypeNBR    LanguageType = "nbr"      // Numana-Nunku-Gbantu-Numbu
	LanguageTypeNBS    LanguageType = "nbs"      // Namibian Sign Language
	LanguageTypeNBT    LanguageType = "nbt"      // Na
	LanguageTypeNBU    LanguageType = "nbu"      // Rongmei Naga
	LanguageTypeNBV    LanguageType = "nbv"      // Ngamambo
	LanguageTypeNBW    LanguageType = "nbw"      // Southern Ngbandi
	LanguageTypeNBX    LanguageType = "nbx"      // Ngura
	LanguageTypeNBY    LanguageType = "nby"      // Ningera
	LanguageTypeNCA    LanguageType = "nca"      // Iyo
	LanguageTypeNCB    LanguageType = "ncb"      // Central Nicobarese
	LanguageTypeNCC    LanguageType = "ncc"      // Ponam
	LanguageTypeNCD    LanguageType = "ncd"      // Nachering
	LanguageTypeNCE    LanguageType = "nce"      // Yale
	LanguageTypeNCF    LanguageType = "ncf"      // Notsi
	LanguageTypeNCG    LanguageType = "ncg"      // Nisga'a
	LanguageTypeNCH    LanguageType = "nch"      // Central Huasteca Nahuatl
	LanguageTypeNCI    LanguageType = "nci"      // Classical Nahuatl
	LanguageTypeNCJ    LanguageType = "ncj"      // Northern Puebla Nahuatl
	LanguageTypeNCK    LanguageType = "nck"      // Nakara
	LanguageTypeNCL    LanguageType = "ncl"      // Michoacán Nahuatl
	LanguageTypeNCM    LanguageType = "ncm"      // Nambo
	LanguageTypeNCN    LanguageType = "ncn"      // Nauna
	LanguageTypeNCO    LanguageType = "nco"      // Sibe
	LanguageTypeNCP    LanguageType = "ncp"      // Ndaktup
	LanguageTypeNCR    LanguageType = "ncr"      // Ncane
	LanguageTypeNCS    LanguageType = "ncs"      // Nicaraguan Sign Language
	LanguageTypeNCT    LanguageType = "nct"      // Chothe Naga
	LanguageTypeNCU    LanguageType = "ncu"      // Chumburung
	LanguageTypeNCX    LanguageType = "ncx"      // Central Puebla Nahuatl
	LanguageTypeNCZ    LanguageType = "ncz"      // Natchez
	LanguageTypeNDA    LanguageType = "nda"      // Ndasa
	LanguageTypeNDB    LanguageType = "ndb"      // Kenswei Nsei
	LanguageTypeNDC    LanguageType = "ndc"      // Ndau
	LanguageTypeNDD    LanguageType = "ndd"      // Nde-Nsele-Nta
	LanguageTypeNDF    LanguageType = "ndf"      // Nadruvian
	LanguageTypeNDG    LanguageType = "ndg"      // Ndengereko
	LanguageTypeNDH    LanguageType = "ndh"      // Ndali
	LanguageTypeNDI    LanguageType = "ndi"      // Samba Leko
	LanguageTypeNDJ    LanguageType = "ndj"      // Ndamba
	LanguageTypeNDK    LanguageType = "ndk"      // Ndaka
	LanguageTypeNDL    LanguageType = "ndl"      // Ndolo
	LanguageTypeNDM    LanguageType = "ndm"      // Ndam
	LanguageTypeNDN    LanguageType = "ndn"      // Ngundi
	LanguageTypeNDP    LanguageType = "ndp"      // Ndo
	LanguageTypeNDQ    LanguageType = "ndq"      // Ndombe
	LanguageTypeNDR    LanguageType = "ndr"      // Ndoola
	LanguageTypeNDS    LanguageType = "nds"      // Low German and Low Saxon
	LanguageTypeNDT    LanguageType = "ndt"      // Ndunga
	LanguageTypeNDU    LanguageType = "ndu"      // Dugun
	LanguageTypeNDV    LanguageType = "ndv"      // Ndut
	LanguageTypeNDW    LanguageType = "ndw"      // Ndobo
	LanguageTypeNDX    LanguageType = "ndx"      // Nduga
	LanguageTypeNDY    LanguageType = "ndy"      // Lutos
	LanguageTypeNDZ    LanguageType = "ndz"      // Ndogo
	LanguageTypeNEA    LanguageType = "nea"      // Eastern Ngad'a
	LanguageTypeNEB    LanguageType = "neb"      // Toura (Côte d'Ivoire)
	LanguageTypeNEC    LanguageType = "nec"      // Nedebang
	LanguageTypeNED    LanguageType = "ned"      // Nde-Gbite
	LanguageTypeNEE    LanguageType = "nee"      // Nêlêmwa-Nixumwak
	LanguageTypeNEF    LanguageType = "nef"      // Nefamese
	LanguageTypeNEG    LanguageType = "neg"      // Negidal
	LanguageTypeNEH    LanguageType = "neh"      // Nyenkha
	LanguageTypeNEI    LanguageType = "nei"      // Neo-Hittite
	LanguageTypeNEJ    LanguageType = "nej"      // Neko
	LanguageTypeNEK    LanguageType = "nek"      // Neku
	LanguageTypeNEM    LanguageType = "nem"      // Nemi
	LanguageTypeNEN    LanguageType = "nen"      // Nengone
	LanguageTypeNEO    LanguageType = "neo"      // Ná-Meo
	LanguageTypeNEQ    LanguageType = "neq"      // North Central Mixe
	LanguageTypeNER    LanguageType = "ner"      // Yahadian
	LanguageTypeNES    LanguageType = "nes"      // Bhoti Kinnauri
	LanguageTypeNET    LanguageType = "net"      // Nete
	LanguageTypeNEU    LanguageType = "neu"      // Neo
	LanguageTypeNEV    LanguageType = "nev"      // Nyaheun
	LanguageTypeNEW    LanguageType = "new"      // Newari and Nepal Bhasa
	LanguageTypeNEX    LanguageType = "nex"      // Neme
	LanguageTypeNEY    LanguageType = "ney"      // Neyo
	LanguageTypeNEZ    LanguageType = "nez"      // Nez Perce
	LanguageTypeNFA    LanguageType = "nfa"      // Dhao
	LanguageTypeNFD    LanguageType = "nfd"      // Ahwai
	LanguageTypeNFL    LanguageType = "nfl"      // Ayiwo and Äiwoo
	LanguageTypeNFR    LanguageType = "nfr"      // Nafaanra
	LanguageTypeNFU    LanguageType = "nfu"      // Mfumte
	LanguageTypeNGA    LanguageType = "nga"      // Ngbaka
	LanguageTypeNGB    LanguageType = "ngb"      // Northern Ngbandi
	LanguageTypeNGC    LanguageType = "ngc"      // Ngombe (Democratic Republic of Congo)
	LanguageTypeNGD    LanguageType = "ngd"      // Ngando (Central African Republic)
	LanguageTypeNGE    LanguageType = "nge"      // Ngemba
	LanguageTypeNGF    LanguageType = "ngf"      // Trans-New Guinea languages
	LanguageTypeNGG    LanguageType = "ngg"      // Ngbaka Manza
	LanguageTypeNGH    LanguageType = "ngh"      // N/u
	LanguageTypeNGI    LanguageType = "ngi"      // Ngizim
	LanguageTypeNGJ    LanguageType = "ngj"      // Ngie
	LanguageTypeNGK    LanguageType = "ngk"      // Dalabon
	LanguageTypeNGL    LanguageType = "ngl"      // Lomwe
	LanguageTypeNGM    LanguageType = "ngm"      // Ngatik Men's Creole
	LanguageTypeNGN    LanguageType = "ngn"      // Ngwo
	LanguageTypeNGO    LanguageType = "ngo"      // Ngoni
	LanguageTypeNGP    LanguageType = "ngp"      // Ngulu
	LanguageTypeNGQ    LanguageType = "ngq"      // Ngurimi and Ngoreme
	LanguageTypeNGR    LanguageType = "ngr"      // Engdewu
	LanguageTypeNGS    LanguageType = "ngs"      // Gvoko
	LanguageTypeNGT    LanguageType = "ngt"      // Ngeq
	LanguageTypeNGU    LanguageType = "ngu"      // Guerrero Nahuatl
	LanguageTypeNGV    LanguageType = "ngv"      // Nagumi
	LanguageTypeNGW    LanguageType = "ngw"      // Ngwaba
	LanguageTypeNGX    LanguageType = "ngx"      // Nggwahyi
	LanguageTypeNGY    LanguageType = "ngy"      // Tibea
	LanguageTypeNGZ    LanguageType = "ngz"      // Ngungwel
	LanguageTypeNHA    LanguageType = "nha"      // Nhanda
	LanguageTypeNHB    LanguageType = "nhb"      // Beng
	LanguageTypeNHC    LanguageType = "nhc"      // Tabasco Nahuatl
	LanguageTypeNHD    LanguageType = "nhd"      // Chiripá and Ava Guaraní
	LanguageTypeNHE    LanguageType = "nhe"      // Eastern Huasteca Nahuatl
	LanguageTypeNHF    LanguageType = "nhf"      // Nhuwala
	LanguageTypeNHG    LanguageType = "nhg"      // Tetelcingo Nahuatl
	LanguageTypeNHH    LanguageType = "nhh"      // Nahari
	LanguageTypeNHI    LanguageType = "nhi"      // Zacatlán-Ahuacatlán-Tepetzintla Nahuatl
	LanguageTypeNHK    LanguageType = "nhk"      // Isthmus-Cosoleacaque Nahuatl
	LanguageTypeNHM    LanguageType = "nhm"      // Morelos Nahuatl
	LanguageTypeNHN    LanguageType = "nhn"      // Central Nahuatl
	LanguageTypeNHO    LanguageType = "nho"      // Takuu
	LanguageTypeNHP    LanguageType = "nhp"      // Isthmus-Pajapan Nahuatl
	LanguageTypeNHQ    LanguageType = "nhq"      // Huaxcaleca Nahuatl
	LanguageTypeNHR    LanguageType = "nhr"      // Naro
	LanguageTypeNHT    LanguageType = "nht"      // Ometepec Nahuatl
	LanguageTypeNHU    LanguageType = "nhu"      // Noone
	LanguageTypeNHV    LanguageType = "nhv"      // Temascaltepec Nahuatl
	LanguageTypeNHW    LanguageType = "nhw"      // Western Huasteca Nahuatl
	LanguageTypeNHX    LanguageType = "nhx"      // Isthmus-Mecayapan Nahuatl
	LanguageTypeNHY    LanguageType = "nhy"      // Northern Oaxaca Nahuatl
	LanguageTypeNHZ    LanguageType = "nhz"      // Santa María La Alta Nahuatl
	LanguageTypeNIA    LanguageType = "nia"      // Nias
	LanguageTypeNIB    LanguageType = "nib"      // Nakame
	LanguageTypeNIC    LanguageType = "nic"      // Niger-Kordofanian languages
	LanguageTypeNID    LanguageType = "nid"      // Ngandi
	LanguageTypeNIE    LanguageType = "nie"      // Niellim
	LanguageTypeNIF    LanguageType = "nif"      // Nek
	LanguageTypeNIG    LanguageType = "nig"      // Ngalakan
	LanguageTypeNIH    LanguageType = "nih"      // Nyiha (Tanzania)
	LanguageTypeNII    LanguageType = "nii"      // Nii
	LanguageTypeNIJ    LanguageType = "nij"      // Ngaju
	LanguageTypeNIK    LanguageType = "nik"      // Southern Nicobarese
	LanguageTypeNIL    LanguageType = "nil"      // Nila
	LanguageTypeNIM    LanguageType = "nim"      // Nilamba
	LanguageTypeNIN    LanguageType = "nin"      // Ninzo
	LanguageTypeNIO    LanguageType = "nio"      // Nganasan
	LanguageTypeNIQ    LanguageType = "niq"      // Nandi
	LanguageTypeNIR    LanguageType = "nir"      // Nimboran
	LanguageTypeNIS    LanguageType = "nis"      // Nimi
	LanguageTypeNIT    LanguageType = "nit"      // Southeastern Kolami
	LanguageTypeNIU    LanguageType = "niu"      // Niuean
	LanguageTypeNIV    LanguageType = "niv"      // Gilyak
	LanguageTypeNIW    LanguageType = "niw"      // Nimo
	LanguageTypeNIX    LanguageType = "nix"      // Hema
	LanguageTypeNIY    LanguageType = "niy"      // Ngiti
	LanguageTypeNIZ    LanguageType = "niz"      // Ningil
	LanguageTypeNJA    LanguageType = "nja"      // Nzanyi
	LanguageTypeNJB    LanguageType = "njb"      // Nocte Naga
	LanguageTypeNJD    LanguageType = "njd"      // Ndonde Hamba
	LanguageTypeNJH    LanguageType = "njh"      // Lotha Naga
	LanguageTypeNJI    LanguageType = "nji"      // Gudanji
	LanguageTypeNJJ    LanguageType = "njj"      // Njen
	LanguageTypeNJL    LanguageType = "njl"      // Njalgulgule
	LanguageTypeNJM    LanguageType = "njm"      // Angami Naga
	LanguageTypeNJN    LanguageType = "njn"      // Liangmai Naga
	LanguageTypeNJO    LanguageType = "njo"      // Ao Naga
	LanguageTypeNJR    LanguageType = "njr"      // Njerep
	LanguageTypeNJS    LanguageType = "njs"      // Nisa
	LanguageTypeNJT    LanguageType = "njt"      // Ndyuka-Trio Pidgin
	LanguageTypeNJU    LanguageType = "nju"      // Ngadjunmaya
	LanguageTypeNJX    LanguageType = "njx"      // Kunyi
	LanguageTypeNJY    LanguageType = "njy"      // Njyem
	LanguageTypeNJZ    LanguageType = "njz"      // Nyishi
	LanguageTypeNKA    LanguageType = "nka"      // Nkoya
	LanguageTypeNKB    LanguageType = "nkb"      // Khoibu Naga
	LanguageTypeNKC    LanguageType = "nkc"      // Nkongho
	LanguageTypeNKD    LanguageType = "nkd"      // Koireng
	LanguageTypeNKE    LanguageType = "nke"      // Duke
	LanguageTypeNKF    LanguageType = "nkf"      // Inpui Naga
	LanguageTypeNKG    LanguageType = "nkg"      // Nekgini
	LanguageTypeNKH    LanguageType = "nkh"      // Khezha Naga
	LanguageTypeNKI    LanguageType = "nki"      // Thangal Naga
	LanguageTypeNKJ    LanguageType = "nkj"      // Nakai
	LanguageTypeNKK    LanguageType = "nkk"      // Nokuku
	LanguageTypeNKM    LanguageType = "nkm"      // Namat
	LanguageTypeNKN    LanguageType = "nkn"      // Nkangala
	LanguageTypeNKO    LanguageType = "nko"      // Nkonya
	LanguageTypeNKP    LanguageType = "nkp"      // Niuatoputapu
	LanguageTypeNKQ    LanguageType = "nkq"      // Nkami
	LanguageTypeNKR    LanguageType = "nkr"      // Nukuoro
	LanguageTypeNKS    LanguageType = "nks"      // North Asmat
	LanguageTypeNKT    LanguageType = "nkt"      // Nyika (Tanzania)
	LanguageTypeNKU    LanguageType = "nku"      // Bouna Kulango
	LanguageTypeNKV    LanguageType = "nkv"      // Nyika (Malawi and Zambia)
	LanguageTypeNKW    LanguageType = "nkw"      // Nkutu
	LanguageTypeNKX    LanguageType = "nkx"      // Nkoroo
	LanguageTypeNKZ    LanguageType = "nkz"      // Nkari
	LanguageTypeNLA    LanguageType = "nla"      // Ngombale
	LanguageTypeNLC    LanguageType = "nlc"      // Nalca
	LanguageTypeNLE    LanguageType = "nle"      // East Nyala
	LanguageTypeNLG    LanguageType = "nlg"      // Gela
	LanguageTypeNLI    LanguageType = "nli"      // Grangali
	LanguageTypeNLJ    LanguageType = "nlj"      // Nyali
	LanguageTypeNLK    LanguageType = "nlk"      // Ninia Yali
	LanguageTypeNLL    LanguageType = "nll"      // Nihali
	LanguageTypeNLN    LanguageType = "nln"      // Durango Nahuatl
	LanguageTypeNLO    LanguageType = "nlo"      // Ngul
	LanguageTypeNLQ    LanguageType = "nlq"      // Lao Naga
	LanguageTypeNLR    LanguageType = "nlr"      // Ngarla
	LanguageTypeNLU    LanguageType = "nlu"      // Nchumbulu
	LanguageTypeNLV    LanguageType = "nlv"      // Orizaba Nahuatl
	LanguageTypeNLW    LanguageType = "nlw"      // Walangama
	LanguageTypeNLX    LanguageType = "nlx"      // Nahali
	LanguageTypeNLY    LanguageType = "nly"      // Nyamal
	LanguageTypeNLZ    LanguageType = "nlz"      // Nalögo
	LanguageTypeNMA    LanguageType = "nma"      // Maram Naga
	LanguageTypeNMB    LanguageType = "nmb"      // Big Nambas and V'ënen Taut
	LanguageTypeNMC    LanguageType = "nmc"      // Ngam
	LanguageTypeNMD    LanguageType = "nmd"      // Ndumu
	LanguageTypeNME    LanguageType = "nme"      // Mzieme Naga
	LanguageTypeNMF    LanguageType = "nmf"      // Tangkhul Naga (India)
	LanguageTypeNMG    LanguageType = "nmg"      // Kwasio
	LanguageTypeNMH    LanguageType = "nmh"      // Monsang Naga
	LanguageTypeNMI    LanguageType = "nmi"      // Nyam
	LanguageTypeNMJ    LanguageType = "nmj"      // Ngombe (Central African Republic)
	LanguageTypeNMK    LanguageType = "nmk"      // Namakura
	LanguageTypeNML    LanguageType = "nml"      // Ndemli
	LanguageTypeNMM    LanguageType = "nmm"      // Manangba
	LanguageTypeNMN    LanguageType = "nmn"      // !Xóõ
	LanguageTypeNMO    LanguageType = "nmo"      // Moyon Naga
	LanguageTypeNMP    LanguageType = "nmp"      // Nimanbur
	LanguageTypeNMQ    LanguageType = "nmq"      // Nambya
	LanguageTypeNMR    LanguageType = "nmr"      // Nimbari
	LanguageTypeNMS    LanguageType = "nms"      // Letemboi
	LanguageTypeNMT    LanguageType = "nmt"      // Namonuito
	LanguageTypeNMU    LanguageType = "nmu"      // Northeast Maidu
	LanguageTypeNMV    LanguageType = "nmv"      // Ngamini
	LanguageTypeNMW    LanguageType = "nmw"      // Nimoa
	LanguageTypeNMX    LanguageType = "nmx"      // Nama (Papua New Guinea)
	LanguageTypeNMY    LanguageType = "nmy"      // Namuyi
	LanguageTypeNMZ    LanguageType = "nmz"      // Nawdm
	LanguageTypeNNA    LanguageType = "nna"      // Nyangumarta
	LanguageTypeNNB    LanguageType = "nnb"      // Nande
	LanguageTypeNNC    LanguageType = "nnc"      // Nancere
	LanguageTypeNND    LanguageType = "nnd"      // West Ambae
	LanguageTypeNNE    LanguageType = "nne"      // Ngandyera
	LanguageTypeNNF    LanguageType = "nnf"      // Ngaing
	LanguageTypeNNG    LanguageType = "nng"      // Maring Naga
	LanguageTypeNNH    LanguageType = "nnh"      // Ngiemboon
	LanguageTypeNNI    LanguageType = "nni"      // North Nuaulu
	LanguageTypeNNJ    LanguageType = "nnj"      // Nyangatom
	LanguageTypeNNK    LanguageType = "nnk"      // Nankina
	LanguageTypeNNL    LanguageType = "nnl"      // Northern Rengma Naga
	LanguageTypeNNM    LanguageType = "nnm"      // Namia
	LanguageTypeNNN    LanguageType = "nnn"      // Ngete
	LanguageTypeNNP    LanguageType = "nnp"      // Wancho Naga
	LanguageTypeNNQ    LanguageType = "nnq"      // Ngindo
	LanguageTypeNNR    LanguageType = "nnr"      // Narungga
	LanguageTypeNNS    LanguageType = "nns"      // Ningye
	LanguageTypeNNT    LanguageType = "nnt"      // Nanticoke
	LanguageTypeNNU    LanguageType = "nnu"      // Dwang
	LanguageTypeNNV    LanguageType = "nnv"      // Nugunu (Australia)
	LanguageTypeNNW    LanguageType = "nnw"      // Southern Nuni
	LanguageTypeNNX    LanguageType = "nnx"      // Ngong
	LanguageTypeNNY    LanguageType = "nny"      // Nyangga
	LanguageTypeNNZ    LanguageType = "nnz"      // Nda'nda'
	LanguageTypeNOA    LanguageType = "noa"      // Woun Meu
	LanguageTypeNOC    LanguageType = "noc"      // Nuk
	LanguageTypeNOD    LanguageType = "nod"      // Northern Thai
	LanguageTypeNOE    LanguageType = "noe"      // Nimadi
	LanguageTypeNOF    LanguageType = "nof"      // Nomane
	LanguageTypeNOG    LanguageType = "nog"      // Nogai
	LanguageTypeNOH    LanguageType = "noh"      // Nomu
	LanguageTypeNOI    LanguageType = "noi"      // Noiri
	LanguageTypeNOJ    LanguageType = "noj"      // Nonuya
	LanguageTypeNOK    LanguageType = "nok"      // Nooksack
	LanguageTypeNOL    LanguageType = "nol"      // Nomlaki
	LanguageTypeNOM    LanguageType = "nom"      // Nocamán
	LanguageTypeNON    LanguageType = "non"      // Old Norse
	LanguageTypeNOO    LanguageType = "noo"      // Nootka
	LanguageTypeNOP    LanguageType = "nop"      // Numanggang
	LanguageTypeNOQ    LanguageType = "noq"      // Ngongo
	LanguageTypeNOS    LanguageType = "nos"      // Eastern Nisu
	LanguageTypeNOT    LanguageType = "not"      // Nomatsiguenga
	LanguageTypeNOU    LanguageType = "nou"      // Ewage-Notu
	LanguageTypeNOV    LanguageType = "nov"      // Novial
	LanguageTypeNOW    LanguageType = "now"      // Nyambo
	LanguageTypeNOY    LanguageType = "noy"      // Noy
	LanguageTypeNOZ    LanguageType = "noz"      // Nayi
	LanguageTypeNPA    LanguageType = "npa"      // Nar Phu
	LanguageTypeNPB    LanguageType = "npb"      // Nupbikha
	LanguageTypeNPG    LanguageType = "npg"      // Ponyo-Gongwang Naga
	LanguageTypeNPH    LanguageType = "nph"      // Phom Naga
	LanguageTypeNPI    LanguageType = "npi"      // Nepali (individual language)
	LanguageTypeNPL    LanguageType = "npl"      // Southeastern Puebla Nahuatl
	LanguageTypeNPN    LanguageType = "npn"      // Mondropolon
	LanguageTypeNPO    LanguageType = "npo"      // Pochuri Naga
	LanguageTypeNPS    LanguageType = "nps"      // Nipsan
	LanguageTypeNPU    LanguageType = "npu"      // Puimei Naga
	LanguageTypeNPY    LanguageType = "npy"      // Napu
	LanguageTypeNQG    LanguageType = "nqg"      // Southern Nago
	LanguageTypeNQK    LanguageType = "nqk"      // Kura Ede Nago
	LanguageTypeNQM    LanguageType = "nqm"      // Ndom
	LanguageTypeNQN    LanguageType = "nqn"      // Nen
	LanguageTypeNQO    LanguageType = "nqo"      // N'Ko and N’Ko
	LanguageTypeNQQ    LanguageType = "nqq"      // Kyan-Karyaw Naga
	LanguageTypeNQY    LanguageType = "nqy"      // Akyaung Ari Naga
	LanguageTypeNRA    LanguageType = "nra"      // Ngom
	LanguageTypeNRB    LanguageType = "nrb"      // Nara
	LanguageTypeNRC    LanguageType = "nrc"      // Noric
	LanguageTypeNRE    LanguageType = "nre"      // Southern Rengma Naga
	LanguageTypeNRG    LanguageType = "nrg"      // Narango
	LanguageTypeNRI    LanguageType = "nri"      // Chokri Naga
	LanguageTypeNRK    LanguageType = "nrk"      // Ngarla
	LanguageTypeNRL    LanguageType = "nrl"      // Ngarluma
	LanguageTypeNRM    LanguageType = "nrm"      // Narom
	LanguageTypeNRN    LanguageType = "nrn"      // Norn
	LanguageTypeNRP    LanguageType = "nrp"      // North Picene
	LanguageTypeNRR    LanguageType = "nrr"      // Norra and Nora
	LanguageTypeNRT    LanguageType = "nrt"      // Northern Kalapuya
	LanguageTypeNRU    LanguageType = "nru"      // Narua
	LanguageTypeNRX    LanguageType = "nrx"      // Ngurmbur
	LanguageTypeNRZ    LanguageType = "nrz"      // Lala
	LanguageTypeNSA    LanguageType = "nsa"      // Sangtam Naga
	LanguageTypeNSC    LanguageType = "nsc"      // Nshi
	LanguageTypeNSD    LanguageType = "nsd"      // Southern Nisu
	LanguageTypeNSE    LanguageType = "nse"      // Nsenga
	LanguageTypeNSF    LanguageType = "nsf"      // Northwestern Nisu
	LanguageTypeNSG    LanguageType = "nsg"      // Ngasa
	LanguageTypeNSH    LanguageType = "nsh"      // Ngoshie
	LanguageTypeNSI    LanguageType = "nsi"      // Nigerian Sign Language
	LanguageTypeNSK    LanguageType = "nsk"      // Naskapi
	LanguageTypeNSL    LanguageType = "nsl"      // Norwegian Sign Language
	LanguageTypeNSM    LanguageType = "nsm"      // Sumi Naga
	LanguageTypeNSN    LanguageType = "nsn"      // Nehan
	LanguageTypeNSO    LanguageType = "nso"      // Pedi and Northern Sotho and Sepedi
	LanguageTypeNSP    LanguageType = "nsp"      // Nepalese Sign Language
	LanguageTypeNSQ    LanguageType = "nsq"      // Northern Sierra Miwok
	LanguageTypeNSR    LanguageType = "nsr"      // Maritime Sign Language
	LanguageTypeNSS    LanguageType = "nss"      // Nali
	LanguageTypeNST    LanguageType = "nst"      // Tase Naga
	LanguageTypeNSU    LanguageType = "nsu"      // Sierra Negra Nahuatl
	LanguageTypeNSV    LanguageType = "nsv"      // Southwestern Nisu
	LanguageTypeNSW    LanguageType = "nsw"      // Navut
	LanguageTypeNSX    LanguageType = "nsx"      // Nsongo
	LanguageTypeNSY    LanguageType = "nsy"      // Nasal
	LanguageTypeNSZ    LanguageType = "nsz"      // Nisenan
	LanguageTypeNTE    LanguageType = "nte"      // Nathembo
	LanguageTypeNTG    LanguageType = "ntg"      // Ngantangarra
	LanguageTypeNTI    LanguageType = "nti"      // Natioro
	LanguageTypeNTJ    LanguageType = "ntj"      // Ngaanyatjarra
	LanguageTypeNTK    LanguageType = "ntk"      // Ikoma-Nata-Isenye
	LanguageTypeNTM    LanguageType = "ntm"      // Nateni
	LanguageTypeNTO    LanguageType = "nto"      // Ntomba
	LanguageTypeNTP    LanguageType = "ntp"      // Northern Tepehuan
	LanguageTypeNTR    LanguageType = "ntr"      // Delo
	LanguageTypeNTS    LanguageType = "nts"      // Natagaimas
	LanguageTypeNTU    LanguageType = "ntu"      // Natügu
	LanguageTypeNTW    LanguageType = "ntw"      // Nottoway
	LanguageTypeNTX    LanguageType = "ntx"      // Tangkhul Naga (Myanmar)
	LanguageTypeNTY    LanguageType = "nty"      // Mantsi
	LanguageTypeNTZ    LanguageType = "ntz"      // Natanzi
	LanguageTypeNUA    LanguageType = "nua"      // Yuanga
	LanguageTypeNUB    LanguageType = "nub"      // Nubian languages
	LanguageTypeNUC    LanguageType = "nuc"      // Nukuini
	LanguageTypeNUD    LanguageType = "nud"      // Ngala
	LanguageTypeNUE    LanguageType = "nue"      // Ngundu
	LanguageTypeNUF    LanguageType = "nuf"      // Nusu
	LanguageTypeNUG    LanguageType = "nug"      // Nungali
	LanguageTypeNUH    LanguageType = "nuh"      // Ndunda
	LanguageTypeNUI    LanguageType = "nui"      // Ngumbi
	LanguageTypeNUJ    LanguageType = "nuj"      // Nyole
	LanguageTypeNUK    LanguageType = "nuk"      // Nuu-chah-nulth and Nuuchahnulth
	LanguageTypeNUL    LanguageType = "nul"      // Nusa Laut
	LanguageTypeNUM    LanguageType = "num"      // Niuafo'ou
	LanguageTypeNUN    LanguageType = "nun"      // Anong
	LanguageTypeNUO    LanguageType = "nuo"      // Nguôn
	LanguageTypeNUP    LanguageType = "nup"      // Nupe-Nupe-Tako
	LanguageTypeNUQ    LanguageType = "nuq"      // Nukumanu
	LanguageTypeNUR    LanguageType = "nur"      // Nukuria
	LanguageTypeNUS    LanguageType = "nus"      // Nuer
	LanguageTypeNUT    LanguageType = "nut"      // Nung (Viet Nam)
	LanguageTypeNUU    LanguageType = "nuu"      // Ngbundu
	LanguageTypeNUV    LanguageType = "nuv"      // Northern Nuni
	LanguageTypeNUW    LanguageType = "nuw"      // Nguluwan
	LanguageTypeNUX    LanguageType = "nux"      // Mehek
	LanguageTypeNUY    LanguageType = "nuy"      // Nunggubuyu
	LanguageTypeNUZ    LanguageType = "nuz"      // Tlamacazapa Nahuatl
	LanguageTypeNVH    LanguageType = "nvh"      // Nasarian
	LanguageTypeNVM    LanguageType = "nvm"      // Namiae
	LanguageTypeNVO    LanguageType = "nvo"      // Nyokon
	LanguageTypeNWA    LanguageType = "nwa"      // Nawathinehena
	LanguageTypeNWB    LanguageType = "nwb"      // Nyabwa
	LanguageTypeNWC    LanguageType = "nwc"      // Classical Newari and Classical Nepal Bhasa and Old Newari
	LanguageTypeNWE    LanguageType = "nwe"      // Ngwe
	LanguageTypeNWG    LanguageType = "nwg"      // Ngayawung
	LanguageTypeNWI    LanguageType = "nwi"      // Southwest Tanna
	LanguageTypeNWM    LanguageType = "nwm"      // Nyamusa-Molo
	LanguageTypeNWO    LanguageType = "nwo"      // Nauo
	LanguageTypeNWR    LanguageType = "nwr"      // Nawaru
	LanguageTypeNWX    LanguageType = "nwx"      // Middle Newar
	LanguageTypeNWY    LanguageType = "nwy"      // Nottoway-Meherrin
	LanguageTypeNXA    LanguageType = "nxa"      // Nauete
	LanguageTypeNXD    LanguageType = "nxd"      // Ngando (Democratic Republic of Congo)
	LanguageTypeNXE    LanguageType = "nxe"      // Nage
	LanguageTypeNXG    LanguageType = "nxg"      // Ngad'a
	LanguageTypeNXI    LanguageType = "nxi"      // Nindi
	LanguageTypeNXK    LanguageType = "nxk"      // Koki Naga
	LanguageTypeNXL    LanguageType = "nxl"      // South Nuaulu
	LanguageTypeNXM    LanguageType = "nxm"      // Numidian
	LanguageTypeNXN    LanguageType = "nxn"      // Ngawun
	LanguageTypeNXQ    LanguageType = "nxq"      // Naxi
	LanguageTypeNXR    LanguageType = "nxr"      // Ninggerum
	LanguageTypeNXU    LanguageType = "nxu"      // Narau
	LanguageTypeNXX    LanguageType = "nxx"      // Nafri
	LanguageTypeNYB    LanguageType = "nyb"      // Nyangbo
	LanguageTypeNYC    LanguageType = "nyc"      // Nyanga-li
	LanguageTypeNYD    LanguageType = "nyd"      // Nyore and Olunyole
	LanguageTypeNYE    LanguageType = "nye"      // Nyengo
	LanguageTypeNYF    LanguageType = "nyf"      // Giryama and Kigiryama
	LanguageTypeNYG    LanguageType = "nyg"      // Nyindu
	LanguageTypeNYH    LanguageType = "nyh"      // Nyigina
	LanguageTypeNYI    LanguageType = "nyi"      // Ama (Sudan)
	LanguageTypeNYJ    LanguageType = "nyj"      // Nyanga
	LanguageTypeNYK    LanguageType = "nyk"      // Nyaneka
	LanguageTypeNYL    LanguageType = "nyl"      // Nyeu
	LanguageTypeNYM    LanguageType = "nym"      // Nyamwezi
	LanguageTypeNYN    LanguageType = "nyn"      // Nyankole
	LanguageTypeNYO    LanguageType = "nyo"      // Nyoro
	LanguageTypeNYP    LanguageType = "nyp"      // Nyang'i
	LanguageTypeNYQ    LanguageType = "nyq"      // Nayini
	LanguageTypeNYR    LanguageType = "nyr"      // Nyiha (Malawi)
	LanguageTypeNYS    LanguageType = "nys"      // Nyunga
	LanguageTypeNYT    LanguageType = "nyt"      // Nyawaygi
	LanguageTypeNYU    LanguageType = "nyu"      // Nyungwe
	LanguageTypeNYV    LanguageType = "nyv"      // Nyulnyul
	LanguageTypeNYW    LanguageType = "nyw"      // Nyaw
	LanguageTypeNYX    LanguageType = "nyx"      // Nganyaywana
	LanguageTypeNYY    LanguageType = "nyy"      // Nyakyusa-Ngonde
	LanguageTypeNZA    LanguageType = "nza"      // Tigon Mbembe
	LanguageTypeNZB    LanguageType = "nzb"      // Njebi
	LanguageTypeNZI    LanguageType = "nzi"      // Nzima
	LanguageTypeNZK    LanguageType = "nzk"      // Nzakara
	LanguageTypeNZM    LanguageType = "nzm"      // Zeme Naga
	LanguageTypeNZS    LanguageType = "nzs"      // New Zealand Sign Language
	LanguageTypeNZU    LanguageType = "nzu"      // Teke-Nzikou
	LanguageTypeNZY    LanguageType = "nzy"      // Nzakambay
	LanguageTypeNZZ    LanguageType = "nzz"      // Nanga Dama Dogon
	LanguageTypeOAA    LanguageType = "oaa"      // Orok
	LanguageTypeOAC    LanguageType = "oac"      // Oroch
	LanguageTypeOAR    LanguageType = "oar"      // Old Aramaic (up to 700 BCE) and Ancient Aramaic (up to 700 BCE)
	LanguageTypeOAV    LanguageType = "oav"      // Old Avar
	LanguageTypeOBI    LanguageType = "obi"      // Obispeño
	LanguageTypeOBK    LanguageType = "obk"      // Southern Bontok
	LanguageTypeOBL    LanguageType = "obl"      // Oblo
	LanguageTypeOBM    LanguageType = "obm"      // Moabite
	LanguageTypeOBO    LanguageType = "obo"      // Obo Manobo
	LanguageTypeOBR    LanguageType = "obr"      // Old Burmese
	LanguageTypeOBT    LanguageType = "obt"      // Old Breton
	LanguageTypeOBU    LanguageType = "obu"      // Obulom
	LanguageTypeOCA    LanguageType = "oca"      // Ocaina
	LanguageTypeOCH    LanguageType = "och"      // Old Chinese
	LanguageTypeOCO    LanguageType = "oco"      // Old Cornish
	LanguageTypeOCU    LanguageType = "ocu"      // Atzingo Matlatzinca
	LanguageTypeODA    LanguageType = "oda"      // Odut
	LanguageTypeODK    LanguageType = "odk"      // Od
	LanguageTypeODT    LanguageType = "odt"      // Old Dutch
	LanguageTypeODU    LanguageType = "odu"      // Odual
	LanguageTypeOFO    LanguageType = "ofo"      // Ofo
	LanguageTypeOFS    LanguageType = "ofs"      // Old Frisian
	LanguageTypeOFU    LanguageType = "ofu"      // Efutop
	LanguageTypeOGB    LanguageType = "ogb"      // Ogbia
	LanguageTypeOGC    LanguageType = "ogc"      // Ogbah
	LanguageTypeOGE    LanguageType = "oge"      // Old Georgian
	LanguageTypeOGG    LanguageType = "ogg"      // Ogbogolo
	LanguageTypeOGO    LanguageType = "ogo"      // Khana
	LanguageTypeOGU    LanguageType = "ogu"      // Ogbronuagum
	LanguageTypeOHT    LanguageType = "oht"      // Old Hittite
	LanguageTypeOHU    LanguageType = "ohu"      // Old Hungarian
	LanguageTypeOIA    LanguageType = "oia"      // Oirata
	LanguageTypeOIN    LanguageType = "oin"      // Inebu One
	LanguageTypeOJB    LanguageType = "ojb"      // Northwestern Ojibwa
	LanguageTypeOJC    LanguageType = "ojc"      // Central Ojibwa
	LanguageTypeOJG    LanguageType = "ojg"      // Eastern Ojibwa
	LanguageTypeOJP    LanguageType = "ojp"      // Old Japanese
	LanguageTypeOJS    LanguageType = "ojs"      // Severn Ojibwa
	LanguageTypeOJV    LanguageType = "ojv"      // Ontong Java
	LanguageTypeOJW    LanguageType = "ojw"      // Western Ojibwa
	LanguageTypeOKA    LanguageType = "oka"      // Okanagan
	LanguageTypeOKB    LanguageType = "okb"      // Okobo
	LanguageTypeOKD    LanguageType = "okd"      // Okodia
	LanguageTypeOKE    LanguageType = "oke"      // Okpe (Southwestern Edo)
	LanguageTypeOKG    LanguageType = "okg"      // Koko Babangk
	LanguageTypeOKH    LanguageType = "okh"      // Koresh-e Rostam
	LanguageTypeOKI    LanguageType = "oki"      // Okiek
	LanguageTypeOKJ    LanguageType = "okj"      // Oko-Juwoi
	LanguageTypeOKK    LanguageType = "okk"      // Kwamtim One
	LanguageTypeOKL    LanguageType = "okl"      // Old Kentish Sign Language
	LanguageTypeOKM    LanguageType = "okm"      // Middle Korean (10th-16th cent.)
	LanguageTypeOKN    LanguageType = "okn"      // Oki-No-Erabu
	LanguageTypeOKO    LanguageType = "oko"      // Old Korean (3rd-9th cent.)
	LanguageTypeOKR    LanguageType = "okr"      // Kirike
	LanguageTypeOKS    LanguageType = "oks"      // Oko-Eni-Osayen
	LanguageTypeOKU    LanguageType = "oku"      // Oku
	LanguageTypeOKV    LanguageType = "okv"      // Orokaiva
	LanguageTypeOKX    LanguageType = "okx"      // Okpe (Northwestern Edo)
	LanguageTypeOLA    LanguageType = "ola"      // Walungge
	LanguageTypeOLD    LanguageType = "old"      // Mochi
	LanguageTypeOLE    LanguageType = "ole"      // Olekha
	LanguageTypeOLK    LanguageType = "olk"      // Olkol
	LanguageTypeOLM    LanguageType = "olm"      // Oloma
	LanguageTypeOLO    LanguageType = "olo"      // Livvi
	LanguageTypeOLR    LanguageType = "olr"      // Olrat
	LanguageTypeOMA    LanguageType = "oma"      // Omaha-Ponca
	LanguageTypeOMB    LanguageType = "omb"      // East Ambae
	LanguageTypeOMC    LanguageType = "omc"      // Mochica
	LanguageTypeOME    LanguageType = "ome"      // Omejes
	LanguageTypeOMG    LanguageType = "omg"      // Omagua
	LanguageTypeOMI    LanguageType = "omi"      // Omi
	LanguageTypeOMK    LanguageType = "omk"      // Omok
	LanguageTypeOML    LanguageType = "oml"      // Ombo
	LanguageTypeOMN    LanguageType = "omn"      // Minoan
	LanguageTypeOMO    LanguageType = "omo"      // Utarmbung
	LanguageTypeOMP    LanguageType = "omp"      // Old Manipuri
	LanguageTypeOMQ    LanguageType = "omq"      // Oto-Manguean languages
	LanguageTypeOMR    LanguageType = "omr"      // Old Marathi
	LanguageTypeOMT    LanguageType = "omt"      // Omotik
	LanguageTypeOMU    LanguageType = "omu"      // Omurano
	LanguageTypeOMV    LanguageType = "omv"      // Omotic languages
	LanguageTypeOMW    LanguageType = "omw"      // South Tairora
	LanguageTypeOMX    LanguageType = "omx"      // Old Mon
	LanguageTypeONA    LanguageType = "ona"      // Ona
	LanguageTypeONB    LanguageType = "onb"      // Lingao
	LanguageTypeONE    LanguageType = "one"      // Oneida
	LanguageTypeONG    LanguageType = "ong"      // Olo
	LanguageTypeONI    LanguageType = "oni"      // Onin
	LanguageTypeONJ    LanguageType = "onj"      // Onjob
	LanguageTypeONK    LanguageType = "onk"      // Kabore One
	LanguageTypeONN    LanguageType = "onn"      // Onobasulu
	LanguageTypeONO    LanguageType = "ono"      // Onondaga
	LanguageTypeONP    LanguageType = "onp"      // Sartang
	LanguageTypeONR    LanguageType = "onr"      // Northern One
	LanguageTypeONS    LanguageType = "ons"      // Ono
	LanguageTypeONT    LanguageType = "ont"      // Ontenu
	LanguageTypeONU    LanguageType = "onu"      // Unua
	LanguageTypeONW    LanguageType = "onw"      // Old Nubian
	LanguageTypeONX    LanguageType = "onx"      // Onin Based Pidgin
	LanguageTypeOOD    LanguageType = "ood"      // Tohono O'odham
	LanguageTypeOOG    LanguageType = "oog"      // Ong
	LanguageTypeOON    LanguageType = "oon"      // Önge
	LanguageTypeOOR    LanguageType = "oor"      // Oorlams
	LanguageTypeOOS    LanguageType = "oos"      // Old Ossetic
	LanguageTypeOPA    LanguageType = "opa"      // Okpamheri
	LanguageTypeOPK    LanguageType = "opk"      // Kopkaka
	LanguageTypeOPM    LanguageType = "opm"      // Oksapmin
	LanguageTypeOPO    LanguageType = "opo"      // Opao
	LanguageTypeOPT    LanguageType = "opt"      // Opata
	LanguageTypeOPY    LanguageType = "opy"      // Ofayé
	LanguageTypeORA    LanguageType = "ora"      // Oroha
	LanguageTypeORC    LanguageType = "orc"      // Orma
	LanguageTypeORE    LanguageType = "ore"      // Orejón
	LanguageTypeORG    LanguageType = "org"      // Oring
	LanguageTypeORH    LanguageType = "orh"      // Oroqen
	LanguageTypeORN    LanguageType = "orn"      // Orang Kanaq
	LanguageTypeORO    LanguageType = "oro"      // Orokolo
	LanguageTypeORR    LanguageType = "orr"      // Oruma
	LanguageTypeORS    LanguageType = "ors"      // Orang Seletar
	LanguageTypeORT    LanguageType = "ort"      // Adivasi Oriya
	LanguageTypeORU    LanguageType = "oru"      // Ormuri
	LanguageTypeORV    LanguageType = "orv"      // Old Russian
	LanguageTypeORW    LanguageType = "orw"      // Oro Win
	LanguageTypeORX    LanguageType = "orx"      // Oro
	LanguageTypeORY    LanguageType = "ory"      // Oriya (individual language)
	LanguageTypeORZ    LanguageType = "orz"      // Ormu
	LanguageTypeOSA    LanguageType = "osa"      // Osage
	LanguageTypeOSC    LanguageType = "osc"      // Oscan
	LanguageTypeOSI    LanguageType = "osi"      // Osing
	LanguageTypeOSO    LanguageType = "oso"      // Ososo
	LanguageTypeOSP    LanguageType = "osp"      // Old Spanish
	LanguageTypeOST    LanguageType = "ost"      // Osatu
	LanguageTypeOSU    LanguageType = "osu"      // Southern One
	LanguageTypeOSX    LanguageType = "osx"      // Old Saxon
	LanguageTypeOTA    LanguageType = "ota"      // Ottoman Turkish (1500-1928)
	LanguageTypeOTB    LanguageType = "otb"      // Old Tibetan
	LanguageTypeOTD    LanguageType = "otd"      // Ot Danum
	LanguageTypeOTE    LanguageType = "ote"      // Mezquital Otomi
	LanguageTypeOTI    LanguageType = "oti"      // Oti
	LanguageTypeOTK    LanguageType = "otk"      // Old Turkish
	LanguageTypeOTL    LanguageType = "otl"      // Tilapa Otomi
	LanguageTypeOTM    LanguageType = "otm"      // Eastern Highland Otomi
	LanguageTypeOTN    LanguageType = "otn"      // Tenango Otomi
	LanguageTypeOTO    LanguageType = "oto"      // Otomian languages
	LanguageTypeOTQ    LanguageType = "otq"      // Querétaro Otomi
	LanguageTypeOTR    LanguageType = "otr"      // Otoro
	LanguageTypeOTS    LanguageType = "ots"      // Estado de México Otomi
	LanguageTypeOTT    LanguageType = "ott"      // Temoaya Otomi
	LanguageTypeOTU    LanguageType = "otu"      // Otuke
	LanguageTypeOTW    LanguageType = "otw"      // Ottawa
	LanguageTypeOTX    LanguageType = "otx"      // Texcatepec Otomi
	LanguageTypeOTY    LanguageType = "oty"      // Old Tamil
	LanguageTypeOTZ    LanguageType = "otz"      // Ixtenco Otomi
	LanguageTypeOUA    LanguageType = "oua"      // Tagargrent
	LanguageTypeOUB    LanguageType = "oub"      // Glio-Oubi
	LanguageTypeOUE    LanguageType = "oue"      // Oune
	LanguageTypeOUI    LanguageType = "oui"      // Old Uighur
	LanguageTypeOUM    LanguageType = "oum"      // Ouma
	LanguageTypeOUN    LanguageType = "oun"      // !O!ung
	LanguageTypeOWI    LanguageType = "owi"      // Owiniga
	LanguageTypeOWL    LanguageType = "owl"      // Old Welsh
	LanguageTypeOYB    LanguageType = "oyb"      // Oy
	LanguageTypeOYD    LanguageType = "oyd"      // Oyda
	LanguageTypeOYM    LanguageType = "oym"      // Wayampi
	LanguageTypeOYY    LanguageType = "oyy"      // Oya'oya
	LanguageTypeOZM    LanguageType = "ozm"      // Koonzime
	LanguageTypePAA    LanguageType = "paa"      // Papuan languages
	LanguageTypePAB    LanguageType = "pab"      // Parecís
	LanguageTypePAC    LanguageType = "pac"      // Pacoh
	LanguageTypePAD    LanguageType = "pad"      // Paumarí
	LanguageTypePAE    LanguageType = "pae"      // Pagibete
	LanguageTypePAF    LanguageType = "paf"      // Paranawát
	LanguageTypePAG    LanguageType = "pag"      // Pangasinan
	LanguageTypePAH    LanguageType = "pah"      // Tenharim
	LanguageTypePAI    LanguageType = "pai"      // Pe
	LanguageTypePAK    LanguageType = "pak"      // Parakanã
	LanguageTypePAL    LanguageType = "pal"      // Pahlavi
	LanguageTypePAM    LanguageType = "pam"      // Pampanga and Kapampangan
	LanguageTypePAO    LanguageType = "pao"      // Northern Paiute
	LanguageTypePAP    LanguageType = "pap"      // Papiamento
	LanguageTypePAQ    LanguageType = "paq"      // Parya
	LanguageTypePAR    LanguageType = "par"      // Panamint and Timbisha
	LanguageTypePAS    LanguageType = "pas"      // Papasena
	LanguageTypePAT    LanguageType = "pat"      // Papitalai
	LanguageTypePAU    LanguageType = "pau"      // Palauan
	LanguageTypePAV    LanguageType = "pav"      // Pakaásnovos
	LanguageTypePAW    LanguageType = "paw"      // Pawnee
	LanguageTypePAX    LanguageType = "pax"      // Pankararé
	LanguageTypePAY    LanguageType = "pay"      // Pech
	LanguageTypePAZ    LanguageType = "paz"      // Pankararú
	LanguageTypePBB    LanguageType = "pbb"      // Páez
	LanguageTypePBC    LanguageType = "pbc"      // Patamona
	LanguageTypePBE    LanguageType = "pbe"      // Mezontla Popoloca
	LanguageTypePBF    LanguageType = "pbf"      // Coyotepec Popoloca
	LanguageTypePBG    LanguageType = "pbg"      // Paraujano
	LanguageTypePBH    LanguageType = "pbh"      // E'ñapa Woromaipu
	LanguageTypePBI    LanguageType = "pbi"      // Parkwa
	LanguageTypePBL    LanguageType = "pbl"      // Mak (Nigeria)
	LanguageTypePBN    LanguageType = "pbn"      // Kpasam
	LanguageTypePBO    LanguageType = "pbo"      // Papel
	LanguageTypePBP    LanguageType = "pbp"      // Badyara
	LanguageTypePBR    LanguageType = "pbr"      // Pangwa
	LanguageTypePBS    LanguageType = "pbs"      // Central Pame
	LanguageTypePBT    LanguageType = "pbt"      // Southern Pashto
	LanguageTypePBU    LanguageType = "pbu"      // Northern Pashto
	LanguageTypePBV    LanguageType = "pbv"      // Pnar
	LanguageTypePBY    LanguageType = "pby"      // Pyu
	LanguageTypePBZ    LanguageType = "pbz"      // Palu
	LanguageTypePCA    LanguageType = "pca"      // Santa Inés Ahuatempan Popoloca
	LanguageTypePCB    LanguageType = "pcb"      // Pear
	LanguageTypePCC    LanguageType = "pcc"      // Bouyei
	LanguageTypePCD    LanguageType = "pcd"      // Picard
	LanguageTypePCE    LanguageType = "pce"      // Ruching Palaung
	LanguageTypePCF    LanguageType = "pcf"      // Paliyan
	LanguageTypePCG    LanguageType = "pcg"      // Paniya
	LanguageTypePCH    LanguageType = "pch"      // Pardhan
	LanguageTypePCI    LanguageType = "pci"      // Duruwa
	LanguageTypePCJ    LanguageType = "pcj"      // Parenga
	LanguageTypePCK    LanguageType = "pck"      // Paite Chin
	LanguageTypePCL    LanguageType = "pcl"      // Pardhi
	LanguageTypePCM    LanguageType = "pcm"      // Nigerian Pidgin
	LanguageTypePCN    LanguageType = "pcn"      // Piti
	LanguageTypePCP    LanguageType = "pcp"      // Pacahuara
	LanguageTypePCR    LanguageType = "pcr"      // Panang
	LanguageTypePCW    LanguageType = "pcw"      // Pyapun
	LanguageTypePDA    LanguageType = "pda"      // Anam
	LanguageTypePDC    LanguageType = "pdc"      // Pennsylvania German
	LanguageTypePDI    LanguageType = "pdi"      // Pa Di
	LanguageTypePDN    LanguageType = "pdn"      // Podena and Fedan
	LanguageTypePDO    LanguageType = "pdo"      // Padoe
	LanguageTypePDT    LanguageType = "pdt"      // Plautdietsch
	LanguageTypePDU    LanguageType = "pdu"      // Kayan
	LanguageTypePEA    LanguageType = "pea"      // Peranakan Indonesian
	LanguageTypePEB    LanguageType = "peb"      // Eastern Pomo
	LanguageTypePED    LanguageType = "ped"      // Mala (Papua New Guinea)
	LanguageTypePEE    LanguageType = "pee"      // Taje
	LanguageTypePEF    LanguageType = "pef"      // Northeastern Pomo
	LanguageTypePEG    LanguageType = "peg"      // Pengo
	LanguageTypePEH    LanguageType = "peh"      // Bonan
	LanguageTypePEI    LanguageType = "pei"      // Chichimeca-Jonaz
	LanguageTypePEJ    LanguageType = "pej"      // Northern Pomo
	LanguageTypePEK    LanguageType = "pek"      // Penchal
	LanguageTypePEL    LanguageType = "pel"      // Pekal
	LanguageTypePEM    LanguageType = "pem"      // Phende
	LanguageTypePEO    LanguageType = "peo"      // Old Persian (ca. 600-400 B.C.)
	LanguageTypePEP    LanguageType = "pep"      // Kunja
	LanguageTypePEQ    LanguageType = "peq"      // Southern Pomo
	LanguageTypePES    LanguageType = "pes"      // Iranian Persian
	LanguageTypePEV    LanguageType = "pev"      // Pémono
	LanguageTypePEX    LanguageType = "pex"      // Petats
	LanguageTypePEY    LanguageType = "pey"      // Petjo
	LanguageTypePEZ    LanguageType = "pez"      // Eastern Penan
	LanguageTypePFA    LanguageType = "pfa"      // Pááfang
	LanguageTypePFE    LanguageType = "pfe"      // Peere
	LanguageTypePFL    LanguageType = "pfl"      // Pfaelzisch
	LanguageTypePGA    LanguageType = "pga"      // Sudanese Creole Arabic
	LanguageTypePGG    LanguageType = "pgg"      // Pangwali
	LanguageTypePGI    LanguageType = "pgi"      // Pagi
	LanguageTypePGK    LanguageType = "pgk"      // Rerep
	LanguageTypePGL    LanguageType = "pgl"      // Primitive Irish
	LanguageTypePGN    LanguageType = "pgn"      // Paelignian
	LanguageTypePGS    LanguageType = "pgs"      // Pangseng
	LanguageTypePGU    LanguageType = "pgu"      // Pagu
	LanguageTypePGY    LanguageType = "pgy"      // Pongyong
	LanguageTypePHA    LanguageType = "pha"      // Pa-Hng
	LanguageTypePHD    LanguageType = "phd"      // Phudagi
	LanguageTypePHG    LanguageType = "phg"      // Phuong
	LanguageTypePHH    LanguageType = "phh"      // Phukha
	LanguageTypePHI    LanguageType = "phi"      // Philippine languages
	LanguageTypePHK    LanguageType = "phk"      // Phake
	LanguageTypePHL    LanguageType = "phl"      // Phalura and Palula
	LanguageTypePHM    LanguageType = "phm"      // Phimbi
	LanguageTypePHN    LanguageType = "phn"      // Phoenician
	LanguageTypePHO    LanguageType = "pho"      // Phunoi
	LanguageTypePHQ    LanguageType = "phq"      // Phana'
	LanguageTypePHR    LanguageType = "phr"      // Pahari-Potwari
	LanguageTypePHT    LanguageType = "pht"      // Phu Thai
	LanguageTypePHU    LanguageType = "phu"      // Phuan
	LanguageTypePHV    LanguageType = "phv"      // Pahlavani
	LanguageTypePHW    LanguageType = "phw"      // Phangduwali
	LanguageTypePIA    LanguageType = "pia"      // Pima Bajo
	LanguageTypePIB    LanguageType = "pib"      // Yine
	LanguageTypePIC    LanguageType = "pic"      // Pinji
	LanguageTypePID    LanguageType = "pid"      // Piaroa
	LanguageTypePIE    LanguageType = "pie"      // Piro
	LanguageTypePIF    LanguageType = "pif"      // Pingelapese
	LanguageTypePIG    LanguageType = "pig"      // Pisabo
	LanguageTypePIH    LanguageType = "pih"      // Pitcairn-Norfolk
	LanguageTypePII    LanguageType = "pii"      // Pini
	LanguageTypePIJ    LanguageType = "pij"      // Pijao
	LanguageTypePIL    LanguageType = "pil"      // Yom
	LanguageTypePIM    LanguageType = "pim"      // Powhatan
	LanguageTypePIN    LanguageType = "pin"      // Piame
	LanguageTypePIO    LanguageType = "pio"      // Piapoco
	LanguageTypePIP    LanguageType = "pip"      // Pero
	LanguageTypePIR    LanguageType = "pir"      // Piratapuyo
	LanguageTypePIS    LanguageType = "pis"      // Pijin
	LanguageTypePIT    LanguageType = "pit"      // Pitta Pitta
	LanguageTypePIU    LanguageType = "piu"      // Pintupi-Luritja
	LanguageTypePIV    LanguageType = "piv"      // Pileni and Vaeakau-Taumako
	LanguageTypePIW    LanguageType = "piw"      // Pimbwe
	LanguageTypePIX    LanguageType = "pix"      // Piu
	LanguageTypePIY    LanguageType = "piy"      // Piya-Kwonci
	LanguageTypePIZ    LanguageType = "piz"      // Pije
	LanguageTypePJT    LanguageType = "pjt"      // Pitjantjatjara
	LanguageTypePKA    LanguageType = "pka"      // Ardhamāgadhī Prākrit
	LanguageTypePKB    LanguageType = "pkb"      // Pokomo and Kipfokomo
	LanguageTypePKC    LanguageType = "pkc"      // Paekche
	LanguageTypePKG    LanguageType = "pkg"      // Pak-Tong
	LanguageTypePKH    LanguageType = "pkh"      // Pankhu
	LanguageTypePKN    LanguageType = "pkn"      // Pakanha
	LanguageTypePKO    LanguageType = "pko"      // Pökoot
	LanguageTypePKP    LanguageType = "pkp"      // Pukapuka
	LanguageTypePKR    LanguageType = "pkr"      // Attapady Kurumba
	LanguageTypePKS    LanguageType = "pks"      // Pakistan Sign Language
	LanguageTypePKT    LanguageType = "pkt"      // Maleng
	LanguageTypePKU    LanguageType = "pku"      // Paku
	LanguageTypePLA    LanguageType = "pla"      // Miani
	LanguageTypePLB    LanguageType = "plb"      // Polonombauk
	LanguageTypePLC    LanguageType = "plc"      // Central Palawano
	LanguageTypePLD    LanguageType = "pld"      // Polari
	LanguageTypePLE    LanguageType = "ple"      // Palu'e
	LanguageTypePLF    LanguageType = "plf"      // Central Malayo-Polynesian languages
	LanguageTypePLG    LanguageType = "plg"      // Pilagá
	LanguageTypePLH    LanguageType = "plh"      // Paulohi
	LanguageTypePLJ    LanguageType = "plj"      // Polci
	LanguageTypePLK    LanguageType = "plk"      // Kohistani Shina
	LanguageTypePLL    LanguageType = "pll"      // Shwe Palaung
	LanguageTypePLN    LanguageType = "pln"      // Palenquero
	LanguageTypePLO    LanguageType = "plo"      // Oluta Popoluca
	LanguageTypePLP    LanguageType = "plp"      // Palpa
	LanguageTypePLQ    LanguageType = "plq"      // Palaic
	LanguageTypePLR    LanguageType = "plr"      // Palaka Senoufo
	LanguageTypePLS    LanguageType = "pls"      // San Marcos Tlalcoyalco Popoloca
	LanguageTypePLT    LanguageType = "plt"      // Plateau Malagasy
	LanguageTypePLU    LanguageType = "plu"      // Palikúr
	LanguageTypePLV    LanguageType = "plv"      // Southwest Palawano
	LanguageTypePLW    LanguageType = "plw"      // Brooke's Point Palawano
	LanguageTypePLY    LanguageType = "ply"      // Bolyu
	LanguageTypePLZ    LanguageType = "plz"      // Paluan
	LanguageTypePMA    LanguageType = "pma"      // Paama
	LanguageTypePMB    LanguageType = "pmb"      // Pambia
	LanguageTypePMC    LanguageType = "pmc"      // Palumata
	LanguageTypePMD    LanguageType = "pmd"      // Pallanganmiddang
	LanguageTypePME    LanguageType = "pme"      // Pwaamei
	LanguageTypePMF    LanguageType = "pmf"      // Pamona
	LanguageTypePMH    LanguageType = "pmh"      // Māhārāṣṭri Prākrit
	LanguageTypePMI    LanguageType = "pmi"      // Northern Pumi
	LanguageTypePMJ    LanguageType = "pmj"      // Southern Pumi
	LanguageTypePMK    LanguageType = "pmk"      // Pamlico
	LanguageTypePML    LanguageType = "pml"      // Lingua Franca
	LanguageTypePMM    LanguageType = "pmm"      // Pomo
	LanguageTypePMN    LanguageType = "pmn"      // Pam
	LanguageTypePMO    LanguageType = "pmo"      // Pom
	LanguageTypePMQ    LanguageType = "pmq"      // Northern Pame
	LanguageTypePMR    LanguageType = "pmr"      // Paynamar
	LanguageTypePMS    LanguageType = "pms"      // Piemontese
	LanguageTypePMT    LanguageType = "pmt"      // Tuamotuan
	LanguageTypePMU    LanguageType = "pmu"      // Mirpur Panjabi
	LanguageTypePMW    LanguageType = "pmw"      // Plains Miwok
	LanguageTypePMX    LanguageType = "pmx"      // Poumei Naga
	LanguageTypePMY    LanguageType = "pmy"      // Papuan Malay
	LanguageTypePMZ    LanguageType = "pmz"      // Southern Pame
	LanguageTypePNA    LanguageType = "pna"      // Punan Bah-Biau
	LanguageTypePNB    LanguageType = "pnb"      // Western Panjabi
	LanguageTypePNC    LanguageType = "pnc"      // Pannei
	LanguageTypePNE    LanguageType = "pne"      // Western Penan
	LanguageTypePNG    LanguageType = "png"      // Pongu
	LanguageTypePNH    LanguageType = "pnh"      // Penrhyn
	LanguageTypePNI    LanguageType = "pni"      // Aoheng
	LanguageTypePNJ    LanguageType = "pnj"      // Pinjarup
	LanguageTypePNK    LanguageType = "pnk"      // Paunaka
	LanguageTypePNL    LanguageType = "pnl"      // Paleni
	LanguageTypePNM    LanguageType = "pnm"      // Punan Batu 1
	LanguageTypePNN    LanguageType = "pnn"      // Pinai-Hagahai
	LanguageTypePNO    LanguageType = "pno"      // Panobo
	LanguageTypePNP    LanguageType = "pnp"      // Pancana
	LanguageTypePNQ    LanguageType = "pnq"      // Pana (Burkina Faso)
	LanguageTypePNR    LanguageType = "pnr"      // Panim
	LanguageTypePNS    LanguageType = "pns"      // Ponosakan
	LanguageTypePNT    LanguageType = "pnt"      // Pontic
	LanguageTypePNU    LanguageType = "pnu"      // Jiongnai Bunu
	LanguageTypePNV    LanguageType = "pnv"      // Pinigura
	LanguageTypePNW    LanguageType = "pnw"      // Panytyima
	LanguageTypePNX    LanguageType = "pnx"      // Phong-Kniang
	LanguageTypePNY    LanguageType = "pny"      // Pinyin
	LanguageTypePNZ    LanguageType = "pnz"      // Pana (Central African Republic)
	LanguageTypePOC    LanguageType = "poc"      // Poqomam
	LanguageTypePOD    LanguageType = "pod"      // Ponares
	LanguageTypePOE    LanguageType = "poe"      // San Juan Atzingo Popoloca
	LanguageTypePOF    LanguageType = "pof"      // Poke
	LanguageTypePOG    LanguageType = "pog"      // Potiguára
	LanguageTypePOH    LanguageType = "poh"      // Poqomchi'
	LanguageTypePOI    LanguageType = "poi"      // Highland Popoluca
	LanguageTypePOK    LanguageType = "pok"      // Pokangá
	LanguageTypePOM    LanguageType = "pom"      // Southeastern Pomo
	LanguageTypePON    LanguageType = "pon"      // Pohnpeian
	LanguageTypePOO    LanguageType = "poo"      // Central Pomo
	LanguageTypePOP    LanguageType = "pop"      // Pwapwâ
	LanguageTypePOQ    LanguageType = "poq"      // Texistepec Popoluca
	LanguageTypePOS    LanguageType = "pos"      // Sayula Popoluca
	LanguageTypePOT    LanguageType = "pot"      // Potawatomi
	LanguageTypePOV    LanguageType = "pov"      // Upper Guinea Crioulo
	LanguageTypePOW    LanguageType = "pow"      // San Felipe Otlaltepec Popoloca
	LanguageTypePOX    LanguageType = "pox"      // Polabian
	LanguageTypePOY    LanguageType = "poy"      // Pogolo
	LanguageTypePOZ    LanguageType = "poz"      // Malayo-Polynesian languages
	LanguageTypePPA    LanguageType = "ppa"      // Pao
	LanguageTypePPE    LanguageType = "ppe"      // Papi
	LanguageTypePPI    LanguageType = "ppi"      // Paipai
	LanguageTypePPK    LanguageType = "ppk"      // Uma
	LanguageTypePPL    LanguageType = "ppl"      // Pipil and Nicarao
	LanguageTypePPM    LanguageType = "ppm"      // Papuma
	LanguageTypePPN    LanguageType = "ppn"      // Papapana
	LanguageTypePPO    LanguageType = "ppo"      // Folopa
	LanguageTypePPP    LanguageType = "ppp"      // Pelende
	LanguageTypePPQ    LanguageType = "ppq"      // Pei
	LanguageTypePPR    LanguageType = "ppr"      // Piru
	LanguageTypePPS    LanguageType = "pps"      // San Luís Temalacayuca Popoloca
	LanguageTypePPT    LanguageType = "ppt"      // Pare
	LanguageTypePPU    LanguageType = "ppu"      // Papora
	LanguageTypePQA    LanguageType = "pqa"      // Pa'a
	LanguageTypePQE    LanguageType = "pqe"      // Eastern Malayo-Polynesian languages
	LanguageTypePQM    LanguageType = "pqm"      // Malecite-Passamaquoddy
	LanguageTypePQW    LanguageType = "pqw"      // Western Malayo-Polynesian languages
	LanguageTypePRA    LanguageType = "pra"      // Prakrit languages
	LanguageTypePRB    LanguageType = "prb"      // Lua'
	LanguageTypePRC    LanguageType = "prc"      // Parachi
	LanguageTypePRD    LanguageType = "prd"      // Parsi-Dari
	LanguageTypePRE    LanguageType = "pre"      // Principense
	LanguageTypePRF    LanguageType = "prf"      // Paranan
	LanguageTypePRG    LanguageType = "prg"      // Prussian
	LanguageTypePRH    LanguageType = "prh"      // Porohanon
	LanguageTypePRI    LanguageType = "pri"      // Paicî
	LanguageTypePRK    LanguageType = "prk"      // Parauk
	LanguageTypePRL    LanguageType = "prl"      // Peruvian Sign Language
	LanguageTypePRM    LanguageType = "prm"      // Kibiri
	LanguageTypePRN    LanguageType = "prn"      // Prasuni
	LanguageTypePRO    LanguageType = "pro"      // Old Provençal (to 1500) and Old Occitan (to 1500)
	LanguageTypePRP    LanguageType = "prp"      // Parsi
	LanguageTypePRQ    LanguageType = "prq"      // Ashéninka Perené
	LanguageTypePRR    LanguageType = "prr"      // Puri
	LanguageTypePRS    LanguageType = "prs"      // Dari and Afghan Persian
	LanguageTypePRT    LanguageType = "prt"      // Phai
	LanguageTypePRU    LanguageType = "pru"      // Puragi
	LanguageTypePRW    LanguageType = "prw"      // Parawen
	LanguageTypePRX    LanguageType = "prx"      // Purik
	LanguageTypePRY    LanguageType = "pry"      // Pray 3
	LanguageTypePRZ    LanguageType = "prz"      // Providencia Sign Language
	LanguageTypePSA    LanguageType = "psa"      // Asue Awyu
	LanguageTypePSC    LanguageType = "psc"      // Persian Sign Language
	LanguageTypePSD    LanguageType = "psd"      // Plains Indian Sign Language
	LanguageTypePSE    LanguageType = "pse"      // Central Malay
	LanguageTypePSG    LanguageType = "psg"      // Penang Sign Language
	LanguageTypePSH    LanguageType = "psh"      // Southwest Pashayi
	LanguageTypePSI    LanguageType = "psi"      // Southeast Pashayi
	LanguageTypePSL    LanguageType = "psl"      // Puerto Rican Sign Language
	LanguageTypePSM    LanguageType = "psm"      // Pauserna
	LanguageTypePSN    LanguageType = "psn"      // Panasuan
	LanguageTypePSO    LanguageType = "pso"      // Polish Sign Language
	LanguageTypePSP    LanguageType = "psp"      // Philippine Sign Language
	LanguageTypePSQ    LanguageType = "psq"      // Pasi
	LanguageTypePSR    LanguageType = "psr"      // Portuguese Sign Language
	LanguageTypePSS    LanguageType = "pss"      // Kaulong
	LanguageTypePST    LanguageType = "pst"      // Central Pashto
	LanguageTypePSU    LanguageType = "psu"      // Sauraseni Prākrit
	LanguageTypePSW    LanguageType = "psw"      // Port Sandwich
	LanguageTypePSY    LanguageType = "psy"      // Piscataway
	LanguageTypePTA    LanguageType = "pta"      // Pai Tavytera
	LanguageTypePTH    LanguageType = "pth"      // Pataxó Hã-Ha-Hãe
	LanguageTypePTI    LanguageType = "pti"      // Pintiini
	LanguageTypePTN    LanguageType = "ptn"      // Patani
	LanguageTypePTO    LanguageType = "pto"      // Zo'é
	LanguageTypePTP    LanguageType = "ptp"      // Patep
	LanguageTypePTR    LanguageType = "ptr"      // Piamatsina
	LanguageTypePTT    LanguageType = "ptt"      // Enrekang
	LanguageTypePTU    LanguageType = "ptu"      // Bambam
	LanguageTypePTV    LanguageType = "ptv"      // Port Vato
	LanguageTypePTW    LanguageType = "ptw"      // Pentlatch
	LanguageTypePTY    LanguageType = "pty"      // Pathiya
	LanguageTypePUA    LanguageType = "pua"      // Western Highland Purepecha
	LanguageTypePUB    LanguageType = "pub"      // Purum
	LanguageTypePUC    LanguageType = "puc"      // Punan Merap
	LanguageTypePUD    LanguageType = "pud"      // Punan Aput
	LanguageTypePUE    LanguageType = "pue"      // Puelche
	LanguageTypePUF    LanguageType = "puf"      // Punan Merah
	LanguageTypePUG    LanguageType = "pug"      // Phuie
	LanguageTypePUI    LanguageType = "pui"      // Puinave
	LanguageTypePUJ    LanguageType = "puj"      // Punan Tubu
	LanguageTypePUK    LanguageType = "puk"      // Pu Ko
	LanguageTypePUM    LanguageType = "pum"      // Puma
	LanguageTypePUO    LanguageType = "puo"      // Puoc
	LanguageTypePUP    LanguageType = "pup"      // Pulabu
	LanguageTypePUQ    LanguageType = "puq"      // Puquina
	LanguageTypePUR    LanguageType = "pur"      // Puruborá
	LanguageTypePUT    LanguageType = "put"      // Putoh
	LanguageTypePUU    LanguageType = "puu"      // Punu
	LanguageTypePUW    LanguageType = "puw"      // Puluwatese
	LanguageTypePUX    LanguageType = "pux"      // Puare
	LanguageTypePUY    LanguageType = "puy"      // Purisimeño
	LanguageTypePUZ    LanguageType = "puz"      // Purum Naga
	LanguageTypePWA    LanguageType = "pwa"      // Pawaia
	LanguageTypePWB    LanguageType = "pwb"      // Panawa
	LanguageTypePWG    LanguageType = "pwg"      // Gapapaiwa
	LanguageTypePWI    LanguageType = "pwi"      // Patwin
	LanguageTypePWM    LanguageType = "pwm"      // Molbog
	LanguageTypePWN    LanguageType = "pwn"      // Paiwan
	LanguageTypePWO    LanguageType = "pwo"      // Pwo Western Karen
	LanguageTypePWR    LanguageType = "pwr"      // Powari
	LanguageTypePWW    LanguageType = "pww"      // Pwo Northern Karen
	LanguageTypePXM    LanguageType = "pxm"      // Quetzaltepec Mixe
	LanguageTypePYE    LanguageType = "pye"      // Pye Krumen
	LanguageTypePYM    LanguageType = "pym"      // Fyam
	LanguageTypePYN    LanguageType = "pyn"      // Poyanáwa
	LanguageTypePYS    LanguageType = "pys"      // Paraguayan Sign Language and Lengua de Señas del Paraguay
	LanguageTypePYU    LanguageType = "pyu"      // Puyuma
	LanguageTypePYX    LanguageType = "pyx"      // Pyu (Myanmar)
	LanguageTypePYY    LanguageType = "pyy"      // Pyen
	LanguageTypePZN    LanguageType = "pzn"      // Para Naga
	LanguageTypeQAAQTZ LanguageType = "qaa..qtz" // Private use
	LanguageTypeQUA    LanguageType = "qua"      // Quapaw
	LanguageTypeQUB    LanguageType = "qub"      // Huallaga Huánuco Quechua
	LanguageTypeQUC    LanguageType = "quc"      // K'iche' and Quiché
	LanguageTypeQUD    LanguageType = "qud"      // Calderón Highland Quichua
	LanguageTypeQUF    LanguageType = "quf"      // Lambayeque Quechua
	LanguageTypeQUG    LanguageType = "qug"      // Chimborazo Highland Quichua
	LanguageTypeQUH    LanguageType = "quh"      // South Bolivian Quechua
	LanguageTypeQUI    LanguageType = "qui"      // Quileute
	LanguageTypeQUK    LanguageType = "quk"      // Chachapoyas Quechua
	LanguageTypeQUL    LanguageType = "qul"      // North Bolivian Quechua
	LanguageTypeQUM    LanguageType = "qum"      // Sipacapense
	LanguageTypeQUN    LanguageType = "qun"      // Quinault
	LanguageTypeQUP    LanguageType = "qup"      // Southern Pastaza Quechua
	LanguageTypeQUQ    LanguageType = "quq"      // Quinqui
	LanguageTypeQUR    LanguageType = "qur"      // Yanahuanca Pasco Quechua
	LanguageTypeQUS    LanguageType = "qus"      // Santiago del Estero Quichua
	LanguageTypeQUV    LanguageType = "quv"      // Sacapulteco
	LanguageTypeQUW    LanguageType = "quw"      // Tena Lowland Quichua
	LanguageTypeQUX    LanguageType = "qux"      // Yauyos Quechua
	LanguageTypeQUY    LanguageType = "quy"      // Ayacucho Quechua
	LanguageTypeQUZ    LanguageType = "quz"      // Cusco Quechua
	LanguageTypeQVA    LanguageType = "qva"      // Ambo-Pasco Quechua
	LanguageTypeQVC    LanguageType = "qvc"      // Cajamarca Quechua
	LanguageTypeQVE    LanguageType = "qve"      // Eastern Apurímac Quechua
	LanguageTypeQVH    LanguageType = "qvh"      // Huamalíes-Dos de Mayo Huánuco Quechua
	LanguageTypeQVI    LanguageType = "qvi"      // Imbabura Highland Quichua
	LanguageTypeQVJ    LanguageType = "qvj"      // Loja Highland Quichua
	LanguageTypeQVL    LanguageType = "qvl"      // Cajatambo North Lima Quechua
	LanguageTypeQVM    LanguageType = "qvm"      // Margos-Yarowilca-Lauricocha Quechua
	LanguageTypeQVN    LanguageType = "qvn"      // North Junín Quechua
	LanguageTypeQVO    LanguageType = "qvo"      // Napo Lowland Quechua
	LanguageTypeQVP    LanguageType = "qvp"      // Pacaraos Quechua
	LanguageTypeQVS    LanguageType = "qvs"      // San Martín Quechua
	LanguageTypeQVW    LanguageType = "qvw"      // Huaylla Wanca Quechua
	LanguageTypeQVY    LanguageType = "qvy"      // Queyu
	LanguageTypeQVZ    LanguageType = "qvz"      // Northern Pastaza Quichua
	LanguageTypeQWA    LanguageType = "qwa"      // Corongo Ancash Quechua
	LanguageTypeQWC    LanguageType = "qwc"      // Classical Quechua
	LanguageTypeQWE    LanguageType = "qwe"      // Quechuan (family)
	LanguageTypeQWH    LanguageType = "qwh"      // Huaylas Ancash Quechua
	LanguageTypeQWM    LanguageType = "qwm"      // Kuman (Russia)
	LanguageTypeQWS    LanguageType = "qws"      // Sihuas Ancash Quechua
	LanguageTypeQWT    LanguageType = "qwt"      // Kwalhioqua-Tlatskanai
	LanguageTypeQXA    LanguageType = "qxa"      // Chiquián Ancash Quechua
	LanguageTypeQXC    LanguageType = "qxc"      // Chincha Quechua
	LanguageTypeQXH    LanguageType = "qxh"      // Panao Huánuco Quechua
	LanguageTypeQXL    LanguageType = "qxl"      // Salasaca Highland Quichua
	LanguageTypeQXN    LanguageType = "qxn"      // Northern Conchucos Ancash Quechua
	LanguageTypeQXO    LanguageType = "qxo"      // Southern Conchucos Ancash Quechua
	LanguageTypeQXP    LanguageType = "qxp"      // Puno Quechua
	LanguageTypeQXQ    LanguageType = "qxq"      // Qashqa'i
	LanguageTypeQXR    LanguageType = "qxr"      // Cañar Highland Quichua
	LanguageTypeQXS    LanguageType = "qxs"      // Southern Qiang
	LanguageTypeQXT    LanguageType = "qxt"      // Santa Ana de Tusi Pasco Quechua
	LanguageTypeQXU    LanguageType = "qxu"      // Arequipa-La Unión Quechua
	LanguageTypeQXW    LanguageType = "qxw"      // Jauja Wanca Quechua
	LanguageTypeQYA    LanguageType = "qya"      // Quenya
	LanguageTypeQYP    LanguageType = "qyp"      // Quiripi
	LanguageTypeRAA    LanguageType = "raa"      // Dungmali
	LanguageTypeRAB    LanguageType = "rab"      // Camling
	LanguageTypeRAC    LanguageType = "rac"      // Rasawa
	LanguageTypeRAD    LanguageType = "rad"      // Rade
	LanguageTypeRAF    LanguageType = "raf"      // Western Meohang
	LanguageTypeRAG    LanguageType = "rag"      // Logooli and Lulogooli
	LanguageTypeRAH    LanguageType = "rah"      // Rabha
	LanguageTypeRAI    LanguageType = "rai"      // Ramoaaina
	LanguageTypeRAJ    LanguageType = "raj"      // Rajasthani
	LanguageTypeRAK    LanguageType = "rak"      // Tulu-Bohuai
	LanguageTypeRAL    LanguageType = "ral"      // Ralte
	LanguageTypeRAM    LanguageType = "ram"      // Canela
	LanguageTypeRAN    LanguageType = "ran"      // Riantana
	LanguageTypeRAO    LanguageType = "rao"      // Rao
	LanguageTypeRAP    LanguageType = "rap"      // Rapanui
	LanguageTypeRAQ    LanguageType = "raq"      // Saam
	LanguageTypeRAR    LanguageType = "rar"      // Rarotongan and Cook Islands Maori
	LanguageTypeRAS    LanguageType = "ras"      // Tegali
	LanguageTypeRAT    LanguageType = "rat"      // Razajerdi
	LanguageTypeRAU    LanguageType = "rau"      // Raute
	LanguageTypeRAV    LanguageType = "rav"      // Sampang
	LanguageTypeRAW    LanguageType = "raw"      // Rawang
	LanguageTypeRAX    LanguageType = "rax"      // Rang
	LanguageTypeRAY    LanguageType = "ray"      // Rapa
	LanguageTypeRAZ    LanguageType = "raz"      // Rahambuu
	LanguageTypeRBB    LanguageType = "rbb"      // Rumai Palaung
	LanguageTypeRBK    LanguageType = "rbk"      // Northern Bontok
	LanguageTypeRBL    LanguageType = "rbl"      // Miraya Bikol
	LanguageTypeRBP    LanguageType = "rbp"      // Barababaraba
	LanguageTypeRCF    LanguageType = "rcf"      // Réunion Creole French
	LanguageTypeRDB    LanguageType = "rdb"      // Rudbari
	LanguageTypeREA    LanguageType = "rea"      // Rerau
	LanguageTypeREB    LanguageType = "reb"      // Rembong
	LanguageTypeREE    LanguageType = "ree"      // Rejang Kayan
	LanguageTypeREG    LanguageType = "reg"      // Kara (Tanzania)
	LanguageTypeREI    LanguageType = "rei"      // Reli
	LanguageTypeREJ    LanguageType = "rej"      // Rejang
	LanguageTypeREL    LanguageType = "rel"      // Rendille
	LanguageTypeREM    LanguageType = "rem"      // Remo
	LanguageTypeREN    LanguageType = "ren"      // Rengao
	LanguageTypeRER    LanguageType = "rer"      // Rer Bare
	LanguageTypeRES    LanguageType = "res"      // Reshe
	LanguageTypeRET    LanguageType = "ret"      // Retta
	LanguageTypeREY    LanguageType = "rey"      // Reyesano
	LanguageTypeRGA    LanguageType = "rga"      // Roria
	LanguageTypeRGE    LanguageType = "rge"      // Romano-Greek
	LanguageTypeRGK    LanguageType = "rgk"      // Rangkas
	LanguageTypeRGN    LanguageType = "rgn"      // Romagnol
	LanguageTypeRGR    LanguageType = "rgr"      // Resígaro
	LanguageTypeRGS    LanguageType = "rgs"      // Southern Roglai
	LanguageTypeRGU    LanguageType = "rgu"      // Ringgou
	LanguageTypeRHG    LanguageType = "rhg"      // Rohingya
	LanguageTypeRHP    LanguageType = "rhp"      // Yahang
	LanguageTypeRIA    LanguageType = "ria"      // Riang (India)
	LanguageTypeRIE    LanguageType = "rie"      // Rien
	LanguageTypeRIF    LanguageType = "rif"      // Tarifit
	LanguageTypeRIL    LanguageType = "ril"      // Riang (Myanmar)
	LanguageTypeRIM    LanguageType = "rim"      // Nyaturu
	LanguageTypeRIN    LanguageType = "rin"      // Nungu
	LanguageTypeRIR    LanguageType = "rir"      // Ribun
	LanguageTypeRIT    LanguageType = "rit"      // Ritarungo
	LanguageTypeRIU    LanguageType = "riu"      // Riung
	LanguageTypeRJG    LanguageType = "rjg"      // Rajong
	LanguageTypeRJI    LanguageType = "rji"      // Raji
	LanguageTypeRJS    LanguageType = "rjs"      // Rajbanshi
	LanguageTypeRKA    LanguageType = "rka"      // Kraol
	LanguageTypeRKB    LanguageType = "rkb"      // Rikbaktsa
	LanguageTypeRKH    LanguageType = "rkh"      // Rakahanga-Manihiki
	LanguageTypeRKI    LanguageType = "rki"      // Rakhine
	LanguageTypeRKM    LanguageType = "rkm"      // Marka
	LanguageTypeRKT    LanguageType = "rkt"      // Rangpuri and Kamta
	LanguageTypeRKW    LanguageType = "rkw"      // Arakwal
	LanguageTypeRMA    LanguageType = "rma"      // Rama
	LanguageTypeRMB    LanguageType = "rmb"      // Rembarunga
	LanguageTypeRMC    LanguageType = "rmc"      // Carpathian Romani
	LanguageTypeRMD    LanguageType = "rmd"      // Traveller Danish
	LanguageTypeRME    LanguageType = "rme"      // Angloromani
	LanguageTypeRMF    LanguageType = "rmf"      // Kalo Finnish Romani
	LanguageTypeRMG    LanguageType = "rmg"      // Traveller Norwegian
	LanguageTypeRMH    LanguageType = "rmh"      // Murkim
	LanguageTypeRMI    LanguageType = "rmi"      // Lomavren
	LanguageTypeRMK    LanguageType = "rmk"      // Romkun
	LanguageTypeRML    LanguageType = "rml"      // Baltic Romani
	LanguageTypeRMM    LanguageType = "rmm"      // Roma
	LanguageTypeRMN    LanguageType = "rmn"      // Balkan Romani
	LanguageTypeRMO    LanguageType = "rmo"      // Sinte Romani
	LanguageTypeRMP    LanguageType = "rmp"      // Rempi
	LanguageTypeRMQ    LanguageType = "rmq"      // Caló
	LanguageTypeRMR    LanguageType = "rmr"      // Caló
	LanguageTypeRMS    LanguageType = "rms"      // Romanian Sign Language
	LanguageTypeRMT    LanguageType = "rmt"      // Domari
	LanguageTypeRMU    LanguageType = "rmu"      // Tavringer Romani
	LanguageTypeRMV    LanguageType = "rmv"      // Romanova
	LanguageTypeRMW    LanguageType = "rmw"      // Welsh Romani
	LanguageTypeRMX    LanguageType = "rmx"      // Romam
	LanguageTypeRMY    LanguageType = "rmy"      // Vlax Romani
	LanguageTypeRMZ    LanguageType = "rmz"      // Marma
	LanguageTypeRNA    LanguageType = "rna"      // Runa
	LanguageTypeRND    LanguageType = "rnd"      // Ruund
	LanguageTypeRNG    LanguageType = "rng"      // Ronga
	LanguageTypeRNL    LanguageType = "rnl"      // Ranglong
	LanguageTypeRNN    LanguageType = "rnn"      // Roon
	LanguageTypeRNP    LanguageType = "rnp"      // Rongpo
	LanguageTypeRNR    LanguageType = "rnr"      // Nari Nari
	LanguageTypeRNW    LanguageType = "rnw"      // Rungwa
	LanguageTypeROA    LanguageType = "roa"      // Romance languages
	LanguageTypeROB    LanguageType = "rob"      // Tae'
	LanguageTypeROC    LanguageType = "roc"      // Cacgia Roglai
	LanguageTypeROD    LanguageType = "rod"      // Rogo
	LanguageTypeROE    LanguageType = "roe"      // Ronji
	LanguageTypeROF    LanguageType = "rof"      // Rombo
	LanguageTypeROG    LanguageType = "rog"      // Northern Roglai
	LanguageTypeROL    LanguageType = "rol"      // Romblomanon
	LanguageTypeROM    LanguageType = "rom"      // Romany
	LanguageTypeROO    LanguageType = "roo"      // Rotokas
	LanguageTypeROP    LanguageType = "rop"      // Kriol
	LanguageTypeROR    LanguageType = "ror"      // Rongga
	LanguageTypeROU    LanguageType = "rou"      // Runga
	LanguageTypeROW    LanguageType = "row"      // Dela-Oenale
	LanguageTypeRPN    LanguageType = "rpn"      // Repanbitip
	LanguageTypeRPT    LanguageType = "rpt"      // Rapting
	LanguageTypeRRI    LanguageType = "rri"      // Ririo
	LanguageTypeRRO    LanguageType = "rro"      // Waima
	LanguageTypeRRT    LanguageType = "rrt"      // Arritinngithigh
	LanguageTypeRSB    LanguageType = "rsb"      // Romano-Serbian
	LanguageTypeRSI    LanguageType = "rsi"      // Rennellese Sign Language
	LanguageTypeRSL    LanguageType = "rsl"      // Russian Sign Language
	LanguageTypeRTC    LanguageType = "rtc"      // Rungtu Chin
	LanguageTypeRTH    LanguageType = "rth"      // Ratahan
	LanguageTypeRTM    LanguageType = "rtm"      // Rotuman
	LanguageTypeRTW    LanguageType = "rtw"      // Rathawi
	LanguageTypeRUB    LanguageType = "rub"      // Gungu
	LanguageTypeRUC    LanguageType = "ruc"      // Ruuli
	LanguageTypeRUE    LanguageType = "rue"      // Rusyn
	LanguageTypeRUF    LanguageType = "ruf"      // Luguru
	LanguageTypeRUG    LanguageType = "rug"      // Roviana
	LanguageTypeRUH    LanguageType = "ruh"      // Ruga
	LanguageTypeRUI    LanguageType = "rui"      // Rufiji
	LanguageTypeRUK    LanguageType = "ruk"      // Che
	LanguageTypeRUO    LanguageType = "ruo"      // Istro Romanian
	LanguageTypeRUP    LanguageType = "rup"      // Macedo-Romanian and Aromanian and Arumanian
	LanguageTypeRUQ    LanguageType = "ruq"      // Megleno Romanian
	LanguageTypeRUT    LanguageType = "rut"      // Rutul
	LanguageTypeRUU    LanguageType = "ruu"      // Lanas Lobu
	LanguageTypeRUY    LanguageType = "ruy"      // Mala (Nigeria)
	LanguageTypeRUZ    LanguageType = "ruz"      // Ruma
	LanguageTypeRWA    LanguageType = "rwa"      // Rawo
	LanguageTypeRWK    LanguageType = "rwk"      // Rwa
	LanguageTypeRWM    LanguageType = "rwm"      // Amba (Uganda)
	LanguageTypeRWO    LanguageType = "rwo"      // Rawa
	LanguageTypeRWR    LanguageType = "rwr"      // Marwari (India)
	LanguageTypeRXD    LanguageType = "rxd"      // Ngardi
	LanguageTypeRXW    LanguageType = "rxw"      // Karuwali
	LanguageTypeRYN    LanguageType = "ryn"      // Northern Amami-Oshima
	LanguageTypeRYS    LanguageType = "rys"      // Yaeyama
	LanguageTypeRYU    LanguageType = "ryu"      // Central Okinawan
	LanguageTypeSAA    LanguageType = "saa"      // Saba
	LanguageTypeSAB    LanguageType = "sab"      // Buglere
	LanguageTypeSAC    LanguageType = "sac"      // Meskwaki
	LanguageTypeSAD    LanguageType = "sad"      // Sandawe
	LanguageTypeSAE    LanguageType = "sae"      // Sabanê
	LanguageTypeSAF    LanguageType = "saf"      // Safaliba
	LanguageTypeSAH    LanguageType = "sah"      // Yakut
	LanguageTypeSAI    LanguageType = "sai"      // South American Indian languages
	LanguageTypeSAJ    LanguageType = "saj"      // Sahu
	LanguageTypeSAK    LanguageType = "sak"      // Sake
	LanguageTypeSAL    LanguageType = "sal"      // Salishan languages
	LanguageTypeSAM    LanguageType = "sam"      // Samaritan Aramaic
	LanguageTypeSAO    LanguageType = "sao"      // Sause
	LanguageTypeSAP    LanguageType = "sap"      // Sanapaná
	LanguageTypeSAQ    LanguageType = "saq"      // Samburu
	LanguageTypeSAR    LanguageType = "sar"      // Saraveca
	LanguageTypeSAS    LanguageType = "sas"      // Sasak
	LanguageTypeSAT    LanguageType = "sat"      // Santali
	LanguageTypeSAU    LanguageType = "sau"      // Saleman
	LanguageTypeSAV    LanguageType = "sav"      // Saafi-Saafi
	LanguageTypeSAW    LanguageType = "saw"      // Sawi
	LanguageTypeSAX    LanguageType = "sax"      // Sa
	LanguageTypeSAY    LanguageType = "say"      // Saya
	LanguageTypeSAZ    LanguageType = "saz"      // Saurashtra
	LanguageTypeSBA    LanguageType = "sba"      // Ngambay
	LanguageTypeSBB    LanguageType = "sbb"      // Simbo
	LanguageTypeSBC    LanguageType = "sbc"      // Kele (Papua New Guinea)
	LanguageTypeSBD    LanguageType = "sbd"      // Southern Samo
	LanguageTypeSBE    LanguageType = "sbe"      // Saliba
	LanguageTypeSBF    LanguageType = "sbf"      // Shabo
	LanguageTypeSBG    LanguageType = "sbg"      // Seget
	LanguageTypeSBH    LanguageType = "sbh"      // Sori-Harengan
	LanguageTypeSBI    LanguageType = "sbi"      // Seti
	LanguageTypeSBJ    LanguageType = "sbj"      // Surbakhal
	LanguageTypeSBK    LanguageType = "sbk"      // Safwa
	LanguageTypeSBL    LanguageType = "sbl"      // Botolan Sambal
	LanguageTypeSBM    LanguageType = "sbm"      // Sagala
	LanguageTypeSBN    LanguageType = "sbn"      // Sindhi Bhil
	LanguageTypeSBO    LanguageType = "sbo"      // Sabüm
	LanguageTypeSBP    LanguageType = "sbp"      // Sangu (Tanzania)
	LanguageTypeSBQ    LanguageType = "sbq"      // Sileibi
	LanguageTypeSBR    LanguageType = "sbr"      // Sembakung Murut
	LanguageTypeSBS    LanguageType = "sbs"      // Subiya
	LanguageTypeSBT    LanguageType = "sbt"      // Kimki
	LanguageTypeSBU    LanguageType = "sbu"      // Stod Bhoti
	LanguageTypeSBV    LanguageType = "sbv"      // Sabine
	LanguageTypeSBW    LanguageType = "sbw"      // Simba
	LanguageTypeSBX    LanguageType = "sbx"      // Seberuang
	LanguageTypeSBY    LanguageType = "sby"      // Soli
	LanguageTypeSBZ    LanguageType = "sbz"      // Sara Kaba
	LanguageTypeSCA    LanguageType = "sca"      // Sansu
	LanguageTypeSCB    LanguageType = "scb"      // Chut
	LanguageTypeSCE    LanguageType = "sce"      // Dongxiang
	LanguageTypeSCF    LanguageType = "scf"      // San Miguel Creole French
	LanguageTypeSCG    LanguageType = "scg"      // Sanggau
	LanguageTypeSCH    LanguageType = "sch"      // Sakachep
	LanguageTypeSCI    LanguageType = "sci"      // Sri Lankan Creole Malay
	LanguageTypeSCK    LanguageType = "sck"      // Sadri
	LanguageTypeSCL    LanguageType = "scl"      // Shina
	LanguageTypeSCN    LanguageType = "scn"      // Sicilian
	LanguageTypeSCO    LanguageType = "sco"      // Scots
	LanguageTypeSCP    LanguageType = "scp"      // Helambu Sherpa
	LanguageTypeSCQ    LanguageType = "scq"      // Sa'och
	LanguageTypeSCS    LanguageType = "scs"      // North Slavey
	LanguageTypeSCU    LanguageType = "scu"      // Shumcho
	LanguageTypeSCV    LanguageType = "scv"      // Sheni
	LanguageTypeSCW    LanguageType = "scw"      // Sha
	LanguageTypeSCX    LanguageType = "scx"      // Sicel
	LanguageTypeSDA    LanguageType = "sda"      // Toraja-Sa'dan
	LanguageTypeSDB    LanguageType = "sdb"      // Shabak
	LanguageTypeSDC    LanguageType = "sdc"      // Sassarese Sardinian
	LanguageTypeSDE    LanguageType = "sde"      // Surubu
	LanguageTypeSDF    LanguageType = "sdf"      // Sarli
	LanguageTypeSDG    LanguageType = "sdg"      // Savi
	LanguageTypeSDH    LanguageType = "sdh"      // Southern Kurdish
	LanguageTypeSDJ    LanguageType = "sdj"      // Suundi
	LanguageTypeSDK    LanguageType = "sdk"      // Sos Kundi
	LanguageTypeSDL    LanguageType = "sdl"      // Saudi Arabian Sign Language
	LanguageTypeSDM    LanguageType = "sdm"      // Semandang
	LanguageTypeSDN    LanguageType = "sdn"      // Gallurese Sardinian
	LanguageTypeSDO    LanguageType = "sdo"      // Bukar-Sadung Bidayuh
	LanguageTypeSDP    LanguageType = "sdp"      // Sherdukpen
	LanguageTypeSDR    LanguageType = "sdr"      // Oraon Sadri
	LanguageTypeSDS    LanguageType = "sds"      // Sened
	LanguageTypeSDT    LanguageType = "sdt"      // Shuadit
	LanguageTypeSDU    LanguageType = "sdu"      // Sarudu
	LanguageTypeSDV    LanguageType = "sdv"      // Eastern Sudanic languages
	LanguageTypeSDX    LanguageType = "sdx"      // Sibu Melanau
	LanguageTypeSDZ    LanguageType = "sdz"      // Sallands
	LanguageTypeSEA    LanguageType = "sea"      // Semai
	LanguageTypeSEB    LanguageType = "seb"      // Shempire Senoufo
	LanguageTypeSEC    LanguageType = "sec"      // Sechelt
	LanguageTypeSED    LanguageType = "sed"      // Sedang
	LanguageTypeSEE    LanguageType = "see"      // Seneca
	LanguageTypeSEF    LanguageType = "sef"      // Cebaara Senoufo
	LanguageTypeSEG    LanguageType = "seg"      // Segeju
	LanguageTypeSEH    LanguageType = "seh"      // Sena
	LanguageTypeSEI    LanguageType = "sei"      // Seri
	LanguageTypeSEJ    LanguageType = "sej"      // Sene
	LanguageTypeSEK    LanguageType = "sek"      // Sekani
	LanguageTypeSEL    LanguageType = "sel"      // Selkup
	LanguageTypeSEM    LanguageType = "sem"      // Semitic languages
	LanguageTypeSEN    LanguageType = "sen"      // Nanerigé Sénoufo
	LanguageTypeSEO    LanguageType = "seo"      // Suarmin
	LanguageTypeSEP    LanguageType = "sep"      // Sìcìté Sénoufo
	LanguageTypeSEQ    LanguageType = "seq"      // Senara Sénoufo
	LanguageTypeSER    LanguageType = "ser"      // Serrano
	LanguageTypeSES    LanguageType = "ses"      // Koyraboro Senni Songhai
	LanguageTypeSET    LanguageType = "set"      // Sentani
	LanguageTypeSEU    LanguageType = "seu"      // Serui-Laut
	LanguageTypeSEV    LanguageType = "sev"      // Nyarafolo Senoufo
	LanguageTypeSEW    LanguageType = "sew"      // Sewa Bay
	LanguageTypeSEY    LanguageType = "sey"      // Secoya
	LanguageTypeSEZ    LanguageType = "sez"      // Senthang Chin
	LanguageTypeSFB    LanguageType = "sfb"      // Langue des signes de Belgique Francophone and French Belgian Sign Language
	LanguageTypeSFE    LanguageType = "sfe"      // Eastern Subanen
	LanguageTypeSFM    LanguageType = "sfm"      // Small Flowery Miao
	LanguageTypeSFS    LanguageType = "sfs"      // South African Sign Language
	LanguageTypeSFW    LanguageType = "sfw"      // Sehwi
	LanguageTypeSGA    LanguageType = "sga"      // Old Irish (to 900)
	LanguageTypeSGB    LanguageType = "sgb"      // Mag-antsi Ayta
	LanguageTypeSGC    LanguageType = "sgc"      // Kipsigis
	LanguageTypeSGD    LanguageType = "sgd"      // Surigaonon
	LanguageTypeSGE    LanguageType = "sge"      // Segai
	LanguageTypeSGG    LanguageType = "sgg"      // Swiss-German Sign Language
	LanguageTypeSGH    LanguageType = "sgh"      // Shughni
	LanguageTypeSGI    LanguageType = "sgi"      // Suga
	LanguageTypeSGJ    LanguageType = "sgj"      // Surgujia
	LanguageTypeSGK    LanguageType = "sgk"      // Sangkong
	LanguageTypeSGL    LanguageType = "sgl"      // Sanglechi-Ishkashimi
	LanguageTypeSGM    LanguageType = "sgm"      // Singa
	LanguageTypeSGN    LanguageType = "sgn"      // Sign languages
	LanguageTypeSGO    LanguageType = "sgo"      // Songa
	LanguageTypeSGP    LanguageType = "sgp"      // Singpho
	LanguageTypeSGR    LanguageType = "sgr"      // Sangisari
	LanguageTypeSGS    LanguageType = "sgs"      // Samogitian
	LanguageTypeSGT    LanguageType = "sgt"      // Brokpake
	LanguageTypeSGU    LanguageType = "sgu"      // Salas
	LanguageTypeSGW    LanguageType = "sgw"      // Sebat Bet Gurage
	LanguageTypeSGX    LanguageType = "sgx"      // Sierra Leone Sign Language
	LanguageTypeSGY    LanguageType = "sgy"      // Sanglechi
	LanguageTypeSGZ    LanguageType = "sgz"      // Sursurunga
	LanguageTypeSHA    LanguageType = "sha"      // Shall-Zwall
	LanguageTypeSHB    LanguageType = "shb"      // Ninam
	LanguageTypeSHC    LanguageType = "shc"      // Sonde
	LanguageTypeSHD    LanguageType = "shd"      // Kundal Shahi
	LanguageTypeSHE    LanguageType = "she"      // Sheko
	LanguageTypeSHG    LanguageType = "shg"      // Shua
	LanguageTypeSHH    LanguageType = "shh"      // Shoshoni
	LanguageTypeSHI    LanguageType = "shi"      // Tachelhit
	LanguageTypeSHJ    LanguageType = "shj"      // Shatt
	LanguageTypeSHK    LanguageType = "shk"      // Shilluk
	LanguageTypeSHL    LanguageType = "shl"      // Shendu
	LanguageTypeSHM    LanguageType = "shm"      // Shahrudi
	LanguageTypeSHN    LanguageType = "shn"      // Shan
	LanguageTypeSHO    LanguageType = "sho"      // Shanga
	LanguageTypeSHP    LanguageType = "shp"      // Shipibo-Conibo
	LanguageTypeSHQ    LanguageType = "shq"      // Sala
	LanguageTypeSHR    LanguageType = "shr"      // Shi
	LanguageTypeSHS    LanguageType = "shs"      // Shuswap
	LanguageTypeSHT    LanguageType = "sht"      // Shasta
	LanguageTypeSHU    LanguageType = "shu"      // Chadian Arabic
	LanguageTypeSHV    LanguageType = "shv"      // Shehri
	LanguageTypeSHW    LanguageType = "shw"      // Shwai
	LanguageTypeSHX    LanguageType = "shx"      // She
	LanguageTypeSHY    LanguageType = "shy"      // Tachawit
	LanguageTypeSHZ    LanguageType = "shz"      // Syenara Senoufo
	LanguageTypeSIA    LanguageType = "sia"      // Akkala Sami
	LanguageTypeSIB    LanguageType = "sib"      // Sebop
	LanguageTypeSID    LanguageType = "sid"      // Sidamo
	LanguageTypeSIE    LanguageType = "sie"      // Simaa
	LanguageTypeSIF    LanguageType = "sif"      // Siamou
	LanguageTypeSIG    LanguageType = "sig"      // Paasaal
	LanguageTypeSIH    LanguageType = "sih"      // Zire and Sîshëë
	LanguageTypeSII    LanguageType = "sii"      // Shom Peng
	LanguageTypeSIJ    LanguageType = "sij"      // Numbami
	LanguageTypeSIK    LanguageType = "sik"      // Sikiana
	LanguageTypeSIL    LanguageType = "sil"      // Tumulung Sisaala
	LanguageTypeSIM    LanguageType = "sim"      // Mende (Papua New Guinea)
	LanguageTypeSIO    LanguageType = "sio"      // Siouan languages
	LanguageTypeSIP    LanguageType = "sip"      // Sikkimese
	LanguageTypeSIQ    LanguageType = "siq"      // Sonia
	LanguageTypeSIR    LanguageType = "sir"      // Siri
	LanguageTypeSIS    LanguageType = "sis"      // Siuslaw
	LanguageTypeSIT    LanguageType = "sit"      // Sino-Tibetan languages
	LanguageTypeSIU    LanguageType = "siu"      // Sinagen
	LanguageTypeSIV    LanguageType = "siv"      // Sumariup
	LanguageTypeSIW    LanguageType = "siw"      // Siwai
	LanguageTypeSIX    LanguageType = "six"      // Sumau
	LanguageTypeSIY    LanguageType = "siy"      // Sivandi
	LanguageTypeSIZ    LanguageType = "siz"      // Siwi
	LanguageTypeSJA    LanguageType = "sja"      // Epena
	LanguageTypeSJB    LanguageType = "sjb"      // Sajau Basap
	LanguageTypeSJD    LanguageType = "sjd"      // Kildin Sami
	LanguageTypeSJE    LanguageType = "sje"      // Pite Sami
	LanguageTypeSJG    LanguageType = "sjg"      // Assangori
	LanguageTypeSJK    LanguageType = "sjk"      // Kemi Sami
	LanguageTypeSJL    LanguageType = "sjl"      // Sajalong and Miji
	LanguageTypeSJM    LanguageType = "sjm"      // Mapun
	LanguageTypeSJN    LanguageType = "sjn"      // Sindarin
	LanguageTypeSJO    LanguageType = "sjo"      // Xibe
	LanguageTypeSJP    LanguageType = "sjp"      // Surjapuri
	LanguageTypeSJR    LanguageType = "sjr"      // Siar-Lak
	LanguageTypeSJS    LanguageType = "sjs"      // Senhaja De Srair
	LanguageTypeSJT    LanguageType = "sjt"      // Ter Sami
	LanguageTypeSJU    LanguageType = "sju"      // Ume Sami
	LanguageTypeSJW    LanguageType = "sjw"      // Shawnee
	LanguageTypeSKA    LanguageType = "ska"      // Skagit
	LanguageTypeSKB    LanguageType = "skb"      // Saek
	LanguageTypeSKC    LanguageType = "skc"      // Ma Manda
	LanguageTypeSKD    LanguageType = "skd"      // Southern Sierra Miwok
	LanguageTypeSKE    LanguageType = "ske"      // Seke (Vanuatu)
	LanguageTypeSKF    LanguageType = "skf"      // Sakirabiá
	LanguageTypeSKG    LanguageType = "skg"      // Sakalava Malagasy
	LanguageTypeSKH    LanguageType = "skh"      // Sikule
	LanguageTypeSKI    LanguageType = "ski"      // Sika
	LanguageTypeSKJ    LanguageType = "skj"      // Seke (Nepal)
	LanguageTypeSKK    LanguageType = "skk"      // Sok
	LanguageTypeSKM    LanguageType = "skm"      // Kutong
	LanguageTypeSKN    LanguageType = "skn"      // Kolibugan Subanon
	LanguageTypeSKO    LanguageType = "sko"      // Seko Tengah
	LanguageTypeSKP    LanguageType = "skp"      // Sekapan
	LanguageTypeSKQ    LanguageType = "skq"      // Sininkere
	LanguageTypeSKR    LanguageType = "skr"      // Seraiki
	LanguageTypeSKS    LanguageType = "sks"      // Maia
	LanguageTypeSKT    LanguageType = "skt"      // Sakata
	LanguageTypeSKU    LanguageType = "sku"      // Sakao
	LanguageTypeSKV    LanguageType = "skv"      // Skou
	LanguageTypeSKW    LanguageType = "skw"      // Skepi Creole Dutch
	LanguageTypeSKX    LanguageType = "skx"      // Seko Padang
	LanguageTypeSKY    LanguageType = "sky"      // Sikaiana
	LanguageTypeSKZ    LanguageType = "skz"      // Sekar
	LanguageTypeSLA    LanguageType = "sla"      // Slavic languages
	LanguageTypeSLC    LanguageType = "slc"      // Sáliba
	LanguageTypeSLD    LanguageType = "sld"      // Sissala
	LanguageTypeSLE    LanguageType = "sle"      // Sholaga
	LanguageTypeSLF    LanguageType = "slf"      // Swiss-Italian Sign Language
	LanguageTypeSLG    LanguageType = "slg"      // Selungai Murut
	LanguageTypeSLH    LanguageType = "slh"      // Southern Puget Sound Salish
	LanguageTypeSLI    LanguageType = "sli"      // Lower Silesian
	LanguageTypeSLJ    LanguageType = "slj"      // Salumá
	LanguageTypeSLL    LanguageType = "sll"      // Salt-Yui
	LanguageTypeSLM    LanguageType = "slm"      // Pangutaran Sama
	LanguageTypeSLN    LanguageType = "sln"      // Salinan
	LanguageTypeSLP    LanguageType = "slp"      // Lamaholot
	LanguageTypeSLQ    LanguageType = "slq"      // Salchuq
	LanguageTypeSLR    LanguageType = "slr"      // Salar
	LanguageTypeSLS    LanguageType = "sls"      // Singapore Sign Language
	LanguageTypeSLT    LanguageType = "slt"      // Sila
	LanguageTypeSLU    LanguageType = "slu"      // Selaru
	LanguageTypeSLW    LanguageType = "slw"      // Sialum
	LanguageTypeSLX    LanguageType = "slx"      // Salampasu
	LanguageTypeSLY    LanguageType = "sly"      // Selayar
	LanguageTypeSLZ    LanguageType = "slz"      // Ma'ya
	LanguageTypeSMA    LanguageType = "sma"      // Southern Sami
	LanguageTypeSMB    LanguageType = "smb"      // Simbari
	LanguageTypeSMC    LanguageType = "smc"      // Som
	LanguageTypeSMD    LanguageType = "smd"      // Sama
	LanguageTypeSMF    LanguageType = "smf"      // Auwe
	LanguageTypeSMG    LanguageType = "smg"      // Simbali
	LanguageTypeSMH    LanguageType = "smh"      // Samei
	LanguageTypeSMI    LanguageType = "smi"      // Sami languages
	LanguageTypeSMJ    LanguageType = "smj"      // Lule Sami
	LanguageTypeSMK    LanguageType = "smk"      // Bolinao
	LanguageTypeSML    LanguageType = "sml"      // Central Sama
	LanguageTypeSMM    LanguageType = "smm"      // Musasa
	LanguageTypeSMN    LanguageType = "smn"      // Inari Sami
	LanguageTypeSMP    LanguageType = "smp"      // Samaritan
	LanguageTypeSMQ    LanguageType = "smq"      // Samo
	LanguageTypeSMR    LanguageType = "smr"      // Simeulue
	LanguageTypeSMS    LanguageType = "sms"      // Skolt Sami
	LanguageTypeSMT    LanguageType = "smt"      // Simte
	LanguageTypeSMU    LanguageType = "smu"      // Somray
	LanguageTypeSMV    LanguageType = "smv"      // Samvedi
	LanguageTypeSMW    LanguageType = "smw"      // Sumbawa
	LanguageTypeSMX    LanguageType = "smx"      // Samba
	LanguageTypeSMY    LanguageType = "smy"      // Semnani
	LanguageTypeSMZ    LanguageType = "smz"      // Simeku
	LanguageTypeSNB    LanguageType = "snb"      // Sebuyau
	LanguageTypeSNC    LanguageType = "snc"      // Sinaugoro
	LanguageTypeSNE    LanguageType = "sne"      // Bau Bidayuh
	LanguageTypeSNF    LanguageType = "snf"      // Noon
	LanguageTypeSNG    LanguageType = "sng"      // Sanga (Democratic Republic of Congo)
	LanguageTypeSNH    LanguageType = "snh"      // Shinabo
	LanguageTypeSNI    LanguageType = "sni"      // Sensi
	LanguageTypeSNJ    LanguageType = "snj"      // Riverain Sango
	LanguageTypeSNK    LanguageType = "snk"      // Soninke
	LanguageTypeSNL    LanguageType = "snl"      // Sangil
	LanguageTypeSNM    LanguageType = "snm"      // Southern Ma'di
	LanguageTypeSNN    LanguageType = "snn"      // Siona
	LanguageTypeSNO    LanguageType = "sno"      // Snohomish
	LanguageTypeSNP    LanguageType = "snp"      // Siane
	LanguageTypeSNQ    LanguageType = "snq"      // Sangu (Gabon)
	LanguageTypeSNR    LanguageType = "snr"      // Sihan
	LanguageTypeSNS    LanguageType = "sns"      // South West Bay and Nahavaq
	LanguageTypeSNU    LanguageType = "snu"      // Senggi and Viid
	LanguageTypeSNV    LanguageType = "snv"      // Sa'ban
	LanguageTypeSNW    LanguageType = "snw"      // Selee
	LanguageTypeSNX    LanguageType = "snx"      // Sam
	LanguageTypeSNY    LanguageType = "sny"      // Saniyo-Hiyewe
	LanguageTypeSNZ    LanguageType = "snz"      // Sinsauru
	LanguageTypeSOA    LanguageType = "soa"      // Thai Song
	LanguageTypeSOB    LanguageType = "sob"      // Sobei
	LanguageTypeSOC    LanguageType = "soc"      // So (Democratic Republic of Congo)
	LanguageTypeSOD    LanguageType = "sod"      // Songoora
	LanguageTypeSOE    LanguageType = "soe"      // Songomeno
	LanguageTypeSOG    LanguageType = "sog"      // Sogdian
	LanguageTypeSOH    LanguageType = "soh"      // Aka
	LanguageTypeSOI    LanguageType = "soi"      // Sonha
	LanguageTypeSOJ    LanguageType = "soj"      // Soi
	LanguageTypeSOK    LanguageType = "sok"      // Sokoro
	LanguageTypeSOL    LanguageType = "sol"      // Solos
	LanguageTypeSON    LanguageType = "son"      // Songhai languages
	LanguageTypeSOO    LanguageType = "soo"      // Songo
	LanguageTypeSOP    LanguageType = "sop"      // Songe
	LanguageTypeSOQ    LanguageType = "soq"      // Kanasi
	LanguageTypeSOR    LanguageType = "sor"      // Somrai
	LanguageTypeSOS    LanguageType = "sos"      // Seeku
	LanguageTypeSOU    LanguageType = "sou"      // Southern Thai
	LanguageTypeSOV    LanguageType = "sov"      // Sonsorol
	LanguageTypeSOW    LanguageType = "sow"      // Sowanda
	LanguageTypeSOX    LanguageType = "sox"      // Swo
	LanguageTypeSOY    LanguageType = "soy"      // Miyobe
	LanguageTypeSOZ    LanguageType = "soz"      // Temi
	LanguageTypeSPB    LanguageType = "spb"      // Sepa (Indonesia)
	LanguageTypeSPC    LanguageType = "spc"      // Sapé
	LanguageTypeSPD    LanguageType = "spd"      // Saep
	LanguageTypeSPE    LanguageType = "spe"      // Sepa (Papua New Guinea)
	LanguageTypeSPG    LanguageType = "spg"      // Sian
	LanguageTypeSPI    LanguageType = "spi"      // Saponi
	LanguageTypeSPK    LanguageType = "spk"      // Sengo
	LanguageTypeSPL    LanguageType = "spl"      // Selepet
	LanguageTypeSPM    LanguageType = "spm"      // Akukem
	LanguageTypeSPO    LanguageType = "spo"      // Spokane
	LanguageTypeSPP    LanguageType = "spp"      // Supyire Senoufo
	LanguageTypeSPQ    LanguageType = "spq"      // Loreto-Ucayali Spanish
	LanguageTypeSPR    LanguageType = "spr"      // Saparua
	LanguageTypeSPS    LanguageType = "sps"      // Saposa
	LanguageTypeSPT    LanguageType = "spt"      // Spiti Bhoti
	LanguageTypeSPU    LanguageType = "spu"      // Sapuan
	LanguageTypeSPV    LanguageType = "spv"      // Sambalpuri and Kosli
	LanguageTypeSPX    LanguageType = "spx"      // South Picene
	LanguageTypeSPY    LanguageType = "spy"      // Sabaot
	LanguageTypeSQA    LanguageType = "sqa"      // Shama-Sambuga
	LanguageTypeSQH    LanguageType = "sqh"      // Shau
	LanguageTypeSQJ    LanguageType = "sqj"      // Albanian languages
	LanguageTypeSQK    LanguageType = "sqk"      // Albanian Sign Language
	LanguageTypeSQM    LanguageType = "sqm"      // Suma
	LanguageTypeSQN    LanguageType = "sqn"      // Susquehannock
	LanguageTypeSQO    LanguageType = "sqo"      // Sorkhei
	LanguageTypeSQQ    LanguageType = "sqq"      // Sou
	LanguageTypeSQR    LanguageType = "sqr"      // Siculo Arabic
	LanguageTypeSQS    LanguageType = "sqs"      // Sri Lankan Sign Language
	LanguageTypeSQT    LanguageType = "sqt"      // Soqotri
	LanguageTypeSQU    LanguageType = "squ"      // Squamish
	LanguageTypeSRA    LanguageType = "sra"      // Saruga
	LanguageTypeSRB    LanguageType = "srb"      // Sora
	LanguageTypeSRC    LanguageType = "src"      // Logudorese Sardinian
	LanguageTypeSRE    LanguageType = "sre"      // Sara
	LanguageTypeSRF    LanguageType = "srf"      // Nafi
	LanguageTypeSRG    LanguageType = "srg"      // Sulod
	LanguageTypeSRH    LanguageType = "srh"      // Sarikoli
	LanguageTypeSRI    LanguageType = "sri"      // Siriano
	LanguageTypeSRK    LanguageType = "srk"      // Serudung Murut
	LanguageTypeSRL    LanguageType = "srl"      // Isirawa
	LanguageTypeSRM    LanguageType = "srm"      // Saramaccan
	LanguageTypeSRN    LanguageType = "srn"      // Sranan Tongo
	LanguageTypeSRO    LanguageType = "sro"      // Campidanese Sardinian
	LanguageTypeSRQ    LanguageType = "srq"      // Sirionó
	LanguageTypeSRR    LanguageType = "srr"      // Serer
	LanguageTypeSRS    LanguageType = "srs"      // Sarsi
	LanguageTypeSRT    LanguageType = "srt"      // Sauri
	LanguageTypeSRU    LanguageType = "sru"      // Suruí
	LanguageTypeSRV    LanguageType = "srv"      // Southern Sorsoganon
	LanguageTypeSRW    LanguageType = "srw"      // Serua
	LanguageTypeSRX    LanguageType = "srx"      // Sirmauri
	LanguageTypeSRY    LanguageType = "sry"      // Sera
	LanguageTypeSRZ    LanguageType = "srz"      // Shahmirzadi
	LanguageTypeSSA    LanguageType = "ssa"      // Nilo-Saharan languages
	LanguageTypeSSB    LanguageType = "ssb"      // Southern Sama
	LanguageTypeSSC    LanguageType = "ssc"      // Suba-Simbiti
	LanguageTypeSSD    LanguageType = "ssd"      // Siroi
	LanguageTypeSSE    LanguageType = "sse"      // Balangingi and Bangingih Sama
	LanguageTypeSSF    LanguageType = "ssf"      // Thao
	LanguageTypeSSG    LanguageType = "ssg"      // Seimat
	LanguageTypeSSH    LanguageType = "ssh"      // Shihhi Arabic
	LanguageTypeSSI    LanguageType = "ssi"      // Sansi
	LanguageTypeSSJ    LanguageType = "ssj"      // Sausi
	LanguageTypeSSK    LanguageType = "ssk"      // Sunam
	LanguageTypeSSL    LanguageType = "ssl"      // Western Sisaala
	LanguageTypeSSM    LanguageType = "ssm"      // Semnam
	LanguageTypeSSN    LanguageType = "ssn"      // Waata
	LanguageTypeSSO    LanguageType = "sso"      // Sissano
	LanguageTypeSSP    LanguageType = "ssp"      // Spanish Sign Language
	LanguageTypeSSQ    LanguageType = "ssq"      // So'a
	LanguageTypeSSR    LanguageType = "ssr"      // Swiss-French Sign Language
	LanguageTypeSSS    LanguageType = "sss"      // Sô
	LanguageTypeSST    LanguageType = "sst"      // Sinasina
	LanguageTypeSSU    LanguageType = "ssu"      // Susuami
	LanguageTypeSSV    LanguageType = "ssv"      // Shark Bay
	LanguageTypeSSX    LanguageType = "ssx"      // Samberigi
	LanguageTypeSSY    LanguageType = "ssy"      // Saho
	LanguageTypeSSZ    LanguageType = "ssz"      // Sengseng
	LanguageTypeSTA    LanguageType = "sta"      // Settla
	LanguageTypeSTB    LanguageType = "stb"      // Northern Subanen
	LanguageTypeSTD    LanguageType = "std"      // Sentinel
	LanguageTypeSTE    LanguageType = "ste"      // Liana-Seti
	LanguageTypeSTF    LanguageType = "stf"      // Seta
	LanguageTypeSTG    LanguageType = "stg"      // Trieng
	LanguageTypeSTH    LanguageType = "sth"      // Shelta
	LanguageTypeSTI    LanguageType = "sti"      // Bulo Stieng
	LanguageTypeSTJ    LanguageType = "stj"      // Matya Samo
	LanguageTypeSTK    LanguageType = "stk"      // Arammba
	LanguageTypeSTL    LanguageType = "stl"      // Stellingwerfs
	LanguageTypeSTM    LanguageType = "stm"      // Setaman
	LanguageTypeSTN    LanguageType = "stn"      // Owa
	LanguageTypeSTO    LanguageType = "sto"      // Stoney
	LanguageTypeSTP    LanguageType = "stp"      // Southeastern Tepehuan
	LanguageTypeSTQ    LanguageType = "stq"      // Saterfriesisch
	LanguageTypeSTR    LanguageType = "str"      // Straits Salish
	LanguageTypeSTS    LanguageType = "sts"      // Shumashti
	LanguageTypeSTT    LanguageType = "stt"      // Budeh Stieng
	LanguageTypeSTU    LanguageType = "stu"      // Samtao
	LanguageTypeSTV    LanguageType = "stv"      // Silt'e
	LanguageTypeSTW    LanguageType = "stw"      // Satawalese
	LanguageTypeSTY    LanguageType = "sty"      // Siberian Tatar
	LanguageTypeSUA    LanguageType = "sua"      // Sulka
	LanguageTypeSUB    LanguageType = "sub"      // Suku
	LanguageTypeSUC    LanguageType = "suc"      // Western Subanon
	LanguageTypeSUE    LanguageType = "sue"      // Suena
	LanguageTypeSUG    LanguageType = "sug"      // Suganga
	LanguageTypeSUI    LanguageType = "sui"      // Suki
	LanguageTypeSUJ    LanguageType = "suj"      // Shubi
	LanguageTypeSUK    LanguageType = "suk"      // Sukuma
	LanguageTypeSUL    LanguageType = "sul"      // Surigaonon
	LanguageTypeSUM    LanguageType = "sum"      // Sumo-Mayangna
	LanguageTypeSUQ    LanguageType = "suq"      // Suri
	LanguageTypeSUR    LanguageType = "sur"      // Mwaghavul
	LanguageTypeSUS    LanguageType = "sus"      // Susu
	LanguageTypeSUT    LanguageType = "sut"      // Subtiaba
	LanguageTypeSUV    LanguageType = "suv"      // Puroik
	LanguageTypeSUW    LanguageType = "suw"      // Sumbwa
	LanguageTypeSUX    LanguageType = "sux"      // Sumerian
	LanguageTypeSUY    LanguageType = "suy"      // Suyá
	LanguageTypeSUZ    LanguageType = "suz"      // Sunwar
	LanguageTypeSVA    LanguageType = "sva"      // Svan
	LanguageTypeSVB    LanguageType = "svb"      // Ulau-Suain
	LanguageTypeSVC    LanguageType = "svc"      // Vincentian Creole English
	LanguageTypeSVE    LanguageType = "sve"      // Serili
	LanguageTypeSVK    LanguageType = "svk"      // Slovakian Sign Language
	LanguageTypeSVM    LanguageType = "svm"      // Slavomolisano
	LanguageTypeSVR    LanguageType = "svr"      // Savara
	LanguageTypeSVS    LanguageType = "svs"      // Savosavo
	LanguageTypeSVX    LanguageType = "svx"      // Skalvian
	LanguageTypeSWB    LanguageType = "swb"      // Maore Comorian
	LanguageTypeSWC    LanguageType = "swc"      // Congo Swahili
	LanguageTypeSWF    LanguageType = "swf"      // Sere
	LanguageTypeSWG    LanguageType = "swg"      // Swabian
	LanguageTypeSWH    LanguageType = "swh"      // Swahili (individual language) and Kiswahili
	LanguageTypeSWI    LanguageType = "swi"      // Sui
	LanguageTypeSWJ    LanguageType = "swj"      // Sira
	LanguageTypeSWK    LanguageType = "swk"      // Malawi Sena
	LanguageTypeSWL    LanguageType = "swl"      // Swedish Sign Language
	LanguageTypeSWM    LanguageType = "swm"      // Samosa
	LanguageTypeSWN    LanguageType = "swn"      // Sawknah
	LanguageTypeSWO    LanguageType = "swo"      // Shanenawa
	LanguageTypeSWP    LanguageType = "swp"      // Suau
	LanguageTypeSWQ    LanguageType = "swq"      // Sharwa
	LanguageTypeSWR    LanguageType = "swr"      // Saweru
	LanguageTypeSWS    LanguageType = "sws"      // Seluwasan
	LanguageTypeSWT    LanguageType = "swt"      // Sawila
	LanguageTypeSWU    LanguageType = "swu"      // Suwawa
	LanguageTypeSWV    LanguageType = "swv"      // Shekhawati
	LanguageTypeSWW    LanguageType = "sww"      // Sowa
	LanguageTypeSWX    LanguageType = "swx"      // Suruahá
	LanguageTypeSWY    LanguageType = "swy"      // Sarua
	LanguageTypeSXB    LanguageType = "sxb"      // Suba
	LanguageTypeSXC    LanguageType = "sxc"      // Sicanian
	LanguageTypeSXE    LanguageType = "sxe"      // Sighu
	LanguageTypeSXG    LanguageType = "sxg"      // Shixing
	LanguageTypeSXK    LanguageType = "sxk"      // Southern Kalapuya
	LanguageTypeSXL    LanguageType = "sxl"      // Selian
	LanguageTypeSXM    LanguageType = "sxm"      // Samre
	LanguageTypeSXN    LanguageType = "sxn"      // Sangir
	LanguageTypeSXO    LanguageType = "sxo"      // Sorothaptic
	LanguageTypeSXR    LanguageType = "sxr"      // Saaroa
	LanguageTypeSXS    LanguageType = "sxs"      // Sasaru
	LanguageTypeSXU    LanguageType = "sxu"      // Upper Saxon
	LanguageTypeSXW    LanguageType = "sxw"      // Saxwe Gbe
	LanguageTypeSYA    LanguageType = "sya"      // Siang
	LanguageTypeSYB    LanguageType = "syb"      // Central Subanen
	LanguageTypeSYC    LanguageType = "syc"      // Classical Syriac
	LanguageTypeSYD    LanguageType = "syd"      // Samoyedic languages
	LanguageTypeSYI    LanguageType = "syi"      // Seki
	LanguageTypeSYK    LanguageType = "syk"      // Sukur
	LanguageTypeSYL    LanguageType = "syl"      // Sylheti
	LanguageTypeSYM    LanguageType = "sym"      // Maya Samo
	LanguageTypeSYN    LanguageType = "syn"      // Senaya
	LanguageTypeSYO    LanguageType = "syo"      // Suoy
	LanguageTypeSYR    LanguageType = "syr"      // Syriac
	LanguageTypeSYS    LanguageType = "sys"      // Sinyar
	LanguageTypeSYW    LanguageType = "syw"      // Kagate
	LanguageTypeSYY    LanguageType = "syy"      // Al-Sayyid Bedouin Sign Language
	LanguageTypeSZA    LanguageType = "sza"      // Semelai
	LanguageTypeSZB    LanguageType = "szb"      // Ngalum
	LanguageTypeSZC    LanguageType = "szc"      // Semaq Beri
	LanguageTypeSZD    LanguageType = "szd"      // Seru
	LanguageTypeSZE    LanguageType = "sze"      // Seze
	LanguageTypeSZG    LanguageType = "szg"      // Sengele
	LanguageTypeSZL    LanguageType = "szl"      // Silesian
	LanguageTypeSZN    LanguageType = "szn"      // Sula
	LanguageTypeSZP    LanguageType = "szp"      // Suabo
	LanguageTypeSZV    LanguageType = "szv"      // Isu (Fako Division)
	LanguageTypeSZW    LanguageType = "szw"      // Sawai
	LanguageTypeTAA    LanguageType = "taa"      // Lower Tanana
	LanguageTypeTAB    LanguageType = "tab"      // Tabassaran
	LanguageTypeTAC    LanguageType = "tac"      // Lowland Tarahumara
	LanguageTypeTAD    LanguageType = "tad"      // Tause
	LanguageTypeTAE    LanguageType = "tae"      // Tariana
	LanguageTypeTAF    LanguageType = "taf"      // Tapirapé
	LanguageTypeTAG    LanguageType = "tag"      // Tagoi
	LanguageTypeTAI    LanguageType = "tai"      // Tai languages
	LanguageTypeTAJ    LanguageType = "taj"      // Eastern Tamang
	LanguageTypeTAK    LanguageType = "tak"      // Tala
	LanguageTypeTAL    LanguageType = "tal"      // Tal
	LanguageTypeTAN    LanguageType = "tan"      // Tangale
	LanguageTypeTAO    LanguageType = "tao"      // Yami
	LanguageTypeTAP    LanguageType = "tap"      // Taabwa
	LanguageTypeTAQ    LanguageType = "taq"      // Tamasheq
	LanguageTypeTAR    LanguageType = "tar"      // Central Tarahumara
	LanguageTypeTAS    LanguageType = "tas"      // Tay Boi
	LanguageTypeTAU    LanguageType = "tau"      // Upper Tanana
	LanguageTypeTAV    LanguageType = "tav"      // Tatuyo
	LanguageTypeTAW    LanguageType = "taw"      // Tai
	LanguageTypeTAX    LanguageType = "tax"      // Tamki
	LanguageTypeTAY    LanguageType = "tay"      // Atayal
	LanguageTypeTAZ    LanguageType = "taz"      // Tocho
	LanguageTypeTBA    LanguageType = "tba"      // Aikanã
	LanguageTypeTBB    LanguageType = "tbb"      // Tapeba
	LanguageTypeTBC    LanguageType = "tbc"      // Takia
	LanguageTypeTBD    LanguageType = "tbd"      // Kaki Ae
	LanguageTypeTBE    LanguageType = "tbe"      // Tanimbili
	LanguageTypeTBF    LanguageType = "tbf"      // Mandara
	LanguageTypeTBG    LanguageType = "tbg"      // North Tairora
	LanguageTypeTBH    LanguageType = "tbh"      // Thurawal
	LanguageTypeTBI    LanguageType = "tbi"      // Gaam
	LanguageTypeTBJ    LanguageType = "tbj"      // Tiang
	LanguageTypeTBK    LanguageType = "tbk"      // Calamian Tagbanwa
	LanguageTypeTBL    LanguageType = "tbl"      // Tboli
	LanguageTypeTBM    LanguageType = "tbm"      // Tagbu
	LanguageTypeTBN    LanguageType = "tbn"      // Barro Negro Tunebo
	LanguageTypeTBO    LanguageType = "tbo"      // Tawala
	LanguageTypeTBP    LanguageType = "tbp"      // Taworta and Diebroud
	LanguageTypeTBQ    LanguageType = "tbq"      // Tibeto-Burman languages
	LanguageTypeTBR    LanguageType = "tbr"      // Tumtum
	LanguageTypeTBS    LanguageType = "tbs"      // Tanguat
	LanguageTypeTBT    LanguageType = "tbt"      // Tembo (Kitembo)
	LanguageTypeTBU    LanguageType = "tbu"      // Tubar
	LanguageTypeTBV    LanguageType = "tbv"      // Tobo
	LanguageTypeTBW    LanguageType = "tbw"      // Tagbanwa
	LanguageTypeTBX    LanguageType = "tbx"      // Kapin
	LanguageTypeTBY    LanguageType = "tby"      // Tabaru
	LanguageTypeTBZ    LanguageType = "tbz"      // Ditammari
	LanguageTypeTCA    LanguageType = "tca"      // Ticuna
	LanguageTypeTCB    LanguageType = "tcb"      // Tanacross
	LanguageTypeTCC    LanguageType = "tcc"      // Datooga
	LanguageTypeTCD    LanguageType = "tcd"      // Tafi
	LanguageTypeTCE    LanguageType = "tce"      // Southern Tutchone
	LanguageTypeTCF    LanguageType = "tcf"      // Malinaltepec Me'phaa and Malinaltepec Tlapanec
	LanguageTypeTCG    LanguageType = "tcg"      // Tamagario
	LanguageTypeTCH    LanguageType = "tch"      // Turks And Caicos Creole English
	LanguageTypeTCI    LanguageType = "tci"      // Wára
	LanguageTypeTCK    LanguageType = "tck"      // Tchitchege
	LanguageTypeTCL    LanguageType = "tcl"      // Taman (Myanmar)
	LanguageTypeTCM    LanguageType = "tcm"      // Tanahmerah
	LanguageTypeTCN    LanguageType = "tcn"      // Tichurong
	LanguageTypeTCO    LanguageType = "tco"      // Taungyo
	LanguageTypeTCP    LanguageType = "tcp"      // Tawr Chin
	LanguageTypeTCQ    LanguageType = "tcq"      // Kaiy
	LanguageTypeTCS    LanguageType = "tcs"      // Torres Strait Creole
	LanguageTypeTCT    LanguageType = "tct"      // T'en
	LanguageTypeTCU    LanguageType = "tcu"      // Southeastern Tarahumara
	LanguageTypeTCW    LanguageType = "tcw"      // Tecpatlán Totonac
	LanguageTypeTCX    LanguageType = "tcx"      // Toda
	LanguageTypeTCY    LanguageType = "tcy"      // Tulu
	LanguageTypeTCZ    LanguageType = "tcz"      // Thado Chin
	LanguageTypeTDA    LanguageType = "tda"      // Tagdal
	LanguageTypeTDB    LanguageType = "tdb"      // Panchpargania
	LanguageTypeTDC    LanguageType = "tdc"      // Emberá-Tadó
	LanguageTypeTDD    LanguageType = "tdd"      // Tai Nüa
	LanguageTypeTDE    LanguageType = "tde"      // Tiranige Diga Dogon
	LanguageTypeTDF    LanguageType = "tdf"      // Talieng
	LanguageTypeTDG    LanguageType = "tdg"      // Western Tamang
	LanguageTypeTDH    LanguageType = "tdh"      // Thulung
	LanguageTypeTDI    LanguageType = "tdi"      // Tomadino
	LanguageTypeTDJ    LanguageType = "tdj"      // Tajio
	LanguageTypeTDK    LanguageType = "tdk"      // Tambas
	LanguageTypeTDL    LanguageType = "tdl"      // Sur
	LanguageTypeTDN    LanguageType = "tdn"      // Tondano
	LanguageTypeTDO    LanguageType = "tdo"      // Teme
	LanguageTypeTDQ    LanguageType = "tdq"      // Tita
	LanguageTypeTDR    LanguageType = "tdr"      // Todrah
	LanguageTypeTDS    LanguageType = "tds"      // Doutai
	LanguageTypeTDT    LanguageType = "tdt"      // Tetun Dili
	LanguageTypeTDU    LanguageType = "tdu"      // Tempasuk Dusun
	LanguageTypeTDV    LanguageType = "tdv"      // Toro
	LanguageTypeTDX    LanguageType = "tdx"      // Tandroy-Mahafaly Malagasy
	LanguageTypeTDY    LanguageType = "tdy"      // Tadyawan
	LanguageTypeTEA    LanguageType = "tea"      // Temiar
	LanguageTypeTEB    LanguageType = "teb"      // Tetete
	LanguageTypeTEC    LanguageType = "tec"      // Terik
	LanguageTypeTED    LanguageType = "ted"      // Tepo Krumen
	LanguageTypeTEE    LanguageType = "tee"      // Huehuetla Tepehua
	LanguageTypeTEF    LanguageType = "tef"      // Teressa
	LanguageTypeTEG    LanguageType = "teg"      // Teke-Tege
	LanguageTypeTEH    LanguageType = "teh"      // Tehuelche
	LanguageTypeTEI    LanguageType = "tei"      // Torricelli
	LanguageTypeTEK    LanguageType = "tek"      // Ibali Teke
	LanguageTypeTEM    LanguageType = "tem"      // Timne
	LanguageTypeTEN    LanguageType = "ten"      // Tama (Colombia)
	LanguageTypeTEO    LanguageType = "teo"      // Teso
	LanguageTypeTEP    LanguageType = "tep"      // Tepecano
	LanguageTypeTEQ    LanguageType = "teq"      // Temein
	LanguageTypeTER    LanguageType = "ter"      // Tereno
	LanguageTypeTES    LanguageType = "tes"      // Tengger
	LanguageTypeTET    LanguageType = "tet"      // Tetum
	LanguageTypeTEU    LanguageType = "teu"      // Soo
	LanguageTypeTEV    LanguageType = "tev"      // Teor
	LanguageTypeTEW    LanguageType = "tew"      // Tewa (USA)
	LanguageTypeTEX    LanguageType = "tex"      // Tennet
	LanguageTypeTEY    LanguageType = "tey"      // Tulishi
	LanguageTypeTFI    LanguageType = "tfi"      // Tofin Gbe
	LanguageTypeTFN    LanguageType = "tfn"      // Tanaina
	LanguageTypeTFO    LanguageType = "tfo"      // Tefaro
	LanguageTypeTFR    LanguageType = "tfr"      // Teribe
	LanguageTypeTFT    LanguageType = "tft"      // Ternate
	LanguageTypeTGA    LanguageType = "tga"      // Sagalla
	LanguageTypeTGB    LanguageType = "tgb"      // Tobilung
	LanguageTypeTGC    LanguageType = "tgc"      // Tigak
	LanguageTypeTGD    LanguageType = "tgd"      // Ciwogai
	LanguageTypeTGE    LanguageType = "tge"      // Eastern Gorkha Tamang
	LanguageTypeTGF    LanguageType = "tgf"      // Chalikha
	LanguageTypeTGG    LanguageType = "tgg"      // Tangga
	LanguageTypeTGH    LanguageType = "tgh"      // Tobagonian Creole English
	LanguageTypeTGI    LanguageType = "tgi"      // Lawunuia
	LanguageTypeTGJ    LanguageType = "tgj"      // Tagin
	LanguageTypeTGN    LanguageType = "tgn"      // Tandaganon
	LanguageTypeTGO    LanguageType = "tgo"      // Sudest
	LanguageTypeTGP    LanguageType = "tgp"      // Tangoa
	LanguageTypeTGQ    LanguageType = "tgq"      // Tring
	LanguageTypeTGR    LanguageType = "tgr"      // Tareng
	LanguageTypeTGS    LanguageType = "tgs"      // Nume
	LanguageTypeTGT    LanguageType = "tgt"      // Central Tagbanwa
	LanguageTypeTGU    LanguageType = "tgu"      // Tanggu
	LanguageTypeTGV    LanguageType = "tgv"      // Tingui-Boto
	LanguageTypeTGW    LanguageType = "tgw"      // Tagwana Senoufo
	LanguageTypeTGX    LanguageType = "tgx"      // Tagish
	LanguageTypeTGY    LanguageType = "tgy"      // Togoyo
	LanguageTypeTGZ    LanguageType = "tgz"      // Tagalaka
	LanguageTypeTHC    LanguageType = "thc"      // Tai Hang Tong
	LanguageTypeTHD    LanguageType = "thd"      // Thayore
	LanguageTypeTHE    LanguageType = "the"      // Chitwania Tharu
	LanguageTypeTHF    LanguageType = "thf"      // Thangmi
	LanguageTypeTHH    LanguageType = "thh"      // Northern Tarahumara
	LanguageTypeTHI    LanguageType = "thi"      // Tai Long
	LanguageTypeTHK    LanguageType = "thk"      // Tharaka and Kitharaka
	LanguageTypeTHL    LanguageType = "thl"      // Dangaura Tharu
	LanguageTypeTHM    LanguageType = "thm"      // Aheu
	LanguageTypeTHN    LanguageType = "thn"      // Thachanadan
	LanguageTypeTHP    LanguageType = "thp"      // Thompson
	LanguageTypeTHQ    LanguageType = "thq"      // Kochila Tharu
	LanguageTypeTHR    LanguageType = "thr"      // Rana Tharu
	LanguageTypeTHS    LanguageType = "ths"      // Thakali
	LanguageTypeTHT    LanguageType = "tht"      // Tahltan
	LanguageTypeTHU    LanguageType = "thu"      // Thuri
	LanguageTypeTHV    LanguageType = "thv"      // Tahaggart Tamahaq
	LanguageTypeTHW    LanguageType = "thw"      // Thudam
	LanguageTypeTHX    LanguageType = "thx"      // The
	LanguageTypeTHY    LanguageType = "thy"      // Tha
	LanguageTypeTHZ    LanguageType = "thz"      // Tayart Tamajeq
	LanguageTypeTIA    LanguageType = "tia"      // Tidikelt Tamazight
	LanguageTypeTIC    LanguageType = "tic"      // Tira
	LanguageTypeTID    LanguageType = "tid"      // Tidong
	LanguageTypeTIE    LanguageType = "tie"      // Tingal
	LanguageTypeTIF    LanguageType = "tif"      // Tifal
	LanguageTypeTIG    LanguageType = "tig"      // Tigre
	LanguageTypeTIH    LanguageType = "tih"      // Timugon Murut
	LanguageTypeTII    LanguageType = "tii"      // Tiene
	LanguageTypeTIJ    LanguageType = "tij"      // Tilung
	LanguageTypeTIK    LanguageType = "tik"      // Tikar
	LanguageTypeTIL    LanguageType = "til"      // Tillamook
	LanguageTypeTIM    LanguageType = "tim"      // Timbe
	LanguageTypeTIN    LanguageType = "tin"      // Tindi
	LanguageTypeTIO    LanguageType = "tio"      // Teop
	LanguageTypeTIP    LanguageType = "tip"      // Trimuris
	LanguageTypeTIQ    LanguageType = "tiq"      // Tiéfo
	LanguageTypeTIS    LanguageType = "tis"      // Masadiit Itneg
	LanguageTypeTIT    LanguageType = "tit"      // Tinigua
	LanguageTypeTIU    LanguageType = "tiu"      // Adasen
	LanguageTypeTIV    LanguageType = "tiv"      // Tiv
	LanguageTypeTIW    LanguageType = "tiw"      // Tiwi
	LanguageTypeTIX    LanguageType = "tix"      // Southern Tiwa
	LanguageTypeTIY    LanguageType = "tiy"      // Tiruray
	LanguageTypeTIZ    LanguageType = "tiz"      // Tai Hongjin
	LanguageTypeTJA    LanguageType = "tja"      // Tajuasohn
	LanguageTypeTJG    LanguageType = "tjg"      // Tunjung
	LanguageTypeTJI    LanguageType = "tji"      // Northern Tujia
	LanguageTypeTJL    LanguageType = "tjl"      // Tai Laing
	LanguageTypeTJM    LanguageType = "tjm"      // Timucua
	LanguageTypeTJN    LanguageType = "tjn"      // Tonjon
	LanguageTypeTJO    LanguageType = "tjo"      // Temacine Tamazight
	LanguageTypeTJS    LanguageType = "tjs"      // Southern Tujia
	LanguageTypeTJU    LanguageType = "tju"      // Tjurruru
	LanguageTypeTJW    LanguageType = "tjw"      // Djabwurrung
	LanguageTypeTKA    LanguageType = "tka"      // Truká
	LanguageTypeTKB    LanguageType = "tkb"      // Buksa
	LanguageTypeTKD    LanguageType = "tkd"      // Tukudede
	LanguageTypeTKE    LanguageType = "tke"      // Takwane
	LanguageTypeTKF    LanguageType = "tkf"      // Tukumanféd
	LanguageTypeTKG    LanguageType = "tkg"      // Tesaka Malagasy
	LanguageTypeTKK    LanguageType = "tkk"      // Takpa
	LanguageTypeTKL    LanguageType = "tkl"      // Tokelau
	LanguageTypeTKM    LanguageType = "tkm"      // Takelma
	LanguageTypeTKN    LanguageType = "tkn"      // Toku-No-Shima
	LanguageTypeTKP    LanguageType = "tkp"      // Tikopia
	LanguageTypeTKQ    LanguageType = "tkq"      // Tee
	LanguageTypeTKR    LanguageType = "tkr"      // Tsakhur
	LanguageTypeTKS    LanguageType = "tks"      // Takestani
	LanguageTypeTKT    LanguageType = "tkt"      // Kathoriya Tharu
	LanguageTypeTKU    LanguageType = "tku"      // Upper Necaxa Totonac
	LanguageTypeTKW    LanguageType = "tkw"      // Teanu
	LanguageTypeTKX    LanguageType = "tkx"      // Tangko
	LanguageTypeTKZ    LanguageType = "tkz"      // Takua
	LanguageTypeTLA    LanguageType = "tla"      // Southwestern Tepehuan
	LanguageTypeTLB    LanguageType = "tlb"      // Tobelo
	LanguageTypeTLC    LanguageType = "tlc"      // Yecuatla Totonac
	LanguageTypeTLD    LanguageType = "tld"      // Talaud
	LanguageTypeTLF    LanguageType = "tlf"      // Telefol
	LanguageTypeTLG    LanguageType = "tlg"      // Tofanma
	LanguageTypeTLH    LanguageType = "tlh"      // Klingon and tlhIngan-Hol
	LanguageTypeTLI    LanguageType = "tli"      // Tlingit
	LanguageTypeTLJ    LanguageType = "tlj"      // Talinga-Bwisi
	LanguageTypeTLK    LanguageType = "tlk"      // Taloki
	LanguageTypeTLL    LanguageType = "tll"      // Tetela
	LanguageTypeTLM    LanguageType = "tlm"      // Tolomako
	LanguageTypeTLN    LanguageType = "tln"      // Talondo'
	LanguageTypeTLO    LanguageType = "tlo"      // Talodi
	LanguageTypeTLP    LanguageType = "tlp"      // Filomena Mata-Coahuitlán Totonac
	LanguageTypeTLQ    LanguageType = "tlq"      // Tai Loi
	LanguageTypeTLR    LanguageType = "tlr"      // Talise
	LanguageTypeTLS    LanguageType = "tls"      // Tambotalo
	LanguageTypeTLT    LanguageType = "tlt"      // Teluti
	LanguageTypeTLU    LanguageType = "tlu"      // Tulehu
	LanguageTypeTLV    LanguageType = "tlv"      // Taliabu
	LanguageTypeTLW    LanguageType = "tlw"      // South Wemale
	LanguageTypeTLX    LanguageType = "tlx"      // Khehek
	LanguageTypeTLY    LanguageType = "tly"      // Talysh
	LanguageTypeTMA    LanguageType = "tma"      // Tama (Chad)
	LanguageTypeTMB    LanguageType = "tmb"      // Katbol and Avava
	LanguageTypeTMC    LanguageType = "tmc"      // Tumak
	LanguageTypeTMD    LanguageType = "tmd"      // Haruai
	LanguageTypeTME    LanguageType = "tme"      // Tremembé
	LanguageTypeTMF    LanguageType = "tmf"      // Toba-Maskoy
	LanguageTypeTMG    LanguageType = "tmg"      // Ternateño
	LanguageTypeTMH    LanguageType = "tmh"      // Tamashek
	LanguageTypeTMI    LanguageType = "tmi"      // Tutuba
	LanguageTypeTMJ    LanguageType = "tmj"      // Samarokena
	LanguageTypeTMK    LanguageType = "tmk"      // Northwestern Tamang
	LanguageTypeTML    LanguageType = "tml"      // Tamnim Citak
	LanguageTypeTMM    LanguageType = "tmm"      // Tai Thanh
	LanguageTypeTMN    LanguageType = "tmn"      // Taman (Indonesia)
	LanguageTypeTMO    LanguageType = "tmo"      // Temoq
	LanguageTypeTMP    LanguageType = "tmp"      // Tai Mène
	LanguageTypeTMQ    LanguageType = "tmq"      // Tumleo
	LanguageTypeTMR    LanguageType = "tmr"      // Jewish Babylonian Aramaic (ca. 200-1200 CE)
	LanguageTypeTMS    LanguageType = "tms"      // Tima
	LanguageTypeTMT    LanguageType = "tmt"      // Tasmate
	LanguageTypeTMU    LanguageType = "tmu"      // Iau
	LanguageTypeTMV    LanguageType = "tmv"      // Tembo (Motembo)
	LanguageTypeTMW    LanguageType = "tmw"      // Temuan
	LanguageTypeTMY    LanguageType = "tmy"      // Tami
	LanguageTypeTMZ    LanguageType = "tmz"      // Tamanaku
	LanguageTypeTNA    LanguageType = "tna"      // Tacana
	LanguageTypeTNB    LanguageType = "tnb"      // Western Tunebo
	LanguageTypeTNC    LanguageType = "tnc"      // Tanimuca-Retuarã
	LanguageTypeTND    LanguageType = "tnd"      // Angosturas Tunebo
	LanguageTypeTNE    LanguageType = "tne"      // Tinoc Kallahan
	LanguageTypeTNF    LanguageType = "tnf"      // Tangshewi
	LanguageTypeTNG    LanguageType = "tng"      // Tobanga
	LanguageTypeTNH    LanguageType = "tnh"      // Maiani
	LanguageTypeTNI    LanguageType = "tni"      // Tandia
	LanguageTypeTNK    LanguageType = "tnk"      // Kwamera
	LanguageTypeTNL    LanguageType = "tnl"      // Lenakel
	LanguageTypeTNM    LanguageType = "tnm"      // Tabla
	LanguageTypeTNN    LanguageType = "tnn"      // North Tanna
	LanguageTypeTNO    LanguageType = "tno"      // Toromono
	LanguageTypeTNP    LanguageType = "tnp"      // Whitesands
	LanguageTypeTNQ    LanguageType = "tnq"      // Taino
	LanguageTypeTNR    LanguageType = "tnr"      // Ménik
	LanguageTypeTNS    LanguageType = "tns"      // Tenis
	LanguageTypeTNT    LanguageType = "tnt"      // Tontemboan
	LanguageTypeTNU    LanguageType = "tnu"      // Tay Khang
	LanguageTypeTNV    LanguageType = "tnv"      // Tangchangya
	LanguageTypeTNW    LanguageType = "tnw"      // Tonsawang
	LanguageTypeTNX    LanguageType = "tnx"      // Tanema
	LanguageTypeTNY    LanguageType = "tny"      // Tongwe
	LanguageTypeTNZ    LanguageType = "tnz"      // Tonga (Thailand)
	LanguageTypeTOB    LanguageType = "tob"      // Toba
	LanguageTypeTOC    LanguageType = "toc"      // Coyutla Totonac
	LanguageTypeTOD    LanguageType = "tod"      // Toma
	LanguageTypeTOE    LanguageType = "toe"      // Tomedes
	LanguageTypeTOF    LanguageType = "tof"      // Gizrra
	LanguageTypeTOG    LanguageType = "tog"      // Tonga (Nyasa)
	LanguageTypeTOH    LanguageType = "toh"      // Gitonga
	LanguageTypeTOI    LanguageType = "toi"      // Tonga (Zambia)
	LanguageTypeTOJ    LanguageType = "toj"      // Tojolabal
	LanguageTypeTOL    LanguageType = "tol"      // Tolowa
	LanguageTypeTOM    LanguageType = "tom"      // Tombulu
	LanguageTypeTOO    LanguageType = "too"      // Xicotepec De Juárez Totonac
	LanguageTypeTOP    LanguageType = "top"      // Papantla Totonac
	LanguageTypeTOQ    LanguageType = "toq"      // Toposa
	LanguageTypeTOR    LanguageType = "tor"      // Togbo-Vara Banda
	LanguageTypeTOS    LanguageType = "tos"      // Highland Totonac
	LanguageTypeTOU    LanguageType = "tou"      // Tho
	LanguageTypeTOV    LanguageType = "tov"      // Upper Taromi
	LanguageTypeTOW    LanguageType = "tow"      // Jemez
	LanguageTypeTOX    LanguageType = "tox"      // Tobian
	LanguageTypeTOY    LanguageType = "toy"      // Topoiyo
	LanguageTypeTOZ    LanguageType = "toz"      // To
	LanguageTypeTPA    LanguageType = "tpa"      // Taupota
	LanguageTypeTPC    LanguageType = "tpc"      // Azoyú Me'phaa and Azoyú Tlapanec
	LanguageTypeTPE    LanguageType = "tpe"      // Tippera
	LanguageTypeTPF    LanguageType = "tpf"      // Tarpia
	LanguageTypeTPG    LanguageType = "tpg"      // Kula
	LanguageTypeTPI    LanguageType = "tpi"      // Tok Pisin
	LanguageTypeTPJ    LanguageType = "tpj"      // Tapieté
	LanguageTypeTPK    LanguageType = "tpk"      // Tupinikin
	LanguageTypeTPL    LanguageType = "tpl"      // Tlacoapa Me'phaa and Tlacoapa Tlapanec
	LanguageTypeTPM    LanguageType = "tpm"      // Tampulma
	LanguageTypeTPN    LanguageType = "tpn"      // Tupinambá
	LanguageTypeTPO    LanguageType = "tpo"      // Tai Pao
	LanguageTypeTPP    LanguageType = "tpp"      // Pisaflores Tepehua
	LanguageTypeTPQ    LanguageType = "tpq"      // Tukpa
	LanguageTypeTPR    LanguageType = "tpr"      // Tuparí
	LanguageTypeTPT    LanguageType = "tpt"      // Tlachichilco Tepehua
	LanguageTypeTPU    LanguageType = "tpu"      // Tampuan
	LanguageTypeTPV    LanguageType = "tpv"      // Tanapag
	LanguageTypeTPW    LanguageType = "tpw"      // Tupí
	LanguageTypeTPX    LanguageType = "tpx"      // Acatepec Me'phaa and Acatepec Tlapanec
	LanguageTypeTPY    LanguageType = "tpy"      // Trumai
	LanguageTypeTPZ    LanguageType = "tpz"      // Tinputz
	LanguageTypeTQB    LanguageType = "tqb"      // Tembé
	LanguageTypeTQL    LanguageType = "tql"      // Lehali
	LanguageTypeTQM    LanguageType = "tqm"      // Turumsa
	LanguageTypeTQN    LanguageType = "tqn"      // Tenino
	LanguageTypeTQO    LanguageType = "tqo"      // Toaripi
	LanguageTypeTQP    LanguageType = "tqp"      // Tomoip
	LanguageTypeTQQ    LanguageType = "tqq"      // Tunni
	LanguageTypeTQR    LanguageType = "tqr"      // Torona
	LanguageTypeTQT    LanguageType = "tqt"      // Western Totonac
	LanguageTypeTQU    LanguageType = "tqu"      // Touo
	LanguageTypeTQW    LanguageType = "tqw"      // Tonkawa
	LanguageTypeTRA    LanguageType = "tra"      // Tirahi
	LanguageTypeTRB    LanguageType = "trb"      // Terebu
	LanguageTypeTRC    LanguageType = "trc"      // Copala Triqui
	LanguageTypeTRD    LanguageType = "trd"      // Turi
	LanguageTypeTRE    LanguageType = "tre"      // East Tarangan
	LanguageTypeTRF    LanguageType = "trf"      // Trinidadian Creole English
	LanguageTypeTRG    LanguageType = "trg"      // Lishán Didán
	LanguageTypeTRH    LanguageType = "trh"      // Turaka
	LanguageTypeTRI    LanguageType = "tri"      // Trió
	LanguageTypeTRJ    LanguageType = "trj"      // Toram
	LanguageTypeTRK    LanguageType = "trk"      // Turkic languages
	LanguageTypeTRL    LanguageType = "trl"      // Traveller Scottish
	LanguageTypeTRM    LanguageType = "trm"      // Tregami
	LanguageTypeTRN    LanguageType = "trn"      // Trinitario
	LanguageTypeTRO    LanguageType = "tro"      // Tarao Naga
	LanguageTypeTRP    LanguageType = "trp"      // Kok Borok
	LanguageTypeTRQ    LanguageType = "trq"      // San Martín Itunyoso Triqui
	LanguageTypeTRR    LanguageType = "trr"      // Taushiro
	LanguageTypeTRS    LanguageType = "trs"      // Chicahuaxtla Triqui
	LanguageTypeTRT    LanguageType = "trt"      // Tunggare
	LanguageTypeTRU    LanguageType = "tru"      // Turoyo and Surayt
	LanguageTypeTRV    LanguageType = "trv"      // Taroko
	LanguageTypeTRW    LanguageType = "trw"      // Torwali
	LanguageTypeTRX    LanguageType = "trx"      // Tringgus-Sembaan Bidayuh
	LanguageTypeTRY    LanguageType = "try"      // Turung
	LanguageTypeTRZ    LanguageType = "trz"      // Torá
	LanguageTypeTSA    LanguageType = "tsa"      // Tsaangi
	LanguageTypeTSB    LanguageType = "tsb"      // Tsamai
	LanguageTypeTSC    LanguageType = "tsc"      // Tswa
	LanguageTypeTSD    LanguageType = "tsd"      // Tsakonian
	LanguageTypeTSE    LanguageType = "tse"      // Tunisian Sign Language
	LanguageTypeTSF    LanguageType = "tsf"      // Southwestern Tamang
	LanguageTypeTSG    LanguageType = "tsg"      // Tausug
	LanguageTypeTSH    LanguageType = "tsh"      // Tsuvan
	LanguageTypeTSI    LanguageType = "tsi"      // Tsimshian
	LanguageTypeTSJ    LanguageType = "tsj"      // Tshangla
	LanguageTypeTSK    LanguageType = "tsk"      // Tseku
	LanguageTypeTSL    LanguageType = "tsl"      // Ts'ün-Lao
	LanguageTypeTSM    LanguageType = "tsm"      // Turkish Sign Language and Türk İşaret Dili
	LanguageTypeTSP    LanguageType = "tsp"      // Northern Toussian
	LanguageTypeTSQ    LanguageType = "tsq"      // Thai Sign Language
	LanguageTypeTSR    LanguageType = "tsr"      // Akei
	LanguageTypeTSS    LanguageType = "tss"      // Taiwan Sign Language
	LanguageTypeTST    LanguageType = "tst"      // Tondi Songway Kiini
	LanguageTypeTSU    LanguageType = "tsu"      // Tsou
	LanguageTypeTSV    LanguageType = "tsv"      // Tsogo
	LanguageTypeTSW    LanguageType = "tsw"      // Tsishingini
	LanguageTypeTSX    LanguageType = "tsx"      // Mubami
	LanguageTypeTSY    LanguageType = "tsy"      // Tebul Sign Language
	LanguageTypeTSZ    LanguageType = "tsz"      // Purepecha
	LanguageTypeTTA    LanguageType = "tta"      // Tutelo
	LanguageTypeTTB    LanguageType = "ttb"      // Gaa
	LanguageTypeTTC    LanguageType = "ttc"      // Tektiteko
	LanguageTypeTTD    LanguageType = "ttd"      // Tauade
	LanguageTypeTTE    LanguageType = "tte"      // Bwanabwana
	LanguageTypeTTF    LanguageType = "ttf"      // Tuotomb
	LanguageTypeTTG    LanguageType = "ttg"      // Tutong
	LanguageTypeTTH    LanguageType = "tth"      // Upper Ta'oih
	LanguageTypeTTI    LanguageType = "tti"      // Tobati
	LanguageTypeTTJ    LanguageType = "ttj"      // Tooro
	LanguageTypeTTK    LanguageType = "ttk"      // Totoro
	LanguageTypeTTL    LanguageType = "ttl"      // Totela
	LanguageTypeTTM    LanguageType = "ttm"      // Northern Tutchone
	LanguageTypeTTN    LanguageType = "ttn"      // Towei
	LanguageTypeTTO    LanguageType = "tto"      // Lower Ta'oih
	LanguageTypeTTP    LanguageType = "ttp"      // Tombelala
	LanguageTypeTTQ    LanguageType = "ttq"      // Tawallammat Tamajaq
	LanguageTypeTTR    LanguageType = "ttr"      // Tera
	LanguageTypeTTS    LanguageType = "tts"      // Northeastern Thai
	LanguageTypeTTT    LanguageType = "ttt"      // Muslim Tat
	LanguageTypeTTU    LanguageType = "ttu"      // Torau
	LanguageTypeTTV    LanguageType = "ttv"      // Titan
	LanguageTypeTTW    LanguageType = "ttw"      // Long Wat
	LanguageTypeTTY    LanguageType = "tty"      // Sikaritai
	LanguageTypeTTZ    LanguageType = "ttz"      // Tsum
	LanguageTypeTUA    LanguageType = "tua"      // Wiarumus
	LanguageTypeTUB    LanguageType = "tub"      // Tübatulabal
	LanguageTypeTUC    LanguageType = "tuc"      // Mutu
	LanguageTypeTUD    LanguageType = "tud"      // Tuxá
	LanguageTypeTUE    LanguageType = "tue"      // Tuyuca
	LanguageTypeTUF    LanguageType = "tuf"      // Central Tunebo
	LanguageTypeTUG    LanguageType = "tug"      // Tunia
	LanguageTypeTUH    LanguageType = "tuh"      // Taulil
	LanguageTypeTUI    LanguageType = "tui"      // Tupuri
	LanguageTypeTUJ    LanguageType = "tuj"      // Tugutil
	LanguageTypeTUL    LanguageType = "tul"      // Tula
	LanguageTypeTUM    LanguageType = "tum"      // Tumbuka
	LanguageTypeTUN    LanguageType = "tun"      // Tunica
	LanguageTypeTUO    LanguageType = "tuo"      // Tucano
	LanguageTypeTUP    LanguageType = "tup"      // Tupi languages
	LanguageTypeTUQ    LanguageType = "tuq"      // Tedaga
	LanguageTypeTUS    LanguageType = "tus"      // Tuscarora
	LanguageTypeTUT    LanguageType = "tut"      // Altaic languages
	LanguageTypeTUU    LanguageType = "tuu"      // Tututni
	LanguageTypeTUV    LanguageType = "tuv"      // Turkana
	LanguageTypeTUW    LanguageType = "tuw"      // Tungus languages
	LanguageTypeTUX    LanguageType = "tux"      // Tuxináwa
	LanguageTypeTUY    LanguageType = "tuy"      // Tugen
	LanguageTypeTUZ    LanguageType = "tuz"      // Turka
	LanguageTypeTVA    LanguageType = "tva"      // Vaghua
	LanguageTypeTVD    LanguageType = "tvd"      // Tsuvadi
	LanguageTypeTVE    LanguageType = "tve"      // Te'un
	LanguageTypeTVK    LanguageType = "tvk"      // Southeast Ambrym
	LanguageTypeTVL    LanguageType = "tvl"      // Tuvalu
	LanguageTypeTVM    LanguageType = "tvm"      // Tela-Masbuar
	LanguageTypeTVN    LanguageType = "tvn"      // Tavoyan
	LanguageTypeTVO    LanguageType = "tvo"      // Tidore
	LanguageTypeTVS    LanguageType = "tvs"      // Taveta
	LanguageTypeTVT    LanguageType = "tvt"      // Tutsa Naga
	LanguageTypeTVU    LanguageType = "tvu"      // Tunen
	LanguageTypeTVW    LanguageType = "tvw"      // Sedoa
	LanguageTypeTVY    LanguageType = "tvy"      // Timor Pidgin
	LanguageTypeTWA    LanguageType = "twa"      // Twana
	LanguageTypeTWB    LanguageType = "twb"      // Western Tawbuid
	LanguageTypeTWC    LanguageType = "twc"      // Teshenawa
	LanguageTypeTWD    LanguageType = "twd"      // Twents
	LanguageTypeTWE    LanguageType = "twe"      // Tewa (Indonesia)
	LanguageTypeTWF    LanguageType = "twf"      // Northern Tiwa
	LanguageTypeTWG    LanguageType = "twg"      // Tereweng
	LanguageTypeTWH    LanguageType = "twh"      // Tai Dón
	LanguageTypeTWL    LanguageType = "twl"      // Tawara
	LanguageTypeTWM    LanguageType = "twm"      // Tawang Monpa
	LanguageTypeTWN    LanguageType = "twn"      // Twendi
	LanguageTypeTWO    LanguageType = "two"      // Tswapong
	LanguageTypeTWP    LanguageType = "twp"      // Ere
	LanguageTypeTWQ    LanguageType = "twq"      // Tasawaq
	LanguageTypeTWR    LanguageType = "twr"      // Southwestern Tarahumara
	LanguageTypeTWT    LanguageType = "twt"      // Turiwára
	LanguageTypeTWU    LanguageType = "twu"      // Termanu
	LanguageTypeTWW    LanguageType = "tww"      // Tuwari
	LanguageTypeTWX    LanguageType = "twx"      // Tewe
	LanguageTypeTWY    LanguageType = "twy"      // Tawoyan
	LanguageTypeTXA    LanguageType = "txa"      // Tombonuo
	LanguageTypeTXB    LanguageType = "txb"      // Tokharian B
	LanguageTypeTXC    LanguageType = "txc"      // Tsetsaut
	LanguageTypeTXE    LanguageType = "txe"      // Totoli
	LanguageTypeTXG    LanguageType = "txg"      // Tangut
	LanguageTypeTXH    LanguageType = "txh"      // Thracian
	LanguageTypeTXI    LanguageType = "txi"      // Ikpeng
	LanguageTypeTXM    LanguageType = "txm"      // Tomini
	LanguageTypeTXN    LanguageType = "txn"      // West Tarangan
	LanguageTypeTXO    LanguageType = "txo"      // Toto
	LanguageTypeTXQ    LanguageType = "txq"      // Tii
	LanguageTypeTXR    LanguageType = "txr"      // Tartessian
	LanguageTypeTXS    LanguageType = "txs"      // Tonsea
	LanguageTypeTXT    LanguageType = "txt"      // Citak
	LanguageTypeTXU    LanguageType = "txu"      // Kayapó
	LanguageTypeTXX    LanguageType = "txx"      // Tatana
	LanguageTypeTXY    LanguageType = "txy"      // Tanosy Malagasy
	LanguageTypeTYA    LanguageType = "tya"      // Tauya
	LanguageTypeTYE    LanguageType = "tye"      // Kyanga
	LanguageTypeTYH    LanguageType = "tyh"      // O'du
	LanguageTypeTYI    LanguageType = "tyi"      // Teke-Tsaayi
	LanguageTypeTYJ    LanguageType = "tyj"      // Tai Do
	LanguageTypeTYL    LanguageType = "tyl"      // Thu Lao
	LanguageTypeTYN    LanguageType = "tyn"      // Kombai
	LanguageTypeTYP    LanguageType = "typ"      // Thaypan
	LanguageTypeTYR    LanguageType = "tyr"      // Tai Daeng
	LanguageTypeTYS    LanguageType = "tys"      // Tày Sa Pa
	LanguageTypeTYT    LanguageType = "tyt"      // Tày Tac
	LanguageTypeTYU    LanguageType = "tyu"      // Kua
	LanguageTypeTYV    LanguageType = "tyv"      // Tuvinian
	LanguageTypeTYX    LanguageType = "tyx"      // Teke-Tyee
	LanguageTypeTYZ    LanguageType = "tyz"      // Tày
	LanguageTypeTZA    LanguageType = "tza"      // Tanzanian Sign Language
	LanguageTypeTZH    LanguageType = "tzh"      // Tzeltal
	LanguageTypeTZJ    LanguageType = "tzj"      // Tz'utujil
	LanguageTypeTZL    LanguageType = "tzl"      // Talossan
	LanguageTypeTZM    LanguageType = "tzm"      // Central Atlas Tamazight
	LanguageTypeTZN    LanguageType = "tzn"      // Tugun
	LanguageTypeTZO    LanguageType = "tzo"      // Tzotzil
	LanguageTypeTZX    LanguageType = "tzx"      // Tabriak
	LanguageTypeUAM    LanguageType = "uam"      // Uamué
	LanguageTypeUAN    LanguageType = "uan"      // Kuan
	LanguageTypeUAR    LanguageType = "uar"      // Tairuma
	LanguageTypeUBA    LanguageType = "uba"      // Ubang
	LanguageTypeUBI    LanguageType = "ubi"      // Ubi
	LanguageTypeUBL    LanguageType = "ubl"      // Buhi'non Bikol
	LanguageTypeUBR    LanguageType = "ubr"      // Ubir
	LanguageTypeUBU    LanguageType = "ubu"      // Umbu-Ungu
	LanguageTypeUBY    LanguageType = "uby"      // Ubykh
	LanguageTypeUDA    LanguageType = "uda"      // Uda
	LanguageTypeUDE    LanguageType = "ude"      // Udihe
	LanguageTypeUDG    LanguageType = "udg"      // Muduga
	LanguageTypeUDI    LanguageType = "udi"      // Udi
	LanguageTypeUDJ    LanguageType = "udj"      // Ujir
	LanguageTypeUDL    LanguageType = "udl"      // Wuzlam
	LanguageTypeUDM    LanguageType = "udm"      // Udmurt
	LanguageTypeUDU    LanguageType = "udu"      // Uduk
	LanguageTypeUES    LanguageType = "ues"      // Kioko
	LanguageTypeUFI    LanguageType = "ufi"      // Ufim
	LanguageTypeUGA    LanguageType = "uga"      // Ugaritic
	LanguageTypeUGB    LanguageType = "ugb"      // Kuku-Ugbanh
	LanguageTypeUGE    LanguageType = "uge"      // Ughele
	LanguageTypeUGN    LanguageType = "ugn"      // Ugandan Sign Language
	LanguageTypeUGO    LanguageType = "ugo"      // Ugong
	LanguageTypeUGY    LanguageType = "ugy"      // Uruguayan Sign Language
	LanguageTypeUHA    LanguageType = "uha"      // Uhami
	LanguageTypeUHN    LanguageType = "uhn"      // Damal
	LanguageTypeUIS    LanguageType = "uis"      // Uisai
	LanguageTypeUIV    LanguageType = "uiv"      // Iyive
	LanguageTypeUJI    LanguageType = "uji"      // Tanjijili
	LanguageTypeUKA    LanguageType = "uka"      // Kaburi
	LanguageTypeUKG    LanguageType = "ukg"      // Ukuriguma
	LanguageTypeUKH    LanguageType = "ukh"      // Ukhwejo
	LanguageTypeUKL    LanguageType = "ukl"      // Ukrainian Sign Language
	LanguageTypeUKP    LanguageType = "ukp"      // Ukpe-Bayobiri
	LanguageTypeUKQ    LanguageType = "ukq"      // Ukwa
	LanguageTypeUKS    LanguageType = "uks"      // Urubú-Kaapor Sign Language and Kaapor Sign Language
	LanguageTypeUKU    LanguageType = "uku"      // Ukue
	LanguageTypeUKW    LanguageType = "ukw"      // Ukwuani-Aboh-Ndoni
	LanguageTypeUKY    LanguageType = "uky"      // Kuuk-Yak
	LanguageTypeULA    LanguageType = "ula"      // Fungwa
	LanguageTypeULB    LanguageType = "ulb"      // Ulukwumi
	LanguageTypeULC    LanguageType = "ulc"      // Ulch
	LanguageTypeULE    LanguageType = "ule"      // Lule
	LanguageTypeULF    LanguageType = "ulf"      // Usku and Afra
	LanguageTypeULI    LanguageType = "uli"      // Ulithian
	LanguageTypeULK    LanguageType = "ulk"      // Meriam
	LanguageTypeULL    LanguageType = "ull"      // Ullatan
	LanguageTypeULM    LanguageType = "ulm"      // Ulumanda'
	LanguageTypeULN    LanguageType = "uln"      // Unserdeutsch
	LanguageTypeULU    LanguageType = "ulu"      // Uma' Lung
	LanguageTypeULW    LanguageType = "ulw"      // Ulwa
	LanguageTypeUMA    LanguageType = "uma"      // Umatilla
	LanguageTypeUMB    LanguageType = "umb"      // Umbundu
	LanguageTypeUMC    LanguageType = "umc"      // Marrucinian
	LanguageTypeUMD    LanguageType = "umd"      // Umbindhamu
	LanguageTypeUMG    LanguageType = "umg"      // Umbuygamu
	LanguageTypeUMI    LanguageType = "umi"      // Ukit
	LanguageTypeUMM    LanguageType = "umm"      // Umon
	LanguageTypeUMN    LanguageType = "umn"      // Makyan Naga
	LanguageTypeUMO    LanguageType = "umo"      // Umotína
	LanguageTypeUMP    LanguageType = "ump"      // Umpila
	LanguageTypeUMR    LanguageType = "umr"      // Umbugarla
	LanguageTypeUMS    LanguageType = "ums"      // Pendau
	LanguageTypeUMU    LanguageType = "umu"      // Munsee
	LanguageTypeUNA    LanguageType = "una"      // North Watut
	LanguageTypeUND    LanguageType = "und"      // Undetermined
	LanguageTypeUNE    LanguageType = "une"      // Uneme
	LanguageTypeUNG    LanguageType = "ung"      // Ngarinyin
	LanguageTypeUNK    LanguageType = "unk"      // Enawené-Nawé
	LanguageTypeUNM    LanguageType = "unm"      // Unami
	LanguageTypeUNN    LanguageType = "unn"      // Kurnai
	LanguageTypeUNP    LanguageType = "unp"      // Worora
	LanguageTypeUNR    LanguageType = "unr"      // Mundari
	LanguageTypeUNU    LanguageType = "unu"      // Unubahe
	LanguageTypeUNX    LanguageType = "unx"      // Munda
	LanguageTypeUNZ    LanguageType = "unz"      // Unde Kaili
	LanguageTypeUOK    LanguageType = "uok"      // Uokha
	LanguageTypeUPI    LanguageType = "upi"      // Umeda
	LanguageTypeUPV    LanguageType = "upv"      // Uripiv-Wala-Rano-Atchin
	LanguageTypeURA    LanguageType = "ura"      // Urarina
	LanguageTypeURB    LanguageType = "urb"      // Urubú-Kaapor and Kaapor
	LanguageTypeURC    LanguageType = "urc"      // Urningangg
	LanguageTypeURE    LanguageType = "ure"      // Uru
	LanguageTypeURF    LanguageType = "urf"      // Uradhi
	LanguageTypeURG    LanguageType = "urg"      // Urigina
	LanguageTypeURH    LanguageType = "urh"      // Urhobo
	LanguageTypeURI    LanguageType = "uri"      // Urim
	LanguageTypeURJ    LanguageType = "urj"      // Uralic languages
	LanguageTypeURK    LanguageType = "urk"      // Urak Lawoi'
	LanguageTypeURL    LanguageType = "url"      // Urali
	LanguageTypeURM    LanguageType = "urm"      // Urapmin
	LanguageTypeURN    LanguageType = "urn"      // Uruangnirin
	LanguageTypeURO    LanguageType = "uro"      // Ura (Papua New Guinea)
	LanguageTypeURP    LanguageType = "urp"      // Uru-Pa-In
	LanguageTypeURR    LanguageType = "urr"      // Lehalurup and Löyöp
	LanguageTypeURT    LanguageType = "urt"      // Urat
	LanguageTypeURU    LanguageType = "uru"      // Urumi
	LanguageTypeURV    LanguageType = "urv"      // Uruava
	LanguageTypeURW    LanguageType = "urw"      // Sop
	LanguageTypeURX    LanguageType = "urx"      // Urimo
	LanguageTypeURY    LanguageType = "ury"      // Orya
	LanguageTypeURZ    LanguageType = "urz"      // Uru-Eu-Wau-Wau
	LanguageTypeUSA    LanguageType = "usa"      // Usarufa
	LanguageTypeUSH    LanguageType = "ush"      // Ushojo
	LanguageTypeUSI    LanguageType = "usi"      // Usui
	LanguageTypeUSK    LanguageType = "usk"      // Usaghade
	LanguageTypeUSP    LanguageType = "usp"      // Uspanteco
	LanguageTypeUSU    LanguageType = "usu"      // Uya
	LanguageTypeUTA    LanguageType = "uta"      // Otank
	LanguageTypeUTE    LanguageType = "ute"      // Ute-Southern Paiute
	LanguageTypeUTP    LanguageType = "utp"      // Amba (Solomon Islands)
	LanguageTypeUTR    LanguageType = "utr"      // Etulo
	LanguageTypeUTU    LanguageType = "utu"      // Utu
	LanguageTypeUUM    LanguageType = "uum"      // Urum
	LanguageTypeUUN    LanguageType = "uun"      // Kulon-Pazeh
	LanguageTypeUUR    LanguageType = "uur"      // Ura (Vanuatu)
	LanguageTypeUUU    LanguageType = "uuu"      // U
	LanguageTypeUVE    LanguageType = "uve"      // West Uvean and Fagauvea
	LanguageTypeUVH    LanguageType = "uvh"      // Uri
	LanguageTypeUVL    LanguageType = "uvl"      // Lote
	LanguageTypeUWA    LanguageType = "uwa"      // Kuku-Uwanh
	LanguageTypeUYA    LanguageType = "uya"      // Doko-Uyanga
	LanguageTypeUZN    LanguageType = "uzn"      // Northern Uzbek
	LanguageTypeUZS    LanguageType = "uzs"      // Southern Uzbek
	LanguageTypeVAA    LanguageType = "vaa"      // Vaagri Booli
	LanguageTypeVAE    LanguageType = "vae"      // Vale
	LanguageTypeVAF    LanguageType = "vaf"      // Vafsi
	LanguageTypeVAG    LanguageType = "vag"      // Vagla
	LanguageTypeVAH    LanguageType = "vah"      // Varhadi-Nagpuri
	LanguageTypeVAI    LanguageType = "vai"      // Vai
	LanguageTypeVAJ    LanguageType = "vaj"      // Vasekela Bushman
	LanguageTypeVAL    LanguageType = "val"      // Vehes
	LanguageTypeVAM    LanguageType = "vam"      // Vanimo
	LanguageTypeVAN    LanguageType = "van"      // Valman
	LanguageTypeVAO    LanguageType = "vao"      // Vao
	LanguageTypeVAP    LanguageType = "vap"      // Vaiphei
	LanguageTypeVAR    LanguageType = "var"      // Huarijio
	LanguageTypeVAS    LanguageType = "vas"      // Vasavi
	LanguageTypeVAU    LanguageType = "vau"      // Vanuma
	LanguageTypeVAV    LanguageType = "vav"      // Varli
	LanguageTypeVAY    LanguageType = "vay"      // Wayu
	LanguageTypeVBB    LanguageType = "vbb"      // Southeast Babar
	LanguageTypeVBK    LanguageType = "vbk"      // Southwestern Bontok
	LanguageTypeVEC    LanguageType = "vec"      // Venetian
	LanguageTypeVED    LanguageType = "ved"      // Veddah
	LanguageTypeVEL    LanguageType = "vel"      // Veluws
	LanguageTypeVEM    LanguageType = "vem"      // Vemgo-Mabas
	LanguageTypeVEO    LanguageType = "veo"      // Ventureño
	LanguageTypeVEP    LanguageType = "vep"      // Veps
	LanguageTypeVER    LanguageType = "ver"      // Mom Jango
	LanguageTypeVGR    LanguageType = "vgr"      // Vaghri
	LanguageTypeVGT    LanguageType = "vgt"      // Vlaamse Gebarentaal and Flemish Sign Language
	LanguageTypeVIC    LanguageType = "vic"      // Virgin Islands Creole English
	LanguageTypeVID    LanguageType = "vid"      // Vidunda
	LanguageTypeVIF    LanguageType = "vif"      // Vili
	LanguageTypeVIG    LanguageType = "vig"      // Viemo
	LanguageTypeVIL    LanguageType = "vil"      // Vilela
	LanguageTypeVIN    LanguageType = "vin"      // Vinza
	LanguageTypeVIS    LanguageType = "vis"      // Vishavan
	LanguageTypeVIT    LanguageType = "vit"      // Viti
	LanguageTypeVIV    LanguageType = "viv"      // Iduna
	LanguageTypeVKA    LanguageType = "vka"      // Kariyarra
	LanguageTypeVKI    LanguageType = "vki"      // Ija-Zuba
	LanguageTypeVKJ    LanguageType = "vkj"      // Kujarge
	LanguageTypeVKK    LanguageType = "vkk"      // Kaur
	LanguageTypeVKL    LanguageType = "vkl"      // Kulisusu
	LanguageTypeVKM    LanguageType = "vkm"      // Kamakan
	LanguageTypeVKO    LanguageType = "vko"      // Kodeoha
	LanguageTypeVKP    LanguageType = "vkp"      // Korlai Creole Portuguese
	LanguageTypeVKT    LanguageType = "vkt"      // Tenggarong Kutai Malay
	LanguageTypeVKU    LanguageType = "vku"      // Kurrama
	LanguageTypeVLP    LanguageType = "vlp"      // Valpei
	LanguageTypeVLS    LanguageType = "vls"      // Vlaams
	LanguageTypeVMA    LanguageType = "vma"      // Martuyhunira
	LanguageTypeVMB    LanguageType = "vmb"      // Barbaram
	LanguageTypeVMC    LanguageType = "vmc"      // Juxtlahuaca Mixtec
	LanguageTypeVMD    LanguageType = "vmd"      // Mudu Koraga
	LanguageTypeVME    LanguageType = "vme"      // East Masela
	LanguageTypeVMF    LanguageType = "vmf"      // Mainfränkisch
	LanguageTypeVMG    LanguageType = "vmg"      // Lungalunga
	LanguageTypeVMH    LanguageType = "vmh"      // Maraghei
	LanguageTypeVMI    LanguageType = "vmi"      // Miwa
	LanguageTypeVMJ    LanguageType = "vmj"      // Ixtayutla Mixtec
	LanguageTypeVMK    LanguageType = "vmk"      // Makhuwa-Shirima
	LanguageTypeVML    LanguageType = "vml"      // Malgana
	LanguageTypeVMM    LanguageType = "vmm"      // Mitlatongo Mixtec
	LanguageTypeVMP    LanguageType = "vmp"      // Soyaltepec Mazatec
	LanguageTypeVMQ    LanguageType = "vmq"      // Soyaltepec Mixtec
	LanguageTypeVMR    LanguageType = "vmr"      // Marenje
	LanguageTypeVMS    LanguageType = "vms"      // Moksela
	LanguageTypeVMU    LanguageType = "vmu"      // Muluridyi
	LanguageTypeVMV    LanguageType = "vmv"      // Valley Maidu
	LanguageTypeVMW    LanguageType = "vmw"      // Makhuwa
	LanguageTypeVMX    LanguageType = "vmx"      // Tamazola Mixtec
	LanguageTypeVMY    LanguageType = "vmy"      // Ayautla Mazatec
	LanguageTypeVMZ    LanguageType = "vmz"      // Mazatlán Mazatec
	LanguageTypeVNK    LanguageType = "vnk"      // Vano and Lovono
	LanguageTypeVNM    LanguageType = "vnm"      // Vinmavis and Neve'ei
	LanguageTypeVNP    LanguageType = "vnp"      // Vunapu
	LanguageTypeVOR    LanguageType = "vor"      // Voro
	LanguageTypeVOT    LanguageType = "vot"      // Votic
	LanguageTypeVRA    LanguageType = "vra"      // Vera'a
	LanguageTypeVRO    LanguageType = "vro"      // Võro
	LanguageTypeVRS    LanguageType = "vrs"      // Varisi
	LanguageTypeVRT    LanguageType = "vrt"      // Burmbar and Banam Bay
	LanguageTypeVSI    LanguageType = "vsi"      // Moldova Sign Language
	LanguageTypeVSL    LanguageType = "vsl"      // Venezuelan Sign Language
	LanguageTypeVSV    LanguageType = "vsv"      // Valencian Sign Language and Llengua de signes valenciana
	LanguageTypeVTO    LanguageType = "vto"      // Vitou
	LanguageTypeVUM    LanguageType = "vum"      // Vumbu
	LanguageTypeVUN    LanguageType = "vun"      // Vunjo
	LanguageTypeVUT    LanguageType = "vut"      // Vute
	LanguageTypeVWA    LanguageType = "vwa"      // Awa (China)
	LanguageTypeWAA    LanguageType = "waa"      // Walla Walla
	LanguageTypeWAB    LanguageType = "wab"      // Wab
	LanguageTypeWAC    LanguageType = "wac"      // Wasco-Wishram
	LanguageTypeWAD    LanguageType = "wad"      // Wandamen
	LanguageTypeWAE    LanguageType = "wae"      // Walser
	LanguageTypeWAF    LanguageType = "waf"      // Wakoná
	LanguageTypeWAG    LanguageType = "wag"      // Wa'ema
	LanguageTypeWAH    LanguageType = "wah"      // Watubela
	LanguageTypeWAI    LanguageType = "wai"      // Wares
	LanguageTypeWAJ    LanguageType = "waj"      // Waffa
	LanguageTypeWAK    LanguageType = "wak"      // Wakashan languages
	LanguageTypeWAL    LanguageType = "wal"      // Wolaytta and Wolaitta
	LanguageTypeWAM    LanguageType = "wam"      // Wampanoag
	LanguageTypeWAN    LanguageType = "wan"      // Wan
	LanguageTypeWAO    LanguageType = "wao"      // Wappo
	LanguageTypeWAP    LanguageType = "wap"      // Wapishana
	LanguageTypeWAQ    LanguageType = "waq"      // Wageman
	LanguageTypeWAR    LanguageType = "war"      // Waray (Philippines)
	LanguageTypeWAS    LanguageType = "was"      // Washo
	LanguageTypeWAT    LanguageType = "wat"      // Kaninuwa
	LanguageTypeWAU    LanguageType = "wau"      // Waurá
	LanguageTypeWAV    LanguageType = "wav"      // Waka
	LanguageTypeWAW    LanguageType = "waw"      // Waiwai
	LanguageTypeWAX    LanguageType = "wax"      // Watam and Marangis
	LanguageTypeWAY    LanguageType = "way"      // Wayana
	LanguageTypeWAZ    LanguageType = "waz"      // Wampur
	LanguageTypeWBA    LanguageType = "wba"      // Warao
	LanguageTypeWBB    LanguageType = "wbb"      // Wabo
	LanguageTypeWBE    LanguageType = "wbe"      // Waritai
	LanguageTypeWBF    LanguageType = "wbf"      // Wara
	LanguageTypeWBH    LanguageType = "wbh"      // Wanda
	LanguageTypeWBI    LanguageType = "wbi"      // Vwanji
	LanguageTypeWBJ    LanguageType = "wbj"      // Alagwa
	LanguageTypeWBK    LanguageType = "wbk"      // Waigali
	LanguageTypeWBL    LanguageType = "wbl"      // Wakhi
	LanguageTypeWBM    LanguageType = "wbm"      // Wa
	LanguageTypeWBP    LanguageType = "wbp"      // Warlpiri
	LanguageTypeWBQ    LanguageType = "wbq"      // Waddar
	LanguageTypeWBR    LanguageType = "wbr"      // Wagdi
	LanguageTypeWBT    LanguageType = "wbt"      // Wanman
	LanguageTypeWBV    LanguageType = "wbv"      // Wajarri
	LanguageTypeWBW    LanguageType = "wbw"      // Woi
	LanguageTypeWCA    LanguageType = "wca"      // Yanomámi
	LanguageTypeWCI    LanguageType = "wci"      // Waci Gbe
	LanguageTypeWDD    LanguageType = "wdd"      // Wandji
	LanguageTypeWDG    LanguageType = "wdg"      // Wadaginam
	LanguageTypeWDJ    LanguageType = "wdj"      // Wadjiginy
	LanguageTypeWDK    LanguageType = "wdk"      // Wadikali
	LanguageTypeWDU    LanguageType = "wdu"      // Wadjigu
	LanguageTypeWDY    LanguageType = "wdy"      // Wadjabangayi
	LanguageTypeWEA    LanguageType = "wea"      // Wewaw
	LanguageTypeWEC    LanguageType = "wec"      // Wè Western
	LanguageTypeWED    LanguageType = "wed"      // Wedau
	LanguageTypeWEG    LanguageType = "weg"      // Wergaia
	LanguageTypeWEH    LanguageType = "weh"      // Weh
	LanguageTypeWEI    LanguageType = "wei"      // Kiunum
	LanguageTypeWEM    LanguageType = "wem"      // Weme Gbe
	LanguageTypeWEN    LanguageType = "wen"      // Sorbian languages
	LanguageTypeWEO    LanguageType = "weo"      // Wemale
	LanguageTypeWEP    LanguageType = "wep"      // Westphalien
	LanguageTypeWER    LanguageType = "wer"      // Weri
	LanguageTypeWES    LanguageType = "wes"      // Cameroon Pidgin
	LanguageTypeWET    LanguageType = "wet"      // Perai
	LanguageTypeWEU    LanguageType = "weu"      // Rawngtu Chin
	LanguageTypeWEW    LanguageType = "wew"      // Wejewa
	LanguageTypeWFG    LanguageType = "wfg"      // Yafi and Zorop
	LanguageTypeWGA    LanguageType = "wga"      // Wagaya
	LanguageTypeWGB    LanguageType = "wgb"      // Wagawaga
	LanguageTypeWGG    LanguageType = "wgg"      // Wangganguru
	LanguageTypeWGI    LanguageType = "wgi"      // Wahgi
	LanguageTypeWGO    LanguageType = "wgo"      // Waigeo
	LanguageTypeWGU    LanguageType = "wgu"      // Wirangu
	LanguageTypeWGW    LanguageType = "wgw"      // Wagawaga
	LanguageTypeWGY    LanguageType = "wgy"      // Warrgamay
	LanguageTypeWHA    LanguageType = "wha"      // Manusela
	LanguageTypeWHG    LanguageType = "whg"      // North Wahgi
	LanguageTypeWHK    LanguageType = "whk"      // Wahau Kenyah
	LanguageTypeWHU    LanguageType = "whu"      // Wahau Kayan
	LanguageTypeWIB    LanguageType = "wib"      // Southern Toussian
	LanguageTypeWIC    LanguageType = "wic"      // Wichita
	LanguageTypeWIE    LanguageType = "wie"      // Wik-Epa
	LanguageTypeWIF    LanguageType = "wif"      // Wik-Keyangan
	LanguageTypeWIG    LanguageType = "wig"      // Wik-Ngathana
	LanguageTypeWIH    LanguageType = "wih"      // Wik-Me'anha
	LanguageTypeWII    LanguageType = "wii"      // Minidien
	LanguageTypeWIJ    LanguageType = "wij"      // Wik-Iiyanh
	LanguageTypeWIK    LanguageType = "wik"      // Wikalkan
	LanguageTypeWIL    LanguageType = "wil"      // Wilawila
	LanguageTypeWIM    LanguageType = "wim"      // Wik-Mungkan
	LanguageTypeWIN    LanguageType = "win"      // Ho-Chunk
	LanguageTypeWIR    LanguageType = "wir"      // Wiraféd
	LanguageTypeWIT    LanguageType = "wit"      // Wintu
	LanguageTypeWIU    LanguageType = "wiu"      // Wiru
	LanguageTypeWIV    LanguageType = "wiv"      // Vitu
	LanguageTypeWIW    LanguageType = "wiw"      // Wirangu
	LanguageTypeWIY    LanguageType = "wiy"      // Wiyot
	LanguageTypeWJA    LanguageType = "wja"      // Waja
	LanguageTypeWJI    LanguageType = "wji"      // Warji
	LanguageTypeWKA    LanguageType = "wka"      // Kw'adza
	LanguageTypeWKB    LanguageType = "wkb"      // Kumbaran
	LanguageTypeWKD    LanguageType = "wkd"      // Wakde and Mo
	LanguageTypeWKL    LanguageType = "wkl"      // Kalanadi
	LanguageTypeWKU    LanguageType = "wku"      // Kunduvadi
	LanguageTypeWKW    LanguageType = "wkw"      // Wakawaka
	LanguageTypeWKY    LanguageType = "wky"      // Wangkayutyuru
	LanguageTypeWLA    LanguageType = "wla"      // Walio
	LanguageTypeWLC    LanguageType = "wlc"      // Mwali Comorian
	LanguageTypeWLE    LanguageType = "wle"      // Wolane
	LanguageTypeWLG    LanguageType = "wlg"      // Kunbarlang
	LanguageTypeWLI    LanguageType = "wli"      // Waioli
	LanguageTypeWLK    LanguageType = "wlk"      // Wailaki
	LanguageTypeWLL    LanguageType = "wll"      // Wali (Sudan)
	LanguageTypeWLM    LanguageType = "wlm"      // Middle Welsh
	LanguageTypeWLO    LanguageType = "wlo"      // Wolio
	LanguageTypeWLR    LanguageType = "wlr"      // Wailapa
	LanguageTypeWLS    LanguageType = "wls"      // Wallisian
	LanguageTypeWLU    LanguageType = "wlu"      // Wuliwuli
	LanguageTypeWLV    LanguageType = "wlv"      // Wichí Lhamtés Vejoz
	LanguageTypeWLW    LanguageType = "wlw"      // Walak
	LanguageTypeWLX    LanguageType = "wlx"      // Wali (Ghana)
	LanguageTypeWLY    LanguageType = "wly"      // Waling
	LanguageTypeWMA    LanguageType = "wma"      // Mawa (Nigeria)
	LanguageTypeWMB    LanguageType = "wmb"      // Wambaya
	LanguageTypeWMC    LanguageType = "wmc"      // Wamas
	LanguageTypeWMD    LanguageType = "wmd"      // Mamaindé
	LanguageTypeWME    LanguageType = "wme"      // Wambule
	LanguageTypeWMH    LanguageType = "wmh"      // Waima'a
	LanguageTypeWMI    LanguageType = "wmi"      // Wamin
	LanguageTypeWMM    LanguageType = "wmm"      // Maiwa (Indonesia)
	LanguageTypeWMN    LanguageType = "wmn"      // Waamwang
	LanguageTypeWMO    LanguageType = "wmo"      // Wom (Papua New Guinea)
	LanguageTypeWMS    LanguageType = "wms"      // Wambon
	LanguageTypeWMT    LanguageType = "wmt"      // Walmajarri
	LanguageTypeWMW    LanguageType = "wmw"      // Mwani
	LanguageTypeWMX    LanguageType = "wmx"      // Womo
	LanguageTypeWNB    LanguageType = "wnb"      // Wanambre
	LanguageTypeWNC    LanguageType = "wnc"      // Wantoat
	LanguageTypeWND    LanguageType = "wnd"      // Wandarang
	LanguageTypeWNE    LanguageType = "wne"      // Waneci
	LanguageTypeWNG    LanguageType = "wng"      // Wanggom
	LanguageTypeWNI    LanguageType = "wni"      // Ndzwani Comorian
	LanguageTypeWNK    LanguageType = "wnk"      // Wanukaka
	LanguageTypeWNM    LanguageType = "wnm"      // Wanggamala
	LanguageTypeWNN    LanguageType = "wnn"      // Wunumara
	LanguageTypeWNO    LanguageType = "wno"      // Wano
	LanguageTypeWNP    LanguageType = "wnp"      // Wanap
	LanguageTypeWNU    LanguageType = "wnu"      // Usan
	LanguageTypeWNW    LanguageType = "wnw"      // Wintu
	LanguageTypeWNY    LanguageType = "wny"      // Wanyi
	LanguageTypeWOA    LanguageType = "woa"      // Tyaraity
	LanguageTypeWOB    LanguageType = "wob"      // Wè Northern
	LanguageTypeWOC    LanguageType = "woc"      // Wogeo
	LanguageTypeWOD    LanguageType = "wod"      // Wolani
	LanguageTypeWOE    LanguageType = "woe"      // Woleaian
	LanguageTypeWOF    LanguageType = "wof"      // Gambian Wolof
	LanguageTypeWOG    LanguageType = "wog"      // Wogamusin
	LanguageTypeWOI    LanguageType = "woi"      // Kamang
	LanguageTypeWOK    LanguageType = "wok"      // Longto
	LanguageTypeWOM    LanguageType = "wom"      // Wom (Nigeria)
	LanguageTypeWON    LanguageType = "won"      // Wongo
	LanguageTypeWOO    LanguageType = "woo"      // Manombai
	LanguageTypeWOR    LanguageType = "wor"      // Woria
	LanguageTypeWOS    LanguageType = "wos"      // Hanga Hundi
	LanguageTypeWOW    LanguageType = "wow"      // Wawonii
	LanguageTypeWOY    LanguageType = "woy"      // Weyto
	LanguageTypeWPC    LanguageType = "wpc"      // Maco
	LanguageTypeWRA    LanguageType = "wra"      // Warapu
	LanguageTypeWRB    LanguageType = "wrb"      // Warluwara
	LanguageTypeWRD    LanguageType = "wrd"      // Warduji
	LanguageTypeWRG    LanguageType = "wrg"      // Warungu
	LanguageTypeWRH    LanguageType = "wrh"      // Wiradhuri
	LanguageTypeWRI    LanguageType = "wri"      // Wariyangga
	LanguageTypeWRK    LanguageType = "wrk"      // Garrwa
	LanguageTypeWRL    LanguageType = "wrl"      // Warlmanpa
	LanguageTypeWRM    LanguageType = "wrm"      // Warumungu
	LanguageTypeWRN    LanguageType = "wrn"      // Warnang
	LanguageTypeWRO    LanguageType = "wro"      // Worrorra
	LanguageTypeWRP    LanguageType = "wrp"      // Waropen
	LanguageTypeWRR    LanguageType = "wrr"      // Wardaman
	LanguageTypeWRS    LanguageType = "wrs"      // Waris
	LanguageTypeWRU    LanguageType = "wru"      // Waru
	LanguageTypeWRV    LanguageType = "wrv"      // Waruna
	LanguageTypeWRW    LanguageType = "wrw"      // Gugu Warra
	LanguageTypeWRX    LanguageType = "wrx"      // Wae Rana
	LanguageTypeWRY    LanguageType = "wry"      // Merwari
	LanguageTypeWRZ    LanguageType = "wrz"      // Waray (Australia)
	LanguageTypeWSA    LanguageType = "wsa"      // Warembori
	LanguageTypeWSI    LanguageType = "wsi"      // Wusi
	LanguageTypeWSK    LanguageType = "wsk"      // Waskia
	LanguageTypeWSR    LanguageType = "wsr"      // Owenia
	LanguageTypeWSS    LanguageType = "wss"      // Wasa
	LanguageTypeWSU    LanguageType = "wsu"      // Wasu
	LanguageTypeWSV    LanguageType = "wsv"      // Wotapuri-Katarqalai
	LanguageTypeWTF    LanguageType = "wtf"      // Watiwa
	LanguageTypeWTH    LanguageType = "wth"      // Wathawurrung
	LanguageTypeWTI    LanguageType = "wti"      // Berta
	LanguageTypeWTK    LanguageType = "wtk"      // Watakataui
	LanguageTypeWTM    LanguageType = "wtm"      // Mewati
	LanguageTypeWTW    LanguageType = "wtw"      // Wotu
	LanguageTypeWUA    LanguageType = "wua"      // Wikngenchera
	LanguageTypeWUB    LanguageType = "wub"      // Wunambal
	LanguageTypeWUD    LanguageType = "wud"      // Wudu
	LanguageTypeWUH    LanguageType = "wuh"      // Wutunhua
	LanguageTypeWUL    LanguageType = "wul"      // Silimo
	LanguageTypeWUM    LanguageType = "wum"      // Wumbvu
	LanguageTypeWUN    LanguageType = "wun"      // Bungu
	LanguageTypeWUR    LanguageType = "wur"      // Wurrugu
	LanguageTypeWUT    LanguageType = "wut"      // Wutung
	LanguageTypeWUU    LanguageType = "wuu"      // Wu Chinese
	LanguageTypeWUV    LanguageType = "wuv"      // Wuvulu-Aua
	LanguageTypeWUX    LanguageType = "wux"      // Wulna
	LanguageTypeWUY    LanguageType = "wuy"      // Wauyai
	LanguageTypeWWA    LanguageType = "wwa"      // Waama
	LanguageTypeWWB    LanguageType = "wwb"      // Wakabunga
	LanguageTypeWWO    LanguageType = "wwo"      // Wetamut and Dorig
	LanguageTypeWWR    LanguageType = "wwr"      // Warrwa
	LanguageTypeWWW    LanguageType = "www"      // Wawa
	LanguageTypeWXA    LanguageType = "wxa"      // Waxianghua
	LanguageTypeWXW    LanguageType = "wxw"      // Wardandi
	LanguageTypeWYA    LanguageType = "wya"      // Wyandot
	LanguageTypeWYB    LanguageType = "wyb"      // Wangaaybuwan-Ngiyambaa
	LanguageTypeWYI    LanguageType = "wyi"      // Woiwurrung
	LanguageTypeWYM    LanguageType = "wym"      // Wymysorys
	LanguageTypeWYR    LanguageType = "wyr"      // Wayoró
	LanguageTypeWYY    LanguageType = "wyy"      // Western Fijian
	LanguageTypeXAA    LanguageType = "xaa"      // Andalusian Arabic
	LanguageTypeXAB    LanguageType = "xab"      // Sambe
	LanguageTypeXAC    LanguageType = "xac"      // Kachari
	LanguageTypeXAD    LanguageType = "xad"      // Adai
	LanguageTypeXAE    LanguageType = "xae"      // Aequian
	LanguageTypeXAG    LanguageType = "xag"      // Aghwan
	LanguageTypeXAI    LanguageType = "xai"      // Kaimbé
	LanguageTypeXAL    LanguageType = "xal"      // Kalmyk and Oirat
	LanguageTypeXAM    LanguageType = "xam"      // /Xam
	LanguageTypeXAN    LanguageType = "xan"      // Xamtanga
	LanguageTypeXAO    LanguageType = "xao"      // Khao
	LanguageTypeXAP    LanguageType = "xap"      // Apalachee
	LanguageTypeXAQ    LanguageType = "xaq"      // Aquitanian
	LanguageTypeXAR    LanguageType = "xar"      // Karami
	LanguageTypeXAS    LanguageType = "xas"      // Kamas
	LanguageTypeXAT    LanguageType = "xat"      // Katawixi
	LanguageTypeXAU    LanguageType = "xau"      // Kauwera
	LanguageTypeXAV    LanguageType = "xav"      // Xavánte
	LanguageTypeXAW    LanguageType = "xaw"      // Kawaiisu
	LanguageTypeXAY    LanguageType = "xay"      // Kayan Mahakam
	LanguageTypeXBA    LanguageType = "xba"      // Kamba (Brazil)
	LanguageTypeXBB    LanguageType = "xbb"      // Lower Burdekin
	LanguageTypeXBC    LanguageType = "xbc"      // Bactrian
	LanguageTypeXBD    LanguageType = "xbd"      // Bindal
	LanguageTypeXBE    LanguageType = "xbe"      // Bigambal
	LanguageTypeXBG    LanguageType = "xbg"      // Bunganditj
	LanguageTypeXBI    LanguageType = "xbi"      // Kombio
	LanguageTypeXBJ    LanguageType = "xbj"      // Birrpayi
	LanguageTypeXBM    LanguageType = "xbm"      // Middle Breton
	LanguageTypeXBN    LanguageType = "xbn"      // Kenaboi
	LanguageTypeXBO    LanguageType = "xbo"      // Bolgarian
	LanguageTypeXBP    LanguageType = "xbp"      // Bibbulman
	LanguageTypeXBR    LanguageType = "xbr"      // Kambera
	LanguageTypeXBW    LanguageType = "xbw"      // Kambiwá
	LanguageTypeXBX    LanguageType = "xbx"      // Kabixí
	LanguageTypeXBY    LanguageType = "xby"      // Batyala
	LanguageTypeXCB    LanguageType = "xcb"      // Cumbric
	LanguageTypeXCC    LanguageType = "xcc"      // Camunic
	LanguageTypeXCE    LanguageType = "xce"      // Celtiberian
	LanguageTypeXCG    LanguageType = "xcg"      // Cisalpine Gaulish
	LanguageTypeXCH    LanguageType = "xch"      // Chemakum and Chimakum
	LanguageTypeXCL    LanguageType = "xcl"      // Classical Armenian
	LanguageTypeXCM    LanguageType = "xcm"      // Comecrudo
	LanguageTypeXCN    LanguageType = "xcn"      // Cotoname
	LanguageTypeXCO    LanguageType = "xco"      // Chorasmian
	LanguageTypeXCR    LanguageType = "xcr"      // Carian
	LanguageTypeXCT    LanguageType = "xct"      // Classical Tibetan
	LanguageTypeXCU    LanguageType = "xcu"      // Curonian
	LanguageTypeXCV    LanguageType = "xcv"      // Chuvantsy
	LanguageTypeXCW    LanguageType = "xcw"      // Coahuilteco
	LanguageTypeXCY    LanguageType = "xcy"      // Cayuse
	LanguageTypeXDA    LanguageType = "xda"      // Darkinyung
	LanguageTypeXDC    LanguageType = "xdc"      // Dacian
	LanguageTypeXDK    LanguageType = "xdk"      // Dharuk
	LanguageTypeXDM    LanguageType = "xdm"      // Edomite
	LanguageTypeXDY    LanguageType = "xdy"      // Malayic Dayak
	LanguageTypeXEB    LanguageType = "xeb"      // Eblan
	LanguageTypeXED    LanguageType = "xed"      // Hdi
	LanguageTypeXEG    LanguageType = "xeg"      // //Xegwi
	LanguageTypeXEL    LanguageType = "xel"      // Kelo
	LanguageTypeXEM    LanguageType = "xem"      // Kembayan
	LanguageTypeXEP    LanguageType = "xep"      // Epi-Olmec
	LanguageTypeXER    LanguageType = "xer"      // Xerénte
	LanguageTypeXES    LanguageType = "xes"      // Kesawai
	LanguageTypeXET    LanguageType = "xet"      // Xetá
	LanguageTypeXEU    LanguageType = "xeu"      // Keoru-Ahia
	LanguageTypeXFA    LanguageType = "xfa"      // Faliscan
	LanguageTypeXGA    LanguageType = "xga"      // Galatian
	LanguageTypeXGB    LanguageType = "xgb"      // Gbin
	LanguageTypeXGD    LanguageType = "xgd"      // Gudang
	LanguageTypeXGF    LanguageType = "xgf"      // Gabrielino-Fernandeño
	LanguageTypeXGG    LanguageType = "xgg"      // Goreng
	LanguageTypeXGI    LanguageType = "xgi"      // Garingbal
	LanguageTypeXGL    LanguageType = "xgl"      // Galindan
	LanguageTypeXGM    LanguageType = "xgm"      // Guwinmal
	LanguageTypeXGN    LanguageType = "xgn"      // Mongolian languages
	LanguageTypeXGR    LanguageType = "xgr"      // Garza
	LanguageTypeXGU    LanguageType = "xgu"      // Unggumi
	LanguageTypeXGW    LanguageType = "xgw"      // Guwa
	LanguageTypeXHA    LanguageType = "xha"      // Harami
	LanguageTypeXHC    LanguageType = "xhc"      // Hunnic
	LanguageTypeXHD    LanguageType = "xhd"      // Hadrami
	LanguageTypeXHE    LanguageType = "xhe"      // Khetrani
	LanguageTypeXHR    LanguageType = "xhr"      // Hernican
	LanguageTypeXHT    LanguageType = "xht"      // Hattic
	LanguageTypeXHU    LanguageType = "xhu"      // Hurrian
	LanguageTypeXHV    LanguageType = "xhv"      // Khua
	LanguageTypeXIA    LanguageType = "xia"      // Xiandao
	LanguageTypeXIB    LanguageType = "xib"      // Iberian
	LanguageTypeXII    LanguageType = "xii"      // Xiri
	LanguageTypeXIL    LanguageType = "xil"      // Illyrian
	LanguageTypeXIN    LanguageType = "xin"      // Xinca
	LanguageTypeXIP    LanguageType = "xip"      // Xipináwa
	LanguageTypeXIR    LanguageType = "xir"      // Xiriâna
	LanguageTypeXIV    LanguageType = "xiv"      // Indus Valley Language
	LanguageTypeXIY    LanguageType = "xiy"      // Xipaya
	LanguageTypeXJB    LanguageType = "xjb"      // Minjungbal
	LanguageTypeXJT    LanguageType = "xjt"      // Jaitmatang
	LanguageTypeXKA    LanguageType = "xka"      // Kalkoti
	LanguageTypeXKB    LanguageType = "xkb"      // Northern Nago
	LanguageTypeXKC    LanguageType = "xkc"      // Kho'ini
	LanguageTypeXKD    LanguageType = "xkd"      // Mendalam Kayan
	LanguageTypeXKE    LanguageType = "xke"      // Kereho
	LanguageTypeXKF    LanguageType = "xkf"      // Khengkha
	LanguageTypeXKG    LanguageType = "xkg"      // Kagoro
	LanguageTypeXKH    LanguageType = "xkh"      // Karahawyana
	LanguageTypeXKI    LanguageType = "xki"      // Kenyan Sign Language
	LanguageTypeXKJ    LanguageType = "xkj"      // Kajali
	LanguageTypeXKK    LanguageType = "xkk"      // Kaco'
	LanguageTypeXKL    LanguageType = "xkl"      // Mainstream Kenyah
	LanguageTypeXKN    LanguageType = "xkn"      // Kayan River Kayan
	LanguageTypeXKO    LanguageType = "xko"      // Kiorr
	LanguageTypeXKP    LanguageType = "xkp"      // Kabatei
	LanguageTypeXKQ    LanguageType = "xkq"      // Koroni
	LanguageTypeXKR    LanguageType = "xkr"      // Xakriabá
	LanguageTypeXKS    LanguageType = "xks"      // Kumbewaha
	LanguageTypeXKT    LanguageType = "xkt"      // Kantosi
	LanguageTypeXKU    LanguageType = "xku"      // Kaamba
	LanguageTypeXKV    LanguageType = "xkv"      // Kgalagadi
	LanguageTypeXKW    LanguageType = "xkw"      // Kembra
	LanguageTypeXKX    LanguageType = "xkx"      // Karore
	LanguageTypeXKY    LanguageType = "xky"      // Uma' Lasan
	LanguageTypeXKZ    LanguageType = "xkz"      // Kurtokha
	LanguageTypeXLA    LanguageType = "xla"      // Kamula
	LanguageTypeXLB    LanguageType = "xlb"      // Loup B
	LanguageTypeXLC    LanguageType = "xlc"      // Lycian
	LanguageTypeXLD    LanguageType = "xld"      // Lydian
	LanguageTypeXLE    LanguageType = "xle"      // Lemnian
	LanguageTypeXLG    LanguageType = "xlg"      // Ligurian (Ancient)
	LanguageTypeXLI    LanguageType = "xli"      // Liburnian
	LanguageTypeXLN    LanguageType = "xln"      // Alanic
	LanguageTypeXLO    LanguageType = "xlo"      // Loup A
	LanguageTypeXLP    LanguageType = "xlp"      // Lepontic
	LanguageTypeXLS    LanguageType = "xls"      // Lusitanian
	LanguageTypeXLU    LanguageType = "xlu"      // Cuneiform Luwian
	LanguageTypeXLY    LanguageType = "xly"      // Elymian
	LanguageTypeXMA    LanguageType = "xma"      // Mushungulu
	LanguageTypeXMB    LanguageType = "xmb"      // Mbonga
	LanguageTypeXMC    LanguageType = "xmc"      // Makhuwa-Marrevone
	LanguageTypeXMD    LanguageType = "xmd"      // Mbudum
	LanguageTypeXME    LanguageType = "xme"      // Median
	LanguageTypeXMF    LanguageType = "xmf"      // Mingrelian
	LanguageTypeXMG    LanguageType = "xmg"      // Mengaka
	LanguageTypeXMH    LanguageType = "xmh"      // Kuku-Muminh
	LanguageTypeXMJ    LanguageType = "xmj"      // Majera
	LanguageTypeXMK    LanguageType = "xmk"      // Ancient Macedonian
	LanguageTypeXML    LanguageType = "xml"      // Malaysian Sign Language
	LanguageTypeXMM    LanguageType = "xmm"      // Manado Malay
	LanguageTypeXMN    LanguageType = "xmn"      // Manichaean Middle Persian
	LanguageTypeXMO    LanguageType = "xmo"      // Morerebi
	LanguageTypeXMP    LanguageType = "xmp"      // Kuku-Mu'inh
	LanguageTypeXMQ    LanguageType = "xmq"      // Kuku-Mangk
	LanguageTypeXMR    LanguageType = "xmr"      // Meroitic
	LanguageTypeXMS    LanguageType = "xms"      // Moroccan Sign Language
	LanguageTypeXMT    LanguageType = "xmt"      // Matbat
	LanguageTypeXMU    LanguageType = "xmu"      // Kamu
	LanguageTypeXMV    LanguageType = "xmv"      // Antankarana Malagasy and Tankarana Malagasy
	LanguageTypeXMW    LanguageType = "xmw"      // Tsimihety Malagasy
	LanguageTypeXMX    LanguageType = "xmx"      // Maden
	LanguageTypeXMY    LanguageType = "xmy"      // Mayaguduna
	LanguageTypeXMZ    LanguageType = "xmz"      // Mori Bawah
	LanguageTypeXNA    LanguageType = "xna"      // Ancient North Arabian
	LanguageTypeXNB    LanguageType = "xnb"      // Kanakanabu
	LanguageTypeXND    LanguageType = "xnd"      // Na-Dene languages
	LanguageTypeXNG    LanguageType = "xng"      // Middle Mongolian
	LanguageTypeXNH    LanguageType = "xnh"      // Kuanhua
	LanguageTypeXNI    LanguageType = "xni"      // Ngarigu
	LanguageTypeXNK    LanguageType = "xnk"      // Nganakarti
	LanguageTypeXNN    LanguageType = "xnn"      // Northern Kankanay
	LanguageTypeXNO    LanguageType = "xno"      // Anglo-Norman
	LanguageTypeXNR    LanguageType = "xnr"      // Kangri
	LanguageTypeXNS    LanguageType = "xns"      // Kanashi
	LanguageTypeXNT    LanguageType = "xnt"      // Narragansett
	LanguageTypeXNU    LanguageType = "xnu"      // Nukunul
	LanguageTypeXNY    LanguageType = "xny"      // Nyiyaparli
	LanguageTypeXNZ    LanguageType = "xnz"      // Kenzi and Mattoki
	LanguageTypeXOC    LanguageType = "xoc"      // O'chi'chi'
	LanguageTypeXOD    LanguageType = "xod"      // Kokoda
	LanguageTypeXOG    LanguageType = "xog"      // Soga
	LanguageTypeXOI    LanguageType = "xoi"      // Kominimung
	LanguageTypeXOK    LanguageType = "xok"      // Xokleng
	LanguageTypeXOM    LanguageType = "xom"      // Komo (Sudan)
	LanguageTypeXON    LanguageType = "xon"      // Konkomba
	LanguageTypeXOO    LanguageType = "xoo"      // Xukurú
	LanguageTypeXOP    LanguageType = "xop"      // Kopar
	LanguageTypeXOR    LanguageType = "xor"      // Korubo
	LanguageTypeXOW    LanguageType = "xow"      // Kowaki
	LanguageTypeXPA    LanguageType = "xpa"      // Pirriya
	LanguageTypeXPC    LanguageType = "xpc"      // Pecheneg
	LanguageTypeXPE    LanguageType = "xpe"      // Liberia Kpelle
	LanguageTypeXPG    LanguageType = "xpg"      // Phrygian
	LanguageTypeXPI    LanguageType = "xpi"      // Pictish
	LanguageTypeXPJ    LanguageType = "xpj"      // Mpalitjanh
	LanguageTypeXPK    LanguageType = "xpk"      // Kulina Pano
	LanguageTypeXPM    LanguageType = "xpm"      // Pumpokol
	LanguageTypeXPN    LanguageType = "xpn"      // Kapinawá
	LanguageTypeXPO    LanguageType = "xpo"      // Pochutec
	LanguageTypeXPP    LanguageType = "xpp"      // Puyo-Paekche
	LanguageTypeXPQ    LanguageType = "xpq"      // Mohegan-Pequot
	LanguageTypeXPR    LanguageType = "xpr"      // Parthian
	LanguageTypeXPS    LanguageType = "xps"      // Pisidian
	LanguageTypeXPT    LanguageType = "xpt"      // Punthamara
	LanguageTypeXPU    LanguageType = "xpu"      // Punic
	LanguageTypeXPY    LanguageType = "xpy"      // Puyo
	LanguageTypeXQA    LanguageType = "xqa"      // Karakhanid
	LanguageTypeXQT    LanguageType = "xqt"      // Qatabanian
	LanguageTypeXRA    LanguageType = "xra"      // Krahô
	LanguageTypeXRB    LanguageType = "xrb"      // Eastern Karaboro
	LanguageTypeXRD    LanguageType = "xrd"      // Gundungurra
	LanguageTypeXRE    LanguageType = "xre"      // Kreye
	LanguageTypeXRG    LanguageType = "xrg"      // Minang
	LanguageTypeXRI    LanguageType = "xri"      // Krikati-Timbira
	LanguageTypeXRM    LanguageType = "xrm"      // Armazic
	LanguageTypeXRN    LanguageType = "xrn"      // Arin
	LanguageTypeXRQ    LanguageType = "xrq"      // Karranga
	LanguageTypeXRR    LanguageType = "xrr"      // Raetic
	LanguageTypeXRT    LanguageType = "xrt"      // Aranama-Tamique
	LanguageTypeXRU    LanguageType = "xru"      // Marriammu
	LanguageTypeXRW    LanguageType = "xrw"      // Karawa
	LanguageTypeXSA    LanguageType = "xsa"      // Sabaean
	LanguageTypeXSB    LanguageType = "xsb"      // Sambal
	LanguageTypeXSC    LanguageType = "xsc"      // Scythian
	LanguageTypeXSD    LanguageType = "xsd"      // Sidetic
	LanguageTypeXSE    LanguageType = "xse"      // Sempan
	LanguageTypeXSH    LanguageType = "xsh"      // Shamang
	LanguageTypeXSI    LanguageType = "xsi"      // Sio
	LanguageTypeXSJ    LanguageType = "xsj"      // Subi
	LanguageTypeXSL    LanguageType = "xsl"      // South Slavey
	LanguageTypeXSM    LanguageType = "xsm"      // Kasem
	LanguageTypeXSN    LanguageType = "xsn"      // Sanga (Nigeria)
	LanguageTypeXSO    LanguageType = "xso"      // Solano
	LanguageTypeXSP    LanguageType = "xsp"      // Silopi
	LanguageTypeXSQ    LanguageType = "xsq"      // Makhuwa-Saka
	LanguageTypeXSR    LanguageType = "xsr"      // Sherpa
	LanguageTypeXSS    LanguageType = "xss"      // Assan
	LanguageTypeXSU    LanguageType = "xsu"      // Sanumá
	LanguageTypeXSV    LanguageType = "xsv"      // Sudovian
	LanguageTypeXSY    LanguageType = "xsy"      // Saisiyat
	LanguageTypeXTA    LanguageType = "xta"      // Alcozauca Mixtec
	LanguageTypeXTB    LanguageType = "xtb"      // Chazumba Mixtec
	LanguageTypeXTC    LanguageType = "xtc"      // Katcha-Kadugli-Miri
	LanguageTypeXTD    LanguageType = "xtd"      // Diuxi-Tilantongo Mixtec
	LanguageTypeXTE    LanguageType = "xte"      // Ketengban
	LanguageTypeXTG    LanguageType = "xtg"      // Transalpine Gaulish
	LanguageTypeXTH    LanguageType = "xth"      // Yitha Yitha
	LanguageTypeXTI    LanguageType = "xti"      // Sinicahua Mixtec
	LanguageTypeXTJ    LanguageType = "xtj"      // San Juan Teita Mixtec
	LanguageTypeXTL    LanguageType = "xtl"      // Tijaltepec Mixtec
	LanguageTypeXTM    LanguageType = "xtm"      // Magdalena Peñasco Mixtec
	LanguageTypeXTN    LanguageType = "xtn"      // Northern Tlaxiaco Mixtec
	LanguageTypeXTO    LanguageType = "xto"      // Tokharian A
	LanguageTypeXTP    LanguageType = "xtp"      // San Miguel Piedras Mixtec
	LanguageTypeXTQ    LanguageType = "xtq"      // Tumshuqese
	LanguageTypeXTR    LanguageType = "xtr"      // Early Tripuri
	LanguageTypeXTS    LanguageType = "xts"      // Sindihui Mixtec
	LanguageTypeXTT    LanguageType = "xtt"      // Tacahua Mixtec
	LanguageTypeXTU    LanguageType = "xtu"      // Cuyamecalco Mixtec
	LanguageTypeXTV    LanguageType = "xtv"      // Thawa
	LanguageTypeXTW    LanguageType = "xtw"      // Tawandê
	LanguageTypeXTY    LanguageType = "xty"      // Yoloxochitl Mixtec
	LanguageTypeXTZ    LanguageType = "xtz"      // Tasmanian
	LanguageTypeXUA    LanguageType = "xua"      // Alu Kurumba
	LanguageTypeXUB    LanguageType = "xub"      // Betta Kurumba
	LanguageTypeXUD    LanguageType = "xud"      // Umiida
	LanguageTypeXUG    LanguageType = "xug"      // Kunigami
	LanguageTypeXUJ    LanguageType = "xuj"      // Jennu Kurumba
	LanguageTypeXUL    LanguageType = "xul"      // Ngunawal
	LanguageTypeXUM    LanguageType = "xum"      // Umbrian
	LanguageTypeXUN    LanguageType = "xun"      // Unggaranggu
	LanguageTypeXUO    LanguageType = "xuo"      // Kuo
	LanguageTypeXUP    LanguageType = "xup"      // Upper Umpqua
	LanguageTypeXUR    LanguageType = "xur"      // Urartian
	LanguageTypeXUT    LanguageType = "xut"      // Kuthant
	LanguageTypeXUU    LanguageType = "xuu"      // Kxoe
	LanguageTypeXVE    LanguageType = "xve"      // Venetic
	LanguageTypeXVI    LanguageType = "xvi"      // Kamviri
	LanguageTypeXVN    LanguageType = "xvn"      // Vandalic
	LanguageTypeXVO    LanguageType = "xvo"      // Volscian
	LanguageTypeXVS    LanguageType = "xvs"      // Vestinian
	LanguageTypeXWA    LanguageType = "xwa"      // Kwaza
	LanguageTypeXWC    LanguageType = "xwc"      // Woccon
	LanguageTypeXWD    LanguageType = "xwd"      // Wadi Wadi
	LanguageTypeXWE    LanguageType = "xwe"      // Xwela Gbe
	LanguageTypeXWG    LanguageType = "xwg"      // Kwegu
	LanguageTypeXWJ    LanguageType = "xwj"      // Wajuk
	LanguageTypeXWK    LanguageType = "xwk"      // Wangkumara
	LanguageTypeXWL    LanguageType = "xwl"      // Western Xwla Gbe
	LanguageTypeXWO    LanguageType = "xwo"      // Written Oirat
	LanguageTypeXWR    LanguageType = "xwr"      // Kwerba Mamberamo
	LanguageTypeXWT    LanguageType = "xwt"      // Wotjobaluk
	LanguageTypeXWW    LanguageType = "xww"      // Wemba Wemba
	LanguageTypeXXB    LanguageType = "xxb"      // Boro (Ghana)
	LanguageTypeXXK    LanguageType = "xxk"      // Ke'o
	LanguageTypeXXM    LanguageType = "xxm"      // Minkin
	LanguageTypeXXR    LanguageType = "xxr"      // Koropó
	LanguageTypeXXT    LanguageType = "xxt"      // Tambora
	LanguageTypeXYA    LanguageType = "xya"      // Yaygir
	LanguageTypeXYB    LanguageType = "xyb"      // Yandjibara
	LanguageTypeXYJ    LanguageType = "xyj"      // Mayi-Yapi
	LanguageTypeXYK    LanguageType = "xyk"      // Mayi-Kulan
	LanguageTypeXYL    LanguageType = "xyl"      // Yalakalore
	LanguageTypeXYT    LanguageType = "xyt"      // Mayi-Thakurti
	LanguageTypeXYY    LanguageType = "xyy"      // Yorta Yorta
	LanguageTypeXZH    LanguageType = "xzh"      // Zhang-Zhung
	LanguageTypeXZM    LanguageType = "xzm"      // Zemgalian
	LanguageTypeXZP    LanguageType = "xzp"      // Ancient Zapotec
	LanguageTypeYAA    LanguageType = "yaa"      // Yaminahua
	LanguageTypeYAB    LanguageType = "yab"      // Yuhup
	LanguageTypeYAC    LanguageType = "yac"      // Pass Valley Yali
	LanguageTypeYAD    LanguageType = "yad"      // Yagua
	LanguageTypeYAE    LanguageType = "yae"      // Pumé
	LanguageTypeYAF    LanguageType = "yaf"      // Yaka (Democratic Republic of Congo)
	LanguageTypeYAG    LanguageType = "yag"      // Yámana
	LanguageTypeYAH    LanguageType = "yah"      // Yazgulyam
	LanguageTypeYAI    LanguageType = "yai"      // Yagnobi
	LanguageTypeYAJ    LanguageType = "yaj"      // Banda-Yangere
	LanguageTypeYAK    LanguageType = "yak"      // Yakama
	LanguageTypeYAL    LanguageType = "yal"      // Yalunka
	LanguageTypeYAM    LanguageType = "yam"      // Yamba
	LanguageTypeYAN    LanguageType = "yan"      // Mayangna
	LanguageTypeYAO    LanguageType = "yao"      // Yao
	LanguageTypeYAP    LanguageType = "yap"      // Yapese
	LanguageTypeYAQ    LanguageType = "yaq"      // Yaqui
	LanguageTypeYAR    LanguageType = "yar"      // Yabarana
	LanguageTypeYAS    LanguageType = "yas"      // Nugunu (Cameroon)
	LanguageTypeYAT    LanguageType = "yat"      // Yambeta
	LanguageTypeYAU    LanguageType = "yau"      // Yuwana
	LanguageTypeYAV    LanguageType = "yav"      // Yangben
	LanguageTypeYAW    LanguageType = "yaw"      // Yawalapití
	LanguageTypeYAX    LanguageType = "yax"      // Yauma
	LanguageTypeYAY    LanguageType = "yay"      // Agwagwune
	LanguageTypeYAZ    LanguageType = "yaz"      // Lokaa
	LanguageTypeYBA    LanguageType = "yba"      // Yala
	LanguageTypeYBB    LanguageType = "ybb"      // Yemba
	LanguageTypeYBD    LanguageType = "ybd"      // Yangbye
	LanguageTypeYBE    LanguageType = "ybe"      // West Yugur
	LanguageTypeYBH    LanguageType = "ybh"      // Yakha
	LanguageTypeYBI    LanguageType = "ybi"      // Yamphu
	LanguageTypeYBJ    LanguageType = "ybj"      // Hasha
	LanguageTypeYBK    LanguageType = "ybk"      // Bokha
	LanguageTypeYBL    LanguageType = "ybl"      // Yukuben
	LanguageTypeYBM    LanguageType = "ybm"      // Yaben
	LanguageTypeYBN    LanguageType = "ybn"      // Yabaâna
	LanguageTypeYBO    LanguageType = "ybo"      // Yabong
	LanguageTypeYBX    LanguageType = "ybx"      // Yawiyo
	LanguageTypeYBY    LanguageType = "yby"      // Yaweyuha
	LanguageTypeYCH    LanguageType = "ych"      // Chesu
	LanguageTypeYCL    LanguageType = "ycl"      // Lolopo
	LanguageTypeYCN    LanguageType = "ycn"      // Yucuna
	LanguageTypeYCP    LanguageType = "ycp"      // Chepya
	LanguageTypeYDA    LanguageType = "yda"      // Yanda
	LanguageTypeYDD    LanguageType = "ydd"      // Eastern Yiddish
	LanguageTypeYDE    LanguageType = "yde"      // Yangum Dey
	LanguageTypeYDG    LanguageType = "ydg"      // Yidgha
	LanguageTypeYDK    LanguageType = "ydk"      // Yoidik
	LanguageTypeYDS    LanguageType = "yds"      // Yiddish Sign Language
	LanguageTypeYEA    LanguageType = "yea"      // Ravula
	LanguageTypeYEC    LanguageType = "yec"      // Yeniche
	LanguageTypeYEE    LanguageType = "yee"      // Yimas
	LanguageTypeYEI    LanguageType = "yei"      // Yeni
	LanguageTypeYEJ    LanguageType = "yej"      // Yevanic
	LanguageTypeYEL    LanguageType = "yel"      // Yela
	LanguageTypeYEN    LanguageType = "yen"      // Yendang
	LanguageTypeYER    LanguageType = "yer"      // Tarok
	LanguageTypeYES    LanguageType = "yes"      // Nyankpa
	LanguageTypeYET    LanguageType = "yet"      // Yetfa
	LanguageTypeYEU    LanguageType = "yeu"      // Yerukula
	LanguageTypeYEV    LanguageType = "yev"      // Yapunda
	LanguageTypeYEY    LanguageType = "yey"      // Yeyi
	LanguageTypeYGA    LanguageType = "yga"      // Malyangapa
	LanguageTypeYGI    LanguageType = "ygi"      // Yiningayi
	LanguageTypeYGL    LanguageType = "ygl"      // Yangum Gel
	LanguageTypeYGM    LanguageType = "ygm"      // Yagomi
	LanguageTypeYGP    LanguageType = "ygp"      // Gepo
	LanguageTypeYGR    LanguageType = "ygr"      // Yagaria
	LanguageTypeYGU    LanguageType = "ygu"      // Yugul
	LanguageTypeYGW    LanguageType = "ygw"      // Yagwoia
	LanguageTypeYHA    LanguageType = "yha"      // Baha Buyang
	LanguageTypeYHD    LanguageType = "yhd"      // Judeo-Iraqi Arabic
	LanguageTypeYHL    LanguageType = "yhl"      // Hlepho Phowa
	LanguageTypeYIA    LanguageType = "yia"      // Yinggarda
	LanguageTypeYIF    LanguageType = "yif"      // Ache
	LanguageTypeYIG    LanguageType = "yig"      // Wusa Nasu
	LanguageTypeYIH    LanguageType = "yih"      // Western Yiddish
	LanguageTypeYII    LanguageType = "yii"      // Yidiny
	LanguageTypeYIJ    LanguageType = "yij"      // Yindjibarndi
	LanguageTypeYIK    LanguageType = "yik"      // Dongshanba Lalo
	LanguageTypeYIL    LanguageType = "yil"      // Yindjilandji
	LanguageTypeYIM    LanguageType = "yim"      // Yimchungru Naga
	LanguageTypeYIN    LanguageType = "yin"      // Yinchia
	LanguageTypeYIP    LanguageType = "yip"      // Pholo
	LanguageTypeYIQ    LanguageType = "yiq"      // Miqie
	LanguageTypeYIR    LanguageType = "yir"      // North Awyu
	LanguageTypeYIS    LanguageType = "yis"      // Yis
	LanguageTypeYIT    LanguageType = "yit"      // Eastern Lalu
	LanguageTypeYIU    LanguageType = "yiu"      // Awu
	LanguageTypeYIV    LanguageType = "yiv"      // Northern Nisu
	LanguageTypeYIX    LanguageType = "yix"      // Axi Yi
	LanguageTypeYIY    LanguageType = "yiy"      // Yir Yoront
	LanguageTypeYIZ    LanguageType = "yiz"      // Azhe
	LanguageTypeYKA    LanguageType = "yka"      // Yakan
	LanguageTypeYKG    LanguageType = "ykg"      // Northern Yukaghir
	LanguageTypeYKI    LanguageType = "yki"      // Yoke
	LanguageTypeYKK    LanguageType = "ykk"      // Yakaikeke
	LanguageTypeYKL    LanguageType = "ykl"      // Khlula
	LanguageTypeYKM    LanguageType = "ykm"      // Kap
	LanguageTypeYKN    LanguageType = "ykn"      // Kua-nsi
	LanguageTypeYKO    LanguageType = "yko"      // Yasa
	LanguageTypeYKR    LanguageType = "ykr"      // Yekora
	LanguageTypeYKT    LanguageType = "ykt"      // Kathu
	LanguageTypeYKU    LanguageType = "yku"      // Kuamasi
	LanguageTypeYKY    LanguageType = "yky"      // Yakoma
	LanguageTypeYLA    LanguageType = "yla"      // Yaul
	LanguageTypeYLB    LanguageType = "ylb"      // Yaleba
	LanguageTypeYLE    LanguageType = "yle"      // Yele
	LanguageTypeYLG    LanguageType = "ylg"      // Yelogu
	LanguageTypeYLI    LanguageType = "yli"      // Angguruk Yali
	LanguageTypeYLL    LanguageType = "yll"      // Yil
	LanguageTypeYLM    LanguageType = "ylm"      // Limi
	LanguageTypeYLN    LanguageType = "yln"      // Langnian Buyang
	LanguageTypeYLO    LanguageType = "ylo"      // Naluo Yi
	LanguageTypeYLR    LanguageType = "ylr"      // Yalarnnga
	LanguageTypeYLU    LanguageType = "ylu"      // Aribwaung
	LanguageTypeYLY    LanguageType = "yly"      // Nyâlayu and Nyelâyu
	LanguageTypeYMA    LanguageType = "yma"      // Yamphe
	LanguageTypeYMB    LanguageType = "ymb"      // Yambes
	LanguageTypeYMC    LanguageType = "ymc"      // Southern Muji
	LanguageTypeYMD    LanguageType = "ymd"      // Muda
	LanguageTypeYME    LanguageType = "yme"      // Yameo
	LanguageTypeYMG    LanguageType = "ymg"      // Yamongeri
	LanguageTypeYMH    LanguageType = "ymh"      // Mili
	LanguageTypeYMI    LanguageType = "ymi"      // Moji
	LanguageTypeYMK    LanguageType = "ymk"      // Makwe
	LanguageTypeYML    LanguageType = "yml"      // Iamalele
	LanguageTypeYMM    LanguageType = "ymm"      // Maay
	LanguageTypeYMN    LanguageType = "ymn"      // Yamna and Sunum
	LanguageTypeYMO    LanguageType = "ymo"      // Yangum Mon
	LanguageTypeYMP    LanguageType = "ymp"      // Yamap
	LanguageTypeYMQ    LanguageType = "ymq"      // Qila Muji
	LanguageTypeYMR    LanguageType = "ymr"      // Malasar
	LanguageTypeYMS    LanguageType = "yms"      // Mysian
	LanguageTypeYMT    LanguageType = "ymt"      // Mator-Taygi-Karagas
	LanguageTypeYMX    LanguageType = "ymx"      // Northern Muji
	LanguageTypeYMZ    LanguageType = "ymz"      // Muzi
	LanguageTypeYNA    LanguageType = "yna"      // Aluo
	LanguageTypeYND    LanguageType = "ynd"      // Yandruwandha
	LanguageTypeYNE    LanguageType = "yne"      // Lang'e
	LanguageTypeYNG    LanguageType = "yng"      // Yango
	LanguageTypeYNH    LanguageType = "ynh"      // Yangho
	LanguageTypeYNK    LanguageType = "ynk"      // Naukan Yupik
	LanguageTypeYNL    LanguageType = "ynl"      // Yangulam
	LanguageTypeYNN    LanguageType = "ynn"      // Yana
	LanguageTypeYNO    LanguageType = "yno"      // Yong
	LanguageTypeYNQ    LanguageType = "ynq"      // Yendang
	LanguageTypeYNS    LanguageType = "yns"      // Yansi
	LanguageTypeYNU    LanguageType = "ynu"      // Yahuna
	LanguageTypeYOB    LanguageType = "yob"      // Yoba
	LanguageTypeYOG    LanguageType = "yog"      // Yogad
	LanguageTypeYOI    LanguageType = "yoi"      // Yonaguni
	LanguageTypeYOK    LanguageType = "yok"      // Yokuts
	LanguageTypeYOL    LanguageType = "yol"      // Yola
	LanguageTypeYOM    LanguageType = "yom"      // Yombe
	LanguageTypeYON    LanguageType = "yon"      // Yongkom
	LanguageTypeYOS    LanguageType = "yos"      // Yos
	LanguageTypeYOT    LanguageType = "yot"      // Yotti
	LanguageTypeYOX    LanguageType = "yox"      // Yoron
	LanguageTypeYOY    LanguageType = "yoy"      // Yoy
	LanguageTypeYPA    LanguageType = "ypa"      // Phala
	LanguageTypeYPB    LanguageType = "ypb"      // Labo Phowa
	LanguageTypeYPG    LanguageType = "ypg"      // Phola
	LanguageTypeYPH    LanguageType = "yph"      // Phupha
	LanguageTypeYPK    LanguageType = "ypk"      // Yupik languages
	LanguageTypeYPM    LanguageType = "ypm"      // Phuma
	LanguageTypeYPN    LanguageType = "ypn"      // Ani Phowa
	LanguageTypeYPO    LanguageType = "ypo"      // Alo Phola
	LanguageTypeYPP    LanguageType = "ypp"      // Phupa
	LanguageTypeYPZ    LanguageType = "ypz"      // Phuza
	LanguageTypeYRA    LanguageType = "yra"      // Yerakai
	LanguageTypeYRB    LanguageType = "yrb"      // Yareba
	LanguageTypeYRE    LanguageType = "yre"      // Yaouré
	LanguageTypeYRI    LanguageType = "yri"      // Yarí
	LanguageTypeYRK    LanguageType = "yrk"      // Nenets
	LanguageTypeYRL    LanguageType = "yrl"      // Nhengatu
	LanguageTypeYRM    LanguageType = "yrm"      // Yirrk-Mel
	LanguageTypeYRN    LanguageType = "yrn"      // Yerong
	LanguageTypeYRS    LanguageType = "yrs"      // Yarsun
	LanguageTypeYRW    LanguageType = "yrw"      // Yarawata
	LanguageTypeYRY    LanguageType = "yry"      // Yarluyandi
	LanguageTypeYSC    LanguageType = "ysc"      // Yassic
	LanguageTypeYSD    LanguageType = "ysd"      // Samatao
	LanguageTypeYSG    LanguageType = "ysg"      // Sonaga
	LanguageTypeYSL    LanguageType = "ysl"      // Yugoslavian Sign Language
	LanguageTypeYSN    LanguageType = "ysn"      // Sani
	LanguageTypeYSO    LanguageType = "yso"      // Nisi (China)
	LanguageTypeYSP    LanguageType = "ysp"      // Southern Lolopo
	LanguageTypeYSR    LanguageType = "ysr"      // Sirenik Yupik
	LanguageTypeYSS    LanguageType = "yss"      // Yessan-Mayo
	LanguageTypeYSY    LanguageType = "ysy"      // Sanie
	LanguageTypeYTA    LanguageType = "yta"      // Talu
	LanguageTypeYTL    LanguageType = "ytl"      // Tanglang
	LanguageTypeYTP    LanguageType = "ytp"      // Thopho
	LanguageTypeYTW    LanguageType = "ytw"      // Yout Wam
	LanguageTypeYTY    LanguageType = "yty"      // Yatay
	LanguageTypeYUA    LanguageType = "yua"      // Yucateco and Yucatec Maya
	LanguageTypeYUB    LanguageType = "yub"      // Yugambal
	LanguageTypeYUC    LanguageType = "yuc"      // Yuchi
	LanguageTypeYUD    LanguageType = "yud"      // Judeo-Tripolitanian Arabic
	LanguageTypeYUE    LanguageType = "yue"      // Yue Chinese
	LanguageTypeYUF    LanguageType = "yuf"      // Havasupai-Walapai-Yavapai
	LanguageTypeYUG    LanguageType = "yug"      // Yug
	LanguageTypeYUI    LanguageType = "yui"      // Yurutí
	LanguageTypeYUJ    LanguageType = "yuj"      // Karkar-Yuri
	LanguageTypeYUK    LanguageType = "yuk"      // Yuki
	LanguageTypeYUL    LanguageType = "yul"      // Yulu
	LanguageTypeYUM    LanguageType = "yum"      // Quechan
	LanguageTypeYUN    LanguageType = "yun"      // Bena (Nigeria)
	LanguageTypeYUP    LanguageType = "yup"      // Yukpa
	LanguageTypeYUQ    LanguageType = "yuq"      // Yuqui
	LanguageTypeYUR    LanguageType = "yur"      // Yurok
	LanguageTypeYUT    LanguageType = "yut"      // Yopno
	LanguageTypeYUU    LanguageType = "yuu"      // Yugh
	LanguageTypeYUW    LanguageType = "yuw"      // Yau (Morobe Province)
	LanguageTypeYUX    LanguageType = "yux"      // Southern Yukaghir
	LanguageTypeYUY    LanguageType = "yuy"      // East Yugur
	LanguageTypeYUZ    LanguageType = "yuz"      // Yuracare
	LanguageTypeYVA    LanguageType = "yva"      // Yawa
	LanguageTypeYVT    LanguageType = "yvt"      // Yavitero
	LanguageTypeYWA    LanguageType = "ywa"      // Kalou
	LanguageTypeYWG    LanguageType = "ywg"      // Yinhawangka
	LanguageTypeYWL    LanguageType = "ywl"      // Western Lalu
	LanguageTypeYWN    LanguageType = "ywn"      // Yawanawa
	LanguageTypeYWQ    LanguageType = "ywq"      // Wuding-Luquan Yi
	LanguageTypeYWR    LanguageType = "ywr"      // Yawuru
	LanguageTypeYWT    LanguageType = "ywt"      // Xishanba Lalo and Central Lalo
	LanguageTypeYWU    LanguageType = "ywu"      // Wumeng Nasu
	LanguageTypeYWW    LanguageType = "yww"      // Yawarawarga
	LanguageTypeYXA    LanguageType = "yxa"      // Mayawali
	LanguageTypeYXG    LanguageType = "yxg"      // Yagara
	LanguageTypeYXL    LanguageType = "yxl"      // Yardliyawarra
	LanguageTypeYXM    LanguageType = "yxm"      // Yinwum
	LanguageTypeYXU    LanguageType = "yxu"      // Yuyu
	LanguageTypeYXY    LanguageType = "yxy"      // Yabula Yabula
	LanguageTypeYYR    LanguageType = "yyr"      // Yir Yoront
	LanguageTypeYYU    LanguageType = "yyu"      // Yau (Sandaun Province)
	LanguageTypeYYZ    LanguageType = "yyz"      // Ayizi
	LanguageTypeYZG    LanguageType = "yzg"      // E'ma Buyang
	LanguageTypeYZK    LanguageType = "yzk"      // Zokhuo
	LanguageTypeZAA    LanguageType = "zaa"      // Sierra de Juárez Zapotec
	LanguageTypeZAB    LanguageType = "zab"      // San Juan Guelavía Zapotec
	LanguageTypeZAC    LanguageType = "zac"      // Ocotlán Zapotec
	LanguageTypeZAD    LanguageType = "zad"      // Cajonos Zapotec
	LanguageTypeZAE    LanguageType = "zae"      // Yareni Zapotec
	LanguageTypeZAF    LanguageType = "zaf"      // Ayoquesco Zapotec
	LanguageTypeZAG    LanguageType = "zag"      // Zaghawa
	LanguageTypeZAH    LanguageType = "zah"      // Zangwal
	LanguageTypeZAI    LanguageType = "zai"      // Isthmus Zapotec
	LanguageTypeZAJ    LanguageType = "zaj"      // Zaramo
	LanguageTypeZAK    LanguageType = "zak"      // Zanaki
	LanguageTypeZAL    LanguageType = "zal"      // Zauzou
	LanguageTypeZAM    LanguageType = "zam"      // Miahuatlán Zapotec
	LanguageTypeZAO    LanguageType = "zao"      // Ozolotepec Zapotec
	LanguageTypeZAP    LanguageType = "zap"      // Zapotec
	LanguageTypeZAQ    LanguageType = "zaq"      // Aloápam Zapotec
	LanguageTypeZAR    LanguageType = "zar"      // Rincón Zapotec
	LanguageTypeZAS    LanguageType = "zas"      // Santo Domingo Albarradas Zapotec
	LanguageTypeZAT    LanguageType = "zat"      // Tabaa Zapotec
	LanguageTypeZAU    LanguageType = "zau"      // Zangskari
	LanguageTypeZAV    LanguageType = "zav"      // Yatzachi Zapotec
	LanguageTypeZAW    LanguageType = "zaw"      // Mitla Zapotec
	LanguageTypeZAX    LanguageType = "zax"      // Xadani Zapotec
	LanguageTypeZAY    LanguageType = "zay"      // Zayse-Zergulla and Zaysete
	LanguageTypeZAZ    LanguageType = "zaz"      // Zari
	LanguageTypeZBC    LanguageType = "zbc"      // Central Berawan
	LanguageTypeZBE    LanguageType = "zbe"      // East Berawan
	LanguageTypeZBL    LanguageType = "zbl"      // Blissymbols and Bliss and Blissymbolics
	LanguageTypeZBT    LanguageType = "zbt"      // Batui
	LanguageTypeZBW    LanguageType = "zbw"      // West Berawan
	LanguageTypeZCA    LanguageType = "zca"      // Coatecas Altas Zapotec
	LanguageTypeZCH    LanguageType = "zch"      // Central Hongshuihe Zhuang
	LanguageTypeZDJ    LanguageType = "zdj"      // Ngazidja Comorian
	LanguageTypeZEA    LanguageType = "zea"      // Zeeuws
	LanguageTypeZEG    LanguageType = "zeg"      // Zenag
	LanguageTypeZEH    LanguageType = "zeh"      // Eastern Hongshuihe Zhuang
	LanguageTypeZEN    LanguageType = "zen"      // Zenaga
	LanguageTypeZGA    LanguageType = "zga"      // Kinga
	LanguageTypeZGB    LanguageType = "zgb"      // Guibei Zhuang
	LanguageTypeZGH    LanguageType = "zgh"      // Standard Moroccan Tamazight
	LanguageTypeZGM    LanguageType = "zgm"      // Minz Zhuang
	LanguageTypeZGN    LanguageType = "zgn"      // Guibian Zhuang
	LanguageTypeZGR    LanguageType = "zgr"      // Magori
	LanguageTypeZHB    LanguageType = "zhb"      // Zhaba
	LanguageTypeZHD    LanguageType = "zhd"      // Dai Zhuang
	LanguageTypeZHI    LanguageType = "zhi"      // Zhire
	LanguageTypeZHN    LanguageType = "zhn"      // Nong Zhuang
	LanguageTypeZHW    LanguageType = "zhw"      // Zhoa
	LanguageTypeZHX    LanguageType = "zhx"      // Chinese (family)
	LanguageTypeZIA    LanguageType = "zia"      // Zia
	LanguageTypeZIB    LanguageType = "zib"      // Zimbabwe Sign Language
	LanguageTypeZIK    LanguageType = "zik"      // Zimakani
	LanguageTypeZIL    LanguageType = "zil"      // Zialo
	LanguageTypeZIM    LanguageType = "zim"      // Mesme
	LanguageTypeZIN    LanguageType = "zin"      // Zinza
	LanguageTypeZIR    LanguageType = "zir"      // Ziriya
	LanguageTypeZIW    LanguageType = "ziw"      // Zigula
	LanguageTypeZIZ    LanguageType = "ziz"      // Zizilivakan
	LanguageTypeZKA    LanguageType = "zka"      // Kaimbulawa
	LanguageTypeZKB    LanguageType = "zkb"      // Koibal
	LanguageTypeZKD    LanguageType = "zkd"      // Kadu
	LanguageTypeZKG    LanguageType = "zkg"      // Koguryo
	LanguageTypeZKH    LanguageType = "zkh"      // Khorezmian
	LanguageTypeZKK    LanguageType = "zkk"      // Karankawa
	LanguageTypeZKN    LanguageType = "zkn"      // Kanan
	LanguageTypeZKO    LanguageType = "zko"      // Kott
	LanguageTypeZKP    LanguageType = "zkp"      // São Paulo Kaingáng
	LanguageTypeZKR    LanguageType = "zkr"      // Zakhring
	LanguageTypeZKT    LanguageType = "zkt"      // Kitan
	LanguageTypeZKU    LanguageType = "zku"      // Kaurna
	LanguageTypeZKV    LanguageType = "zkv"      // Krevinian
	LanguageTypeZKZ    LanguageType = "zkz"      // Khazar
	LanguageTypeZLE    LanguageType = "zle"      // East Slavic languages
	LanguageTypeZLJ    LanguageType = "zlj"      // Liujiang Zhuang
	LanguageTypeZLM    LanguageType = "zlm"      // Malay (individual language)
	LanguageTypeZLN    LanguageType = "zln"      // Lianshan Zhuang
	LanguageTypeZLQ    LanguageType = "zlq"      // Liuqian Zhuang
	LanguageTypeZLS    LanguageType = "zls"      // South Slavic languages
	LanguageTypeZLW    LanguageType = "zlw"      // West Slavic languages
	LanguageTypeZMA    LanguageType = "zma"      // Manda (Australia)
	LanguageTypeZMB    LanguageType = "zmb"      // Zimba
	LanguageTypeZMC    LanguageType = "zmc"      // Margany
	LanguageTypeZMD    LanguageType = "zmd"      // Maridan
	LanguageTypeZME    LanguageType = "zme"      // Mangerr
	LanguageTypeZMF    LanguageType = "zmf"      // Mfinu
	LanguageTypeZMG    LanguageType = "zmg"      // Marti Ke
	LanguageTypeZMH    LanguageType = "zmh"      // Makolkol
	LanguageTypeZMI    LanguageType = "zmi"      // Negeri Sembilan Malay
	LanguageTypeZMJ    LanguageType = "zmj"      // Maridjabin
	LanguageTypeZMK    LanguageType = "zmk"      // Mandandanyi
	LanguageTypeZML    LanguageType = "zml"      // Madngele
	LanguageTypeZMM    LanguageType = "zmm"      // Marimanindji
	LanguageTypeZMN    LanguageType = "zmn"      // Mbangwe
	LanguageTypeZMO    LanguageType = "zmo"      // Molo
	LanguageTypeZMP    LanguageType = "zmp"      // Mpuono
	LanguageTypeZMQ    LanguageType = "zmq"      // Mituku
	LanguageTypeZMR    LanguageType = "zmr"      // Maranunggu
	LanguageTypeZMS    LanguageType = "zms"      // Mbesa
	LanguageTypeZMT    LanguageType = "zmt"      // Maringarr
	LanguageTypeZMU    LanguageType = "zmu"      // Muruwari
	LanguageTypeZMV    LanguageType = "zmv"      // Mbariman-Gudhinma
	LanguageTypeZMW    LanguageType = "zmw"      // Mbo (Democratic Republic of Congo)
	LanguageTypeZMX    LanguageType = "zmx"      // Bomitaba
	LanguageTypeZMY    LanguageType = "zmy"      // Mariyedi
	LanguageTypeZMZ    LanguageType = "zmz"      // Mbandja
	LanguageTypeZNA    LanguageType = "zna"      // Zan Gula
	LanguageTypeZND    LanguageType = "znd"      // Zande languages
	LanguageTypeZNE    LanguageType = "zne"      // Zande (individual language)
	LanguageTypeZNG    LanguageType = "zng"      // Mang
	LanguageTypeZNK    LanguageType = "znk"      // Manangkari
	LanguageTypeZNS    LanguageType = "zns"      // Mangas
	LanguageTypeZOC    LanguageType = "zoc"      // Copainalá Zoque
	LanguageTypeZOH    LanguageType = "zoh"      // Chimalapa Zoque
	LanguageTypeZOM    LanguageType = "zom"      // Zou
	LanguageTypeZOO    LanguageType = "zoo"      // Asunción Mixtepec Zapotec
	LanguageTypeZOQ    LanguageType = "zoq"      // Tabasco Zoque
	LanguageTypeZOR    LanguageType = "zor"      // Rayón Zoque
	LanguageTypeZOS    LanguageType = "zos"      // Francisco León Zoque
	LanguageTypeZPA    LanguageType = "zpa"      // Lachiguiri Zapotec
	LanguageTypeZPB    LanguageType = "zpb"      // Yautepec Zapotec
	LanguageTypeZPC    LanguageType = "zpc"      // Choapan Zapotec
	LanguageTypeZPD    LanguageType = "zpd"      // Southeastern Ixtlán Zapotec
	LanguageTypeZPE    LanguageType = "zpe"      // Petapa Zapotec
	LanguageTypeZPF    LanguageType = "zpf"      // San Pedro Quiatoni Zapotec
	LanguageTypeZPG    LanguageType = "zpg"      // Guevea De Humboldt Zapotec
	LanguageTypeZPH    LanguageType = "zph"      // Totomachapan Zapotec
	LanguageTypeZPI    LanguageType = "zpi"      // Santa María Quiegolani Zapotec
	LanguageTypeZPJ    LanguageType = "zpj"      // Quiavicuzas Zapotec
	LanguageTypeZPK    LanguageType = "zpk"      // Tlacolulita Zapotec
	LanguageTypeZPL    LanguageType = "zpl"      // Lachixío Zapotec
	LanguageTypeZPM    LanguageType = "zpm"      // Mixtepec Zapotec
	LanguageTypeZPN    LanguageType = "zpn"      // Santa Inés Yatzechi Zapotec
	LanguageTypeZPO    LanguageType = "zpo"      // Amatlán Zapotec
	LanguageTypeZPP    LanguageType = "zpp"      // El Alto Zapotec
	LanguageTypeZPQ    LanguageType = "zpq"      // Zoogocho Zapotec
	LanguageTypeZPR    LanguageType = "zpr"      // Santiago Xanica Zapotec
	LanguageTypeZPS    LanguageType = "zps"      // Coatlán Zapotec
	LanguageTypeZPT    LanguageType = "zpt"      // San Vicente Coatlán Zapotec
	LanguageTypeZPU    LanguageType = "zpu"      // Yalálag Zapotec
	LanguageTypeZPV    LanguageType = "zpv"      // Chichicapan Zapotec
	LanguageTypeZPW    LanguageType = "zpw"      // Zaniza Zapotec
	LanguageTypeZPX    LanguageType = "zpx"      // San Baltazar Loxicha Zapotec
	LanguageTypeZPY    LanguageType = "zpy"      // Mazaltepec Zapotec
	LanguageTypeZPZ    LanguageType = "zpz"      // Texmelucan Zapotec
	LanguageTypeZQE    LanguageType = "zqe"      // Qiubei Zhuang
	LanguageTypeZRA    LanguageType = "zra"      // Kara (Korea)
	LanguageTypeZRG    LanguageType = "zrg"      // Mirgan
	LanguageTypeZRN    LanguageType = "zrn"      // Zerenkel
	LanguageTypeZRO    LanguageType = "zro"      // Záparo
	LanguageTypeZRP    LanguageType = "zrp"      // Zarphatic
	LanguageTypeZRS    LanguageType = "zrs"      // Mairasi
	LanguageTypeZSA    LanguageType = "zsa"      // Sarasira
	LanguageTypeZSK    LanguageType = "zsk"      // Kaskean
	LanguageTypeZSL    LanguageType = "zsl"      // Zambian Sign Language
	LanguageTypeZSM    LanguageType = "zsm"      // Standard Malay
	LanguageTypeZSR    LanguageType = "zsr"      // Southern Rincon Zapotec
	LanguageTypeZSU    LanguageType = "zsu"      // Sukurum
	LanguageTypeZTE    LanguageType = "zte"      // Elotepec Zapotec
	LanguageTypeZTG    LanguageType = "ztg"      // Xanaguía Zapotec
	LanguageTypeZTL    LanguageType = "ztl"      // Lapaguía-Guivini Zapotec
	LanguageTypeZTM    LanguageType = "ztm"      // San Agustín Mixtepec Zapotec
	LanguageTypeZTN    LanguageType = "ztn"      // Santa Catarina Albarradas Zapotec
	LanguageTypeZTP    LanguageType = "ztp"      // Loxicha Zapotec
	LanguageTypeZTQ    LanguageType = "ztq"      // Quioquitani-Quierí Zapotec
	LanguageTypeZTS    LanguageType = "zts"      // Tilquiapan Zapotec
	LanguageTypeZTT    LanguageType = "ztt"      // Tejalapan Zapotec
	LanguageTypeZTU    LanguageType = "ztu"      // Güilá Zapotec
	LanguageTypeZTX    LanguageType = "ztx"      // Zaachila Zapotec
	LanguageTypeZTY    LanguageType = "zty"      // Yatee Zapotec
	LanguageTypeZUA    LanguageType = "zua"      // Zeem
	LanguageTypeZUH    LanguageType = "zuh"      // Tokano
	LanguageTypeZUM    LanguageType = "zum"      // Kumzari
	LanguageTypeZUN    LanguageType = "zun"      // Zuni
	LanguageTypeZUY    LanguageType = "zuy"      // Zumaya
	LanguageTypeZWA    LanguageType = "zwa"      // Zay
	LanguageTypeZXX    LanguageType = "zxx"      // No linguistic content and Not applicable
	LanguageTypeZYB    LanguageType = "zyb"      // Yongbei Zhuang
	LanguageTypeZYG    LanguageType = "zyg"      // Yang Zhuang
	LanguageTypeZYJ    LanguageType = "zyj"      // Youjiang Zhuang
	LanguageTypeZYN    LanguageType = "zyn"      // Yongnan Zhuang
	LanguageTypeZYP    LanguageType = "zyp"      // Zyphe Chin
	LanguageTypeZZA    LanguageType = "zza"      // Zaza and Dimili and Dimli (macrolanguage) and Kirdki and Kirmanjki (macrolanguage) and Zazaki
	LanguageTypeZZJ    LanguageType = "zzj"      // Zuojiang Zhuang
)

// Used in the system everytime that we need to determinate a preferred language, the LanguageType
// is an enumerate that describe all possible languages
type LanguageType string

// Structure used to identify if a language exists or not
var (
	languageTypes map[LanguageType]bool = map[LanguageType]bool{
		LanguageTypeAA:     true,
		LanguageTypeAB:     true,
		LanguageTypeAE:     true,
		LanguageTypeAF:     true,
		LanguageTypeAK:     true,
		LanguageTypeAM:     true,
		LanguageTypeAN:     true,
		LanguageTypeAR:     true,
		LanguageTypeAS:     true,
		LanguageTypeAV:     true,
		LanguageTypeAY:     true,
		LanguageTypeAZ:     true,
		LanguageTypeBA:     true,
		LanguageTypeBE:     true,
		LanguageTypeBG:     true,
		LanguageTypeBH:     true,
		LanguageTypeBI:     true,
		LanguageTypeBM:     true,
		LanguageTypeBN:     true,
		LanguageTypeBO:     true,
		LanguageTypeBR:     true,
		LanguageTypeBS:     true,
		LanguageTypeCA:     true,
		LanguageTypeCE:     true,
		LanguageTypeCH:     true,
		LanguageTypeCO:     true,
		LanguageTypeCR:     true,
		LanguageTypeCS:     true,
		LanguageTypeCU:     true,
		LanguageTypeCV:     true,
		LanguageTypeCY:     true,
		LanguageTypeDA:     true,
		LanguageTypeDE:     true,
		LanguageTypeDV:     true,
		LanguageTypeDZ:     true,
		LanguageTypeEE:     true,
		LanguageTypeEL:     true,
		LanguageTypeEN:     true,
		LanguageTypeEO:     true,
		LanguageTypeES:     true,
		LanguageTypeET:     true,
		LanguageTypeEU:     true,
		LanguageTypeFA:     true,
		LanguageTypeFF:     true,
		LanguageTypeFI:     true,
		LanguageTypeFJ:     true,
		LanguageTypeFO:     true,
		LanguageTypeFR:     true,
		LanguageTypeFY:     true,
		LanguageTypeGA:     true,
		LanguageTypeGD:     true,
		LanguageTypeGL:     true,
		LanguageTypeGN:     true,
		LanguageTypeGU:     true,
		LanguageTypeGV:     true,
		LanguageTypeHA:     true,
		LanguageTypeHE:     true,
		LanguageTypeHI:     true,
		LanguageTypeHO:     true,
		LanguageTypeHR:     true,
		LanguageTypeHT:     true,
		LanguageTypeHU:     true,
		LanguageTypeHY:     true,
		LanguageTypeHZ:     true,
		LanguageTypeIA:     true,
		LanguageTypeID:     true,
		LanguageTypeIE:     true,
		LanguageTypeIG:     true,
		LanguageTypeII:     true,
		LanguageTypeIK:     true,
		LanguageTypeIN:     true,
		LanguageTypeIO:     true,
		LanguageTypeIS:     true,
		LanguageTypeIT:     true,
		LanguageTypeIU:     true,
		LanguageTypeIW:     true,
		LanguageTypeJA:     true,
		LanguageTypeJI:     true,
		LanguageTypeJV:     true,
		LanguageTypeJW:     true,
		LanguageTypeKA:     true,
		LanguageTypeKG:     true,
		LanguageTypeKI:     true,
		LanguageTypeKJ:     true,
		LanguageTypeKK:     true,
		LanguageTypeKL:     true,
		LanguageTypeKM:     true,
		LanguageTypeKN:     true,
		LanguageTypeKO:     true,
		LanguageTypeKR:     true,
		LanguageTypeKS:     true,
		LanguageTypeKU:     true,
		LanguageTypeKV:     true,
		LanguageTypeKW:     true,
		LanguageTypeKY:     true,
		LanguageTypeLA:     true,
		LanguageTypeLB:     true,
		LanguageTypeLG:     true,
		LanguageTypeLI:     true,
		LanguageTypeLN:     true,
		LanguageTypeLO:     true,
		LanguageTypeLT:     true,
		LanguageTypeLU:     true,
		LanguageTypeLV:     true,
		LanguageTypeMG:     true,
		LanguageTypeMH:     true,
		LanguageTypeMI:     true,
		LanguageTypeMK:     true,
		LanguageTypeML:     true,
		LanguageTypeMN:     true,
		LanguageTypeMO:     true,
		LanguageTypeMR:     true,
		LanguageTypeMS:     true,
		LanguageTypeMT:     true,
		LanguageTypeMY:     true,
		LanguageTypeNA:     true,
		LanguageTypeNB:     true,
		LanguageTypeND:     true,
		LanguageTypeNE:     true,
		LanguageTypeNG:     true,
		LanguageTypeNL:     true,
		LanguageTypeNN:     true,
		LanguageTypeNO:     true,
		LanguageTypeNR:     true,
		LanguageTypeNV:     true,
		LanguageTypeNY:     true,
		LanguageTypeOC:     true,
		LanguageTypeOJ:     true,
		LanguageTypeOM:     true,
		LanguageTypeOR:     true,
		LanguageTypeOS:     true,
		LanguageTypePA:     true,
		LanguageTypePI:     true,
		LanguageTypePL:     true,
		LanguageTypePS:     true,
		LanguageTypePT:     true,
		LanguageTypeQU:     true,
		LanguageTypeRM:     true,
		LanguageTypeRN:     true,
		LanguageTypeRO:     true,
		LanguageTypeRU:     true,
		LanguageTypeRW:     true,
		LanguageTypeSA:     true,
		LanguageTypeSC:     true,
		LanguageTypeSD:     true,
		LanguageTypeSE:     true,
		LanguageTypeSG:     true,
		LanguageTypeSH:     true,
		LanguageTypeSI:     true,
		LanguageTypeSK:     true,
		LanguageTypeSL:     true,
		LanguageTypeSM:     true,
		LanguageTypeSN:     true,
		LanguageTypeSO:     true,
		LanguageTypeSQ:     true,
		LanguageTypeSR:     true,
		LanguageTypeSS:     true,
		LanguageTypeST:     true,
		LanguageTypeSU:     true,
		LanguageTypeSV:     true,
		LanguageTypeSW:     true,
		LanguageTypeTA:     true,
		LanguageTypeTE:     true,
		LanguageTypeTG:     true,
		LanguageTypeTH:     true,
		LanguageTypeTI:     true,
		LanguageTypeTK:     true,
		LanguageTypeTL:     true,
		LanguageTypeTN:     true,
		LanguageTypeTO:     true,
		LanguageTypeTR:     true,
		LanguageTypeTS:     true,
		LanguageTypeTT:     true,
		LanguageTypeTW:     true,
		LanguageTypeTY:     true,
		LanguageTypeUG:     true,
		LanguageTypeUK:     true,
		LanguageTypeUR:     true,
		LanguageTypeUZ:     true,
		LanguageTypeVE:     true,
		LanguageTypeVI:     true,
		LanguageTypeVO:     true,
		LanguageTypeWA:     true,
		LanguageTypeWO:     true,
		LanguageTypeXH:     true,
		LanguageTypeYI:     true,
		LanguageTypeYO:     true,
		LanguageTypeZA:     true,
		LanguageTypeZH:     true,
		LanguageTypeZU:     true,
		LanguageTypeAAA:    true,
		LanguageTypeAAB:    true,
		LanguageTypeAAC:    true,
		LanguageTypeAAD:    true,
		LanguageTypeAAE:    true,
		LanguageTypeAAF:    true,
		LanguageTypeAAG:    true,
		LanguageTypeAAH:    true,
		LanguageTypeAAI:    true,
		LanguageTypeAAK:    true,
		LanguageTypeAAL:    true,
		LanguageTypeAAM:    true,
		LanguageTypeAAN:    true,
		LanguageTypeAAO:    true,
		LanguageTypeAAP:    true,
		LanguageTypeAAQ:    true,
		LanguageTypeAAS:    true,
		LanguageTypeAAT:    true,
		LanguageTypeAAU:    true,
		LanguageTypeAAV:    true,
		LanguageTypeAAW:    true,
		LanguageTypeAAX:    true,
		LanguageTypeAAZ:    true,
		LanguageTypeABA:    true,
		LanguageTypeABB:    true,
		LanguageTypeABC:    true,
		LanguageTypeABD:    true,
		LanguageTypeABE:    true,
		LanguageTypeABF:    true,
		LanguageTypeABG:    true,
		LanguageTypeABH:    true,
		LanguageTypeABI:    true,
		LanguageTypeABJ:    true,
		LanguageTypeABL:    true,
		LanguageTypeABM:    true,
		LanguageTypeABN:    true,
		LanguageTypeABO:    true,
		LanguageTypeABP:    true,
		LanguageTypeABQ:    true,
		LanguageTypeABR:    true,
		LanguageTypeABS:    true,
		LanguageTypeABT:    true,
		LanguageTypeABU:    true,
		LanguageTypeABV:    true,
		LanguageTypeABW:    true,
		LanguageTypeABX:    true,
		LanguageTypeABY:    true,
		LanguageTypeABZ:    true,
		LanguageTypeACA:    true,
		LanguageTypeACB:    true,
		LanguageTypeACD:    true,
		LanguageTypeACE:    true,
		LanguageTypeACF:    true,
		LanguageTypeACH:    true,
		LanguageTypeACI:    true,
		LanguageTypeACK:    true,
		LanguageTypeACL:    true,
		LanguageTypeACM:    true,
		LanguageTypeACN:    true,
		LanguageTypeACP:    true,
		LanguageTypeACQ:    true,
		LanguageTypeACR:    true,
		LanguageTypeACS:    true,
		LanguageTypeACT:    true,
		LanguageTypeACU:    true,
		LanguageTypeACV:    true,
		LanguageTypeACW:    true,
		LanguageTypeACX:    true,
		LanguageTypeACY:    true,
		LanguageTypeACZ:    true,
		LanguageTypeADA:    true,
		LanguageTypeADB:    true,
		LanguageTypeADD:    true,
		LanguageTypeADE:    true,
		LanguageTypeADF:    true,
		LanguageTypeADG:    true,
		LanguageTypeADH:    true,
		LanguageTypeADI:    true,
		LanguageTypeADJ:    true,
		LanguageTypeADL:    true,
		LanguageTypeADN:    true,
		LanguageTypeADO:    true,
		LanguageTypeADP:    true,
		LanguageTypeADQ:    true,
		LanguageTypeADR:    true,
		LanguageTypeADS:    true,
		LanguageTypeADT:    true,
		LanguageTypeADU:    true,
		LanguageTypeADW:    true,
		LanguageTypeADX:    true,
		LanguageTypeADY:    true,
		LanguageTypeADZ:    true,
		LanguageTypeAEA:    true,
		LanguageTypeAEB:    true,
		LanguageTypeAEC:    true,
		LanguageTypeAED:    true,
		LanguageTypeAEE:    true,
		LanguageTypeAEK:    true,
		LanguageTypeAEL:    true,
		LanguageTypeAEM:    true,
		LanguageTypeAEN:    true,
		LanguageTypeAEQ:    true,
		LanguageTypeAER:    true,
		LanguageTypeAES:    true,
		LanguageTypeAEU:    true,
		LanguageTypeAEW:    true,
		LanguageTypeAEY:    true,
		LanguageTypeAEZ:    true,
		LanguageTypeAFA:    true,
		LanguageTypeAFB:    true,
		LanguageTypeAFD:    true,
		LanguageTypeAFE:    true,
		LanguageTypeAFG:    true,
		LanguageTypeAFH:    true,
		LanguageTypeAFI:    true,
		LanguageTypeAFK:    true,
		LanguageTypeAFN:    true,
		LanguageTypeAFO:    true,
		LanguageTypeAFP:    true,
		LanguageTypeAFS:    true,
		LanguageTypeAFT:    true,
		LanguageTypeAFU:    true,
		LanguageTypeAFZ:    true,
		LanguageTypeAGA:    true,
		LanguageTypeAGB:    true,
		LanguageTypeAGC:    true,
		LanguageTypeAGD:    true,
		LanguageTypeAGE:    true,
		LanguageTypeAGF:    true,
		LanguageTypeAGG:    true,
		LanguageTypeAGH:    true,
		LanguageTypeAGI:    true,
		LanguageTypeAGJ:    true,
		LanguageTypeAGK:    true,
		LanguageTypeAGL:    true,
		LanguageTypeAGM:    true,
		LanguageTypeAGN:    true,
		LanguageTypeAGO:    true,
		LanguageTypeAGP:    true,
		LanguageTypeAGQ:    true,
		LanguageTypeAGR:    true,
		LanguageTypeAGS:    true,
		LanguageTypeAGT:    true,
		LanguageTypeAGU:    true,
		LanguageTypeAGV:    true,
		LanguageTypeAGW:    true,
		LanguageTypeAGX:    true,
		LanguageTypeAGY:    true,
		LanguageTypeAGZ:    true,
		LanguageTypeAHA:    true,
		LanguageTypeAHB:    true,
		LanguageTypeAHG:    true,
		LanguageTypeAHH:    true,
		LanguageTypeAHI:    true,
		LanguageTypeAHK:    true,
		LanguageTypeAHL:    true,
		LanguageTypeAHM:    true,
		LanguageTypeAHN:    true,
		LanguageTypeAHO:    true,
		LanguageTypeAHP:    true,
		LanguageTypeAHR:    true,
		LanguageTypeAHS:    true,
		LanguageTypeAHT:    true,
		LanguageTypeAIA:    true,
		LanguageTypeAIB:    true,
		LanguageTypeAIC:    true,
		LanguageTypeAID:    true,
		LanguageTypeAIE:    true,
		LanguageTypeAIF:    true,
		LanguageTypeAIG:    true,
		LanguageTypeAIH:    true,
		LanguageTypeAII:    true,
		LanguageTypeAIJ:    true,
		LanguageTypeAIK:    true,
		LanguageTypeAIL:    true,
		LanguageTypeAIM:    true,
		LanguageTypeAIN:    true,
		LanguageTypeAIO:    true,
		LanguageTypeAIP:    true,
		LanguageTypeAIQ:    true,
		LanguageTypeAIR:    true,
		LanguageTypeAIS:    true,
		LanguageTypeAIT:    true,
		LanguageTypeAIW:    true,
		LanguageTypeAIX:    true,
		LanguageTypeAIY:    true,
		LanguageTypeAJA:    true,
		LanguageTypeAJG:    true,
		LanguageTypeAJI:    true,
		LanguageTypeAJN:    true,
		LanguageTypeAJP:    true,
		LanguageTypeAJT:    true,
		LanguageTypeAJU:    true,
		LanguageTypeAJW:    true,
		LanguageTypeAJZ:    true,
		LanguageTypeAKB:    true,
		LanguageTypeAKC:    true,
		LanguageTypeAKD:    true,
		LanguageTypeAKE:    true,
		LanguageTypeAKF:    true,
		LanguageTypeAKG:    true,
		LanguageTypeAKH:    true,
		LanguageTypeAKI:    true,
		LanguageTypeAKJ:    true,
		LanguageTypeAKK:    true,
		LanguageTypeAKL:    true,
		LanguageTypeAKM:    true,
		LanguageTypeAKO:    true,
		LanguageTypeAKP:    true,
		LanguageTypeAKQ:    true,
		LanguageTypeAKR:    true,
		LanguageTypeAKS:    true,
		LanguageTypeAKT:    true,
		LanguageTypeAKU:    true,
		LanguageTypeAKV:    true,
		LanguageTypeAKW:    true,
		LanguageTypeAKX:    true,
		LanguageTypeAKY:    true,
		LanguageTypeAKZ:    true,
		LanguageTypeALA:    true,
		LanguageTypeALC:    true,
		LanguageTypeALD:    true,
		LanguageTypeALE:    true,
		LanguageTypeALF:    true,
		LanguageTypeALG:    true,
		LanguageTypeALH:    true,
		LanguageTypeALI:    true,
		LanguageTypeALJ:    true,
		LanguageTypeALK:    true,
		LanguageTypeALL:    true,
		LanguageTypeALM:    true,
		LanguageTypeALN:    true,
		LanguageTypeALO:    true,
		LanguageTypeALP:    true,
		LanguageTypeALQ:    true,
		LanguageTypeALR:    true,
		LanguageTypeALS:    true,
		LanguageTypeALT:    true,
		LanguageTypeALU:    true,
		LanguageTypeALV:    true,
		LanguageTypeALW:    true,
		LanguageTypeALX:    true,
		LanguageTypeALY:    true,
		LanguageTypeALZ:    true,
		LanguageTypeAMA:    true,
		LanguageTypeAMB:    true,
		LanguageTypeAMC:    true,
		LanguageTypeAME:    true,
		LanguageTypeAMF:    true,
		LanguageTypeAMG:    true,
		LanguageTypeAMI:    true,
		LanguageTypeAMJ:    true,
		LanguageTypeAMK:    true,
		LanguageTypeAML:    true,
		LanguageTypeAMM:    true,
		LanguageTypeAMN:    true,
		LanguageTypeAMO:    true,
		LanguageTypeAMP:    true,
		LanguageTypeAMQ:    true,
		LanguageTypeAMR:    true,
		LanguageTypeAMS:    true,
		LanguageTypeAMT:    true,
		LanguageTypeAMU:    true,
		LanguageTypeAMV:    true,
		LanguageTypeAMW:    true,
		LanguageTypeAMX:    true,
		LanguageTypeAMY:    true,
		LanguageTypeAMZ:    true,
		LanguageTypeANA:    true,
		LanguageTypeANB:    true,
		LanguageTypeANC:    true,
		LanguageTypeAND:    true,
		LanguageTypeANE:    true,
		LanguageTypeANF:    true,
		LanguageTypeANG:    true,
		LanguageTypeANH:    true,
		LanguageTypeANI:    true,
		LanguageTypeANJ:    true,
		LanguageTypeANK:    true,
		LanguageTypeANL:    true,
		LanguageTypeANM:    true,
		LanguageTypeANN:    true,
		LanguageTypeANO:    true,
		LanguageTypeANP:    true,
		LanguageTypeANQ:    true,
		LanguageTypeANR:    true,
		LanguageTypeANS:    true,
		LanguageTypeANT:    true,
		LanguageTypeANU:    true,
		LanguageTypeANV:    true,
		LanguageTypeANW:    true,
		LanguageTypeANX:    true,
		LanguageTypeANY:    true,
		LanguageTypeANZ:    true,
		LanguageTypeAOA:    true,
		LanguageTypeAOB:    true,
		LanguageTypeAOC:    true,
		LanguageTypeAOD:    true,
		LanguageTypeAOE:    true,
		LanguageTypeAOF:    true,
		LanguageTypeAOG:    true,
		LanguageTypeAOH:    true,
		LanguageTypeAOI:    true,
		LanguageTypeAOJ:    true,
		LanguageTypeAOK:    true,
		LanguageTypeAOL:    true,
		LanguageTypeAOM:    true,
		LanguageTypeAON:    true,
		LanguageTypeAOR:    true,
		LanguageTypeAOS:    true,
		LanguageTypeAOT:    true,
		LanguageTypeAOU:    true,
		LanguageTypeAOX:    true,
		LanguageTypeAOZ:    true,
		LanguageTypeAPA:    true,
		LanguageTypeAPB:    true,
		LanguageTypeAPC:    true,
		LanguageTypeAPD:    true,
		LanguageTypeAPE:    true,
		LanguageTypeAPF:    true,
		LanguageTypeAPG:    true,
		LanguageTypeAPH:    true,
		LanguageTypeAPI:    true,
		LanguageTypeAPJ:    true,
		LanguageTypeAPK:    true,
		LanguageTypeAPL:    true,
		LanguageTypeAPM:    true,
		LanguageTypeAPN:    true,
		LanguageTypeAPO:    true,
		LanguageTypeAPP:    true,
		LanguageTypeAPQ:    true,
		LanguageTypeAPR:    true,
		LanguageTypeAPS:    true,
		LanguageTypeAPT:    true,
		LanguageTypeAPU:    true,
		LanguageTypeAPV:    true,
		LanguageTypeAPW:    true,
		LanguageTypeAPX:    true,
		LanguageTypeAPY:    true,
		LanguageTypeAPZ:    true,
		LanguageTypeAQA:    true,
		LanguageTypeAQC:    true,
		LanguageTypeAQD:    true,
		LanguageTypeAQG:    true,
		LanguageTypeAQL:    true,
		LanguageTypeAQM:    true,
		LanguageTypeAQN:    true,
		LanguageTypeAQP:    true,
		LanguageTypeAQR:    true,
		LanguageTypeAQZ:    true,
		LanguageTypeARB:    true,
		LanguageTypeARC:    true,
		LanguageTypeARD:    true,
		LanguageTypeARE:    true,
		LanguageTypeARH:    true,
		LanguageTypeARI:    true,
		LanguageTypeARJ:    true,
		LanguageTypeARK:    true,
		LanguageTypeARL:    true,
		LanguageTypeARN:    true,
		LanguageTypeARO:    true,
		LanguageTypeARP:    true,
		LanguageTypeARQ:    true,
		LanguageTypeARR:    true,
		LanguageTypeARS:    true,
		LanguageTypeART:    true,
		LanguageTypeARU:    true,
		LanguageTypeARV:    true,
		LanguageTypeARW:    true,
		LanguageTypeARX:    true,
		LanguageTypeARY:    true,
		LanguageTypeARZ:    true,
		LanguageTypeASA:    true,
		LanguageTypeASB:    true,
		LanguageTypeASC:    true,
		LanguageTypeASD:    true,
		LanguageTypeASE:    true,
		LanguageTypeASF:    true,
		LanguageTypeASG:    true,
		LanguageTypeASH:    true,
		LanguageTypeASI:    true,
		LanguageTypeASJ:    true,
		LanguageTypeASK:    true,
		LanguageTypeASL:    true,
		LanguageTypeASN:    true,
		LanguageTypeASO:    true,
		LanguageTypeASP:    true,
		LanguageTypeASQ:    true,
		LanguageTypeASR:    true,
		LanguageTypeASS:    true,
		LanguageTypeAST:    true,
		LanguageTypeASU:    true,
		LanguageTypeASV:    true,
		LanguageTypeASW:    true,
		LanguageTypeASX:    true,
		LanguageTypeASY:    true,
		LanguageTypeASZ:    true,
		LanguageTypeATA:    true,
		LanguageTypeATB:    true,
		LanguageTypeATC:    true,
		LanguageTypeATD:    true,
		LanguageTypeATE:    true,
		LanguageTypeATG:    true,
		LanguageTypeATH:    true,
		LanguageTypeATI:    true,
		LanguageTypeATJ:    true,
		LanguageTypeATK:    true,
		LanguageTypeATL:    true,
		LanguageTypeATM:    true,
		LanguageTypeATN:    true,
		LanguageTypeATO:    true,
		LanguageTypeATP:    true,
		LanguageTypeATQ:    true,
		LanguageTypeATR:    true,
		LanguageTypeATS:    true,
		LanguageTypeATT:    true,
		LanguageTypeATU:    true,
		LanguageTypeATV:    true,
		LanguageTypeATW:    true,
		LanguageTypeATX:    true,
		LanguageTypeATY:    true,
		LanguageTypeATZ:    true,
		LanguageTypeAUA:    true,
		LanguageTypeAUB:    true,
		LanguageTypeAUC:    true,
		LanguageTypeAUD:    true,
		LanguageTypeAUE:    true,
		LanguageTypeAUF:    true,
		LanguageTypeAUG:    true,
		LanguageTypeAUH:    true,
		LanguageTypeAUI:    true,
		LanguageTypeAUJ:    true,
		LanguageTypeAUK:    true,
		LanguageTypeAUL:    true,
		LanguageTypeAUM:    true,
		LanguageTypeAUN:    true,
		LanguageTypeAUO:    true,
		LanguageTypeAUP:    true,
		LanguageTypeAUQ:    true,
		LanguageTypeAUR:    true,
		LanguageTypeAUS:    true,
		LanguageTypeAUT:    true,
		LanguageTypeAUU:    true,
		LanguageTypeAUW:    true,
		LanguageTypeAUX:    true,
		LanguageTypeAUY:    true,
		LanguageTypeAUZ:    true,
		LanguageTypeAVB:    true,
		LanguageTypeAVD:    true,
		LanguageTypeAVI:    true,
		LanguageTypeAVK:    true,
		LanguageTypeAVL:    true,
		LanguageTypeAVM:    true,
		LanguageTypeAVN:    true,
		LanguageTypeAVO:    true,
		LanguageTypeAVS:    true,
		LanguageTypeAVT:    true,
		LanguageTypeAVU:    true,
		LanguageTypeAVV:    true,
		LanguageTypeAWA:    true,
		LanguageTypeAWB:    true,
		LanguageTypeAWC:    true,
		LanguageTypeAWD:    true,
		LanguageTypeAWE:    true,
		LanguageTypeAWG:    true,
		LanguageTypeAWH:    true,
		LanguageTypeAWI:    true,
		LanguageTypeAWK:    true,
		LanguageTypeAWM:    true,
		LanguageTypeAWN:    true,
		LanguageTypeAWO:    true,
		LanguageTypeAWR:    true,
		LanguageTypeAWS:    true,
		LanguageTypeAWT:    true,
		LanguageTypeAWU:    true,
		LanguageTypeAWV:    true,
		LanguageTypeAWW:    true,
		LanguageTypeAWX:    true,
		LanguageTypeAWY:    true,
		LanguageTypeAXB:    true,
		LanguageTypeAXE:    true,
		LanguageTypeAXG:    true,
		LanguageTypeAXK:    true,
		LanguageTypeAXL:    true,
		LanguageTypeAXM:    true,
		LanguageTypeAXX:    true,
		LanguageTypeAYA:    true,
		LanguageTypeAYB:    true,
		LanguageTypeAYC:    true,
		LanguageTypeAYD:    true,
		LanguageTypeAYE:    true,
		LanguageTypeAYG:    true,
		LanguageTypeAYH:    true,
		LanguageTypeAYI:    true,
		LanguageTypeAYK:    true,
		LanguageTypeAYL:    true,
		LanguageTypeAYN:    true,
		LanguageTypeAYO:    true,
		LanguageTypeAYP:    true,
		LanguageTypeAYQ:    true,
		LanguageTypeAYR:    true,
		LanguageTypeAYS:    true,
		LanguageTypeAYT:    true,
		LanguageTypeAYU:    true,
		LanguageTypeAYX:    true,
		LanguageTypeAYY:    true,
		LanguageTypeAYZ:    true,
		LanguageTypeAZA:    true,
		LanguageTypeAZB:    true,
		LanguageTypeAZC:    true,
		LanguageTypeAZD:    true,
		LanguageTypeAZG:    true,
		LanguageTypeAZJ:    true,
		LanguageTypeAZM:    true,
		LanguageTypeAZN:    true,
		LanguageTypeAZO:    true,
		LanguageTypeAZT:    true,
		LanguageTypeAZZ:    true,
		LanguageTypeBAA:    true,
		LanguageTypeBAB:    true,
		LanguageTypeBAC:    true,
		LanguageTypeBAD:    true,
		LanguageTypeBAE:    true,
		LanguageTypeBAF:    true,
		LanguageTypeBAG:    true,
		LanguageTypeBAH:    true,
		LanguageTypeBAI:    true,
		LanguageTypeBAJ:    true,
		LanguageTypeBAL:    true,
		LanguageTypeBAN:    true,
		LanguageTypeBAO:    true,
		LanguageTypeBAP:    true,
		LanguageTypeBAR:    true,
		LanguageTypeBAS:    true,
		LanguageTypeBAT:    true,
		LanguageTypeBAU:    true,
		LanguageTypeBAV:    true,
		LanguageTypeBAW:    true,
		LanguageTypeBAX:    true,
		LanguageTypeBAY:    true,
		LanguageTypeBAZ:    true,
		LanguageTypeBBA:    true,
		LanguageTypeBBB:    true,
		LanguageTypeBBC:    true,
		LanguageTypeBBD:    true,
		LanguageTypeBBE:    true,
		LanguageTypeBBF:    true,
		LanguageTypeBBG:    true,
		LanguageTypeBBH:    true,
		LanguageTypeBBI:    true,
		LanguageTypeBBJ:    true,
		LanguageTypeBBK:    true,
		LanguageTypeBBL:    true,
		LanguageTypeBBM:    true,
		LanguageTypeBBN:    true,
		LanguageTypeBBO:    true,
		LanguageTypeBBP:    true,
		LanguageTypeBBQ:    true,
		LanguageTypeBBR:    true,
		LanguageTypeBBS:    true,
		LanguageTypeBBT:    true,
		LanguageTypeBBU:    true,
		LanguageTypeBBV:    true,
		LanguageTypeBBW:    true,
		LanguageTypeBBX:    true,
		LanguageTypeBBY:    true,
		LanguageTypeBBZ:    true,
		LanguageTypeBCA:    true,
		LanguageTypeBCB:    true,
		LanguageTypeBCC:    true,
		LanguageTypeBCD:    true,
		LanguageTypeBCE:    true,
		LanguageTypeBCF:    true,
		LanguageTypeBCG:    true,
		LanguageTypeBCH:    true,
		LanguageTypeBCI:    true,
		LanguageTypeBCJ:    true,
		LanguageTypeBCK:    true,
		LanguageTypeBCL:    true,
		LanguageTypeBCM:    true,
		LanguageTypeBCN:    true,
		LanguageTypeBCO:    true,
		LanguageTypeBCP:    true,
		LanguageTypeBCQ:    true,
		LanguageTypeBCR:    true,
		LanguageTypeBCS:    true,
		LanguageTypeBCT:    true,
		LanguageTypeBCU:    true,
		LanguageTypeBCV:    true,
		LanguageTypeBCW:    true,
		LanguageTypeBCY:    true,
		LanguageTypeBCZ:    true,
		LanguageTypeBDA:    true,
		LanguageTypeBDB:    true,
		LanguageTypeBDC:    true,
		LanguageTypeBDD:    true,
		LanguageTypeBDE:    true,
		LanguageTypeBDF:    true,
		LanguageTypeBDG:    true,
		LanguageTypeBDH:    true,
		LanguageTypeBDI:    true,
		LanguageTypeBDJ:    true,
		LanguageTypeBDK:    true,
		LanguageTypeBDL:    true,
		LanguageTypeBDM:    true,
		LanguageTypeBDN:    true,
		LanguageTypeBDO:    true,
		LanguageTypeBDP:    true,
		LanguageTypeBDQ:    true,
		LanguageTypeBDR:    true,
		LanguageTypeBDS:    true,
		LanguageTypeBDT:    true,
		LanguageTypeBDU:    true,
		LanguageTypeBDV:    true,
		LanguageTypeBDW:    true,
		LanguageTypeBDX:    true,
		LanguageTypeBDY:    true,
		LanguageTypeBDZ:    true,
		LanguageTypeBEA:    true,
		LanguageTypeBEB:    true,
		LanguageTypeBEC:    true,
		LanguageTypeBED:    true,
		LanguageTypeBEE:    true,
		LanguageTypeBEF:    true,
		LanguageTypeBEG:    true,
		LanguageTypeBEH:    true,
		LanguageTypeBEI:    true,
		LanguageTypeBEJ:    true,
		LanguageTypeBEK:    true,
		LanguageTypeBEM:    true,
		LanguageTypeBEO:    true,
		LanguageTypeBEP:    true,
		LanguageTypeBEQ:    true,
		LanguageTypeBER:    true,
		LanguageTypeBES:    true,
		LanguageTypeBET:    true,
		LanguageTypeBEU:    true,
		LanguageTypeBEV:    true,
		LanguageTypeBEW:    true,
		LanguageTypeBEX:    true,
		LanguageTypeBEY:    true,
		LanguageTypeBEZ:    true,
		LanguageTypeBFA:    true,
		LanguageTypeBFB:    true,
		LanguageTypeBFC:    true,
		LanguageTypeBFD:    true,
		LanguageTypeBFE:    true,
		LanguageTypeBFF:    true,
		LanguageTypeBFG:    true,
		LanguageTypeBFH:    true,
		LanguageTypeBFI:    true,
		LanguageTypeBFJ:    true,
		LanguageTypeBFK:    true,
		LanguageTypeBFL:    true,
		LanguageTypeBFM:    true,
		LanguageTypeBFN:    true,
		LanguageTypeBFO:    true,
		LanguageTypeBFP:    true,
		LanguageTypeBFQ:    true,
		LanguageTypeBFR:    true,
		LanguageTypeBFS:    true,
		LanguageTypeBFT:    true,
		LanguageTypeBFU:    true,
		LanguageTypeBFW:    true,
		LanguageTypeBFX:    true,
		LanguageTypeBFY:    true,
		LanguageTypeBFZ:    true,
		LanguageTypeBGA:    true,
		LanguageTypeBGB:    true,
		LanguageTypeBGC:    true,
		LanguageTypeBGD:    true,
		LanguageTypeBGE:    true,
		LanguageTypeBGF:    true,
		LanguageTypeBGG:    true,
		LanguageTypeBGI:    true,
		LanguageTypeBGJ:    true,
		LanguageTypeBGK:    true,
		LanguageTypeBGL:    true,
		LanguageTypeBGM:    true,
		LanguageTypeBGN:    true,
		LanguageTypeBGO:    true,
		LanguageTypeBGP:    true,
		LanguageTypeBGQ:    true,
		LanguageTypeBGR:    true,
		LanguageTypeBGS:    true,
		LanguageTypeBGT:    true,
		LanguageTypeBGU:    true,
		LanguageTypeBGV:    true,
		LanguageTypeBGW:    true,
		LanguageTypeBGX:    true,
		LanguageTypeBGY:    true,
		LanguageTypeBGZ:    true,
		LanguageTypeBHA:    true,
		LanguageTypeBHB:    true,
		LanguageTypeBHC:    true,
		LanguageTypeBHD:    true,
		LanguageTypeBHE:    true,
		LanguageTypeBHF:    true,
		LanguageTypeBHG:    true,
		LanguageTypeBHH:    true,
		LanguageTypeBHI:    true,
		LanguageTypeBHJ:    true,
		LanguageTypeBHK:    true,
		LanguageTypeBHL:    true,
		LanguageTypeBHM:    true,
		LanguageTypeBHN:    true,
		LanguageTypeBHO:    true,
		LanguageTypeBHP:    true,
		LanguageTypeBHQ:    true,
		LanguageTypeBHR:    true,
		LanguageTypeBHS:    true,
		LanguageTypeBHT:    true,
		LanguageTypeBHU:    true,
		LanguageTypeBHV:    true,
		LanguageTypeBHW:    true,
		LanguageTypeBHX:    true,
		LanguageTypeBHY:    true,
		LanguageTypeBHZ:    true,
		LanguageTypeBIA:    true,
		LanguageTypeBIB:    true,
		LanguageTypeBIC:    true,
		LanguageTypeBID:    true,
		LanguageTypeBIE:    true,
		LanguageTypeBIF:    true,
		LanguageTypeBIG:    true,
		LanguageTypeBIJ:    true,
		LanguageTypeBIK:    true,
		LanguageTypeBIL:    true,
		LanguageTypeBIM:    true,
		LanguageTypeBIN:    true,
		LanguageTypeBIO:    true,
		LanguageTypeBIP:    true,
		LanguageTypeBIQ:    true,
		LanguageTypeBIR:    true,
		LanguageTypeBIT:    true,
		LanguageTypeBIU:    true,
		LanguageTypeBIV:    true,
		LanguageTypeBIW:    true,
		LanguageTypeBIX:    true,
		LanguageTypeBIY:    true,
		LanguageTypeBIZ:    true,
		LanguageTypeBJA:    true,
		LanguageTypeBJB:    true,
		LanguageTypeBJC:    true,
		LanguageTypeBJD:    true,
		LanguageTypeBJE:    true,
		LanguageTypeBJF:    true,
		LanguageTypeBJG:    true,
		LanguageTypeBJH:    true,
		LanguageTypeBJI:    true,
		LanguageTypeBJJ:    true,
		LanguageTypeBJK:    true,
		LanguageTypeBJL:    true,
		LanguageTypeBJM:    true,
		LanguageTypeBJN:    true,
		LanguageTypeBJO:    true,
		LanguageTypeBJP:    true,
		LanguageTypeBJQ:    true,
		LanguageTypeBJR:    true,
		LanguageTypeBJS:    true,
		LanguageTypeBJT:    true,
		LanguageTypeBJU:    true,
		LanguageTypeBJV:    true,
		LanguageTypeBJW:    true,
		LanguageTypeBJX:    true,
		LanguageTypeBJY:    true,
		LanguageTypeBJZ:    true,
		LanguageTypeBKA:    true,
		LanguageTypeBKB:    true,
		LanguageTypeBKC:    true,
		LanguageTypeBKD:    true,
		LanguageTypeBKF:    true,
		LanguageTypeBKG:    true,
		LanguageTypeBKH:    true,
		LanguageTypeBKI:    true,
		LanguageTypeBKJ:    true,
		LanguageTypeBKK:    true,
		LanguageTypeBKL:    true,
		LanguageTypeBKM:    true,
		LanguageTypeBKN:    true,
		LanguageTypeBKO:    true,
		LanguageTypeBKP:    true,
		LanguageTypeBKQ:    true,
		LanguageTypeBKR:    true,
		LanguageTypeBKS:    true,
		LanguageTypeBKT:    true,
		LanguageTypeBKU:    true,
		LanguageTypeBKV:    true,
		LanguageTypeBKW:    true,
		LanguageTypeBKX:    true,
		LanguageTypeBKY:    true,
		LanguageTypeBKZ:    true,
		LanguageTypeBLA:    true,
		LanguageTypeBLB:    true,
		LanguageTypeBLC:    true,
		LanguageTypeBLD:    true,
		LanguageTypeBLE:    true,
		LanguageTypeBLF:    true,
		LanguageTypeBLG:    true,
		LanguageTypeBLH:    true,
		LanguageTypeBLI:    true,
		LanguageTypeBLJ:    true,
		LanguageTypeBLK:    true,
		LanguageTypeBLL:    true,
		LanguageTypeBLM:    true,
		LanguageTypeBLN:    true,
		LanguageTypeBLO:    true,
		LanguageTypeBLP:    true,
		LanguageTypeBLQ:    true,
		LanguageTypeBLR:    true,
		LanguageTypeBLS:    true,
		LanguageTypeBLT:    true,
		LanguageTypeBLV:    true,
		LanguageTypeBLW:    true,
		LanguageTypeBLX:    true,
		LanguageTypeBLY:    true,
		LanguageTypeBLZ:    true,
		LanguageTypeBMA:    true,
		LanguageTypeBMB:    true,
		LanguageTypeBMC:    true,
		LanguageTypeBMD:    true,
		LanguageTypeBME:    true,
		LanguageTypeBMF:    true,
		LanguageTypeBMG:    true,
		LanguageTypeBMH:    true,
		LanguageTypeBMI:    true,
		LanguageTypeBMJ:    true,
		LanguageTypeBMK:    true,
		LanguageTypeBML:    true,
		LanguageTypeBMM:    true,
		LanguageTypeBMN:    true,
		LanguageTypeBMO:    true,
		LanguageTypeBMP:    true,
		LanguageTypeBMQ:    true,
		LanguageTypeBMR:    true,
		LanguageTypeBMS:    true,
		LanguageTypeBMT:    true,
		LanguageTypeBMU:    true,
		LanguageTypeBMV:    true,
		LanguageTypeBMW:    true,
		LanguageTypeBMX:    true,
		LanguageTypeBMY:    true,
		LanguageTypeBMZ:    true,
		LanguageTypeBNA:    true,
		LanguageTypeBNB:    true,
		LanguageTypeBNC:    true,
		LanguageTypeBND:    true,
		LanguageTypeBNE:    true,
		LanguageTypeBNF:    true,
		LanguageTypeBNG:    true,
		LanguageTypeBNI:    true,
		LanguageTypeBNJ:    true,
		LanguageTypeBNK:    true,
		LanguageTypeBNL:    true,
		LanguageTypeBNM:    true,
		LanguageTypeBNN:    true,
		LanguageTypeBNO:    true,
		LanguageTypeBNP:    true,
		LanguageTypeBNQ:    true,
		LanguageTypeBNR:    true,
		LanguageTypeBNS:    true,
		LanguageTypeBNT:    true,
		LanguageTypeBNU:    true,
		LanguageTypeBNV:    true,
		LanguageTypeBNW:    true,
		LanguageTypeBNX:    true,
		LanguageTypeBNY:    true,
		LanguageTypeBNZ:    true,
		LanguageTypeBOA:    true,
		LanguageTypeBOB:    true,
		LanguageTypeBOE:    true,
		LanguageTypeBOF:    true,
		LanguageTypeBOG:    true,
		LanguageTypeBOH:    true,
		LanguageTypeBOI:    true,
		LanguageTypeBOJ:    true,
		LanguageTypeBOK:    true,
		LanguageTypeBOL:    true,
		LanguageTypeBOM:    true,
		LanguageTypeBON:    true,
		LanguageTypeBOO:    true,
		LanguageTypeBOP:    true,
		LanguageTypeBOQ:    true,
		LanguageTypeBOR:    true,
		LanguageTypeBOT:    true,
		LanguageTypeBOU:    true,
		LanguageTypeBOV:    true,
		LanguageTypeBOW:    true,
		LanguageTypeBOX:    true,
		LanguageTypeBOY:    true,
		LanguageTypeBOZ:    true,
		LanguageTypeBPA:    true,
		LanguageTypeBPB:    true,
		LanguageTypeBPD:    true,
		LanguageTypeBPG:    true,
		LanguageTypeBPH:    true,
		LanguageTypeBPI:    true,
		LanguageTypeBPJ:    true,
		LanguageTypeBPK:    true,
		LanguageTypeBPL:    true,
		LanguageTypeBPM:    true,
		LanguageTypeBPN:    true,
		LanguageTypeBPO:    true,
		LanguageTypeBPP:    true,
		LanguageTypeBPQ:    true,
		LanguageTypeBPR:    true,
		LanguageTypeBPS:    true,
		LanguageTypeBPT:    true,
		LanguageTypeBPU:    true,
		LanguageTypeBPV:    true,
		LanguageTypeBPW:    true,
		LanguageTypeBPX:    true,
		LanguageTypeBPY:    true,
		LanguageTypeBPZ:    true,
		LanguageTypeBQA:    true,
		LanguageTypeBQB:    true,
		LanguageTypeBQC:    true,
		LanguageTypeBQD:    true,
		LanguageTypeBQF:    true,
		LanguageTypeBQG:    true,
		LanguageTypeBQH:    true,
		LanguageTypeBQI:    true,
		LanguageTypeBQJ:    true,
		LanguageTypeBQK:    true,
		LanguageTypeBQL:    true,
		LanguageTypeBQM:    true,
		LanguageTypeBQN:    true,
		LanguageTypeBQO:    true,
		LanguageTypeBQP:    true,
		LanguageTypeBQQ:    true,
		LanguageTypeBQR:    true,
		LanguageTypeBQS:    true,
		LanguageTypeBQT:    true,
		LanguageTypeBQU:    true,
		LanguageTypeBQV:    true,
		LanguageTypeBQW:    true,
		LanguageTypeBQX:    true,
		LanguageTypeBQY:    true,
		LanguageTypeBQZ:    true,
		LanguageTypeBRA:    true,
		LanguageTypeBRB:    true,
		LanguageTypeBRC:    true,
		LanguageTypeBRD:    true,
		LanguageTypeBRF:    true,
		LanguageTypeBRG:    true,
		LanguageTypeBRH:    true,
		LanguageTypeBRI:    true,
		LanguageTypeBRJ:    true,
		LanguageTypeBRK:    true,
		LanguageTypeBRL:    true,
		LanguageTypeBRM:    true,
		LanguageTypeBRN:    true,
		LanguageTypeBRO:    true,
		LanguageTypeBRP:    true,
		LanguageTypeBRQ:    true,
		LanguageTypeBRR:    true,
		LanguageTypeBRS:    true,
		LanguageTypeBRT:    true,
		LanguageTypeBRU:    true,
		LanguageTypeBRV:    true,
		LanguageTypeBRW:    true,
		LanguageTypeBRX:    true,
		LanguageTypeBRY:    true,
		LanguageTypeBRZ:    true,
		LanguageTypeBSA:    true,
		LanguageTypeBSB:    true,
		LanguageTypeBSC:    true,
		LanguageTypeBSE:    true,
		LanguageTypeBSF:    true,
		LanguageTypeBSG:    true,
		LanguageTypeBSH:    true,
		LanguageTypeBSI:    true,
		LanguageTypeBSJ:    true,
		LanguageTypeBSK:    true,
		LanguageTypeBSL:    true,
		LanguageTypeBSM:    true,
		LanguageTypeBSN:    true,
		LanguageTypeBSO:    true,
		LanguageTypeBSP:    true,
		LanguageTypeBSQ:    true,
		LanguageTypeBSR:    true,
		LanguageTypeBSS:    true,
		LanguageTypeBST:    true,
		LanguageTypeBSU:    true,
		LanguageTypeBSV:    true,
		LanguageTypeBSW:    true,
		LanguageTypeBSX:    true,
		LanguageTypeBSY:    true,
		LanguageTypeBTA:    true,
		LanguageTypeBTB:    true,
		LanguageTypeBTC:    true,
		LanguageTypeBTD:    true,
		LanguageTypeBTE:    true,
		LanguageTypeBTF:    true,
		LanguageTypeBTG:    true,
		LanguageTypeBTH:    true,
		LanguageTypeBTI:    true,
		LanguageTypeBTJ:    true,
		LanguageTypeBTK:    true,
		LanguageTypeBTL:    true,
		LanguageTypeBTM:    true,
		LanguageTypeBTN:    true,
		LanguageTypeBTO:    true,
		LanguageTypeBTP:    true,
		LanguageTypeBTQ:    true,
		LanguageTypeBTR:    true,
		LanguageTypeBTS:    true,
		LanguageTypeBTT:    true,
		LanguageTypeBTU:    true,
		LanguageTypeBTV:    true,
		LanguageTypeBTW:    true,
		LanguageTypeBTX:    true,
		LanguageTypeBTY:    true,
		LanguageTypeBTZ:    true,
		LanguageTypeBUA:    true,
		LanguageTypeBUB:    true,
		LanguageTypeBUC:    true,
		LanguageTypeBUD:    true,
		LanguageTypeBUE:    true,
		LanguageTypeBUF:    true,
		LanguageTypeBUG:    true,
		LanguageTypeBUH:    true,
		LanguageTypeBUI:    true,
		LanguageTypeBUJ:    true,
		LanguageTypeBUK:    true,
		LanguageTypeBUM:    true,
		LanguageTypeBUN:    true,
		LanguageTypeBUO:    true,
		LanguageTypeBUP:    true,
		LanguageTypeBUQ:    true,
		LanguageTypeBUS:    true,
		LanguageTypeBUT:    true,
		LanguageTypeBUU:    true,
		LanguageTypeBUV:    true,
		LanguageTypeBUW:    true,
		LanguageTypeBUX:    true,
		LanguageTypeBUY:    true,
		LanguageTypeBUZ:    true,
		LanguageTypeBVA:    true,
		LanguageTypeBVB:    true,
		LanguageTypeBVC:    true,
		LanguageTypeBVD:    true,
		LanguageTypeBVE:    true,
		LanguageTypeBVF:    true,
		LanguageTypeBVG:    true,
		LanguageTypeBVH:    true,
		LanguageTypeBVI:    true,
		LanguageTypeBVJ:    true,
		LanguageTypeBVK:    true,
		LanguageTypeBVL:    true,
		LanguageTypeBVM:    true,
		LanguageTypeBVN:    true,
		LanguageTypeBVO:    true,
		LanguageTypeBVP:    true,
		LanguageTypeBVQ:    true,
		LanguageTypeBVR:    true,
		LanguageTypeBVT:    true,
		LanguageTypeBVU:    true,
		LanguageTypeBVV:    true,
		LanguageTypeBVW:    true,
		LanguageTypeBVX:    true,
		LanguageTypeBVY:    true,
		LanguageTypeBVZ:    true,
		LanguageTypeBWA:    true,
		LanguageTypeBWB:    true,
		LanguageTypeBWC:    true,
		LanguageTypeBWD:    true,
		LanguageTypeBWE:    true,
		LanguageTypeBWF:    true,
		LanguageTypeBWG:    true,
		LanguageTypeBWH:    true,
		LanguageTypeBWI:    true,
		LanguageTypeBWJ:    true,
		LanguageTypeBWK:    true,
		LanguageTypeBWL:    true,
		LanguageTypeBWM:    true,
		LanguageTypeBWN:    true,
		LanguageTypeBWO:    true,
		LanguageTypeBWP:    true,
		LanguageTypeBWQ:    true,
		LanguageTypeBWR:    true,
		LanguageTypeBWS:    true,
		LanguageTypeBWT:    true,
		LanguageTypeBWU:    true,
		LanguageTypeBWW:    true,
		LanguageTypeBWX:    true,
		LanguageTypeBWY:    true,
		LanguageTypeBWZ:    true,
		LanguageTypeBXA:    true,
		LanguageTypeBXB:    true,
		LanguageTypeBXC:    true,
		LanguageTypeBXD:    true,
		LanguageTypeBXE:    true,
		LanguageTypeBXF:    true,
		LanguageTypeBXG:    true,
		LanguageTypeBXH:    true,
		LanguageTypeBXI:    true,
		LanguageTypeBXJ:    true,
		LanguageTypeBXK:    true,
		LanguageTypeBXL:    true,
		LanguageTypeBXM:    true,
		LanguageTypeBXN:    true,
		LanguageTypeBXO:    true,
		LanguageTypeBXP:    true,
		LanguageTypeBXQ:    true,
		LanguageTypeBXR:    true,
		LanguageTypeBXS:    true,
		LanguageTypeBXU:    true,
		LanguageTypeBXV:    true,
		LanguageTypeBXW:    true,
		LanguageTypeBXX:    true,
		LanguageTypeBXZ:    true,
		LanguageTypeBYA:    true,
		LanguageTypeBYB:    true,
		LanguageTypeBYC:    true,
		LanguageTypeBYD:    true,
		LanguageTypeBYE:    true,
		LanguageTypeBYF:    true,
		LanguageTypeBYG:    true,
		LanguageTypeBYH:    true,
		LanguageTypeBYI:    true,
		LanguageTypeBYJ:    true,
		LanguageTypeBYK:    true,
		LanguageTypeBYL:    true,
		LanguageTypeBYM:    true,
		LanguageTypeBYN:    true,
		LanguageTypeBYO:    true,
		LanguageTypeBYP:    true,
		LanguageTypeBYQ:    true,
		LanguageTypeBYR:    true,
		LanguageTypeBYS:    true,
		LanguageTypeBYT:    true,
		LanguageTypeBYV:    true,
		LanguageTypeBYW:    true,
		LanguageTypeBYX:    true,
		LanguageTypeBYY:    true,
		LanguageTypeBYZ:    true,
		LanguageTypeBZA:    true,
		LanguageTypeBZB:    true,
		LanguageTypeBZC:    true,
		LanguageTypeBZD:    true,
		LanguageTypeBZE:    true,
		LanguageTypeBZF:    true,
		LanguageTypeBZG:    true,
		LanguageTypeBZH:    true,
		LanguageTypeBZI:    true,
		LanguageTypeBZJ:    true,
		LanguageTypeBZK:    true,
		LanguageTypeBZL:    true,
		LanguageTypeBZM:    true,
		LanguageTypeBZN:    true,
		LanguageTypeBZO:    true,
		LanguageTypeBZP:    true,
		LanguageTypeBZQ:    true,
		LanguageTypeBZR:    true,
		LanguageTypeBZS:    true,
		LanguageTypeBZT:    true,
		LanguageTypeBZU:    true,
		LanguageTypeBZV:    true,
		LanguageTypeBZW:    true,
		LanguageTypeBZX:    true,
		LanguageTypeBZY:    true,
		LanguageTypeBZZ:    true,
		LanguageTypeCAA:    true,
		LanguageTypeCAB:    true,
		LanguageTypeCAC:    true,
		LanguageTypeCAD:    true,
		LanguageTypeCAE:    true,
		LanguageTypeCAF:    true,
		LanguageTypeCAG:    true,
		LanguageTypeCAH:    true,
		LanguageTypeCAI:    true,
		LanguageTypeCAJ:    true,
		LanguageTypeCAK:    true,
		LanguageTypeCAL:    true,
		LanguageTypeCAM:    true,
		LanguageTypeCAN:    true,
		LanguageTypeCAO:    true,
		LanguageTypeCAP:    true,
		LanguageTypeCAQ:    true,
		LanguageTypeCAR:    true,
		LanguageTypeCAS:    true,
		LanguageTypeCAU:    true,
		LanguageTypeCAV:    true,
		LanguageTypeCAW:    true,
		LanguageTypeCAX:    true,
		LanguageTypeCAY:    true,
		LanguageTypeCAZ:    true,
		LanguageTypeCBA:    true,
		LanguageTypeCBB:    true,
		LanguageTypeCBC:    true,
		LanguageTypeCBD:    true,
		LanguageTypeCBE:    true,
		LanguageTypeCBG:    true,
		LanguageTypeCBH:    true,
		LanguageTypeCBI:    true,
		LanguageTypeCBJ:    true,
		LanguageTypeCBK:    true,
		LanguageTypeCBL:    true,
		LanguageTypeCBN:    true,
		LanguageTypeCBO:    true,
		LanguageTypeCBR:    true,
		LanguageTypeCBS:    true,
		LanguageTypeCBT:    true,
		LanguageTypeCBU:    true,
		LanguageTypeCBV:    true,
		LanguageTypeCBW:    true,
		LanguageTypeCBY:    true,
		LanguageTypeCCA:    true,
		LanguageTypeCCC:    true,
		LanguageTypeCCD:    true,
		LanguageTypeCCE:    true,
		LanguageTypeCCG:    true,
		LanguageTypeCCH:    true,
		LanguageTypeCCJ:    true,
		LanguageTypeCCL:    true,
		LanguageTypeCCM:    true,
		LanguageTypeCCN:    true,
		LanguageTypeCCO:    true,
		LanguageTypeCCP:    true,
		LanguageTypeCCQ:    true,
		LanguageTypeCCR:    true,
		LanguageTypeCCS:    true,
		LanguageTypeCDA:    true,
		LanguageTypeCDC:    true,
		LanguageTypeCDD:    true,
		LanguageTypeCDE:    true,
		LanguageTypeCDF:    true,
		LanguageTypeCDG:    true,
		LanguageTypeCDH:    true,
		LanguageTypeCDI:    true,
		LanguageTypeCDJ:    true,
		LanguageTypeCDM:    true,
		LanguageTypeCDN:    true,
		LanguageTypeCDO:    true,
		LanguageTypeCDR:    true,
		LanguageTypeCDS:    true,
		LanguageTypeCDY:    true,
		LanguageTypeCDZ:    true,
		LanguageTypeCEA:    true,
		LanguageTypeCEB:    true,
		LanguageTypeCEG:    true,
		LanguageTypeCEK:    true,
		LanguageTypeCEL:    true,
		LanguageTypeCEN:    true,
		LanguageTypeCET:    true,
		LanguageTypeCFA:    true,
		LanguageTypeCFD:    true,
		LanguageTypeCFG:    true,
		LanguageTypeCFM:    true,
		LanguageTypeCGA:    true,
		LanguageTypeCGC:    true,
		LanguageTypeCGG:    true,
		LanguageTypeCGK:    true,
		LanguageTypeCHB:    true,
		LanguageTypeCHC:    true,
		LanguageTypeCHD:    true,
		LanguageTypeCHF:    true,
		LanguageTypeCHG:    true,
		LanguageTypeCHH:    true,
		LanguageTypeCHJ:    true,
		LanguageTypeCHK:    true,
		LanguageTypeCHL:    true,
		LanguageTypeCHM:    true,
		LanguageTypeCHN:    true,
		LanguageTypeCHO:    true,
		LanguageTypeCHP:    true,
		LanguageTypeCHQ:    true,
		LanguageTypeCHR:    true,
		LanguageTypeCHT:    true,
		LanguageTypeCHW:    true,
		LanguageTypeCHX:    true,
		LanguageTypeCHY:    true,
		LanguageTypeCHZ:    true,
		LanguageTypeCIA:    true,
		LanguageTypeCIB:    true,
		LanguageTypeCIC:    true,
		LanguageTypeCID:    true,
		LanguageTypeCIE:    true,
		LanguageTypeCIH:    true,
		LanguageTypeCIK:    true,
		LanguageTypeCIM:    true,
		LanguageTypeCIN:    true,
		LanguageTypeCIP:    true,
		LanguageTypeCIR:    true,
		LanguageTypeCIW:    true,
		LanguageTypeCIY:    true,
		LanguageTypeCJA:    true,
		LanguageTypeCJE:    true,
		LanguageTypeCJH:    true,
		LanguageTypeCJI:    true,
		LanguageTypeCJK:    true,
		LanguageTypeCJM:    true,
		LanguageTypeCJN:    true,
		LanguageTypeCJO:    true,
		LanguageTypeCJP:    true,
		LanguageTypeCJR:    true,
		LanguageTypeCJS:    true,
		LanguageTypeCJV:    true,
		LanguageTypeCJY:    true,
		LanguageTypeCKA:    true,
		LanguageTypeCKB:    true,
		LanguageTypeCKH:    true,
		LanguageTypeCKL:    true,
		LanguageTypeCKN:    true,
		LanguageTypeCKO:    true,
		LanguageTypeCKQ:    true,
		LanguageTypeCKR:    true,
		LanguageTypeCKS:    true,
		LanguageTypeCKT:    true,
		LanguageTypeCKU:    true,
		LanguageTypeCKV:    true,
		LanguageTypeCKX:    true,
		LanguageTypeCKY:    true,
		LanguageTypeCKZ:    true,
		LanguageTypeCLA:    true,
		LanguageTypeCLC:    true,
		LanguageTypeCLD:    true,
		LanguageTypeCLE:    true,
		LanguageTypeCLH:    true,
		LanguageTypeCLI:    true,
		LanguageTypeCLJ:    true,
		LanguageTypeCLK:    true,
		LanguageTypeCLL:    true,
		LanguageTypeCLM:    true,
		LanguageTypeCLO:    true,
		LanguageTypeCLT:    true,
		LanguageTypeCLU:    true,
		LanguageTypeCLW:    true,
		LanguageTypeCLY:    true,
		LanguageTypeCMA:    true,
		LanguageTypeCMC:    true,
		LanguageTypeCME:    true,
		LanguageTypeCMG:    true,
		LanguageTypeCMI:    true,
		LanguageTypeCMK:    true,
		LanguageTypeCML:    true,
		LanguageTypeCMM:    true,
		LanguageTypeCMN:    true,
		LanguageTypeCMO:    true,
		LanguageTypeCMR:    true,
		LanguageTypeCMS:    true,
		LanguageTypeCMT:    true,
		LanguageTypeCNA:    true,
		LanguageTypeCNB:    true,
		LanguageTypeCNC:    true,
		LanguageTypeCNG:    true,
		LanguageTypeCNH:    true,
		LanguageTypeCNI:    true,
		LanguageTypeCNK:    true,
		LanguageTypeCNL:    true,
		LanguageTypeCNO:    true,
		LanguageTypeCNS:    true,
		LanguageTypeCNT:    true,
		LanguageTypeCNU:    true,
		LanguageTypeCNW:    true,
		LanguageTypeCNX:    true,
		LanguageTypeCOA:    true,
		LanguageTypeCOB:    true,
		LanguageTypeCOC:    true,
		LanguageTypeCOD:    true,
		LanguageTypeCOE:    true,
		LanguageTypeCOF:    true,
		LanguageTypeCOG:    true,
		LanguageTypeCOH:    true,
		LanguageTypeCOJ:    true,
		LanguageTypeCOK:    true,
		LanguageTypeCOL:    true,
		LanguageTypeCOM:    true,
		LanguageTypeCON:    true,
		LanguageTypeCOO:    true,
		LanguageTypeCOP:    true,
		LanguageTypeCOQ:    true,
		LanguageTypeCOT:    true,
		LanguageTypeCOU:    true,
		LanguageTypeCOV:    true,
		LanguageTypeCOW:    true,
		LanguageTypeCOX:    true,
		LanguageTypeCOY:    true,
		LanguageTypeCOZ:    true,
		LanguageTypeCPA:    true,
		LanguageTypeCPB:    true,
		LanguageTypeCPC:    true,
		LanguageTypeCPE:    true,
		LanguageTypeCPF:    true,
		LanguageTypeCPG:    true,
		LanguageTypeCPI:    true,
		LanguageTypeCPN:    true,
		LanguageTypeCPO:    true,
		LanguageTypeCPP:    true,
		LanguageTypeCPS:    true,
		LanguageTypeCPU:    true,
		LanguageTypeCPX:    true,
		LanguageTypeCPY:    true,
		LanguageTypeCQD:    true,
		LanguageTypeCQU:    true,
		LanguageTypeCRA:    true,
		LanguageTypeCRB:    true,
		LanguageTypeCRC:    true,
		LanguageTypeCRD:    true,
		LanguageTypeCRF:    true,
		LanguageTypeCRG:    true,
		LanguageTypeCRH:    true,
		LanguageTypeCRI:    true,
		LanguageTypeCRJ:    true,
		LanguageTypeCRK:    true,
		LanguageTypeCRL:    true,
		LanguageTypeCRM:    true,
		LanguageTypeCRN:    true,
		LanguageTypeCRO:    true,
		LanguageTypeCRP:    true,
		LanguageTypeCRQ:    true,
		LanguageTypeCRR:    true,
		LanguageTypeCRS:    true,
		LanguageTypeCRT:    true,
		LanguageTypeCRV:    true,
		LanguageTypeCRW:    true,
		LanguageTypeCRX:    true,
		LanguageTypeCRY:    true,
		LanguageTypeCRZ:    true,
		LanguageTypeCSA:    true,
		LanguageTypeCSB:    true,
		LanguageTypeCSC:    true,
		LanguageTypeCSD:    true,
		LanguageTypeCSE:    true,
		LanguageTypeCSF:    true,
		LanguageTypeCSG:    true,
		LanguageTypeCSH:    true,
		LanguageTypeCSI:    true,
		LanguageTypeCSJ:    true,
		LanguageTypeCSK:    true,
		LanguageTypeCSL:    true,
		LanguageTypeCSM:    true,
		LanguageTypeCSN:    true,
		LanguageTypeCSO:    true,
		LanguageTypeCSQ:    true,
		LanguageTypeCSR:    true,
		LanguageTypeCSS:    true,
		LanguageTypeCST:    true,
		LanguageTypeCSU:    true,
		LanguageTypeCSV:    true,
		LanguageTypeCSW:    true,
		LanguageTypeCSY:    true,
		LanguageTypeCSZ:    true,
		LanguageTypeCTA:    true,
		LanguageTypeCTC:    true,
		LanguageTypeCTD:    true,
		LanguageTypeCTE:    true,
		LanguageTypeCTG:    true,
		LanguageTypeCTH:    true,
		LanguageTypeCTL:    true,
		LanguageTypeCTM:    true,
		LanguageTypeCTN:    true,
		LanguageTypeCTO:    true,
		LanguageTypeCTP:    true,
		LanguageTypeCTS:    true,
		LanguageTypeCTT:    true,
		LanguageTypeCTU:    true,
		LanguageTypeCTZ:    true,
		LanguageTypeCUA:    true,
		LanguageTypeCUB:    true,
		LanguageTypeCUC:    true,
		LanguageTypeCUG:    true,
		LanguageTypeCUH:    true,
		LanguageTypeCUI:    true,
		LanguageTypeCUJ:    true,
		LanguageTypeCUK:    true,
		LanguageTypeCUL:    true,
		LanguageTypeCUM:    true,
		LanguageTypeCUO:    true,
		LanguageTypeCUP:    true,
		LanguageTypeCUQ:    true,
		LanguageTypeCUR:    true,
		LanguageTypeCUS:    true,
		LanguageTypeCUT:    true,
		LanguageTypeCUU:    true,
		LanguageTypeCUV:    true,
		LanguageTypeCUW:    true,
		LanguageTypeCUX:    true,
		LanguageTypeCVG:    true,
		LanguageTypeCVN:    true,
		LanguageTypeCWA:    true,
		LanguageTypeCWB:    true,
		LanguageTypeCWD:    true,
		LanguageTypeCWE:    true,
		LanguageTypeCWG:    true,
		LanguageTypeCWT:    true,
		LanguageTypeCYA:    true,
		LanguageTypeCYB:    true,
		LanguageTypeCYO:    true,
		LanguageTypeCZH:    true,
		LanguageTypeCZK:    true,
		LanguageTypeCZN:    true,
		LanguageTypeCZO:    true,
		LanguageTypeCZT:    true,
		LanguageTypeDAA:    true,
		LanguageTypeDAC:    true,
		LanguageTypeDAD:    true,
		LanguageTypeDAE:    true,
		LanguageTypeDAF:    true,
		LanguageTypeDAG:    true,
		LanguageTypeDAH:    true,
		LanguageTypeDAI:    true,
		LanguageTypeDAJ:    true,
		LanguageTypeDAK:    true,
		LanguageTypeDAL:    true,
		LanguageTypeDAM:    true,
		LanguageTypeDAO:    true,
		LanguageTypeDAP:    true,
		LanguageTypeDAQ:    true,
		LanguageTypeDAR:    true,
		LanguageTypeDAS:    true,
		LanguageTypeDAU:    true,
		LanguageTypeDAV:    true,
		LanguageTypeDAW:    true,
		LanguageTypeDAX:    true,
		LanguageTypeDAY:    true,
		LanguageTypeDAZ:    true,
		LanguageTypeDBA:    true,
		LanguageTypeDBB:    true,
		LanguageTypeDBD:    true,
		LanguageTypeDBE:    true,
		LanguageTypeDBF:    true,
		LanguageTypeDBG:    true,
		LanguageTypeDBI:    true,
		LanguageTypeDBJ:    true,
		LanguageTypeDBL:    true,
		LanguageTypeDBM:    true,
		LanguageTypeDBN:    true,
		LanguageTypeDBO:    true,
		LanguageTypeDBP:    true,
		LanguageTypeDBQ:    true,
		LanguageTypeDBR:    true,
		LanguageTypeDBT:    true,
		LanguageTypeDBU:    true,
		LanguageTypeDBV:    true,
		LanguageTypeDBW:    true,
		LanguageTypeDBY:    true,
		LanguageTypeDCC:    true,
		LanguageTypeDCR:    true,
		LanguageTypeDDA:    true,
		LanguageTypeDDD:    true,
		LanguageTypeDDE:    true,
		LanguageTypeDDG:    true,
		LanguageTypeDDI:    true,
		LanguageTypeDDJ:    true,
		LanguageTypeDDN:    true,
		LanguageTypeDDO:    true,
		LanguageTypeDDR:    true,
		LanguageTypeDDS:    true,
		LanguageTypeDDW:    true,
		LanguageTypeDEC:    true,
		LanguageTypeDED:    true,
		LanguageTypeDEE:    true,
		LanguageTypeDEF:    true,
		LanguageTypeDEG:    true,
		LanguageTypeDEH:    true,
		LanguageTypeDEI:    true,
		LanguageTypeDEK:    true,
		LanguageTypeDEL:    true,
		LanguageTypeDEM:    true,
		LanguageTypeDEN:    true,
		LanguageTypeDEP:    true,
		LanguageTypeDEQ:    true,
		LanguageTypeDER:    true,
		LanguageTypeDES:    true,
		LanguageTypeDEV:    true,
		LanguageTypeDEZ:    true,
		LanguageTypeDGA:    true,
		LanguageTypeDGB:    true,
		LanguageTypeDGC:    true,
		LanguageTypeDGD:    true,
		LanguageTypeDGE:    true,
		LanguageTypeDGG:    true,
		LanguageTypeDGH:    true,
		LanguageTypeDGI:    true,
		LanguageTypeDGK:    true,
		LanguageTypeDGL:    true,
		LanguageTypeDGN:    true,
		LanguageTypeDGO:    true,
		LanguageTypeDGR:    true,
		LanguageTypeDGS:    true,
		LanguageTypeDGT:    true,
		LanguageTypeDGU:    true,
		LanguageTypeDGW:    true,
		LanguageTypeDGX:    true,
		LanguageTypeDGZ:    true,
		LanguageTypeDHA:    true,
		LanguageTypeDHD:    true,
		LanguageTypeDHG:    true,
		LanguageTypeDHI:    true,
		LanguageTypeDHL:    true,
		LanguageTypeDHM:    true,
		LanguageTypeDHN:    true,
		LanguageTypeDHO:    true,
		LanguageTypeDHR:    true,
		LanguageTypeDHS:    true,
		LanguageTypeDHU:    true,
		LanguageTypeDHV:    true,
		LanguageTypeDHW:    true,
		LanguageTypeDHX:    true,
		LanguageTypeDIA:    true,
		LanguageTypeDIB:    true,
		LanguageTypeDIC:    true,
		LanguageTypeDID:    true,
		LanguageTypeDIF:    true,
		LanguageTypeDIG:    true,
		LanguageTypeDIH:    true,
		LanguageTypeDII:    true,
		LanguageTypeDIJ:    true,
		LanguageTypeDIK:    true,
		LanguageTypeDIL:    true,
		LanguageTypeDIM:    true,
		LanguageTypeDIN:    true,
		LanguageTypeDIO:    true,
		LanguageTypeDIP:    true,
		LanguageTypeDIQ:    true,
		LanguageTypeDIR:    true,
		LanguageTypeDIS:    true,
		LanguageTypeDIT:    true,
		LanguageTypeDIU:    true,
		LanguageTypeDIW:    true,
		LanguageTypeDIX:    true,
		LanguageTypeDIY:    true,
		LanguageTypeDIZ:    true,
		LanguageTypeDJA:    true,
		LanguageTypeDJB:    true,
		LanguageTypeDJC:    true,
		LanguageTypeDJD:    true,
		LanguageTypeDJE:    true,
		LanguageTypeDJF:    true,
		LanguageTypeDJI:    true,
		LanguageTypeDJJ:    true,
		LanguageTypeDJK:    true,
		LanguageTypeDJL:    true,
		LanguageTypeDJM:    true,
		LanguageTypeDJN:    true,
		LanguageTypeDJO:    true,
		LanguageTypeDJR:    true,
		LanguageTypeDJU:    true,
		LanguageTypeDJW:    true,
		LanguageTypeDKA:    true,
		LanguageTypeDKK:    true,
		LanguageTypeDKL:    true,
		LanguageTypeDKR:    true,
		LanguageTypeDKS:    true,
		LanguageTypeDKX:    true,
		LanguageTypeDLG:    true,
		LanguageTypeDLK:    true,
		LanguageTypeDLM:    true,
		LanguageTypeDLN:    true,
		LanguageTypeDMA:    true,
		LanguageTypeDMB:    true,
		LanguageTypeDMC:    true,
		LanguageTypeDMD:    true,
		LanguageTypeDME:    true,
		LanguageTypeDMG:    true,
		LanguageTypeDMK:    true,
		LanguageTypeDML:    true,
		LanguageTypeDMM:    true,
		LanguageTypeDMN:    true,
		LanguageTypeDMO:    true,
		LanguageTypeDMR:    true,
		LanguageTypeDMS:    true,
		LanguageTypeDMU:    true,
		LanguageTypeDMV:    true,
		LanguageTypeDMW:    true,
		LanguageTypeDMX:    true,
		LanguageTypeDMY:    true,
		LanguageTypeDNA:    true,
		LanguageTypeDND:    true,
		LanguageTypeDNE:    true,
		LanguageTypeDNG:    true,
		LanguageTypeDNI:    true,
		LanguageTypeDNJ:    true,
		LanguageTypeDNK:    true,
		LanguageTypeDNN:    true,
		LanguageTypeDNR:    true,
		LanguageTypeDNT:    true,
		LanguageTypeDNU:    true,
		LanguageTypeDNV:    true,
		LanguageTypeDNW:    true,
		LanguageTypeDNY:    true,
		LanguageTypeDOA:    true,
		LanguageTypeDOB:    true,
		LanguageTypeDOC:    true,
		LanguageTypeDOE:    true,
		LanguageTypeDOF:    true,
		LanguageTypeDOH:    true,
		LanguageTypeDOI:    true,
		LanguageTypeDOK:    true,
		LanguageTypeDOL:    true,
		LanguageTypeDON:    true,
		LanguageTypeDOO:    true,
		LanguageTypeDOP:    true,
		LanguageTypeDOQ:    true,
		LanguageTypeDOR:    true,
		LanguageTypeDOS:    true,
		LanguageTypeDOT:    true,
		LanguageTypeDOV:    true,
		LanguageTypeDOW:    true,
		LanguageTypeDOX:    true,
		LanguageTypeDOY:    true,
		LanguageTypeDOZ:    true,
		LanguageTypeDPP:    true,
		LanguageTypeDRA:    true,
		LanguageTypeDRB:    true,
		LanguageTypeDRC:    true,
		LanguageTypeDRD:    true,
		LanguageTypeDRE:    true,
		LanguageTypeDRG:    true,
		LanguageTypeDRH:    true,
		LanguageTypeDRI:    true,
		LanguageTypeDRL:    true,
		LanguageTypeDRN:    true,
		LanguageTypeDRO:    true,
		LanguageTypeDRQ:    true,
		LanguageTypeDRR:    true,
		LanguageTypeDRS:    true,
		LanguageTypeDRT:    true,
		LanguageTypeDRU:    true,
		LanguageTypeDRW:    true,
		LanguageTypeDRY:    true,
		LanguageTypeDSB:    true,
		LanguageTypeDSE:    true,
		LanguageTypeDSH:    true,
		LanguageTypeDSI:    true,
		LanguageTypeDSL:    true,
		LanguageTypeDSN:    true,
		LanguageTypeDSO:    true,
		LanguageTypeDSQ:    true,
		LanguageTypeDTA:    true,
		LanguageTypeDTB:    true,
		LanguageTypeDTD:    true,
		LanguageTypeDTH:    true,
		LanguageTypeDTI:    true,
		LanguageTypeDTK:    true,
		LanguageTypeDTM:    true,
		LanguageTypeDTO:    true,
		LanguageTypeDTP:    true,
		LanguageTypeDTR:    true,
		LanguageTypeDTS:    true,
		LanguageTypeDTT:    true,
		LanguageTypeDTU:    true,
		LanguageTypeDTY:    true,
		LanguageTypeDUA:    true,
		LanguageTypeDUB:    true,
		LanguageTypeDUC:    true,
		LanguageTypeDUD:    true,
		LanguageTypeDUE:    true,
		LanguageTypeDUF:    true,
		LanguageTypeDUG:    true,
		LanguageTypeDUH:    true,
		LanguageTypeDUI:    true,
		LanguageTypeDUJ:    true,
		LanguageTypeDUK:    true,
		LanguageTypeDUL:    true,
		LanguageTypeDUM:    true,
		LanguageTypeDUN:    true,
		LanguageTypeDUO:    true,
		LanguageTypeDUP:    true,
		LanguageTypeDUQ:    true,
		LanguageTypeDUR:    true,
		LanguageTypeDUS:    true,
		LanguageTypeDUU:    true,
		LanguageTypeDUV:    true,
		LanguageTypeDUW:    true,
		LanguageTypeDUX:    true,
		LanguageTypeDUY:    true,
		LanguageTypeDUZ:    true,
		LanguageTypeDVA:    true,
		LanguageTypeDWA:    true,
		LanguageTypeDWL:    true,
		LanguageTypeDWR:    true,
		LanguageTypeDWS:    true,
		LanguageTypeDWW:    true,
		LanguageTypeDYA:    true,
		LanguageTypeDYB:    true,
		LanguageTypeDYD:    true,
		LanguageTypeDYG:    true,
		LanguageTypeDYI:    true,
		LanguageTypeDYM:    true,
		LanguageTypeDYN:    true,
		LanguageTypeDYO:    true,
		LanguageTypeDYU:    true,
		LanguageTypeDYY:    true,
		LanguageTypeDZA:    true,
		LanguageTypeDZD:    true,
		LanguageTypeDZE:    true,
		LanguageTypeDZG:    true,
		LanguageTypeDZL:    true,
		LanguageTypeDZN:    true,
		LanguageTypeEAA:    true,
		LanguageTypeEBG:    true,
		LanguageTypeEBK:    true,
		LanguageTypeEBO:    true,
		LanguageTypeEBR:    true,
		LanguageTypeEBU:    true,
		LanguageTypeECR:    true,
		LanguageTypeECS:    true,
		LanguageTypeECY:    true,
		LanguageTypeEEE:    true,
		LanguageTypeEFA:    true,
		LanguageTypeEFE:    true,
		LanguageTypeEFI:    true,
		LanguageTypeEGA:    true,
		LanguageTypeEGL:    true,
		LanguageTypeEGO:    true,
		LanguageTypeEGX:    true,
		LanguageTypeEGY:    true,
		LanguageTypeEHU:    true,
		LanguageTypeEIP:    true,
		LanguageTypeEIT:    true,
		LanguageTypeEIV:    true,
		LanguageTypeEJA:    true,
		LanguageTypeEKA:    true,
		LanguageTypeEKC:    true,
		LanguageTypeEKE:    true,
		LanguageTypeEKG:    true,
		LanguageTypeEKI:    true,
		LanguageTypeEKK:    true,
		LanguageTypeEKL:    true,
		LanguageTypeEKM:    true,
		LanguageTypeEKO:    true,
		LanguageTypeEKP:    true,
		LanguageTypeEKR:    true,
		LanguageTypeEKY:    true,
		LanguageTypeELE:    true,
		LanguageTypeELH:    true,
		LanguageTypeELI:    true,
		LanguageTypeELK:    true,
		LanguageTypeELM:    true,
		LanguageTypeELO:    true,
		LanguageTypeELP:    true,
		LanguageTypeELU:    true,
		LanguageTypeELX:    true,
		LanguageTypeEMA:    true,
		LanguageTypeEMB:    true,
		LanguageTypeEME:    true,
		LanguageTypeEMG:    true,
		LanguageTypeEMI:    true,
		LanguageTypeEMK:    true,
		LanguageTypeEMM:    true,
		LanguageTypeEMN:    true,
		LanguageTypeEMO:    true,
		LanguageTypeEMP:    true,
		LanguageTypeEMS:    true,
		LanguageTypeEMU:    true,
		LanguageTypeEMW:    true,
		LanguageTypeEMX:    true,
		LanguageTypeEMY:    true,
		LanguageTypeENA:    true,
		LanguageTypeENB:    true,
		LanguageTypeENC:    true,
		LanguageTypeEND:    true,
		LanguageTypeENF:    true,
		LanguageTypeENH:    true,
		LanguageTypeENM:    true,
		LanguageTypeENN:    true,
		LanguageTypeENO:    true,
		LanguageTypeENQ:    true,
		LanguageTypeENR:    true,
		LanguageTypeENU:    true,
		LanguageTypeENV:    true,
		LanguageTypeENW:    true,
		LanguageTypeEOT:    true,
		LanguageTypeEPI:    true,
		LanguageTypeERA:    true,
		LanguageTypeERG:    true,
		LanguageTypeERH:    true,
		LanguageTypeERI:    true,
		LanguageTypeERK:    true,
		LanguageTypeERO:    true,
		LanguageTypeERR:    true,
		LanguageTypeERS:    true,
		LanguageTypeERT:    true,
		LanguageTypeERW:    true,
		LanguageTypeESE:    true,
		LanguageTypeESH:    true,
		LanguageTypeESI:    true,
		LanguageTypeESK:    true,
		LanguageTypeESL:    true,
		LanguageTypeESM:    true,
		LanguageTypeESN:    true,
		LanguageTypeESO:    true,
		LanguageTypeESQ:    true,
		LanguageTypeESS:    true,
		LanguageTypeESU:    true,
		LanguageTypeESX:    true,
		LanguageTypeETB:    true,
		LanguageTypeETC:    true,
		LanguageTypeETH:    true,
		LanguageTypeETN:    true,
		LanguageTypeETO:    true,
		LanguageTypeETR:    true,
		LanguageTypeETS:    true,
		LanguageTypeETT:    true,
		LanguageTypeETU:    true,
		LanguageTypeETX:    true,
		LanguageTypeETZ:    true,
		LanguageTypeEUQ:    true,
		LanguageTypeEVE:    true,
		LanguageTypeEVH:    true,
		LanguageTypeEVN:    true,
		LanguageTypeEWO:    true,
		LanguageTypeEXT:    true,
		LanguageTypeEYA:    true,
		LanguageTypeEYO:    true,
		LanguageTypeEZA:    true,
		LanguageTypeEZE:    true,
		LanguageTypeFAA:    true,
		LanguageTypeFAB:    true,
		LanguageTypeFAD:    true,
		LanguageTypeFAF:    true,
		LanguageTypeFAG:    true,
		LanguageTypeFAH:    true,
		LanguageTypeFAI:    true,
		LanguageTypeFAJ:    true,
		LanguageTypeFAK:    true,
		LanguageTypeFAL:    true,
		LanguageTypeFAM:    true,
		LanguageTypeFAN:    true,
		LanguageTypeFAP:    true,
		LanguageTypeFAR:    true,
		LanguageTypeFAT:    true,
		LanguageTypeFAU:    true,
		LanguageTypeFAX:    true,
		LanguageTypeFAY:    true,
		LanguageTypeFAZ:    true,
		LanguageTypeFBL:    true,
		LanguageTypeFCS:    true,
		LanguageTypeFER:    true,
		LanguageTypeFFI:    true,
		LanguageTypeFFM:    true,
		LanguageTypeFGR:    true,
		LanguageTypeFIA:    true,
		LanguageTypeFIE:    true,
		LanguageTypeFIL:    true,
		LanguageTypeFIP:    true,
		LanguageTypeFIR:    true,
		LanguageTypeFIT:    true,
		LanguageTypeFIU:    true,
		LanguageTypeFIW:    true,
		LanguageTypeFKK:    true,
		LanguageTypeFKV:    true,
		LanguageTypeFLA:    true,
		LanguageTypeFLH:    true,
		LanguageTypeFLI:    true,
		LanguageTypeFLL:    true,
		LanguageTypeFLN:    true,
		LanguageTypeFLR:    true,
		LanguageTypeFLY:    true,
		LanguageTypeFMP:    true,
		LanguageTypeFMU:    true,
		LanguageTypeFNG:    true,
		LanguageTypeFNI:    true,
		LanguageTypeFOD:    true,
		LanguageTypeFOI:    true,
		LanguageTypeFOM:    true,
		LanguageTypeFON:    true,
		LanguageTypeFOR:    true,
		LanguageTypeFOS:    true,
		LanguageTypeFOX:    true,
		LanguageTypeFPE:    true,
		LanguageTypeFQS:    true,
		LanguageTypeFRC:    true,
		LanguageTypeFRD:    true,
		LanguageTypeFRK:    true,
		LanguageTypeFRM:    true,
		LanguageTypeFRO:    true,
		LanguageTypeFRP:    true,
		LanguageTypeFRQ:    true,
		LanguageTypeFRR:    true,
		LanguageTypeFRS:    true,
		LanguageTypeFRT:    true,
		LanguageTypeFSE:    true,
		LanguageTypeFSL:    true,
		LanguageTypeFSS:    true,
		LanguageTypeFUB:    true,
		LanguageTypeFUC:    true,
		LanguageTypeFUD:    true,
		LanguageTypeFUE:    true,
		LanguageTypeFUF:    true,
		LanguageTypeFUH:    true,
		LanguageTypeFUI:    true,
		LanguageTypeFUJ:    true,
		LanguageTypeFUM:    true,
		LanguageTypeFUN:    true,
		LanguageTypeFUQ:    true,
		LanguageTypeFUR:    true,
		LanguageTypeFUT:    true,
		LanguageTypeFUU:    true,
		LanguageTypeFUV:    true,
		LanguageTypeFUY:    true,
		LanguageTypeFVR:    true,
		LanguageTypeFWA:    true,
		LanguageTypeFWE:    true,
		LanguageTypeGAA:    true,
		LanguageTypeGAB:    true,
		LanguageTypeGAC:    true,
		LanguageTypeGAD:    true,
		LanguageTypeGAE:    true,
		LanguageTypeGAF:    true,
		LanguageTypeGAG:    true,
		LanguageTypeGAH:    true,
		LanguageTypeGAI:    true,
		LanguageTypeGAJ:    true,
		LanguageTypeGAK:    true,
		LanguageTypeGAL:    true,
		LanguageTypeGAM:    true,
		LanguageTypeGAN:    true,
		LanguageTypeGAO:    true,
		LanguageTypeGAP:    true,
		LanguageTypeGAQ:    true,
		LanguageTypeGAR:    true,
		LanguageTypeGAS:    true,
		LanguageTypeGAT:    true,
		LanguageTypeGAU:    true,
		LanguageTypeGAV:    true,
		LanguageTypeGAW:    true,
		LanguageTypeGAX:    true,
		LanguageTypeGAY:    true,
		LanguageTypeGAZ:    true,
		LanguageTypeGBA:    true,
		LanguageTypeGBB:    true,
		LanguageTypeGBC:    true,
		LanguageTypeGBD:    true,
		LanguageTypeGBE:    true,
		LanguageTypeGBF:    true,
		LanguageTypeGBG:    true,
		LanguageTypeGBH:    true,
		LanguageTypeGBI:    true,
		LanguageTypeGBJ:    true,
		LanguageTypeGBK:    true,
		LanguageTypeGBL:    true,
		LanguageTypeGBM:    true,
		LanguageTypeGBN:    true,
		LanguageTypeGBO:    true,
		LanguageTypeGBP:    true,
		LanguageTypeGBQ:    true,
		LanguageTypeGBR:    true,
		LanguageTypeGBS:    true,
		LanguageTypeGBU:    true,
		LanguageTypeGBV:    true,
		LanguageTypeGBW:    true,
		LanguageTypeGBX:    true,
		LanguageTypeGBY:    true,
		LanguageTypeGBZ:    true,
		LanguageTypeGCC:    true,
		LanguageTypeGCD:    true,
		LanguageTypeGCE:    true,
		LanguageTypeGCF:    true,
		LanguageTypeGCL:    true,
		LanguageTypeGCN:    true,
		LanguageTypeGCR:    true,
		LanguageTypeGCT:    true,
		LanguageTypeGDA:    true,
		LanguageTypeGDB:    true,
		LanguageTypeGDC:    true,
		LanguageTypeGDD:    true,
		LanguageTypeGDE:    true,
		LanguageTypeGDF:    true,
		LanguageTypeGDG:    true,
		LanguageTypeGDH:    true,
		LanguageTypeGDI:    true,
		LanguageTypeGDJ:    true,
		LanguageTypeGDK:    true,
		LanguageTypeGDL:    true,
		LanguageTypeGDM:    true,
		LanguageTypeGDN:    true,
		LanguageTypeGDO:    true,
		LanguageTypeGDQ:    true,
		LanguageTypeGDR:    true,
		LanguageTypeGDS:    true,
		LanguageTypeGDT:    true,
		LanguageTypeGDU:    true,
		LanguageTypeGDX:    true,
		LanguageTypeGEA:    true,
		LanguageTypeGEB:    true,
		LanguageTypeGEC:    true,
		LanguageTypeGED:    true,
		LanguageTypeGEG:    true,
		LanguageTypeGEH:    true,
		LanguageTypeGEI:    true,
		LanguageTypeGEJ:    true,
		LanguageTypeGEK:    true,
		LanguageTypeGEL:    true,
		LanguageTypeGEM:    true,
		LanguageTypeGEQ:    true,
		LanguageTypeGES:    true,
		LanguageTypeGEW:    true,
		LanguageTypeGEX:    true,
		LanguageTypeGEY:    true,
		LanguageTypeGEZ:    true,
		LanguageTypeGFK:    true,
		LanguageTypeGFT:    true,
		LanguageTypeGFX:    true,
		LanguageTypeGGA:    true,
		LanguageTypeGGB:    true,
		LanguageTypeGGD:    true,
		LanguageTypeGGE:    true,
		LanguageTypeGGG:    true,
		LanguageTypeGGK:    true,
		LanguageTypeGGL:    true,
		LanguageTypeGGN:    true,
		LanguageTypeGGO:    true,
		LanguageTypeGGR:    true,
		LanguageTypeGGT:    true,
		LanguageTypeGGU:    true,
		LanguageTypeGGW:    true,
		LanguageTypeGHA:    true,
		LanguageTypeGHC:    true,
		LanguageTypeGHE:    true,
		LanguageTypeGHH:    true,
		LanguageTypeGHK:    true,
		LanguageTypeGHL:    true,
		LanguageTypeGHN:    true,
		LanguageTypeGHO:    true,
		LanguageTypeGHR:    true,
		LanguageTypeGHS:    true,
		LanguageTypeGHT:    true,
		LanguageTypeGIA:    true,
		LanguageTypeGIB:    true,
		LanguageTypeGIC:    true,
		LanguageTypeGID:    true,
		LanguageTypeGIG:    true,
		LanguageTypeGIH:    true,
		LanguageTypeGIL:    true,
		LanguageTypeGIM:    true,
		LanguageTypeGIN:    true,
		LanguageTypeGIO:    true,
		LanguageTypeGIP:    true,
		LanguageTypeGIQ:    true,
		LanguageTypeGIR:    true,
		LanguageTypeGIS:    true,
		LanguageTypeGIT:    true,
		LanguageTypeGIU:    true,
		LanguageTypeGIW:    true,
		LanguageTypeGIX:    true,
		LanguageTypeGIY:    true,
		LanguageTypeGIZ:    true,
		LanguageTypeGJI:    true,
		LanguageTypeGJK:    true,
		LanguageTypeGJM:    true,
		LanguageTypeGJN:    true,
		LanguageTypeGJU:    true,
		LanguageTypeGKA:    true,
		LanguageTypeGKE:    true,
		LanguageTypeGKN:    true,
		LanguageTypeGKO:    true,
		LanguageTypeGKP:    true,
		LanguageTypeGLC:    true,
		LanguageTypeGLD:    true,
		LanguageTypeGLH:    true,
		LanguageTypeGLI:    true,
		LanguageTypeGLJ:    true,
		LanguageTypeGLK:    true,
		LanguageTypeGLL:    true,
		LanguageTypeGLO:    true,
		LanguageTypeGLR:    true,
		LanguageTypeGLU:    true,
		LanguageTypeGLW:    true,
		LanguageTypeGLY:    true,
		LanguageTypeGMA:    true,
		LanguageTypeGMB:    true,
		LanguageTypeGMD:    true,
		LanguageTypeGME:    true,
		LanguageTypeGMH:    true,
		LanguageTypeGML:    true,
		LanguageTypeGMM:    true,
		LanguageTypeGMN:    true,
		LanguageTypeGMQ:    true,
		LanguageTypeGMU:    true,
		LanguageTypeGMV:    true,
		LanguageTypeGMW:    true,
		LanguageTypeGMX:    true,
		LanguageTypeGMY:    true,
		LanguageTypeGMZ:    true,
		LanguageTypeGNA:    true,
		LanguageTypeGNB:    true,
		LanguageTypeGNC:    true,
		LanguageTypeGND:    true,
		LanguageTypeGNE:    true,
		LanguageTypeGNG:    true,
		LanguageTypeGNH:    true,
		LanguageTypeGNI:    true,
		LanguageTypeGNK:    true,
		LanguageTypeGNL:    true,
		LanguageTypeGNM:    true,
		LanguageTypeGNN:    true,
		LanguageTypeGNO:    true,
		LanguageTypeGNQ:    true,
		LanguageTypeGNR:    true,
		LanguageTypeGNT:    true,
		LanguageTypeGNU:    true,
		LanguageTypeGNW:    true,
		LanguageTypeGNZ:    true,
		LanguageTypeGOA:    true,
		LanguageTypeGOB:    true,
		LanguageTypeGOC:    true,
		LanguageTypeGOD:    true,
		LanguageTypeGOE:    true,
		LanguageTypeGOF:    true,
		LanguageTypeGOG:    true,
		LanguageTypeGOH:    true,
		LanguageTypeGOI:    true,
		LanguageTypeGOJ:    true,
		LanguageTypeGOK:    true,
		LanguageTypeGOL:    true,
		LanguageTypeGOM:    true,
		LanguageTypeGON:    true,
		LanguageTypeGOO:    true,
		LanguageTypeGOP:    true,
		LanguageTypeGOQ:    true,
		LanguageTypeGOR:    true,
		LanguageTypeGOS:    true,
		LanguageTypeGOT:    true,
		LanguageTypeGOU:    true,
		LanguageTypeGOW:    true,
		LanguageTypeGOX:    true,
		LanguageTypeGOY:    true,
		LanguageTypeGOZ:    true,
		LanguageTypeGPA:    true,
		LanguageTypeGPE:    true,
		LanguageTypeGPN:    true,
		LanguageTypeGQA:    true,
		LanguageTypeGQI:    true,
		LanguageTypeGQN:    true,
		LanguageTypeGQR:    true,
		LanguageTypeGQU:    true,
		LanguageTypeGRA:    true,
		LanguageTypeGRB:    true,
		LanguageTypeGRC:    true,
		LanguageTypeGRD:    true,
		LanguageTypeGRG:    true,
		LanguageTypeGRH:    true,
		LanguageTypeGRI:    true,
		LanguageTypeGRJ:    true,
		LanguageTypeGRK:    true,
		LanguageTypeGRM:    true,
		LanguageTypeGRO:    true,
		LanguageTypeGRQ:    true,
		LanguageTypeGRR:    true,
		LanguageTypeGRS:    true,
		LanguageTypeGRT:    true,
		LanguageTypeGRU:    true,
		LanguageTypeGRV:    true,
		LanguageTypeGRW:    true,
		LanguageTypeGRX:    true,
		LanguageTypeGRY:    true,
		LanguageTypeGRZ:    true,
		LanguageTypeGSE:    true,
		LanguageTypeGSG:    true,
		LanguageTypeGSL:    true,
		LanguageTypeGSM:    true,
		LanguageTypeGSN:    true,
		LanguageTypeGSO:    true,
		LanguageTypeGSP:    true,
		LanguageTypeGSS:    true,
		LanguageTypeGSW:    true,
		LanguageTypeGTA:    true,
		LanguageTypeGTI:    true,
		LanguageTypeGTU:    true,
		LanguageTypeGUA:    true,
		LanguageTypeGUB:    true,
		LanguageTypeGUC:    true,
		LanguageTypeGUD:    true,
		LanguageTypeGUE:    true,
		LanguageTypeGUF:    true,
		LanguageTypeGUG:    true,
		LanguageTypeGUH:    true,
		LanguageTypeGUI:    true,
		LanguageTypeGUK:    true,
		LanguageTypeGUL:    true,
		LanguageTypeGUM:    true,
		LanguageTypeGUN:    true,
		LanguageTypeGUO:    true,
		LanguageTypeGUP:    true,
		LanguageTypeGUQ:    true,
		LanguageTypeGUR:    true,
		LanguageTypeGUS:    true,
		LanguageTypeGUT:    true,
		LanguageTypeGUU:    true,
		LanguageTypeGUV:    true,
		LanguageTypeGUW:    true,
		LanguageTypeGUX:    true,
		LanguageTypeGUZ:    true,
		LanguageTypeGVA:    true,
		LanguageTypeGVC:    true,
		LanguageTypeGVE:    true,
		LanguageTypeGVF:    true,
		LanguageTypeGVJ:    true,
		LanguageTypeGVL:    true,
		LanguageTypeGVM:    true,
		LanguageTypeGVN:    true,
		LanguageTypeGVO:    true,
		LanguageTypeGVP:    true,
		LanguageTypeGVR:    true,
		LanguageTypeGVS:    true,
		LanguageTypeGVY:    true,
		LanguageTypeGWA:    true,
		LanguageTypeGWB:    true,
		LanguageTypeGWC:    true,
		LanguageTypeGWD:    true,
		LanguageTypeGWE:    true,
		LanguageTypeGWF:    true,
		LanguageTypeGWG:    true,
		LanguageTypeGWI:    true,
		LanguageTypeGWJ:    true,
		LanguageTypeGWM:    true,
		LanguageTypeGWN:    true,
		LanguageTypeGWR:    true,
		LanguageTypeGWT:    true,
		LanguageTypeGWU:    true,
		LanguageTypeGWW:    true,
		LanguageTypeGWX:    true,
		LanguageTypeGXX:    true,
		LanguageTypeGYA:    true,
		LanguageTypeGYB:    true,
		LanguageTypeGYD:    true,
		LanguageTypeGYE:    true,
		LanguageTypeGYF:    true,
		LanguageTypeGYG:    true,
		LanguageTypeGYI:    true,
		LanguageTypeGYL:    true,
		LanguageTypeGYM:    true,
		LanguageTypeGYN:    true,
		LanguageTypeGYR:    true,
		LanguageTypeGYY:    true,
		LanguageTypeGZA:    true,
		LanguageTypeGZI:    true,
		LanguageTypeGZN:    true,
		LanguageTypeHAA:    true,
		LanguageTypeHAB:    true,
		LanguageTypeHAC:    true,
		LanguageTypeHAD:    true,
		LanguageTypeHAE:    true,
		LanguageTypeHAF:    true,
		LanguageTypeHAG:    true,
		LanguageTypeHAH:    true,
		LanguageTypeHAI:    true,
		LanguageTypeHAJ:    true,
		LanguageTypeHAK:    true,
		LanguageTypeHAL:    true,
		LanguageTypeHAM:    true,
		LanguageTypeHAN:    true,
		LanguageTypeHAO:    true,
		LanguageTypeHAP:    true,
		LanguageTypeHAQ:    true,
		LanguageTypeHAR:    true,
		LanguageTypeHAS:    true,
		LanguageTypeHAV:    true,
		LanguageTypeHAW:    true,
		LanguageTypeHAX:    true,
		LanguageTypeHAY:    true,
		LanguageTypeHAZ:    true,
		LanguageTypeHBA:    true,
		LanguageTypeHBB:    true,
		LanguageTypeHBN:    true,
		LanguageTypeHBO:    true,
		LanguageTypeHBU:    true,
		LanguageTypeHCA:    true,
		LanguageTypeHCH:    true,
		LanguageTypeHDN:    true,
		LanguageTypeHDS:    true,
		LanguageTypeHDY:    true,
		LanguageTypeHEA:    true,
		LanguageTypeHED:    true,
		LanguageTypeHEG:    true,
		LanguageTypeHEH:    true,
		LanguageTypeHEI:    true,
		LanguageTypeHEM:    true,
		LanguageTypeHGM:    true,
		LanguageTypeHGW:    true,
		LanguageTypeHHI:    true,
		LanguageTypeHHR:    true,
		LanguageTypeHHY:    true,
		LanguageTypeHIA:    true,
		LanguageTypeHIB:    true,
		LanguageTypeHID:    true,
		LanguageTypeHIF:    true,
		LanguageTypeHIG:    true,
		LanguageTypeHIH:    true,
		LanguageTypeHII:    true,
		LanguageTypeHIJ:    true,
		LanguageTypeHIK:    true,
		LanguageTypeHIL:    true,
		LanguageTypeHIM:    true,
		LanguageTypeHIO:    true,
		LanguageTypeHIR:    true,
		LanguageTypeHIT:    true,
		LanguageTypeHIW:    true,
		LanguageTypeHIX:    true,
		LanguageTypeHJI:    true,
		LanguageTypeHKA:    true,
		LanguageTypeHKE:    true,
		LanguageTypeHKK:    true,
		LanguageTypeHKS:    true,
		LanguageTypeHLA:    true,
		LanguageTypeHLB:    true,
		LanguageTypeHLD:    true,
		LanguageTypeHLE:    true,
		LanguageTypeHLT:    true,
		LanguageTypeHLU:    true,
		LanguageTypeHMA:    true,
		LanguageTypeHMB:    true,
		LanguageTypeHMC:    true,
		LanguageTypeHMD:    true,
		LanguageTypeHME:    true,
		LanguageTypeHMF:    true,
		LanguageTypeHMG:    true,
		LanguageTypeHMH:    true,
		LanguageTypeHMI:    true,
		LanguageTypeHMJ:    true,
		LanguageTypeHMK:    true,
		LanguageTypeHML:    true,
		LanguageTypeHMM:    true,
		LanguageTypeHMN:    true,
		LanguageTypeHMP:    true,
		LanguageTypeHMQ:    true,
		LanguageTypeHMR:    true,
		LanguageTypeHMS:    true,
		LanguageTypeHMT:    true,
		LanguageTypeHMU:    true,
		LanguageTypeHMV:    true,
		LanguageTypeHMW:    true,
		LanguageTypeHMX:    true,
		LanguageTypeHMY:    true,
		LanguageTypeHMZ:    true,
		LanguageTypeHNA:    true,
		LanguageTypeHND:    true,
		LanguageTypeHNE:    true,
		LanguageTypeHNH:    true,
		LanguageTypeHNI:    true,
		LanguageTypeHNJ:    true,
		LanguageTypeHNN:    true,
		LanguageTypeHNO:    true,
		LanguageTypeHNS:    true,
		LanguageTypeHNU:    true,
		LanguageTypeHOA:    true,
		LanguageTypeHOB:    true,
		LanguageTypeHOC:    true,
		LanguageTypeHOD:    true,
		LanguageTypeHOE:    true,
		LanguageTypeHOH:    true,
		LanguageTypeHOI:    true,
		LanguageTypeHOJ:    true,
		LanguageTypeHOK:    true,
		LanguageTypeHOL:    true,
		LanguageTypeHOM:    true,
		LanguageTypeHOO:    true,
		LanguageTypeHOP:    true,
		LanguageTypeHOR:    true,
		LanguageTypeHOS:    true,
		LanguageTypeHOT:    true,
		LanguageTypeHOV:    true,
		LanguageTypeHOW:    true,
		LanguageTypeHOY:    true,
		LanguageTypeHOZ:    true,
		LanguageTypeHPO:    true,
		LanguageTypeHPS:    true,
		LanguageTypeHRA:    true,
		LanguageTypeHRC:    true,
		LanguageTypeHRE:    true,
		LanguageTypeHRK:    true,
		LanguageTypeHRM:    true,
		LanguageTypeHRO:    true,
		LanguageTypeHRP:    true,
		LanguageTypeHRR:    true,
		LanguageTypeHRT:    true,
		LanguageTypeHRU:    true,
		LanguageTypeHRW:    true,
		LanguageTypeHRX:    true,
		LanguageTypeHRZ:    true,
		LanguageTypeHSB:    true,
		LanguageTypeHSH:    true,
		LanguageTypeHSL:    true,
		LanguageTypeHSN:    true,
		LanguageTypeHSS:    true,
		LanguageTypeHTI:    true,
		LanguageTypeHTO:    true,
		LanguageTypeHTS:    true,
		LanguageTypeHTU:    true,
		LanguageTypeHTX:    true,
		LanguageTypeHUB:    true,
		LanguageTypeHUC:    true,
		LanguageTypeHUD:    true,
		LanguageTypeHUE:    true,
		LanguageTypeHUF:    true,
		LanguageTypeHUG:    true,
		LanguageTypeHUH:    true,
		LanguageTypeHUI:    true,
		LanguageTypeHUJ:    true,
		LanguageTypeHUK:    true,
		LanguageTypeHUL:    true,
		LanguageTypeHUM:    true,
		LanguageTypeHUO:    true,
		LanguageTypeHUP:    true,
		LanguageTypeHUQ:    true,
		LanguageTypeHUR:    true,
		LanguageTypeHUS:    true,
		LanguageTypeHUT:    true,
		LanguageTypeHUU:    true,
		LanguageTypeHUV:    true,
		LanguageTypeHUW:    true,
		LanguageTypeHUX:    true,
		LanguageTypeHUY:    true,
		LanguageTypeHUZ:    true,
		LanguageTypeHVC:    true,
		LanguageTypeHVE:    true,
		LanguageTypeHVK:    true,
		LanguageTypeHVN:    true,
		LanguageTypeHVV:    true,
		LanguageTypeHWA:    true,
		LanguageTypeHWC:    true,
		LanguageTypeHWO:    true,
		LanguageTypeHYA:    true,
		LanguageTypeHYX:    true,
		LanguageTypeIAI:    true,
		LanguageTypeIAN:    true,
		LanguageTypeIAP:    true,
		LanguageTypeIAR:    true,
		LanguageTypeIBA:    true,
		LanguageTypeIBB:    true,
		LanguageTypeIBD:    true,
		LanguageTypeIBE:    true,
		LanguageTypeIBG:    true,
		LanguageTypeIBI:    true,
		LanguageTypeIBL:    true,
		LanguageTypeIBM:    true,
		LanguageTypeIBN:    true,
		LanguageTypeIBR:    true,
		LanguageTypeIBU:    true,
		LanguageTypeIBY:    true,
		LanguageTypeICA:    true,
		LanguageTypeICH:    true,
		LanguageTypeICL:    true,
		LanguageTypeICR:    true,
		LanguageTypeIDA:    true,
		LanguageTypeIDB:    true,
		LanguageTypeIDC:    true,
		LanguageTypeIDD:    true,
		LanguageTypeIDE:    true,
		LanguageTypeIDI:    true,
		LanguageTypeIDR:    true,
		LanguageTypeIDS:    true,
		LanguageTypeIDT:    true,
		LanguageTypeIDU:    true,
		LanguageTypeIFA:    true,
		LanguageTypeIFB:    true,
		LanguageTypeIFE:    true,
		LanguageTypeIFF:    true,
		LanguageTypeIFK:    true,
		LanguageTypeIFM:    true,
		LanguageTypeIFU:    true,
		LanguageTypeIFY:    true,
		LanguageTypeIGB:    true,
		LanguageTypeIGE:    true,
		LanguageTypeIGG:    true,
		LanguageTypeIGL:    true,
		LanguageTypeIGM:    true,
		LanguageTypeIGN:    true,
		LanguageTypeIGO:    true,
		LanguageTypeIGS:    true,
		LanguageTypeIGW:    true,
		LanguageTypeIHB:    true,
		LanguageTypeIHI:    true,
		LanguageTypeIHP:    true,
		LanguageTypeIHW:    true,
		LanguageTypeIIN:    true,
		LanguageTypeIIR:    true,
		LanguageTypeIJC:    true,
		LanguageTypeIJE:    true,
		LanguageTypeIJJ:    true,
		LanguageTypeIJN:    true,
		LanguageTypeIJO:    true,
		LanguageTypeIJS:    true,
		LanguageTypeIKE:    true,
		LanguageTypeIKI:    true,
		LanguageTypeIKK:    true,
		LanguageTypeIKL:    true,
		LanguageTypeIKO:    true,
		LanguageTypeIKP:    true,
		LanguageTypeIKR:    true,
		LanguageTypeIKT:    true,
		LanguageTypeIKV:    true,
		LanguageTypeIKW:    true,
		LanguageTypeIKX:    true,
		LanguageTypeIKZ:    true,
		LanguageTypeILA:    true,
		LanguageTypeILB:    true,
		LanguageTypeILG:    true,
		LanguageTypeILI:    true,
		LanguageTypeILK:    true,
		LanguageTypeILL:    true,
		LanguageTypeILO:    true,
		LanguageTypeILS:    true,
		LanguageTypeILU:    true,
		LanguageTypeILV:    true,
		LanguageTypeILW:    true,
		LanguageTypeIMA:    true,
		LanguageTypeIME:    true,
		LanguageTypeIMI:    true,
		LanguageTypeIML:    true,
		LanguageTypeIMN:    true,
		LanguageTypeIMO:    true,
		LanguageTypeIMR:    true,
		LanguageTypeIMS:    true,
		LanguageTypeIMY:    true,
		LanguageTypeINB:    true,
		LanguageTypeINC:    true,
		LanguageTypeINE:    true,
		LanguageTypeING:    true,
		LanguageTypeINH:    true,
		LanguageTypeINJ:    true,
		LanguageTypeINL:    true,
		LanguageTypeINM:    true,
		LanguageTypeINN:    true,
		LanguageTypeINO:    true,
		LanguageTypeINP:    true,
		LanguageTypeINS:    true,
		LanguageTypeINT:    true,
		LanguageTypeINZ:    true,
		LanguageTypeIOR:    true,
		LanguageTypeIOU:    true,
		LanguageTypeIOW:    true,
		LanguageTypeIPI:    true,
		LanguageTypeIPO:    true,
		LanguageTypeIQU:    true,
		LanguageTypeIQW:    true,
		LanguageTypeIRA:    true,
		LanguageTypeIRE:    true,
		LanguageTypeIRH:    true,
		LanguageTypeIRI:    true,
		LanguageTypeIRK:    true,
		LanguageTypeIRN:    true,
		LanguageTypeIRO:    true,
		LanguageTypeIRR:    true,
		LanguageTypeIRU:    true,
		LanguageTypeIRX:    true,
		LanguageTypeIRY:    true,
		LanguageTypeISA:    true,
		LanguageTypeISC:    true,
		LanguageTypeISD:    true,
		LanguageTypeISE:    true,
		LanguageTypeISG:    true,
		LanguageTypeISH:    true,
		LanguageTypeISI:    true,
		LanguageTypeISK:    true,
		LanguageTypeISM:    true,
		LanguageTypeISN:    true,
		LanguageTypeISO:    true,
		LanguageTypeISR:    true,
		LanguageTypeIST:    true,
		LanguageTypeISU:    true,
		LanguageTypeITB:    true,
		LanguageTypeITC:    true,
		LanguageTypeITE:    true,
		LanguageTypeITI:    true,
		LanguageTypeITK:    true,
		LanguageTypeITL:    true,
		LanguageTypeITM:    true,
		LanguageTypeITO:    true,
		LanguageTypeITR:    true,
		LanguageTypeITS:    true,
		LanguageTypeITT:    true,
		LanguageTypeITV:    true,
		LanguageTypeITW:    true,
		LanguageTypeITX:    true,
		LanguageTypeITY:    true,
		LanguageTypeITZ:    true,
		LanguageTypeIUM:    true,
		LanguageTypeIVB:    true,
		LanguageTypeIVV:    true,
		LanguageTypeIWK:    true,
		LanguageTypeIWM:    true,
		LanguageTypeIWO:    true,
		LanguageTypeIWS:    true,
		LanguageTypeIXC:    true,
		LanguageTypeIXL:    true,
		LanguageTypeIYA:    true,
		LanguageTypeIYO:    true,
		LanguageTypeIYX:    true,
		LanguageTypeIZH:    true,
		LanguageTypeIZI:    true,
		LanguageTypeIZR:    true,
		LanguageTypeIZZ:    true,
		LanguageTypeJAA:    true,
		LanguageTypeJAB:    true,
		LanguageTypeJAC:    true,
		LanguageTypeJAD:    true,
		LanguageTypeJAE:    true,
		LanguageTypeJAF:    true,
		LanguageTypeJAH:    true,
		LanguageTypeJAJ:    true,
		LanguageTypeJAK:    true,
		LanguageTypeJAL:    true,
		LanguageTypeJAM:    true,
		LanguageTypeJAN:    true,
		LanguageTypeJAO:    true,
		LanguageTypeJAQ:    true,
		LanguageTypeJAR:    true,
		LanguageTypeJAS:    true,
		LanguageTypeJAT:    true,
		LanguageTypeJAU:    true,
		LanguageTypeJAX:    true,
		LanguageTypeJAY:    true,
		LanguageTypeJAZ:    true,
		LanguageTypeJBE:    true,
		LanguageTypeJBI:    true,
		LanguageTypeJBJ:    true,
		LanguageTypeJBK:    true,
		LanguageTypeJBN:    true,
		LanguageTypeJBO:    true,
		LanguageTypeJBR:    true,
		LanguageTypeJBT:    true,
		LanguageTypeJBU:    true,
		LanguageTypeJBW:    true,
		LanguageTypeJCS:    true,
		LanguageTypeJCT:    true,
		LanguageTypeJDA:    true,
		LanguageTypeJDG:    true,
		LanguageTypeJDT:    true,
		LanguageTypeJEB:    true,
		LanguageTypeJEE:    true,
		LanguageTypeJEG:    true,
		LanguageTypeJEH:    true,
		LanguageTypeJEI:    true,
		LanguageTypeJEK:    true,
		LanguageTypeJEL:    true,
		LanguageTypeJEN:    true,
		LanguageTypeJER:    true,
		LanguageTypeJET:    true,
		LanguageTypeJEU:    true,
		LanguageTypeJGB:    true,
		LanguageTypeJGE:    true,
		LanguageTypeJGK:    true,
		LanguageTypeJGO:    true,
		LanguageTypeJHI:    true,
		LanguageTypeJHS:    true,
		LanguageTypeJIA:    true,
		LanguageTypeJIB:    true,
		LanguageTypeJIC:    true,
		LanguageTypeJID:    true,
		LanguageTypeJIE:    true,
		LanguageTypeJIG:    true,
		LanguageTypeJIH:    true,
		LanguageTypeJII:    true,
		LanguageTypeJIL:    true,
		LanguageTypeJIM:    true,
		LanguageTypeJIO:    true,
		LanguageTypeJIQ:    true,
		LanguageTypeJIT:    true,
		LanguageTypeJIU:    true,
		LanguageTypeJIV:    true,
		LanguageTypeJIY:    true,
		LanguageTypeJJR:    true,
		LanguageTypeJKM:    true,
		LanguageTypeJKO:    true,
		LanguageTypeJKP:    true,
		LanguageTypeJKR:    true,
		LanguageTypeJKU:    true,
		LanguageTypeJLE:    true,
		LanguageTypeJLS:    true,
		LanguageTypeJMA:    true,
		LanguageTypeJMB:    true,
		LanguageTypeJMC:    true,
		LanguageTypeJMD:    true,
		LanguageTypeJMI:    true,
		LanguageTypeJML:    true,
		LanguageTypeJMN:    true,
		LanguageTypeJMR:    true,
		LanguageTypeJMS:    true,
		LanguageTypeJMW:    true,
		LanguageTypeJMX:    true,
		LanguageTypeJNA:    true,
		LanguageTypeJND:    true,
		LanguageTypeJNG:    true,
		LanguageTypeJNI:    true,
		LanguageTypeJNJ:    true,
		LanguageTypeJNL:    true,
		LanguageTypeJNS:    true,
		LanguageTypeJOB:    true,
		LanguageTypeJOD:    true,
		LanguageTypeJOR:    true,
		LanguageTypeJOS:    true,
		LanguageTypeJOW:    true,
		LanguageTypeJPA:    true,
		LanguageTypeJPR:    true,
		LanguageTypeJPX:    true,
		LanguageTypeJQR:    true,
		LanguageTypeJRA:    true,
		LanguageTypeJRB:    true,
		LanguageTypeJRR:    true,
		LanguageTypeJRT:    true,
		LanguageTypeJRU:    true,
		LanguageTypeJSL:    true,
		LanguageTypeJUA:    true,
		LanguageTypeJUB:    true,
		LanguageTypeJUC:    true,
		LanguageTypeJUD:    true,
		LanguageTypeJUH:    true,
		LanguageTypeJUI:    true,
		LanguageTypeJUK:    true,
		LanguageTypeJUL:    true,
		LanguageTypeJUM:    true,
		LanguageTypeJUN:    true,
		LanguageTypeJUO:    true,
		LanguageTypeJUP:    true,
		LanguageTypeJUR:    true,
		LanguageTypeJUS:    true,
		LanguageTypeJUT:    true,
		LanguageTypeJUU:    true,
		LanguageTypeJUW:    true,
		LanguageTypeJUY:    true,
		LanguageTypeJVD:    true,
		LanguageTypeJVN:    true,
		LanguageTypeJWI:    true,
		LanguageTypeJYA:    true,
		LanguageTypeJYE:    true,
		LanguageTypeJYY:    true,
		LanguageTypeKAA:    true,
		LanguageTypeKAB:    true,
		LanguageTypeKAC:    true,
		LanguageTypeKAD:    true,
		LanguageTypeKAE:    true,
		LanguageTypeKAF:    true,
		LanguageTypeKAG:    true,
		LanguageTypeKAH:    true,
		LanguageTypeKAI:    true,
		LanguageTypeKAJ:    true,
		LanguageTypeKAK:    true,
		LanguageTypeKAM:    true,
		LanguageTypeKAO:    true,
		LanguageTypeKAP:    true,
		LanguageTypeKAQ:    true,
		LanguageTypeKAR:    true,
		LanguageTypeKAV:    true,
		LanguageTypeKAW:    true,
		LanguageTypeKAX:    true,
		LanguageTypeKAY:    true,
		LanguageTypeKBA:    true,
		LanguageTypeKBB:    true,
		LanguageTypeKBC:    true,
		LanguageTypeKBD:    true,
		LanguageTypeKBE:    true,
		LanguageTypeKBF:    true,
		LanguageTypeKBG:    true,
		LanguageTypeKBH:    true,
		LanguageTypeKBI:    true,
		LanguageTypeKBJ:    true,
		LanguageTypeKBK:    true,
		LanguageTypeKBL:    true,
		LanguageTypeKBM:    true,
		LanguageTypeKBN:    true,
		LanguageTypeKBO:    true,
		LanguageTypeKBP:    true,
		LanguageTypeKBQ:    true,
		LanguageTypeKBR:    true,
		LanguageTypeKBS:    true,
		LanguageTypeKBT:    true,
		LanguageTypeKBU:    true,
		LanguageTypeKBV:    true,
		LanguageTypeKBW:    true,
		LanguageTypeKBX:    true,
		LanguageTypeKBY:    true,
		LanguageTypeKBZ:    true,
		LanguageTypeKCA:    true,
		LanguageTypeKCB:    true,
		LanguageTypeKCC:    true,
		LanguageTypeKCD:    true,
		LanguageTypeKCE:    true,
		LanguageTypeKCF:    true,
		LanguageTypeKCG:    true,
		LanguageTypeKCH:    true,
		LanguageTypeKCI:    true,
		LanguageTypeKCJ:    true,
		LanguageTypeKCK:    true,
		LanguageTypeKCL:    true,
		LanguageTypeKCM:    true,
		LanguageTypeKCN:    true,
		LanguageTypeKCO:    true,
		LanguageTypeKCP:    true,
		LanguageTypeKCQ:    true,
		LanguageTypeKCR:    true,
		LanguageTypeKCS:    true,
		LanguageTypeKCT:    true,
		LanguageTypeKCU:    true,
		LanguageTypeKCV:    true,
		LanguageTypeKCW:    true,
		LanguageTypeKCX:    true,
		LanguageTypeKCY:    true,
		LanguageTypeKCZ:    true,
		LanguageTypeKDA:    true,
		LanguageTypeKDC:    true,
		LanguageTypeKDD:    true,
		LanguageTypeKDE:    true,
		LanguageTypeKDF:    true,
		LanguageTypeKDG:    true,
		LanguageTypeKDH:    true,
		LanguageTypeKDI:    true,
		LanguageTypeKDJ:    true,
		LanguageTypeKDK:    true,
		LanguageTypeKDL:    true,
		LanguageTypeKDM:    true,
		LanguageTypeKDN:    true,
		LanguageTypeKDO:    true,
		LanguageTypeKDP:    true,
		LanguageTypeKDQ:    true,
		LanguageTypeKDR:    true,
		LanguageTypeKDT:    true,
		LanguageTypeKDU:    true,
		LanguageTypeKDV:    true,
		LanguageTypeKDW:    true,
		LanguageTypeKDX:    true,
		LanguageTypeKDY:    true,
		LanguageTypeKDZ:    true,
		LanguageTypeKEA:    true,
		LanguageTypeKEB:    true,
		LanguageTypeKEC:    true,
		LanguageTypeKED:    true,
		LanguageTypeKEE:    true,
		LanguageTypeKEF:    true,
		LanguageTypeKEG:    true,
		LanguageTypeKEH:    true,
		LanguageTypeKEI:    true,
		LanguageTypeKEJ:    true,
		LanguageTypeKEK:    true,
		LanguageTypeKEL:    true,
		LanguageTypeKEM:    true,
		LanguageTypeKEN:    true,
		LanguageTypeKEO:    true,
		LanguageTypeKEP:    true,
		LanguageTypeKEQ:    true,
		LanguageTypeKER:    true,
		LanguageTypeKES:    true,
		LanguageTypeKET:    true,
		LanguageTypeKEU:    true,
		LanguageTypeKEV:    true,
		LanguageTypeKEW:    true,
		LanguageTypeKEX:    true,
		LanguageTypeKEY:    true,
		LanguageTypeKEZ:    true,
		LanguageTypeKFA:    true,
		LanguageTypeKFB:    true,
		LanguageTypeKFC:    true,
		LanguageTypeKFD:    true,
		LanguageTypeKFE:    true,
		LanguageTypeKFF:    true,
		LanguageTypeKFG:    true,
		LanguageTypeKFH:    true,
		LanguageTypeKFI:    true,
		LanguageTypeKFJ:    true,
		LanguageTypeKFK:    true,
		LanguageTypeKFL:    true,
		LanguageTypeKFM:    true,
		LanguageTypeKFN:    true,
		LanguageTypeKFO:    true,
		LanguageTypeKFP:    true,
		LanguageTypeKFQ:    true,
		LanguageTypeKFR:    true,
		LanguageTypeKFS:    true,
		LanguageTypeKFT:    true,
		LanguageTypeKFU:    true,
		LanguageTypeKFV:    true,
		LanguageTypeKFW:    true,
		LanguageTypeKFX:    true,
		LanguageTypeKFY:    true,
		LanguageTypeKFZ:    true,
		LanguageTypeKGA:    true,
		LanguageTypeKGB:    true,
		LanguageTypeKGC:    true,
		LanguageTypeKGD:    true,
		LanguageTypeKGE:    true,
		LanguageTypeKGF:    true,
		LanguageTypeKGG:    true,
		LanguageTypeKGH:    true,
		LanguageTypeKGI:    true,
		LanguageTypeKGJ:    true,
		LanguageTypeKGK:    true,
		LanguageTypeKGL:    true,
		LanguageTypeKGM:    true,
		LanguageTypeKGN:    true,
		LanguageTypeKGO:    true,
		LanguageTypeKGP:    true,
		LanguageTypeKGQ:    true,
		LanguageTypeKGR:    true,
		LanguageTypeKGS:    true,
		LanguageTypeKGT:    true,
		LanguageTypeKGU:    true,
		LanguageTypeKGV:    true,
		LanguageTypeKGW:    true,
		LanguageTypeKGX:    true,
		LanguageTypeKGY:    true,
		LanguageTypeKHA:    true,
		LanguageTypeKHB:    true,
		LanguageTypeKHC:    true,
		LanguageTypeKHD:    true,
		LanguageTypeKHE:    true,
		LanguageTypeKHF:    true,
		LanguageTypeKHG:    true,
		LanguageTypeKHH:    true,
		LanguageTypeKHI:    true,
		LanguageTypeKHJ:    true,
		LanguageTypeKHK:    true,
		LanguageTypeKHL:    true,
		LanguageTypeKHN:    true,
		LanguageTypeKHO:    true,
		LanguageTypeKHP:    true,
		LanguageTypeKHQ:    true,
		LanguageTypeKHR:    true,
		LanguageTypeKHS:    true,
		LanguageTypeKHT:    true,
		LanguageTypeKHU:    true,
		LanguageTypeKHV:    true,
		LanguageTypeKHW:    true,
		LanguageTypeKHX:    true,
		LanguageTypeKHY:    true,
		LanguageTypeKHZ:    true,
		LanguageTypeKIA:    true,
		LanguageTypeKIB:    true,
		LanguageTypeKIC:    true,
		LanguageTypeKID:    true,
		LanguageTypeKIE:    true,
		LanguageTypeKIF:    true,
		LanguageTypeKIG:    true,
		LanguageTypeKIH:    true,
		LanguageTypeKII:    true,
		LanguageTypeKIJ:    true,
		LanguageTypeKIL:    true,
		LanguageTypeKIM:    true,
		LanguageTypeKIO:    true,
		LanguageTypeKIP:    true,
		LanguageTypeKIQ:    true,
		LanguageTypeKIS:    true,
		LanguageTypeKIT:    true,
		LanguageTypeKIU:    true,
		LanguageTypeKIV:    true,
		LanguageTypeKIW:    true,
		LanguageTypeKIX:    true,
		LanguageTypeKIY:    true,
		LanguageTypeKIZ:    true,
		LanguageTypeKJA:    true,
		LanguageTypeKJB:    true,
		LanguageTypeKJC:    true,
		LanguageTypeKJD:    true,
		LanguageTypeKJE:    true,
		LanguageTypeKJF:    true,
		LanguageTypeKJG:    true,
		LanguageTypeKJH:    true,
		LanguageTypeKJI:    true,
		LanguageTypeKJJ:    true,
		LanguageTypeKJK:    true,
		LanguageTypeKJL:    true,
		LanguageTypeKJM:    true,
		LanguageTypeKJN:    true,
		LanguageTypeKJO:    true,
		LanguageTypeKJP:    true,
		LanguageTypeKJQ:    true,
		LanguageTypeKJR:    true,
		LanguageTypeKJS:    true,
		LanguageTypeKJT:    true,
		LanguageTypeKJU:    true,
		LanguageTypeKJX:    true,
		LanguageTypeKJY:    true,
		LanguageTypeKJZ:    true,
		LanguageTypeKKA:    true,
		LanguageTypeKKB:    true,
		LanguageTypeKKC:    true,
		LanguageTypeKKD:    true,
		LanguageTypeKKE:    true,
		LanguageTypeKKF:    true,
		LanguageTypeKKG:    true,
		LanguageTypeKKH:    true,
		LanguageTypeKKI:    true,
		LanguageTypeKKJ:    true,
		LanguageTypeKKK:    true,
		LanguageTypeKKL:    true,
		LanguageTypeKKM:    true,
		LanguageTypeKKN:    true,
		LanguageTypeKKO:    true,
		LanguageTypeKKP:    true,
		LanguageTypeKKQ:    true,
		LanguageTypeKKR:    true,
		LanguageTypeKKS:    true,
		LanguageTypeKKT:    true,
		LanguageTypeKKU:    true,
		LanguageTypeKKV:    true,
		LanguageTypeKKW:    true,
		LanguageTypeKKX:    true,
		LanguageTypeKKY:    true,
		LanguageTypeKKZ:    true,
		LanguageTypeKLA:    true,
		LanguageTypeKLB:    true,
		LanguageTypeKLC:    true,
		LanguageTypeKLD:    true,
		LanguageTypeKLE:    true,
		LanguageTypeKLF:    true,
		LanguageTypeKLG:    true,
		LanguageTypeKLH:    true,
		LanguageTypeKLI:    true,
		LanguageTypeKLJ:    true,
		LanguageTypeKLK:    true,
		LanguageTypeKLL:    true,
		LanguageTypeKLM:    true,
		LanguageTypeKLN:    true,
		LanguageTypeKLO:    true,
		LanguageTypeKLP:    true,
		LanguageTypeKLQ:    true,
		LanguageTypeKLR:    true,
		LanguageTypeKLS:    true,
		LanguageTypeKLT:    true,
		LanguageTypeKLU:    true,
		LanguageTypeKLV:    true,
		LanguageTypeKLW:    true,
		LanguageTypeKLX:    true,
		LanguageTypeKLY:    true,
		LanguageTypeKLZ:    true,
		LanguageTypeKMA:    true,
		LanguageTypeKMB:    true,
		LanguageTypeKMC:    true,
		LanguageTypeKMD:    true,
		LanguageTypeKME:    true,
		LanguageTypeKMF:    true,
		LanguageTypeKMG:    true,
		LanguageTypeKMH:    true,
		LanguageTypeKMI:    true,
		LanguageTypeKMJ:    true,
		LanguageTypeKMK:    true,
		LanguageTypeKML:    true,
		LanguageTypeKMM:    true,
		LanguageTypeKMN:    true,
		LanguageTypeKMO:    true,
		LanguageTypeKMP:    true,
		LanguageTypeKMQ:    true,
		LanguageTypeKMR:    true,
		LanguageTypeKMS:    true,
		LanguageTypeKMT:    true,
		LanguageTypeKMU:    true,
		LanguageTypeKMV:    true,
		LanguageTypeKMW:    true,
		LanguageTypeKMX:    true,
		LanguageTypeKMY:    true,
		LanguageTypeKMZ:    true,
		LanguageTypeKNA:    true,
		LanguageTypeKNB:    true,
		LanguageTypeKNC:    true,
		LanguageTypeKND:    true,
		LanguageTypeKNE:    true,
		LanguageTypeKNF:    true,
		LanguageTypeKNG:    true,
		LanguageTypeKNI:    true,
		LanguageTypeKNJ:    true,
		LanguageTypeKNK:    true,
		LanguageTypeKNL:    true,
		LanguageTypeKNM:    true,
		LanguageTypeKNN:    true,
		LanguageTypeKNO:    true,
		LanguageTypeKNP:    true,
		LanguageTypeKNQ:    true,
		LanguageTypeKNR:    true,
		LanguageTypeKNS:    true,
		LanguageTypeKNT:    true,
		LanguageTypeKNU:    true,
		LanguageTypeKNV:    true,
		LanguageTypeKNW:    true,
		LanguageTypeKNX:    true,
		LanguageTypeKNY:    true,
		LanguageTypeKNZ:    true,
		LanguageTypeKOA:    true,
		LanguageTypeKOC:    true,
		LanguageTypeKOD:    true,
		LanguageTypeKOE:    true,
		LanguageTypeKOF:    true,
		LanguageTypeKOG:    true,
		LanguageTypeKOH:    true,
		LanguageTypeKOI:    true,
		LanguageTypeKOJ:    true,
		LanguageTypeKOK:    true,
		LanguageTypeKOL:    true,
		LanguageTypeKOO:    true,
		LanguageTypeKOP:    true,
		LanguageTypeKOQ:    true,
		LanguageTypeKOS:    true,
		LanguageTypeKOT:    true,
		LanguageTypeKOU:    true,
		LanguageTypeKOV:    true,
		LanguageTypeKOW:    true,
		LanguageTypeKOX:    true,
		LanguageTypeKOY:    true,
		LanguageTypeKOZ:    true,
		LanguageTypeKPA:    true,
		LanguageTypeKPB:    true,
		LanguageTypeKPC:    true,
		LanguageTypeKPD:    true,
		LanguageTypeKPE:    true,
		LanguageTypeKPF:    true,
		LanguageTypeKPG:    true,
		LanguageTypeKPH:    true,
		LanguageTypeKPI:    true,
		LanguageTypeKPJ:    true,
		LanguageTypeKPK:    true,
		LanguageTypeKPL:    true,
		LanguageTypeKPM:    true,
		LanguageTypeKPN:    true,
		LanguageTypeKPO:    true,
		LanguageTypeKPP:    true,
		LanguageTypeKPQ:    true,
		LanguageTypeKPR:    true,
		LanguageTypeKPS:    true,
		LanguageTypeKPT:    true,
		LanguageTypeKPU:    true,
		LanguageTypeKPV:    true,
		LanguageTypeKPW:    true,
		LanguageTypeKPX:    true,
		LanguageTypeKPY:    true,
		LanguageTypeKPZ:    true,
		LanguageTypeKQA:    true,
		LanguageTypeKQB:    true,
		LanguageTypeKQC:    true,
		LanguageTypeKQD:    true,
		LanguageTypeKQE:    true,
		LanguageTypeKQF:    true,
		LanguageTypeKQG:    true,
		LanguageTypeKQH:    true,
		LanguageTypeKQI:    true,
		LanguageTypeKQJ:    true,
		LanguageTypeKQK:    true,
		LanguageTypeKQL:    true,
		LanguageTypeKQM:    true,
		LanguageTypeKQN:    true,
		LanguageTypeKQO:    true,
		LanguageTypeKQP:    true,
		LanguageTypeKQQ:    true,
		LanguageTypeKQR:    true,
		LanguageTypeKQS:    true,
		LanguageTypeKQT:    true,
		LanguageTypeKQU:    true,
		LanguageTypeKQV:    true,
		LanguageTypeKQW:    true,
		LanguageTypeKQX:    true,
		LanguageTypeKQY:    true,
		LanguageTypeKQZ:    true,
		LanguageTypeKRA:    true,
		LanguageTypeKRB:    true,
		LanguageTypeKRC:    true,
		LanguageTypeKRD:    true,
		LanguageTypeKRE:    true,
		LanguageTypeKRF:    true,
		LanguageTypeKRH:    true,
		LanguageTypeKRI:    true,
		LanguageTypeKRJ:    true,
		LanguageTypeKRK:    true,
		LanguageTypeKRL:    true,
		LanguageTypeKRM:    true,
		LanguageTypeKRN:    true,
		LanguageTypeKRO:    true,
		LanguageTypeKRP:    true,
		LanguageTypeKRR:    true,
		LanguageTypeKRS:    true,
		LanguageTypeKRT:    true,
		LanguageTypeKRU:    true,
		LanguageTypeKRV:    true,
		LanguageTypeKRW:    true,
		LanguageTypeKRX:    true,
		LanguageTypeKRY:    true,
		LanguageTypeKRZ:    true,
		LanguageTypeKSA:    true,
		LanguageTypeKSB:    true,
		LanguageTypeKSC:    true,
		LanguageTypeKSD:    true,
		LanguageTypeKSE:    true,
		LanguageTypeKSF:    true,
		LanguageTypeKSG:    true,
		LanguageTypeKSH:    true,
		LanguageTypeKSI:    true,
		LanguageTypeKSJ:    true,
		LanguageTypeKSK:    true,
		LanguageTypeKSL:    true,
		LanguageTypeKSM:    true,
		LanguageTypeKSN:    true,
		LanguageTypeKSO:    true,
		LanguageTypeKSP:    true,
		LanguageTypeKSQ:    true,
		LanguageTypeKSR:    true,
		LanguageTypeKSS:    true,
		LanguageTypeKST:    true,
		LanguageTypeKSU:    true,
		LanguageTypeKSV:    true,
		LanguageTypeKSW:    true,
		LanguageTypeKSX:    true,
		LanguageTypeKSY:    true,
		LanguageTypeKSZ:    true,
		LanguageTypeKTA:    true,
		LanguageTypeKTB:    true,
		LanguageTypeKTC:    true,
		LanguageTypeKTD:    true,
		LanguageTypeKTE:    true,
		LanguageTypeKTF:    true,
		LanguageTypeKTG:    true,
		LanguageTypeKTH:    true,
		LanguageTypeKTI:    true,
		LanguageTypeKTJ:    true,
		LanguageTypeKTK:    true,
		LanguageTypeKTL:    true,
		LanguageTypeKTM:    true,
		LanguageTypeKTN:    true,
		LanguageTypeKTO:    true,
		LanguageTypeKTP:    true,
		LanguageTypeKTQ:    true,
		LanguageTypeKTR:    true,
		LanguageTypeKTS:    true,
		LanguageTypeKTT:    true,
		LanguageTypeKTU:    true,
		LanguageTypeKTV:    true,
		LanguageTypeKTW:    true,
		LanguageTypeKTX:    true,
		LanguageTypeKTY:    true,
		LanguageTypeKTZ:    true,
		LanguageTypeKUB:    true,
		LanguageTypeKUC:    true,
		LanguageTypeKUD:    true,
		LanguageTypeKUE:    true,
		LanguageTypeKUF:    true,
		LanguageTypeKUG:    true,
		LanguageTypeKUH:    true,
		LanguageTypeKUI:    true,
		LanguageTypeKUJ:    true,
		LanguageTypeKUK:    true,
		LanguageTypeKUL:    true,
		LanguageTypeKUM:    true,
		LanguageTypeKUN:    true,
		LanguageTypeKUO:    true,
		LanguageTypeKUP:    true,
		LanguageTypeKUQ:    true,
		LanguageTypeKUS:    true,
		LanguageTypeKUT:    true,
		LanguageTypeKUU:    true,
		LanguageTypeKUV:    true,
		LanguageTypeKUW:    true,
		LanguageTypeKUX:    true,
		LanguageTypeKUY:    true,
		LanguageTypeKUZ:    true,
		LanguageTypeKVA:    true,
		LanguageTypeKVB:    true,
		LanguageTypeKVC:    true,
		LanguageTypeKVD:    true,
		LanguageTypeKVE:    true,
		LanguageTypeKVF:    true,
		LanguageTypeKVG:    true,
		LanguageTypeKVH:    true,
		LanguageTypeKVI:    true,
		LanguageTypeKVJ:    true,
		LanguageTypeKVK:    true,
		LanguageTypeKVL:    true,
		LanguageTypeKVM:    true,
		LanguageTypeKVN:    true,
		LanguageTypeKVO:    true,
		LanguageTypeKVP:    true,
		LanguageTypeKVQ:    true,
		LanguageTypeKVR:    true,
		LanguageTypeKVS:    true,
		LanguageTypeKVT:    true,
		LanguageTypeKVU:    true,
		LanguageTypeKVV:    true,
		LanguageTypeKVW:    true,
		LanguageTypeKVX:    true,
		LanguageTypeKVY:    true,
		LanguageTypeKVZ:    true,
		LanguageTypeKWA:    true,
		LanguageTypeKWB:    true,
		LanguageTypeKWC:    true,
		LanguageTypeKWD:    true,
		LanguageTypeKWE:    true,
		LanguageTypeKWF:    true,
		LanguageTypeKWG:    true,
		LanguageTypeKWH:    true,
		LanguageTypeKWI:    true,
		LanguageTypeKWJ:    true,
		LanguageTypeKWK:    true,
		LanguageTypeKWL:    true,
		LanguageTypeKWM:    true,
		LanguageTypeKWN:    true,
		LanguageTypeKWO:    true,
		LanguageTypeKWP:    true,
		LanguageTypeKWQ:    true,
		LanguageTypeKWR:    true,
		LanguageTypeKWS:    true,
		LanguageTypeKWT:    true,
		LanguageTypeKWU:    true,
		LanguageTypeKWV:    true,
		LanguageTypeKWW:    true,
		LanguageTypeKWX:    true,
		LanguageTypeKWY:    true,
		LanguageTypeKWZ:    true,
		LanguageTypeKXA:    true,
		LanguageTypeKXB:    true,
		LanguageTypeKXC:    true,
		LanguageTypeKXD:    true,
		LanguageTypeKXE:    true,
		LanguageTypeKXF:    true,
		LanguageTypeKXH:    true,
		LanguageTypeKXI:    true,
		LanguageTypeKXJ:    true,
		LanguageTypeKXK:    true,
		LanguageTypeKXL:    true,
		LanguageTypeKXM:    true,
		LanguageTypeKXN:    true,
		LanguageTypeKXO:    true,
		LanguageTypeKXP:    true,
		LanguageTypeKXQ:    true,
		LanguageTypeKXR:    true,
		LanguageTypeKXS:    true,
		LanguageTypeKXT:    true,
		LanguageTypeKXU:    true,
		LanguageTypeKXV:    true,
		LanguageTypeKXW:    true,
		LanguageTypeKXX:    true,
		LanguageTypeKXY:    true,
		LanguageTypeKXZ:    true,
		LanguageTypeKYA:    true,
		LanguageTypeKYB:    true,
		LanguageTypeKYC:    true,
		LanguageTypeKYD:    true,
		LanguageTypeKYE:    true,
		LanguageTypeKYF:    true,
		LanguageTypeKYG:    true,
		LanguageTypeKYH:    true,
		LanguageTypeKYI:    true,
		LanguageTypeKYJ:    true,
		LanguageTypeKYK:    true,
		LanguageTypeKYL:    true,
		LanguageTypeKYM:    true,
		LanguageTypeKYN:    true,
		LanguageTypeKYO:    true,
		LanguageTypeKYP:    true,
		LanguageTypeKYQ:    true,
		LanguageTypeKYR:    true,
		LanguageTypeKYS:    true,
		LanguageTypeKYT:    true,
		LanguageTypeKYU:    true,
		LanguageTypeKYV:    true,
		LanguageTypeKYW:    true,
		LanguageTypeKYX:    true,
		LanguageTypeKYY:    true,
		LanguageTypeKYZ:    true,
		LanguageTypeKZA:    true,
		LanguageTypeKZB:    true,
		LanguageTypeKZC:    true,
		LanguageTypeKZD:    true,
		LanguageTypeKZE:    true,
		LanguageTypeKZF:    true,
		LanguageTypeKZG:    true,
		LanguageTypeKZH:    true,
		LanguageTypeKZI:    true,
		LanguageTypeKZJ:    true,
		LanguageTypeKZK:    true,
		LanguageTypeKZL:    true,
		LanguageTypeKZM:    true,
		LanguageTypeKZN:    true,
		LanguageTypeKZO:    true,
		LanguageTypeKZP:    true,
		LanguageTypeKZQ:    true,
		LanguageTypeKZR:    true,
		LanguageTypeKZS:    true,
		LanguageTypeKZT:    true,
		LanguageTypeKZU:    true,
		LanguageTypeKZV:    true,
		LanguageTypeKZW:    true,
		LanguageTypeKZX:    true,
		LanguageTypeKZY:    true,
		LanguageTypeKZZ:    true,
		LanguageTypeLAA:    true,
		LanguageTypeLAB:    true,
		LanguageTypeLAC:    true,
		LanguageTypeLAD:    true,
		LanguageTypeLAE:    true,
		LanguageTypeLAF:    true,
		LanguageTypeLAG:    true,
		LanguageTypeLAH:    true,
		LanguageTypeLAI:    true,
		LanguageTypeLAJ:    true,
		LanguageTypeLAK:    true,
		LanguageTypeLAL:    true,
		LanguageTypeLAM:    true,
		LanguageTypeLAN:    true,
		LanguageTypeLAP:    true,
		LanguageTypeLAQ:    true,
		LanguageTypeLAR:    true,
		LanguageTypeLAS:    true,
		LanguageTypeLAU:    true,
		LanguageTypeLAW:    true,
		LanguageTypeLAX:    true,
		LanguageTypeLAY:    true,
		LanguageTypeLAZ:    true,
		LanguageTypeLBA:    true,
		LanguageTypeLBB:    true,
		LanguageTypeLBC:    true,
		LanguageTypeLBE:    true,
		LanguageTypeLBF:    true,
		LanguageTypeLBG:    true,
		LanguageTypeLBI:    true,
		LanguageTypeLBJ:    true,
		LanguageTypeLBK:    true,
		LanguageTypeLBL:    true,
		LanguageTypeLBM:    true,
		LanguageTypeLBN:    true,
		LanguageTypeLBO:    true,
		LanguageTypeLBQ:    true,
		LanguageTypeLBR:    true,
		LanguageTypeLBS:    true,
		LanguageTypeLBT:    true,
		LanguageTypeLBU:    true,
		LanguageTypeLBV:    true,
		LanguageTypeLBW:    true,
		LanguageTypeLBX:    true,
		LanguageTypeLBY:    true,
		LanguageTypeLBZ:    true,
		LanguageTypeLCC:    true,
		LanguageTypeLCD:    true,
		LanguageTypeLCE:    true,
		LanguageTypeLCF:    true,
		LanguageTypeLCH:    true,
		LanguageTypeLCL:    true,
		LanguageTypeLCM:    true,
		LanguageTypeLCP:    true,
		LanguageTypeLCQ:    true,
		LanguageTypeLCS:    true,
		LanguageTypeLDA:    true,
		LanguageTypeLDB:    true,
		LanguageTypeLDD:    true,
		LanguageTypeLDG:    true,
		LanguageTypeLDH:    true,
		LanguageTypeLDI:    true,
		LanguageTypeLDJ:    true,
		LanguageTypeLDK:    true,
		LanguageTypeLDL:    true,
		LanguageTypeLDM:    true,
		LanguageTypeLDN:    true,
		LanguageTypeLDO:    true,
		LanguageTypeLDP:    true,
		LanguageTypeLDQ:    true,
		LanguageTypeLEA:    true,
		LanguageTypeLEB:    true,
		LanguageTypeLEC:    true,
		LanguageTypeLED:    true,
		LanguageTypeLEE:    true,
		LanguageTypeLEF:    true,
		LanguageTypeLEG:    true,
		LanguageTypeLEH:    true,
		LanguageTypeLEI:    true,
		LanguageTypeLEJ:    true,
		LanguageTypeLEK:    true,
		LanguageTypeLEL:    true,
		LanguageTypeLEM:    true,
		LanguageTypeLEN:    true,
		LanguageTypeLEO:    true,
		LanguageTypeLEP:    true,
		LanguageTypeLEQ:    true,
		LanguageTypeLER:    true,
		LanguageTypeLES:    true,
		LanguageTypeLET:    true,
		LanguageTypeLEU:    true,
		LanguageTypeLEV:    true,
		LanguageTypeLEW:    true,
		LanguageTypeLEX:    true,
		LanguageTypeLEY:    true,
		LanguageTypeLEZ:    true,
		LanguageTypeLFA:    true,
		LanguageTypeLFN:    true,
		LanguageTypeLGA:    true,
		LanguageTypeLGB:    true,
		LanguageTypeLGG:    true,
		LanguageTypeLGH:    true,
		LanguageTypeLGI:    true,
		LanguageTypeLGK:    true,
		LanguageTypeLGL:    true,
		LanguageTypeLGM:    true,
		LanguageTypeLGN:    true,
		LanguageTypeLGQ:    true,
		LanguageTypeLGR:    true,
		LanguageTypeLGT:    true,
		LanguageTypeLGU:    true,
		LanguageTypeLGZ:    true,
		LanguageTypeLHA:    true,
		LanguageTypeLHH:    true,
		LanguageTypeLHI:    true,
		LanguageTypeLHL:    true,
		LanguageTypeLHM:    true,
		LanguageTypeLHN:    true,
		LanguageTypeLHP:    true,
		LanguageTypeLHS:    true,
		LanguageTypeLHT:    true,
		LanguageTypeLHU:    true,
		LanguageTypeLIA:    true,
		LanguageTypeLIB:    true,
		LanguageTypeLIC:    true,
		LanguageTypeLID:    true,
		LanguageTypeLIE:    true,
		LanguageTypeLIF:    true,
		LanguageTypeLIG:    true,
		LanguageTypeLIH:    true,
		LanguageTypeLII:    true,
		LanguageTypeLIJ:    true,
		LanguageTypeLIK:    true,
		LanguageTypeLIL:    true,
		LanguageTypeLIO:    true,
		LanguageTypeLIP:    true,
		LanguageTypeLIQ:    true,
		LanguageTypeLIR:    true,
		LanguageTypeLIS:    true,
		LanguageTypeLIU:    true,
		LanguageTypeLIV:    true,
		LanguageTypeLIW:    true,
		LanguageTypeLIX:    true,
		LanguageTypeLIY:    true,
		LanguageTypeLIZ:    true,
		LanguageTypeLJA:    true,
		LanguageTypeLJE:    true,
		LanguageTypeLJI:    true,
		LanguageTypeLJL:    true,
		LanguageTypeLJP:    true,
		LanguageTypeLJW:    true,
		LanguageTypeLJX:    true,
		LanguageTypeLKA:    true,
		LanguageTypeLKB:    true,
		LanguageTypeLKC:    true,
		LanguageTypeLKD:    true,
		LanguageTypeLKE:    true,
		LanguageTypeLKH:    true,
		LanguageTypeLKI:    true,
		LanguageTypeLKJ:    true,
		LanguageTypeLKL:    true,
		LanguageTypeLKM:    true,
		LanguageTypeLKN:    true,
		LanguageTypeLKO:    true,
		LanguageTypeLKR:    true,
		LanguageTypeLKS:    true,
		LanguageTypeLKT:    true,
		LanguageTypeLKU:    true,
		LanguageTypeLKY:    true,
		LanguageTypeLLA:    true,
		LanguageTypeLLB:    true,
		LanguageTypeLLC:    true,
		LanguageTypeLLD:    true,
		LanguageTypeLLE:    true,
		LanguageTypeLLF:    true,
		LanguageTypeLLG:    true,
		LanguageTypeLLH:    true,
		LanguageTypeLLI:    true,
		LanguageTypeLLJ:    true,
		LanguageTypeLLK:    true,
		LanguageTypeLLL:    true,
		LanguageTypeLLM:    true,
		LanguageTypeLLN:    true,
		LanguageTypeLLO:    true,
		LanguageTypeLLP:    true,
		LanguageTypeLLQ:    true,
		LanguageTypeLLS:    true,
		LanguageTypeLLU:    true,
		LanguageTypeLLX:    true,
		LanguageTypeLMA:    true,
		LanguageTypeLMB:    true,
		LanguageTypeLMC:    true,
		LanguageTypeLMD:    true,
		LanguageTypeLME:    true,
		LanguageTypeLMF:    true,
		LanguageTypeLMG:    true,
		LanguageTypeLMH:    true,
		LanguageTypeLMI:    true,
		LanguageTypeLMJ:    true,
		LanguageTypeLMK:    true,
		LanguageTypeLML:    true,
		LanguageTypeLMM:    true,
		LanguageTypeLMN:    true,
		LanguageTypeLMO:    true,
		LanguageTypeLMP:    true,
		LanguageTypeLMQ:    true,
		LanguageTypeLMR:    true,
		LanguageTypeLMU:    true,
		LanguageTypeLMV:    true,
		LanguageTypeLMW:    true,
		LanguageTypeLMX:    true,
		LanguageTypeLMY:    true,
		LanguageTypeLMZ:    true,
		LanguageTypeLNA:    true,
		LanguageTypeLNB:    true,
		LanguageTypeLND:    true,
		LanguageTypeLNG:    true,
		LanguageTypeLNH:    true,
		LanguageTypeLNI:    true,
		LanguageTypeLNJ:    true,
		LanguageTypeLNL:    true,
		LanguageTypeLNM:    true,
		LanguageTypeLNN:    true,
		LanguageTypeLNO:    true,
		LanguageTypeLNS:    true,
		LanguageTypeLNU:    true,
		LanguageTypeLNW:    true,
		LanguageTypeLNZ:    true,
		LanguageTypeLOA:    true,
		LanguageTypeLOB:    true,
		LanguageTypeLOC:    true,
		LanguageTypeLOE:    true,
		LanguageTypeLOF:    true,
		LanguageTypeLOG:    true,
		LanguageTypeLOH:    true,
		LanguageTypeLOI:    true,
		LanguageTypeLOJ:    true,
		LanguageTypeLOK:    true,
		LanguageTypeLOL:    true,
		LanguageTypeLOM:    true,
		LanguageTypeLON:    true,
		LanguageTypeLOO:    true,
		LanguageTypeLOP:    true,
		LanguageTypeLOQ:    true,
		LanguageTypeLOR:    true,
		LanguageTypeLOS:    true,
		LanguageTypeLOT:    true,
		LanguageTypeLOU:    true,
		LanguageTypeLOV:    true,
		LanguageTypeLOW:    true,
		LanguageTypeLOX:    true,
		LanguageTypeLOY:    true,
		LanguageTypeLOZ:    true,
		LanguageTypeLPA:    true,
		LanguageTypeLPE:    true,
		LanguageTypeLPN:    true,
		LanguageTypeLPO:    true,
		LanguageTypeLPX:    true,
		LanguageTypeLRA:    true,
		LanguageTypeLRC:    true,
		LanguageTypeLRE:    true,
		LanguageTypeLRG:    true,
		LanguageTypeLRI:    true,
		LanguageTypeLRK:    true,
		LanguageTypeLRL:    true,
		LanguageTypeLRM:    true,
		LanguageTypeLRN:    true,
		LanguageTypeLRO:    true,
		LanguageTypeLRR:    true,
		LanguageTypeLRT:    true,
		LanguageTypeLRV:    true,
		LanguageTypeLRZ:    true,
		LanguageTypeLSA:    true,
		LanguageTypeLSD:    true,
		LanguageTypeLSE:    true,
		LanguageTypeLSG:    true,
		LanguageTypeLSH:    true,
		LanguageTypeLSI:    true,
		LanguageTypeLSL:    true,
		LanguageTypeLSM:    true,
		LanguageTypeLSO:    true,
		LanguageTypeLSP:    true,
		LanguageTypeLSR:    true,
		LanguageTypeLSS:    true,
		LanguageTypeLST:    true,
		LanguageTypeLSY:    true,
		LanguageTypeLTC:    true,
		LanguageTypeLTG:    true,
		LanguageTypeLTI:    true,
		LanguageTypeLTN:    true,
		LanguageTypeLTO:    true,
		LanguageTypeLTS:    true,
		LanguageTypeLTU:    true,
		LanguageTypeLUA:    true,
		LanguageTypeLUC:    true,
		LanguageTypeLUD:    true,
		LanguageTypeLUE:    true,
		LanguageTypeLUF:    true,
		LanguageTypeLUI:    true,
		LanguageTypeLUJ:    true,
		LanguageTypeLUK:    true,
		LanguageTypeLUL:    true,
		LanguageTypeLUM:    true,
		LanguageTypeLUN:    true,
		LanguageTypeLUO:    true,
		LanguageTypeLUP:    true,
		LanguageTypeLUQ:    true,
		LanguageTypeLUR:    true,
		LanguageTypeLUS:    true,
		LanguageTypeLUT:    true,
		LanguageTypeLUU:    true,
		LanguageTypeLUV:    true,
		LanguageTypeLUW:    true,
		LanguageTypeLUY:    true,
		LanguageTypeLUZ:    true,
		LanguageTypeLVA:    true,
		LanguageTypeLVK:    true,
		LanguageTypeLVS:    true,
		LanguageTypeLVU:    true,
		LanguageTypeLWA:    true,
		LanguageTypeLWE:    true,
		LanguageTypeLWG:    true,
		LanguageTypeLWH:    true,
		LanguageTypeLWL:    true,
		LanguageTypeLWM:    true,
		LanguageTypeLWO:    true,
		LanguageTypeLWT:    true,
		LanguageTypeLWU:    true,
		LanguageTypeLWW:    true,
		LanguageTypeLYA:    true,
		LanguageTypeLYG:    true,
		LanguageTypeLYN:    true,
		LanguageTypeLZH:    true,
		LanguageTypeLZL:    true,
		LanguageTypeLZN:    true,
		LanguageTypeLZZ:    true,
		LanguageTypeMAA:    true,
		LanguageTypeMAB:    true,
		LanguageTypeMAD:    true,
		LanguageTypeMAE:    true,
		LanguageTypeMAF:    true,
		LanguageTypeMAG:    true,
		LanguageTypeMAI:    true,
		LanguageTypeMAJ:    true,
		LanguageTypeMAK:    true,
		LanguageTypeMAM:    true,
		LanguageTypeMAN:    true,
		LanguageTypeMAP:    true,
		LanguageTypeMAQ:    true,
		LanguageTypeMAS:    true,
		LanguageTypeMAT:    true,
		LanguageTypeMAU:    true,
		LanguageTypeMAV:    true,
		LanguageTypeMAW:    true,
		LanguageTypeMAX:    true,
		LanguageTypeMAZ:    true,
		LanguageTypeMBA:    true,
		LanguageTypeMBB:    true,
		LanguageTypeMBC:    true,
		LanguageTypeMBD:    true,
		LanguageTypeMBE:    true,
		LanguageTypeMBF:    true,
		LanguageTypeMBH:    true,
		LanguageTypeMBI:    true,
		LanguageTypeMBJ:    true,
		LanguageTypeMBK:    true,
		LanguageTypeMBL:    true,
		LanguageTypeMBM:    true,
		LanguageTypeMBN:    true,
		LanguageTypeMBO:    true,
		LanguageTypeMBP:    true,
		LanguageTypeMBQ:    true,
		LanguageTypeMBR:    true,
		LanguageTypeMBS:    true,
		LanguageTypeMBT:    true,
		LanguageTypeMBU:    true,
		LanguageTypeMBV:    true,
		LanguageTypeMBW:    true,
		LanguageTypeMBX:    true,
		LanguageTypeMBY:    true,
		LanguageTypeMBZ:    true,
		LanguageTypeMCA:    true,
		LanguageTypeMCB:    true,
		LanguageTypeMCC:    true,
		LanguageTypeMCD:    true,
		LanguageTypeMCE:    true,
		LanguageTypeMCF:    true,
		LanguageTypeMCG:    true,
		LanguageTypeMCH:    true,
		LanguageTypeMCI:    true,
		LanguageTypeMCJ:    true,
		LanguageTypeMCK:    true,
		LanguageTypeMCL:    true,
		LanguageTypeMCM:    true,
		LanguageTypeMCN:    true,
		LanguageTypeMCO:    true,
		LanguageTypeMCP:    true,
		LanguageTypeMCQ:    true,
		LanguageTypeMCR:    true,
		LanguageTypeMCS:    true,
		LanguageTypeMCT:    true,
		LanguageTypeMCU:    true,
		LanguageTypeMCV:    true,
		LanguageTypeMCW:    true,
		LanguageTypeMCX:    true,
		LanguageTypeMCY:    true,
		LanguageTypeMCZ:    true,
		LanguageTypeMDA:    true,
		LanguageTypeMDB:    true,
		LanguageTypeMDC:    true,
		LanguageTypeMDD:    true,
		LanguageTypeMDE:    true,
		LanguageTypeMDF:    true,
		LanguageTypeMDG:    true,
		LanguageTypeMDH:    true,
		LanguageTypeMDI:    true,
		LanguageTypeMDJ:    true,
		LanguageTypeMDK:    true,
		LanguageTypeMDL:    true,
		LanguageTypeMDM:    true,
		LanguageTypeMDN:    true,
		LanguageTypeMDP:    true,
		LanguageTypeMDQ:    true,
		LanguageTypeMDR:    true,
		LanguageTypeMDS:    true,
		LanguageTypeMDT:    true,
		LanguageTypeMDU:    true,
		LanguageTypeMDV:    true,
		LanguageTypeMDW:    true,
		LanguageTypeMDX:    true,
		LanguageTypeMDY:    true,
		LanguageTypeMDZ:    true,
		LanguageTypeMEA:    true,
		LanguageTypeMEB:    true,
		LanguageTypeMEC:    true,
		LanguageTypeMED:    true,
		LanguageTypeMEE:    true,
		LanguageTypeMEF:    true,
		LanguageTypeMEG:    true,
		LanguageTypeMEH:    true,
		LanguageTypeMEI:    true,
		LanguageTypeMEJ:    true,
		LanguageTypeMEK:    true,
		LanguageTypeMEL:    true,
		LanguageTypeMEM:    true,
		LanguageTypeMEN:    true,
		LanguageTypeMEO:    true,
		LanguageTypeMEP:    true,
		LanguageTypeMEQ:    true,
		LanguageTypeMER:    true,
		LanguageTypeMES:    true,
		LanguageTypeMET:    true,
		LanguageTypeMEU:    true,
		LanguageTypeMEV:    true,
		LanguageTypeMEW:    true,
		LanguageTypeMEY:    true,
		LanguageTypeMEZ:    true,
		LanguageTypeMFA:    true,
		LanguageTypeMFB:    true,
		LanguageTypeMFC:    true,
		LanguageTypeMFD:    true,
		LanguageTypeMFE:    true,
		LanguageTypeMFF:    true,
		LanguageTypeMFG:    true,
		LanguageTypeMFH:    true,
		LanguageTypeMFI:    true,
		LanguageTypeMFJ:    true,
		LanguageTypeMFK:    true,
		LanguageTypeMFL:    true,
		LanguageTypeMFM:    true,
		LanguageTypeMFN:    true,
		LanguageTypeMFO:    true,
		LanguageTypeMFP:    true,
		LanguageTypeMFQ:    true,
		LanguageTypeMFR:    true,
		LanguageTypeMFS:    true,
		LanguageTypeMFT:    true,
		LanguageTypeMFU:    true,
		LanguageTypeMFV:    true,
		LanguageTypeMFW:    true,
		LanguageTypeMFX:    true,
		LanguageTypeMFY:    true,
		LanguageTypeMFZ:    true,
		LanguageTypeMGA:    true,
		LanguageTypeMGB:    true,
		LanguageTypeMGC:    true,
		LanguageTypeMGD:    true,
		LanguageTypeMGE:    true,
		LanguageTypeMGF:    true,
		LanguageTypeMGG:    true,
		LanguageTypeMGH:    true,
		LanguageTypeMGI:    true,
		LanguageTypeMGJ:    true,
		LanguageTypeMGK:    true,
		LanguageTypeMGL:    true,
		LanguageTypeMGM:    true,
		LanguageTypeMGN:    true,
		LanguageTypeMGO:    true,
		LanguageTypeMGP:    true,
		LanguageTypeMGQ:    true,
		LanguageTypeMGR:    true,
		LanguageTypeMGS:    true,
		LanguageTypeMGT:    true,
		LanguageTypeMGU:    true,
		LanguageTypeMGV:    true,
		LanguageTypeMGW:    true,
		LanguageTypeMGX:    true,
		LanguageTypeMGY:    true,
		LanguageTypeMGZ:    true,
		LanguageTypeMHA:    true,
		LanguageTypeMHB:    true,
		LanguageTypeMHC:    true,
		LanguageTypeMHD:    true,
		LanguageTypeMHE:    true,
		LanguageTypeMHF:    true,
		LanguageTypeMHG:    true,
		LanguageTypeMHH:    true,
		LanguageTypeMHI:    true,
		LanguageTypeMHJ:    true,
		LanguageTypeMHK:    true,
		LanguageTypeMHL:    true,
		LanguageTypeMHM:    true,
		LanguageTypeMHN:    true,
		LanguageTypeMHO:    true,
		LanguageTypeMHP:    true,
		LanguageTypeMHQ:    true,
		LanguageTypeMHR:    true,
		LanguageTypeMHS:    true,
		LanguageTypeMHT:    true,
		LanguageTypeMHU:    true,
		LanguageTypeMHW:    true,
		LanguageTypeMHX:    true,
		LanguageTypeMHY:    true,
		LanguageTypeMHZ:    true,
		LanguageTypeMIA:    true,
		LanguageTypeMIB:    true,
		LanguageTypeMIC:    true,
		LanguageTypeMID:    true,
		LanguageTypeMIE:    true,
		LanguageTypeMIF:    true,
		LanguageTypeMIG:    true,
		LanguageTypeMIH:    true,
		LanguageTypeMII:    true,
		LanguageTypeMIJ:    true,
		LanguageTypeMIK:    true,
		LanguageTypeMIL:    true,
		LanguageTypeMIM:    true,
		LanguageTypeMIN:    true,
		LanguageTypeMIO:    true,
		LanguageTypeMIP:    true,
		LanguageTypeMIQ:    true,
		LanguageTypeMIR:    true,
		LanguageTypeMIS:    true,
		LanguageTypeMIT:    true,
		LanguageTypeMIU:    true,
		LanguageTypeMIW:    true,
		LanguageTypeMIX:    true,
		LanguageTypeMIY:    true,
		LanguageTypeMIZ:    true,
		LanguageTypeMJA:    true,
		LanguageTypeMJC:    true,
		LanguageTypeMJD:    true,
		LanguageTypeMJE:    true,
		LanguageTypeMJG:    true,
		LanguageTypeMJH:    true,
		LanguageTypeMJI:    true,
		LanguageTypeMJJ:    true,
		LanguageTypeMJK:    true,
		LanguageTypeMJL:    true,
		LanguageTypeMJM:    true,
		LanguageTypeMJN:    true,
		LanguageTypeMJO:    true,
		LanguageTypeMJP:    true,
		LanguageTypeMJQ:    true,
		LanguageTypeMJR:    true,
		LanguageTypeMJS:    true,
		LanguageTypeMJT:    true,
		LanguageTypeMJU:    true,
		LanguageTypeMJV:    true,
		LanguageTypeMJW:    true,
		LanguageTypeMJX:    true,
		LanguageTypeMJY:    true,
		LanguageTypeMJZ:    true,
		LanguageTypeMKA:    true,
		LanguageTypeMKB:    true,
		LanguageTypeMKC:    true,
		LanguageTypeMKE:    true,
		LanguageTypeMKF:    true,
		LanguageTypeMKG:    true,
		LanguageTypeMKH:    true,
		LanguageTypeMKI:    true,
		LanguageTypeMKJ:    true,
		LanguageTypeMKK:    true,
		LanguageTypeMKL:    true,
		LanguageTypeMKM:    true,
		LanguageTypeMKN:    true,
		LanguageTypeMKO:    true,
		LanguageTypeMKP:    true,
		LanguageTypeMKQ:    true,
		LanguageTypeMKR:    true,
		LanguageTypeMKS:    true,
		LanguageTypeMKT:    true,
		LanguageTypeMKU:    true,
		LanguageTypeMKV:    true,
		LanguageTypeMKW:    true,
		LanguageTypeMKX:    true,
		LanguageTypeMKY:    true,
		LanguageTypeMKZ:    true,
		LanguageTypeMLA:    true,
		LanguageTypeMLB:    true,
		LanguageTypeMLC:    true,
		LanguageTypeMLD:    true,
		LanguageTypeMLE:    true,
		LanguageTypeMLF:    true,
		LanguageTypeMLH:    true,
		LanguageTypeMLI:    true,
		LanguageTypeMLJ:    true,
		LanguageTypeMLK:    true,
		LanguageTypeMLL:    true,
		LanguageTypeMLM:    true,
		LanguageTypeMLN:    true,
		LanguageTypeMLO:    true,
		LanguageTypeMLP:    true,
		LanguageTypeMLQ:    true,
		LanguageTypeMLR:    true,
		LanguageTypeMLS:    true,
		LanguageTypeMLU:    true,
		LanguageTypeMLV:    true,
		LanguageTypeMLW:    true,
		LanguageTypeMLX:    true,
		LanguageTypeMLZ:    true,
		LanguageTypeMMA:    true,
		LanguageTypeMMB:    true,
		LanguageTypeMMC:    true,
		LanguageTypeMMD:    true,
		LanguageTypeMME:    true,
		LanguageTypeMMF:    true,
		LanguageTypeMMG:    true,
		LanguageTypeMMH:    true,
		LanguageTypeMMI:    true,
		LanguageTypeMMJ:    true,
		LanguageTypeMMK:    true,
		LanguageTypeMML:    true,
		LanguageTypeMMM:    true,
		LanguageTypeMMN:    true,
		LanguageTypeMMO:    true,
		LanguageTypeMMP:    true,
		LanguageTypeMMQ:    true,
		LanguageTypeMMR:    true,
		LanguageTypeMMT:    true,
		LanguageTypeMMU:    true,
		LanguageTypeMMV:    true,
		LanguageTypeMMW:    true,
		LanguageTypeMMX:    true,
		LanguageTypeMMY:    true,
		LanguageTypeMMZ:    true,
		LanguageTypeMNA:    true,
		LanguageTypeMNB:    true,
		LanguageTypeMNC:    true,
		LanguageTypeMND:    true,
		LanguageTypeMNE:    true,
		LanguageTypeMNF:    true,
		LanguageTypeMNG:    true,
		LanguageTypeMNH:    true,
		LanguageTypeMNI:    true,
		LanguageTypeMNJ:    true,
		LanguageTypeMNK:    true,
		LanguageTypeMNL:    true,
		LanguageTypeMNM:    true,
		LanguageTypeMNN:    true,
		LanguageTypeMNO:    true,
		LanguageTypeMNP:    true,
		LanguageTypeMNQ:    true,
		LanguageTypeMNR:    true,
		LanguageTypeMNS:    true,
		LanguageTypeMNT:    true,
		LanguageTypeMNU:    true,
		LanguageTypeMNV:    true,
		LanguageTypeMNW:    true,
		LanguageTypeMNX:    true,
		LanguageTypeMNY:    true,
		LanguageTypeMNZ:    true,
		LanguageTypeMOA:    true,
		LanguageTypeMOC:    true,
		LanguageTypeMOD:    true,
		LanguageTypeMOE:    true,
		LanguageTypeMOF:    true,
		LanguageTypeMOG:    true,
		LanguageTypeMOH:    true,
		LanguageTypeMOI:    true,
		LanguageTypeMOJ:    true,
		LanguageTypeMOK:    true,
		LanguageTypeMOM:    true,
		LanguageTypeMOO:    true,
		LanguageTypeMOP:    true,
		LanguageTypeMOQ:    true,
		LanguageTypeMOR:    true,
		LanguageTypeMOS:    true,
		LanguageTypeMOT:    true,
		LanguageTypeMOU:    true,
		LanguageTypeMOV:    true,
		LanguageTypeMOW:    true,
		LanguageTypeMOX:    true,
		LanguageTypeMOY:    true,
		LanguageTypeMOZ:    true,
		LanguageTypeMPA:    true,
		LanguageTypeMPB:    true,
		LanguageTypeMPC:    true,
		LanguageTypeMPD:    true,
		LanguageTypeMPE:    true,
		LanguageTypeMPG:    true,
		LanguageTypeMPH:    true,
		LanguageTypeMPI:    true,
		LanguageTypeMPJ:    true,
		LanguageTypeMPK:    true,
		LanguageTypeMPL:    true,
		LanguageTypeMPM:    true,
		LanguageTypeMPN:    true,
		LanguageTypeMPO:    true,
		LanguageTypeMPP:    true,
		LanguageTypeMPQ:    true,
		LanguageTypeMPR:    true,
		LanguageTypeMPS:    true,
		LanguageTypeMPT:    true,
		LanguageTypeMPU:    true,
		LanguageTypeMPV:    true,
		LanguageTypeMPW:    true,
		LanguageTypeMPX:    true,
		LanguageTypeMPY:    true,
		LanguageTypeMPZ:    true,
		LanguageTypeMQA:    true,
		LanguageTypeMQB:    true,
		LanguageTypeMQC:    true,
		LanguageTypeMQE:    true,
		LanguageTypeMQF:    true,
		LanguageTypeMQG:    true,
		LanguageTypeMQH:    true,
		LanguageTypeMQI:    true,
		LanguageTypeMQJ:    true,
		LanguageTypeMQK:    true,
		LanguageTypeMQL:    true,
		LanguageTypeMQM:    true,
		LanguageTypeMQN:    true,
		LanguageTypeMQO:    true,
		LanguageTypeMQP:    true,
		LanguageTypeMQQ:    true,
		LanguageTypeMQR:    true,
		LanguageTypeMQS:    true,
		LanguageTypeMQT:    true,
		LanguageTypeMQU:    true,
		LanguageTypeMQV:    true,
		LanguageTypeMQW:    true,
		LanguageTypeMQX:    true,
		LanguageTypeMQY:    true,
		LanguageTypeMQZ:    true,
		LanguageTypeMRA:    true,
		LanguageTypeMRB:    true,
		LanguageTypeMRC:    true,
		LanguageTypeMRD:    true,
		LanguageTypeMRE:    true,
		LanguageTypeMRF:    true,
		LanguageTypeMRG:    true,
		LanguageTypeMRH:    true,
		LanguageTypeMRJ:    true,
		LanguageTypeMRK:    true,
		LanguageTypeMRL:    true,
		LanguageTypeMRM:    true,
		LanguageTypeMRN:    true,
		LanguageTypeMRO:    true,
		LanguageTypeMRP:    true,
		LanguageTypeMRQ:    true,
		LanguageTypeMRR:    true,
		LanguageTypeMRS:    true,
		LanguageTypeMRT:    true,
		LanguageTypeMRU:    true,
		LanguageTypeMRV:    true,
		LanguageTypeMRW:    true,
		LanguageTypeMRX:    true,
		LanguageTypeMRY:    true,
		LanguageTypeMRZ:    true,
		LanguageTypeMSB:    true,
		LanguageTypeMSC:    true,
		LanguageTypeMSD:    true,
		LanguageTypeMSE:    true,
		LanguageTypeMSF:    true,
		LanguageTypeMSG:    true,
		LanguageTypeMSH:    true,
		LanguageTypeMSI:    true,
		LanguageTypeMSJ:    true,
		LanguageTypeMSK:    true,
		LanguageTypeMSL:    true,
		LanguageTypeMSM:    true,
		LanguageTypeMSN:    true,
		LanguageTypeMSO:    true,
		LanguageTypeMSP:    true,
		LanguageTypeMSQ:    true,
		LanguageTypeMSR:    true,
		LanguageTypeMSS:    true,
		LanguageTypeMST:    true,
		LanguageTypeMSU:    true,
		LanguageTypeMSV:    true,
		LanguageTypeMSW:    true,
		LanguageTypeMSX:    true,
		LanguageTypeMSY:    true,
		LanguageTypeMSZ:    true,
		LanguageTypeMTA:    true,
		LanguageTypeMTB:    true,
		LanguageTypeMTC:    true,
		LanguageTypeMTD:    true,
		LanguageTypeMTE:    true,
		LanguageTypeMTF:    true,
		LanguageTypeMTG:    true,
		LanguageTypeMTH:    true,
		LanguageTypeMTI:    true,
		LanguageTypeMTJ:    true,
		LanguageTypeMTK:    true,
		LanguageTypeMTL:    true,
		LanguageTypeMTM:    true,
		LanguageTypeMTN:    true,
		LanguageTypeMTO:    true,
		LanguageTypeMTP:    true,
		LanguageTypeMTQ:    true,
		LanguageTypeMTR:    true,
		LanguageTypeMTS:    true,
		LanguageTypeMTT:    true,
		LanguageTypeMTU:    true,
		LanguageTypeMTV:    true,
		LanguageTypeMTW:    true,
		LanguageTypeMTX:    true,
		LanguageTypeMTY:    true,
		LanguageTypeMUA:    true,
		LanguageTypeMUB:    true,
		LanguageTypeMUC:    true,
		LanguageTypeMUD:    true,
		LanguageTypeMUE:    true,
		LanguageTypeMUG:    true,
		LanguageTypeMUH:    true,
		LanguageTypeMUI:    true,
		LanguageTypeMUJ:    true,
		LanguageTypeMUK:    true,
		LanguageTypeMUL:    true,
		LanguageTypeMUM:    true,
		LanguageTypeMUN:    true,
		LanguageTypeMUO:    true,
		LanguageTypeMUP:    true,
		LanguageTypeMUQ:    true,
		LanguageTypeMUR:    true,
		LanguageTypeMUS:    true,
		LanguageTypeMUT:    true,
		LanguageTypeMUU:    true,
		LanguageTypeMUV:    true,
		LanguageTypeMUX:    true,
		LanguageTypeMUY:    true,
		LanguageTypeMUZ:    true,
		LanguageTypeMVA:    true,
		LanguageTypeMVB:    true,
		LanguageTypeMVD:    true,
		LanguageTypeMVE:    true,
		LanguageTypeMVF:    true,
		LanguageTypeMVG:    true,
		LanguageTypeMVH:    true,
		LanguageTypeMVI:    true,
		LanguageTypeMVK:    true,
		LanguageTypeMVL:    true,
		LanguageTypeMVM:    true,
		LanguageTypeMVN:    true,
		LanguageTypeMVO:    true,
		LanguageTypeMVP:    true,
		LanguageTypeMVQ:    true,
		LanguageTypeMVR:    true,
		LanguageTypeMVS:    true,
		LanguageTypeMVT:    true,
		LanguageTypeMVU:    true,
		LanguageTypeMVV:    true,
		LanguageTypeMVW:    true,
		LanguageTypeMVX:    true,
		LanguageTypeMVY:    true,
		LanguageTypeMVZ:    true,
		LanguageTypeMWA:    true,
		LanguageTypeMWB:    true,
		LanguageTypeMWC:    true,
		LanguageTypeMWD:    true,
		LanguageTypeMWE:    true,
		LanguageTypeMWF:    true,
		LanguageTypeMWG:    true,
		LanguageTypeMWH:    true,
		LanguageTypeMWI:    true,
		LanguageTypeMWJ:    true,
		LanguageTypeMWK:    true,
		LanguageTypeMWL:    true,
		LanguageTypeMWM:    true,
		LanguageTypeMWN:    true,
		LanguageTypeMWO:    true,
		LanguageTypeMWP:    true,
		LanguageTypeMWQ:    true,
		LanguageTypeMWR:    true,
		LanguageTypeMWS:    true,
		LanguageTypeMWT:    true,
		LanguageTypeMWU:    true,
		LanguageTypeMWV:    true,
		LanguageTypeMWW:    true,
		LanguageTypeMWX:    true,
		LanguageTypeMWY:    true,
		LanguageTypeMWZ:    true,
		LanguageTypeMXA:    true,
		LanguageTypeMXB:    true,
		LanguageTypeMXC:    true,
		LanguageTypeMXD:    true,
		LanguageTypeMXE:    true,
		LanguageTypeMXF:    true,
		LanguageTypeMXG:    true,
		LanguageTypeMXH:    true,
		LanguageTypeMXI:    true,
		LanguageTypeMXJ:    true,
		LanguageTypeMXK:    true,
		LanguageTypeMXL:    true,
		LanguageTypeMXM:    true,
		LanguageTypeMXN:    true,
		LanguageTypeMXO:    true,
		LanguageTypeMXP:    true,
		LanguageTypeMXQ:    true,
		LanguageTypeMXR:    true,
		LanguageTypeMXS:    true,
		LanguageTypeMXT:    true,
		LanguageTypeMXU:    true,
		LanguageTypeMXV:    true,
		LanguageTypeMXW:    true,
		LanguageTypeMXX:    true,
		LanguageTypeMXY:    true,
		LanguageTypeMXZ:    true,
		LanguageTypeMYB:    true,
		LanguageTypeMYC:    true,
		LanguageTypeMYD:    true,
		LanguageTypeMYE:    true,
		LanguageTypeMYF:    true,
		LanguageTypeMYG:    true,
		LanguageTypeMYH:    true,
		LanguageTypeMYI:    true,
		LanguageTypeMYJ:    true,
		LanguageTypeMYK:    true,
		LanguageTypeMYL:    true,
		LanguageTypeMYM:    true,
		LanguageTypeMYN:    true,
		LanguageTypeMYO:    true,
		LanguageTypeMYP:    true,
		LanguageTypeMYQ:    true,
		LanguageTypeMYR:    true,
		LanguageTypeMYS:    true,
		LanguageTypeMYT:    true,
		LanguageTypeMYU:    true,
		LanguageTypeMYV:    true,
		LanguageTypeMYW:    true,
		LanguageTypeMYX:    true,
		LanguageTypeMYY:    true,
		LanguageTypeMYZ:    true,
		LanguageTypeMZA:    true,
		LanguageTypeMZB:    true,
		LanguageTypeMZC:    true,
		LanguageTypeMZD:    true,
		LanguageTypeMZE:    true,
		LanguageTypeMZG:    true,
		LanguageTypeMZH:    true,
		LanguageTypeMZI:    true,
		LanguageTypeMZJ:    true,
		LanguageTypeMZK:    true,
		LanguageTypeMZL:    true,
		LanguageTypeMZM:    true,
		LanguageTypeMZN:    true,
		LanguageTypeMZO:    true,
		LanguageTypeMZP:    true,
		LanguageTypeMZQ:    true,
		LanguageTypeMZR:    true,
		LanguageTypeMZS:    true,
		LanguageTypeMZT:    true,
		LanguageTypeMZU:    true,
		LanguageTypeMZV:    true,
		LanguageTypeMZW:    true,
		LanguageTypeMZX:    true,
		LanguageTypeMZY:    true,
		LanguageTypeMZZ:    true,
		LanguageTypeNAA:    true,
		LanguageTypeNAB:    true,
		LanguageTypeNAC:    true,
		LanguageTypeNAD:    true,
		LanguageTypeNAE:    true,
		LanguageTypeNAF:    true,
		LanguageTypeNAG:    true,
		LanguageTypeNAH:    true,
		LanguageTypeNAI:    true,
		LanguageTypeNAJ:    true,
		LanguageTypeNAK:    true,
		LanguageTypeNAL:    true,
		LanguageTypeNAM:    true,
		LanguageTypeNAN:    true,
		LanguageTypeNAO:    true,
		LanguageTypeNAP:    true,
		LanguageTypeNAQ:    true,
		LanguageTypeNAR:    true,
		LanguageTypeNAS:    true,
		LanguageTypeNAT:    true,
		LanguageTypeNAW:    true,
		LanguageTypeNAX:    true,
		LanguageTypeNAY:    true,
		LanguageTypeNAZ:    true,
		LanguageTypeNBA:    true,
		LanguageTypeNBB:    true,
		LanguageTypeNBC:    true,
		LanguageTypeNBD:    true,
		LanguageTypeNBE:    true,
		LanguageTypeNBF:    true,
		LanguageTypeNBG:    true,
		LanguageTypeNBH:    true,
		LanguageTypeNBI:    true,
		LanguageTypeNBJ:    true,
		LanguageTypeNBK:    true,
		LanguageTypeNBM:    true,
		LanguageTypeNBN:    true,
		LanguageTypeNBO:    true,
		LanguageTypeNBP:    true,
		LanguageTypeNBQ:    true,
		LanguageTypeNBR:    true,
		LanguageTypeNBS:    true,
		LanguageTypeNBT:    true,
		LanguageTypeNBU:    true,
		LanguageTypeNBV:    true,
		LanguageTypeNBW:    true,
		LanguageTypeNBX:    true,
		LanguageTypeNBY:    true,
		LanguageTypeNCA:    true,
		LanguageTypeNCB:    true,
		LanguageTypeNCC:    true,
		LanguageTypeNCD:    true,
		LanguageTypeNCE:    true,
		LanguageTypeNCF:    true,
		LanguageTypeNCG:    true,
		LanguageTypeNCH:    true,
		LanguageTypeNCI:    true,
		LanguageTypeNCJ:    true,
		LanguageTypeNCK:    true,
		LanguageTypeNCL:    true,
		LanguageTypeNCM:    true,
		LanguageTypeNCN:    true,
		LanguageTypeNCO:    true,
		LanguageTypeNCP:    true,
		LanguageTypeNCR:    true,
		LanguageTypeNCS:    true,
		LanguageTypeNCT:    true,
		LanguageTypeNCU:    true,
		LanguageTypeNCX:    true,
		LanguageTypeNCZ:    true,
		LanguageTypeNDA:    true,
		LanguageTypeNDB:    true,
		LanguageTypeNDC:    true,
		LanguageTypeNDD:    true,
		LanguageTypeNDF:    true,
		LanguageTypeNDG:    true,
		LanguageTypeNDH:    true,
		LanguageTypeNDI:    true,
		LanguageTypeNDJ:    true,
		LanguageTypeNDK:    true,
		LanguageTypeNDL:    true,
		LanguageTypeNDM:    true,
		LanguageTypeNDN:    true,
		LanguageTypeNDP:    true,
		LanguageTypeNDQ:    true,
		LanguageTypeNDR:    true,
		LanguageTypeNDS:    true,
		LanguageTypeNDT:    true,
		LanguageTypeNDU:    true,
		LanguageTypeNDV:    true,
		LanguageTypeNDW:    true,
		LanguageTypeNDX:    true,
		LanguageTypeNDY:    true,
		LanguageTypeNDZ:    true,
		LanguageTypeNEA:    true,
		LanguageTypeNEB:    true,
		LanguageTypeNEC:    true,
		LanguageTypeNED:    true,
		LanguageTypeNEE:    true,
		LanguageTypeNEF:    true,
		LanguageTypeNEG:    true,
		LanguageTypeNEH:    true,
		LanguageTypeNEI:    true,
		LanguageTypeNEJ:    true,
		LanguageTypeNEK:    true,
		LanguageTypeNEM:    true,
		LanguageTypeNEN:    true,
		LanguageTypeNEO:    true,
		LanguageTypeNEQ:    true,
		LanguageTypeNER:    true,
		LanguageTypeNES:    true,
		LanguageTypeNET:    true,
		LanguageTypeNEU:    true,
		LanguageTypeNEV:    true,
		LanguageTypeNEW:    true,
		LanguageTypeNEX:    true,
		LanguageTypeNEY:    true,
		LanguageTypeNEZ:    true,
		LanguageTypeNFA:    true,
		LanguageTypeNFD:    true,
		LanguageTypeNFL:    true,
		LanguageTypeNFR:    true,
		LanguageTypeNFU:    true,
		LanguageTypeNGA:    true,
		LanguageTypeNGB:    true,
		LanguageTypeNGC:    true,
		LanguageTypeNGD:    true,
		LanguageTypeNGE:    true,
		LanguageTypeNGF:    true,
		LanguageTypeNGG:    true,
		LanguageTypeNGH:    true,
		LanguageTypeNGI:    true,
		LanguageTypeNGJ:    true,
		LanguageTypeNGK:    true,
		LanguageTypeNGL:    true,
		LanguageTypeNGM:    true,
		LanguageTypeNGN:    true,
		LanguageTypeNGO:    true,
		LanguageTypeNGP:    true,
		LanguageTypeNGQ:    true,
		LanguageTypeNGR:    true,
		LanguageTypeNGS:    true,
		LanguageTypeNGT:    true,
		LanguageTypeNGU:    true,
		LanguageTypeNGV:    true,
		LanguageTypeNGW:    true,
		LanguageTypeNGX:    true,
		LanguageTypeNGY:    true,
		LanguageTypeNGZ:    true,
		LanguageTypeNHA:    true,
		LanguageTypeNHB:    true,
		LanguageTypeNHC:    true,
		LanguageTypeNHD:    true,
		LanguageTypeNHE:    true,
		LanguageTypeNHF:    true,
		LanguageTypeNHG:    true,
		LanguageTypeNHH:    true,
		LanguageTypeNHI:    true,
		LanguageTypeNHK:    true,
		LanguageTypeNHM:    true,
		LanguageTypeNHN:    true,
		LanguageTypeNHO:    true,
		LanguageTypeNHP:    true,
		LanguageTypeNHQ:    true,
		LanguageTypeNHR:    true,
		LanguageTypeNHT:    true,
		LanguageTypeNHU:    true,
		LanguageTypeNHV:    true,
		LanguageTypeNHW:    true,
		LanguageTypeNHX:    true,
		LanguageTypeNHY:    true,
		LanguageTypeNHZ:    true,
		LanguageTypeNIA:    true,
		LanguageTypeNIB:    true,
		LanguageTypeNIC:    true,
		LanguageTypeNID:    true,
		LanguageTypeNIE:    true,
		LanguageTypeNIF:    true,
		LanguageTypeNIG:    true,
		LanguageTypeNIH:    true,
		LanguageTypeNII:    true,
		LanguageTypeNIJ:    true,
		LanguageTypeNIK:    true,
		LanguageTypeNIL:    true,
		LanguageTypeNIM:    true,
		LanguageTypeNIN:    true,
		LanguageTypeNIO:    true,
		LanguageTypeNIQ:    true,
		LanguageTypeNIR:    true,
		LanguageTypeNIS:    true,
		LanguageTypeNIT:    true,
		LanguageTypeNIU:    true,
		LanguageTypeNIV:    true,
		LanguageTypeNIW:    true,
		LanguageTypeNIX:    true,
		LanguageTypeNIY:    true,
		LanguageTypeNIZ:    true,
		LanguageTypeNJA:    true,
		LanguageTypeNJB:    true,
		LanguageTypeNJD:    true,
		LanguageTypeNJH:    true,
		LanguageTypeNJI:    true,
		LanguageTypeNJJ:    true,
		LanguageTypeNJL:    true,
		LanguageTypeNJM:    true,
		LanguageTypeNJN:    true,
		LanguageTypeNJO:    true,
		LanguageTypeNJR:    true,
		LanguageTypeNJS:    true,
		LanguageTypeNJT:    true,
		LanguageTypeNJU:    true,
		LanguageTypeNJX:    true,
		LanguageTypeNJY:    true,
		LanguageTypeNJZ:    true,
		LanguageTypeNKA:    true,
		LanguageTypeNKB:    true,
		LanguageTypeNKC:    true,
		LanguageTypeNKD:    true,
		LanguageTypeNKE:    true,
		LanguageTypeNKF:    true,
		LanguageTypeNKG:    true,
		LanguageTypeNKH:    true,
		LanguageTypeNKI:    true,
		LanguageTypeNKJ:    true,
		LanguageTypeNKK:    true,
		LanguageTypeNKM:    true,
		LanguageTypeNKN:    true,
		LanguageTypeNKO:    true,
		LanguageTypeNKP:    true,
		LanguageTypeNKQ:    true,
		LanguageTypeNKR:    true,
		LanguageTypeNKS:    true,
		LanguageTypeNKT:    true,
		LanguageTypeNKU:    true,
		LanguageTypeNKV:    true,
		LanguageTypeNKW:    true,
		LanguageTypeNKX:    true,
		LanguageTypeNKZ:    true,
		LanguageTypeNLA:    true,
		LanguageTypeNLC:    true,
		LanguageTypeNLE:    true,
		LanguageTypeNLG:    true,
		LanguageTypeNLI:    true,
		LanguageTypeNLJ:    true,
		LanguageTypeNLK:    true,
		LanguageTypeNLL:    true,
		LanguageTypeNLN:    true,
		LanguageTypeNLO:    true,
		LanguageTypeNLQ:    true,
		LanguageTypeNLR:    true,
		LanguageTypeNLU:    true,
		LanguageTypeNLV:    true,
		LanguageTypeNLW:    true,
		LanguageTypeNLX:    true,
		LanguageTypeNLY:    true,
		LanguageTypeNLZ:    true,
		LanguageTypeNMA:    true,
		LanguageTypeNMB:    true,
		LanguageTypeNMC:    true,
		LanguageTypeNMD:    true,
		LanguageTypeNME:    true,
		LanguageTypeNMF:    true,
		LanguageTypeNMG:    true,
		LanguageTypeNMH:    true,
		LanguageTypeNMI:    true,
		LanguageTypeNMJ:    true,
		LanguageTypeNMK:    true,
		LanguageTypeNML:    true,
		LanguageTypeNMM:    true,
		LanguageTypeNMN:    true,
		LanguageTypeNMO:    true,
		LanguageTypeNMP:    true,
		LanguageTypeNMQ:    true,
		LanguageTypeNMR:    true,
		LanguageTypeNMS:    true,
		LanguageTypeNMT:    true,
		LanguageTypeNMU:    true,
		LanguageTypeNMV:    true,
		LanguageTypeNMW:    true,
		LanguageTypeNMX:    true,
		LanguageTypeNMY:    true,
		LanguageTypeNMZ:    true,
		LanguageTypeNNA:    true,
		LanguageTypeNNB:    true,
		LanguageTypeNNC:    true,
		LanguageTypeNND:    true,
		LanguageTypeNNE:    true,
		LanguageTypeNNF:    true,
		LanguageTypeNNG:    true,
		LanguageTypeNNH:    true,
		LanguageTypeNNI:    true,
		LanguageTypeNNJ:    true,
		LanguageTypeNNK:    true,
		LanguageTypeNNL:    true,
		LanguageTypeNNM:    true,
		LanguageTypeNNN:    true,
		LanguageTypeNNP:    true,
		LanguageTypeNNQ:    true,
		LanguageTypeNNR:    true,
		LanguageTypeNNS:    true,
		LanguageTypeNNT:    true,
		LanguageTypeNNU:    true,
		LanguageTypeNNV:    true,
		LanguageTypeNNW:    true,
		LanguageTypeNNX:    true,
		LanguageTypeNNY:    true,
		LanguageTypeNNZ:    true,
		LanguageTypeNOA:    true,
		LanguageTypeNOC:    true,
		LanguageTypeNOD:    true,
		LanguageTypeNOE:    true,
		LanguageTypeNOF:    true,
		LanguageTypeNOG:    true,
		LanguageTypeNOH:    true,
		LanguageTypeNOI:    true,
		LanguageTypeNOJ:    true,
		LanguageTypeNOK:    true,
		LanguageTypeNOL:    true,
		LanguageTypeNOM:    true,
		LanguageTypeNON:    true,
		LanguageTypeNOO:    true,
		LanguageTypeNOP:    true,
		LanguageTypeNOQ:    true,
		LanguageTypeNOS:    true,
		LanguageTypeNOT:    true,
		LanguageTypeNOU:    true,
		LanguageTypeNOV:    true,
		LanguageTypeNOW:    true,
		LanguageTypeNOY:    true,
		LanguageTypeNOZ:    true,
		LanguageTypeNPA:    true,
		LanguageTypeNPB:    true,
		LanguageTypeNPG:    true,
		LanguageTypeNPH:    true,
		LanguageTypeNPI:    true,
		LanguageTypeNPL:    true,
		LanguageTypeNPN:    true,
		LanguageTypeNPO:    true,
		LanguageTypeNPS:    true,
		LanguageTypeNPU:    true,
		LanguageTypeNPY:    true,
		LanguageTypeNQG:    true,
		LanguageTypeNQK:    true,
		LanguageTypeNQM:    true,
		LanguageTypeNQN:    true,
		LanguageTypeNQO:    true,
		LanguageTypeNQQ:    true,
		LanguageTypeNQY:    true,
		LanguageTypeNRA:    true,
		LanguageTypeNRB:    true,
		LanguageTypeNRC:    true,
		LanguageTypeNRE:    true,
		LanguageTypeNRG:    true,
		LanguageTypeNRI:    true,
		LanguageTypeNRK:    true,
		LanguageTypeNRL:    true,
		LanguageTypeNRM:    true,
		LanguageTypeNRN:    true,
		LanguageTypeNRP:    true,
		LanguageTypeNRR:    true,
		LanguageTypeNRT:    true,
		LanguageTypeNRU:    true,
		LanguageTypeNRX:    true,
		LanguageTypeNRZ:    true,
		LanguageTypeNSA:    true,
		LanguageTypeNSC:    true,
		LanguageTypeNSD:    true,
		LanguageTypeNSE:    true,
		LanguageTypeNSF:    true,
		LanguageTypeNSG:    true,
		LanguageTypeNSH:    true,
		LanguageTypeNSI:    true,
		LanguageTypeNSK:    true,
		LanguageTypeNSL:    true,
		LanguageTypeNSM:    true,
		LanguageTypeNSN:    true,
		LanguageTypeNSO:    true,
		LanguageTypeNSP:    true,
		LanguageTypeNSQ:    true,
		LanguageTypeNSR:    true,
		LanguageTypeNSS:    true,
		LanguageTypeNST:    true,
		LanguageTypeNSU:    true,
		LanguageTypeNSV:    true,
		LanguageTypeNSW:    true,
		LanguageTypeNSX:    true,
		LanguageTypeNSY:    true,
		LanguageTypeNSZ:    true,
		LanguageTypeNTE:    true,
		LanguageTypeNTG:    true,
		LanguageTypeNTI:    true,
		LanguageTypeNTJ:    true,
		LanguageTypeNTK:    true,
		LanguageTypeNTM:    true,
		LanguageTypeNTO:    true,
		LanguageTypeNTP:    true,
		LanguageTypeNTR:    true,
		LanguageTypeNTS:    true,
		LanguageTypeNTU:    true,
		LanguageTypeNTW:    true,
		LanguageTypeNTX:    true,
		LanguageTypeNTY:    true,
		LanguageTypeNTZ:    true,
		LanguageTypeNUA:    true,
		LanguageTypeNUB:    true,
		LanguageTypeNUC:    true,
		LanguageTypeNUD:    true,
		LanguageTypeNUE:    true,
		LanguageTypeNUF:    true,
		LanguageTypeNUG:    true,
		LanguageTypeNUH:    true,
		LanguageTypeNUI:    true,
		LanguageTypeNUJ:    true,
		LanguageTypeNUK:    true,
		LanguageTypeNUL:    true,
		LanguageTypeNUM:    true,
		LanguageTypeNUN:    true,
		LanguageTypeNUO:    true,
		LanguageTypeNUP:    true,
		LanguageTypeNUQ:    true,
		LanguageTypeNUR:    true,
		LanguageTypeNUS:    true,
		LanguageTypeNUT:    true,
		LanguageTypeNUU:    true,
		LanguageTypeNUV:    true,
		LanguageTypeNUW:    true,
		LanguageTypeNUX:    true,
		LanguageTypeNUY:    true,
		LanguageTypeNUZ:    true,
		LanguageTypeNVH:    true,
		LanguageTypeNVM:    true,
		LanguageTypeNVO:    true,
		LanguageTypeNWA:    true,
		LanguageTypeNWB:    true,
		LanguageTypeNWC:    true,
		LanguageTypeNWE:    true,
		LanguageTypeNWG:    true,
		LanguageTypeNWI:    true,
		LanguageTypeNWM:    true,
		LanguageTypeNWO:    true,
		LanguageTypeNWR:    true,
		LanguageTypeNWX:    true,
		LanguageTypeNWY:    true,
		LanguageTypeNXA:    true,
		LanguageTypeNXD:    true,
		LanguageTypeNXE:    true,
		LanguageTypeNXG:    true,
		LanguageTypeNXI:    true,
		LanguageTypeNXK:    true,
		LanguageTypeNXL:    true,
		LanguageTypeNXM:    true,
		LanguageTypeNXN:    true,
		LanguageTypeNXQ:    true,
		LanguageTypeNXR:    true,
		LanguageTypeNXU:    true,
		LanguageTypeNXX:    true,
		LanguageTypeNYB:    true,
		LanguageTypeNYC:    true,
		LanguageTypeNYD:    true,
		LanguageTypeNYE:    true,
		LanguageTypeNYF:    true,
		LanguageTypeNYG:    true,
		LanguageTypeNYH:    true,
		LanguageTypeNYI:    true,
		LanguageTypeNYJ:    true,
		LanguageTypeNYK:    true,
		LanguageTypeNYL:    true,
		LanguageTypeNYM:    true,
		LanguageTypeNYN:    true,
		LanguageTypeNYO:    true,
		LanguageTypeNYP:    true,
		LanguageTypeNYQ:    true,
		LanguageTypeNYR:    true,
		LanguageTypeNYS:    true,
		LanguageTypeNYT:    true,
		LanguageTypeNYU:    true,
		LanguageTypeNYV:    true,
		LanguageTypeNYW:    true,
		LanguageTypeNYX:    true,
		LanguageTypeNYY:    true,
		LanguageTypeNZA:    true,
		LanguageTypeNZB:    true,
		LanguageTypeNZI:    true,
		LanguageTypeNZK:    true,
		LanguageTypeNZM:    true,
		LanguageTypeNZS:    true,
		LanguageTypeNZU:    true,
		LanguageTypeNZY:    true,
		LanguageTypeNZZ:    true,
		LanguageTypeOAA:    true,
		LanguageTypeOAC:    true,
		LanguageTypeOAR:    true,
		LanguageTypeOAV:    true,
		LanguageTypeOBI:    true,
		LanguageTypeOBK:    true,
		LanguageTypeOBL:    true,
		LanguageTypeOBM:    true,
		LanguageTypeOBO:    true,
		LanguageTypeOBR:    true,
		LanguageTypeOBT:    true,
		LanguageTypeOBU:    true,
		LanguageTypeOCA:    true,
		LanguageTypeOCH:    true,
		LanguageTypeOCO:    true,
		LanguageTypeOCU:    true,
		LanguageTypeODA:    true,
		LanguageTypeODK:    true,
		LanguageTypeODT:    true,
		LanguageTypeODU:    true,
		LanguageTypeOFO:    true,
		LanguageTypeOFS:    true,
		LanguageTypeOFU:    true,
		LanguageTypeOGB:    true,
		LanguageTypeOGC:    true,
		LanguageTypeOGE:    true,
		LanguageTypeOGG:    true,
		LanguageTypeOGO:    true,
		LanguageTypeOGU:    true,
		LanguageTypeOHT:    true,
		LanguageTypeOHU:    true,
		LanguageTypeOIA:    true,
		LanguageTypeOIN:    true,
		LanguageTypeOJB:    true,
		LanguageTypeOJC:    true,
		LanguageTypeOJG:    true,
		LanguageTypeOJP:    true,
		LanguageTypeOJS:    true,
		LanguageTypeOJV:    true,
		LanguageTypeOJW:    true,
		LanguageTypeOKA:    true,
		LanguageTypeOKB:    true,
		LanguageTypeOKD:    true,
		LanguageTypeOKE:    true,
		LanguageTypeOKG:    true,
		LanguageTypeOKH:    true,
		LanguageTypeOKI:    true,
		LanguageTypeOKJ:    true,
		LanguageTypeOKK:    true,
		LanguageTypeOKL:    true,
		LanguageTypeOKM:    true,
		LanguageTypeOKN:    true,
		LanguageTypeOKO:    true,
		LanguageTypeOKR:    true,
		LanguageTypeOKS:    true,
		LanguageTypeOKU:    true,
		LanguageTypeOKV:    true,
		LanguageTypeOKX:    true,
		LanguageTypeOLA:    true,
		LanguageTypeOLD:    true,
		LanguageTypeOLE:    true,
		LanguageTypeOLK:    true,
		LanguageTypeOLM:    true,
		LanguageTypeOLO:    true,
		LanguageTypeOLR:    true,
		LanguageTypeOMA:    true,
		LanguageTypeOMB:    true,
		LanguageTypeOMC:    true,
		LanguageTypeOME:    true,
		LanguageTypeOMG:    true,
		LanguageTypeOMI:    true,
		LanguageTypeOMK:    true,
		LanguageTypeOML:    true,
		LanguageTypeOMN:    true,
		LanguageTypeOMO:    true,
		LanguageTypeOMP:    true,
		LanguageTypeOMQ:    true,
		LanguageTypeOMR:    true,
		LanguageTypeOMT:    true,
		LanguageTypeOMU:    true,
		LanguageTypeOMV:    true,
		LanguageTypeOMW:    true,
		LanguageTypeOMX:    true,
		LanguageTypeONA:    true,
		LanguageTypeONB:    true,
		LanguageTypeONE:    true,
		LanguageTypeONG:    true,
		LanguageTypeONI:    true,
		LanguageTypeONJ:    true,
		LanguageTypeONK:    true,
		LanguageTypeONN:    true,
		LanguageTypeONO:    true,
		LanguageTypeONP:    true,
		LanguageTypeONR:    true,
		LanguageTypeONS:    true,
		LanguageTypeONT:    true,
		LanguageTypeONU:    true,
		LanguageTypeONW:    true,
		LanguageTypeONX:    true,
		LanguageTypeOOD:    true,
		LanguageTypeOOG:    true,
		LanguageTypeOON:    true,
		LanguageTypeOOR:    true,
		LanguageTypeOOS:    true,
		LanguageTypeOPA:    true,
		LanguageTypeOPK:    true,
		LanguageTypeOPM:    true,
		LanguageTypeOPO:    true,
		LanguageTypeOPT:    true,
		LanguageTypeOPY:    true,
		LanguageTypeORA:    true,
		LanguageTypeORC:    true,
		LanguageTypeORE:    true,
		LanguageTypeORG:    true,
		LanguageTypeORH:    true,
		LanguageTypeORN:    true,
		LanguageTypeORO:    true,
		LanguageTypeORR:    true,
		LanguageTypeORS:    true,
		LanguageTypeORT:    true,
		LanguageTypeORU:    true,
		LanguageTypeORV:    true,
		LanguageTypeORW:    true,
		LanguageTypeORX:    true,
		LanguageTypeORY:    true,
		LanguageTypeORZ:    true,
		LanguageTypeOSA:    true,
		LanguageTypeOSC:    true,
		LanguageTypeOSI:    true,
		LanguageTypeOSO:    true,
		LanguageTypeOSP:    true,
		LanguageTypeOST:    true,
		LanguageTypeOSU:    true,
		LanguageTypeOSX:    true,
		LanguageTypeOTA:    true,
		LanguageTypeOTB:    true,
		LanguageTypeOTD:    true,
		LanguageTypeOTE:    true,
		LanguageTypeOTI:    true,
		LanguageTypeOTK:    true,
		LanguageTypeOTL:    true,
		LanguageTypeOTM:    true,
		LanguageTypeOTN:    true,
		LanguageTypeOTO:    true,
		LanguageTypeOTQ:    true,
		LanguageTypeOTR:    true,
		LanguageTypeOTS:    true,
		LanguageTypeOTT:    true,
		LanguageTypeOTU:    true,
		LanguageTypeOTW:    true,
		LanguageTypeOTX:    true,
		LanguageTypeOTY:    true,
		LanguageTypeOTZ:    true,
		LanguageTypeOUA:    true,
		LanguageTypeOUB:    true,
		LanguageTypeOUE:    true,
		LanguageTypeOUI:    true,
		LanguageTypeOUM:    true,
		LanguageTypeOUN:    true,
		LanguageTypeOWI:    true,
		LanguageTypeOWL:    true,
		LanguageTypeOYB:    true,
		LanguageTypeOYD:    true,
		LanguageTypeOYM:    true,
		LanguageTypeOYY:    true,
		LanguageTypeOZM:    true,
		LanguageTypePAA:    true,
		LanguageTypePAB:    true,
		LanguageTypePAC:    true,
		LanguageTypePAD:    true,
		LanguageTypePAE:    true,
		LanguageTypePAF:    true,
		LanguageTypePAG:    true,
		LanguageTypePAH:    true,
		LanguageTypePAI:    true,
		LanguageTypePAK:    true,
		LanguageTypePAL:    true,
		LanguageTypePAM:    true,
		LanguageTypePAO:    true,
		LanguageTypePAP:    true,
		LanguageTypePAQ:    true,
		LanguageTypePAR:    true,
		LanguageTypePAS:    true,
		LanguageTypePAT:    true,
		LanguageTypePAU:    true,
		LanguageTypePAV:    true,
		LanguageTypePAW:    true,
		LanguageTypePAX:    true,
		LanguageTypePAY:    true,
		LanguageTypePAZ:    true,
		LanguageTypePBB:    true,
		LanguageTypePBC:    true,
		LanguageTypePBE:    true,
		LanguageTypePBF:    true,
		LanguageTypePBG:    true,
		LanguageTypePBH:    true,
		LanguageTypePBI:    true,
		LanguageTypePBL:    true,
		LanguageTypePBN:    true,
		LanguageTypePBO:    true,
		LanguageTypePBP:    true,
		LanguageTypePBR:    true,
		LanguageTypePBS:    true,
		LanguageTypePBT:    true,
		LanguageTypePBU:    true,
		LanguageTypePBV:    true,
		LanguageTypePBY:    true,
		LanguageTypePBZ:    true,
		LanguageTypePCA:    true,
		LanguageTypePCB:    true,
		LanguageTypePCC:    true,
		LanguageTypePCD:    true,
		LanguageTypePCE:    true,
		LanguageTypePCF:    true,
		LanguageTypePCG:    true,
		LanguageTypePCH:    true,
		LanguageTypePCI:    true,
		LanguageTypePCJ:    true,
		LanguageTypePCK:    true,
		LanguageTypePCL:    true,
		LanguageTypePCM:    true,
		LanguageTypePCN:    true,
		LanguageTypePCP:    true,
		LanguageTypePCR:    true,
		LanguageTypePCW:    true,
		LanguageTypePDA:    true,
		LanguageTypePDC:    true,
		LanguageTypePDI:    true,
		LanguageTypePDN:    true,
		LanguageTypePDO:    true,
		LanguageTypePDT:    true,
		LanguageTypePDU:    true,
		LanguageTypePEA:    true,
		LanguageTypePEB:    true,
		LanguageTypePED:    true,
		LanguageTypePEE:    true,
		LanguageTypePEF:    true,
		LanguageTypePEG:    true,
		LanguageTypePEH:    true,
		LanguageTypePEI:    true,
		LanguageTypePEJ:    true,
		LanguageTypePEK:    true,
		LanguageTypePEL:    true,
		LanguageTypePEM:    true,
		LanguageTypePEO:    true,
		LanguageTypePEP:    true,
		LanguageTypePEQ:    true,
		LanguageTypePES:    true,
		LanguageTypePEV:    true,
		LanguageTypePEX:    true,
		LanguageTypePEY:    true,
		LanguageTypePEZ:    true,
		LanguageTypePFA:    true,
		LanguageTypePFE:    true,
		LanguageTypePFL:    true,
		LanguageTypePGA:    true,
		LanguageTypePGG:    true,
		LanguageTypePGI:    true,
		LanguageTypePGK:    true,
		LanguageTypePGL:    true,
		LanguageTypePGN:    true,
		LanguageTypePGS:    true,
		LanguageTypePGU:    true,
		LanguageTypePGY:    true,
		LanguageTypePHA:    true,
		LanguageTypePHD:    true,
		LanguageTypePHG:    true,
		LanguageTypePHH:    true,
		LanguageTypePHI:    true,
		LanguageTypePHK:    true,
		LanguageTypePHL:    true,
		LanguageTypePHM:    true,
		LanguageTypePHN:    true,
		LanguageTypePHO:    true,
		LanguageTypePHQ:    true,
		LanguageTypePHR:    true,
		LanguageTypePHT:    true,
		LanguageTypePHU:    true,
		LanguageTypePHV:    true,
		LanguageTypePHW:    true,
		LanguageTypePIA:    true,
		LanguageTypePIB:    true,
		LanguageTypePIC:    true,
		LanguageTypePID:    true,
		LanguageTypePIE:    true,
		LanguageTypePIF:    true,
		LanguageTypePIG:    true,
		LanguageTypePIH:    true,
		LanguageTypePII:    true,
		LanguageTypePIJ:    true,
		LanguageTypePIL:    true,
		LanguageTypePIM:    true,
		LanguageTypePIN:    true,
		LanguageTypePIO:    true,
		LanguageTypePIP:    true,
		LanguageTypePIR:    true,
		LanguageTypePIS:    true,
		LanguageTypePIT:    true,
		LanguageTypePIU:    true,
		LanguageTypePIV:    true,
		LanguageTypePIW:    true,
		LanguageTypePIX:    true,
		LanguageTypePIY:    true,
		LanguageTypePIZ:    true,
		LanguageTypePJT:    true,
		LanguageTypePKA:    true,
		LanguageTypePKB:    true,
		LanguageTypePKC:    true,
		LanguageTypePKG:    true,
		LanguageTypePKH:    true,
		LanguageTypePKN:    true,
		LanguageTypePKO:    true,
		LanguageTypePKP:    true,
		LanguageTypePKR:    true,
		LanguageTypePKS:    true,
		LanguageTypePKT:    true,
		LanguageTypePKU:    true,
		LanguageTypePLA:    true,
		LanguageTypePLB:    true,
		LanguageTypePLC:    true,
		LanguageTypePLD:    true,
		LanguageTypePLE:    true,
		LanguageTypePLF:    true,
		LanguageTypePLG:    true,
		LanguageTypePLH:    true,
		LanguageTypePLJ:    true,
		LanguageTypePLK:    true,
		LanguageTypePLL:    true,
		LanguageTypePLN:    true,
		LanguageTypePLO:    true,
		LanguageTypePLP:    true,
		LanguageTypePLQ:    true,
		LanguageTypePLR:    true,
		LanguageTypePLS:    true,
		LanguageTypePLT:    true,
		LanguageTypePLU:    true,
		LanguageTypePLV:    true,
		LanguageTypePLW:    true,
		LanguageTypePLY:    true,
		LanguageTypePLZ:    true,
		LanguageTypePMA:    true,
		LanguageTypePMB:    true,
		LanguageTypePMC:    true,
		LanguageTypePMD:    true,
		LanguageTypePME:    true,
		LanguageTypePMF:    true,
		LanguageTypePMH:    true,
		LanguageTypePMI:    true,
		LanguageTypePMJ:    true,
		LanguageTypePMK:    true,
		LanguageTypePML:    true,
		LanguageTypePMM:    true,
		LanguageTypePMN:    true,
		LanguageTypePMO:    true,
		LanguageTypePMQ:    true,
		LanguageTypePMR:    true,
		LanguageTypePMS:    true,
		LanguageTypePMT:    true,
		LanguageTypePMU:    true,
		LanguageTypePMW:    true,
		LanguageTypePMX:    true,
		LanguageTypePMY:    true,
		LanguageTypePMZ:    true,
		LanguageTypePNA:    true,
		LanguageTypePNB:    true,
		LanguageTypePNC:    true,
		LanguageTypePNE:    true,
		LanguageTypePNG:    true,
		LanguageTypePNH:    true,
		LanguageTypePNI:    true,
		LanguageTypePNJ:    true,
		LanguageTypePNK:    true,
		LanguageTypePNL:    true,
		LanguageTypePNM:    true,
		LanguageTypePNN:    true,
		LanguageTypePNO:    true,
		LanguageTypePNP:    true,
		LanguageTypePNQ:    true,
		LanguageTypePNR:    true,
		LanguageTypePNS:    true,
		LanguageTypePNT:    true,
		LanguageTypePNU:    true,
		LanguageTypePNV:    true,
		LanguageTypePNW:    true,
		LanguageTypePNX:    true,
		LanguageTypePNY:    true,
		LanguageTypePNZ:    true,
		LanguageTypePOC:    true,
		LanguageTypePOD:    true,
		LanguageTypePOE:    true,
		LanguageTypePOF:    true,
		LanguageTypePOG:    true,
		LanguageTypePOH:    true,
		LanguageTypePOI:    true,
		LanguageTypePOK:    true,
		LanguageTypePOM:    true,
		LanguageTypePON:    true,
		LanguageTypePOO:    true,
		LanguageTypePOP:    true,
		LanguageTypePOQ:    true,
		LanguageTypePOS:    true,
		LanguageTypePOT:    true,
		LanguageTypePOV:    true,
		LanguageTypePOW:    true,
		LanguageTypePOX:    true,
		LanguageTypePOY:    true,
		LanguageTypePOZ:    true,
		LanguageTypePPA:    true,
		LanguageTypePPE:    true,
		LanguageTypePPI:    true,
		LanguageTypePPK:    true,
		LanguageTypePPL:    true,
		LanguageTypePPM:    true,
		LanguageTypePPN:    true,
		LanguageTypePPO:    true,
		LanguageTypePPP:    true,
		LanguageTypePPQ:    true,
		LanguageTypePPR:    true,
		LanguageTypePPS:    true,
		LanguageTypePPT:    true,
		LanguageTypePPU:    true,
		LanguageTypePQA:    true,
		LanguageTypePQE:    true,
		LanguageTypePQM:    true,
		LanguageTypePQW:    true,
		LanguageTypePRA:    true,
		LanguageTypePRB:    true,
		LanguageTypePRC:    true,
		LanguageTypePRD:    true,
		LanguageTypePRE:    true,
		LanguageTypePRF:    true,
		LanguageTypePRG:    true,
		LanguageTypePRH:    true,
		LanguageTypePRI:    true,
		LanguageTypePRK:    true,
		LanguageTypePRL:    true,
		LanguageTypePRM:    true,
		LanguageTypePRN:    true,
		LanguageTypePRO:    true,
		LanguageTypePRP:    true,
		LanguageTypePRQ:    true,
		LanguageTypePRR:    true,
		LanguageTypePRS:    true,
		LanguageTypePRT:    true,
		LanguageTypePRU:    true,
		LanguageTypePRW:    true,
		LanguageTypePRX:    true,
		LanguageTypePRY:    true,
		LanguageTypePRZ:    true,
		LanguageTypePSA:    true,
		LanguageTypePSC:    true,
		LanguageTypePSD:    true,
		LanguageTypePSE:    true,
		LanguageTypePSG:    true,
		LanguageTypePSH:    true,
		LanguageTypePSI:    true,
		LanguageTypePSL:    true,
		LanguageTypePSM:    true,
		LanguageTypePSN:    true,
		LanguageTypePSO:    true,
		LanguageTypePSP:    true,
		LanguageTypePSQ:    true,
		LanguageTypePSR:    true,
		LanguageTypePSS:    true,
		LanguageTypePST:    true,
		LanguageTypePSU:    true,
		LanguageTypePSW:    true,
		LanguageTypePSY:    true,
		LanguageTypePTA:    true,
		LanguageTypePTH:    true,
		LanguageTypePTI:    true,
		LanguageTypePTN:    true,
		LanguageTypePTO:    true,
		LanguageTypePTP:    true,
		LanguageTypePTR:    true,
		LanguageTypePTT:    true,
		LanguageTypePTU:    true,
		LanguageTypePTV:    true,
		LanguageTypePTW:    true,
		LanguageTypePTY:    true,
		LanguageTypePUA:    true,
		LanguageTypePUB:    true,
		LanguageTypePUC:    true,
		LanguageTypePUD:    true,
		LanguageTypePUE:    true,
		LanguageTypePUF:    true,
		LanguageTypePUG:    true,
		LanguageTypePUI:    true,
		LanguageTypePUJ:    true,
		LanguageTypePUK:    true,
		LanguageTypePUM:    true,
		LanguageTypePUO:    true,
		LanguageTypePUP:    true,
		LanguageTypePUQ:    true,
		LanguageTypePUR:    true,
		LanguageTypePUT:    true,
		LanguageTypePUU:    true,
		LanguageTypePUW:    true,
		LanguageTypePUX:    true,
		LanguageTypePUY:    true,
		LanguageTypePUZ:    true,
		LanguageTypePWA:    true,
		LanguageTypePWB:    true,
		LanguageTypePWG:    true,
		LanguageTypePWI:    true,
		LanguageTypePWM:    true,
		LanguageTypePWN:    true,
		LanguageTypePWO:    true,
		LanguageTypePWR:    true,
		LanguageTypePWW:    true,
		LanguageTypePXM:    true,
		LanguageTypePYE:    true,
		LanguageTypePYM:    true,
		LanguageTypePYN:    true,
		LanguageTypePYS:    true,
		LanguageTypePYU:    true,
		LanguageTypePYX:    true,
		LanguageTypePYY:    true,
		LanguageTypePZN:    true,
		LanguageTypeQAAQTZ: true,
		LanguageTypeQUA:    true,
		LanguageTypeQUB:    true,
		LanguageTypeQUC:    true,
		LanguageTypeQUD:    true,
		LanguageTypeQUF:    true,
		LanguageTypeQUG:    true,
		LanguageTypeQUH:    true,
		LanguageTypeQUI:    true,
		LanguageTypeQUK:    true,
		LanguageTypeQUL:    true,
		LanguageTypeQUM:    true,
		LanguageTypeQUN:    true,
		LanguageTypeQUP:    true,
		LanguageTypeQUQ:    true,
		LanguageTypeQUR:    true,
		LanguageTypeQUS:    true,
		LanguageTypeQUV:    true,
		LanguageTypeQUW:    true,
		LanguageTypeQUX:    true,
		LanguageTypeQUY:    true,
		LanguageTypeQUZ:    true,
		LanguageTypeQVA:    true,
		LanguageTypeQVC:    true,
		LanguageTypeQVE:    true,
		LanguageTypeQVH:    true,
		LanguageTypeQVI:    true,
		LanguageTypeQVJ:    true,
		LanguageTypeQVL:    true,
		LanguageTypeQVM:    true,
		LanguageTypeQVN:    true,
		LanguageTypeQVO:    true,
		LanguageTypeQVP:    true,
		LanguageTypeQVS:    true,
		LanguageTypeQVW:    true,
		LanguageTypeQVY:    true,
		LanguageTypeQVZ:    true,
		LanguageTypeQWA:    true,
		LanguageTypeQWC:    true,
		LanguageTypeQWE:    true,
		LanguageTypeQWH:    true,
		LanguageTypeQWM:    true,
		LanguageTypeQWS:    true,
		LanguageTypeQWT:    true,
		LanguageTypeQXA:    true,
		LanguageTypeQXC:    true,
		LanguageTypeQXH:    true,
		LanguageTypeQXL:    true,
		LanguageTypeQXN:    true,
		LanguageTypeQXO:    true,
		LanguageTypeQXP:    true,
		LanguageTypeQXQ:    true,
		LanguageTypeQXR:    true,
		LanguageTypeQXS:    true,
		LanguageTypeQXT:    true,
		LanguageTypeQXU:    true,
		LanguageTypeQXW:    true,
		LanguageTypeQYA:    true,
		LanguageTypeQYP:    true,
		LanguageTypeRAA:    true,
		LanguageTypeRAB:    true,
		LanguageTypeRAC:    true,
		LanguageTypeRAD:    true,
		LanguageTypeRAF:    true,
		LanguageTypeRAG:    true,
		LanguageTypeRAH:    true,
		LanguageTypeRAI:    true,
		LanguageTypeRAJ:    true,
		LanguageTypeRAK:    true,
		LanguageTypeRAL:    true,
		LanguageTypeRAM:    true,
		LanguageTypeRAN:    true,
		LanguageTypeRAO:    true,
		LanguageTypeRAP:    true,
		LanguageTypeRAQ:    true,
		LanguageTypeRAR:    true,
		LanguageTypeRAS:    true,
		LanguageTypeRAT:    true,
		LanguageTypeRAU:    true,
		LanguageTypeRAV:    true,
		LanguageTypeRAW:    true,
		LanguageTypeRAX:    true,
		LanguageTypeRAY:    true,
		LanguageTypeRAZ:    true,
		LanguageTypeRBB:    true,
		LanguageTypeRBK:    true,
		LanguageTypeRBL:    true,
		LanguageTypeRBP:    true,
		LanguageTypeRCF:    true,
		LanguageTypeRDB:    true,
		LanguageTypeREA:    true,
		LanguageTypeREB:    true,
		LanguageTypeREE:    true,
		LanguageTypeREG:    true,
		LanguageTypeREI:    true,
		LanguageTypeREJ:    true,
		LanguageTypeREL:    true,
		LanguageTypeREM:    true,
		LanguageTypeREN:    true,
		LanguageTypeRER:    true,
		LanguageTypeRES:    true,
		LanguageTypeRET:    true,
		LanguageTypeREY:    true,
		LanguageTypeRGA:    true,
		LanguageTypeRGE:    true,
		LanguageTypeRGK:    true,
		LanguageTypeRGN:    true,
		LanguageTypeRGR:    true,
		LanguageTypeRGS:    true,
		LanguageTypeRGU:    true,
		LanguageTypeRHG:    true,
		LanguageTypeRHP:    true,
		LanguageTypeRIA:    true,
		LanguageTypeRIE:    true,
		LanguageTypeRIF:    true,
		LanguageTypeRIL:    true,
		LanguageTypeRIM:    true,
		LanguageTypeRIN:    true,
		LanguageTypeRIR:    true,
		LanguageTypeRIT:    true,
		LanguageTypeRIU:    true,
		LanguageTypeRJG:    true,
		LanguageTypeRJI:    true,
		LanguageTypeRJS:    true,
		LanguageTypeRKA:    true,
		LanguageTypeRKB:    true,
		LanguageTypeRKH:    true,
		LanguageTypeRKI:    true,
		LanguageTypeRKM:    true,
		LanguageTypeRKT:    true,
		LanguageTypeRKW:    true,
		LanguageTypeRMA:    true,
		LanguageTypeRMB:    true,
		LanguageTypeRMC:    true,
		LanguageTypeRMD:    true,
		LanguageTypeRME:    true,
		LanguageTypeRMF:    true,
		LanguageTypeRMG:    true,
		LanguageTypeRMH:    true,
		LanguageTypeRMI:    true,
		LanguageTypeRMK:    true,
		LanguageTypeRML:    true,
		LanguageTypeRMM:    true,
		LanguageTypeRMN:    true,
		LanguageTypeRMO:    true,
		LanguageTypeRMP:    true,
		LanguageTypeRMQ:    true,
		LanguageTypeRMR:    true,
		LanguageTypeRMS:    true,
		LanguageTypeRMT:    true,
		LanguageTypeRMU:    true,
		LanguageTypeRMV:    true,
		LanguageTypeRMW:    true,
		LanguageTypeRMX:    true,
		LanguageTypeRMY:    true,
		LanguageTypeRMZ:    true,
		LanguageTypeRNA:    true,
		LanguageTypeRND:    true,
		LanguageTypeRNG:    true,
		LanguageTypeRNL:    true,
		LanguageTypeRNN:    true,
		LanguageTypeRNP:    true,
		LanguageTypeRNR:    true,
		LanguageTypeRNW:    true,
		LanguageTypeROA:    true,
		LanguageTypeROB:    true,
		LanguageTypeROC:    true,
		LanguageTypeROD:    true,
		LanguageTypeROE:    true,
		LanguageTypeROF:    true,
		LanguageTypeROG:    true,
		LanguageTypeROL:    true,
		LanguageTypeROM:    true,
		LanguageTypeROO:    true,
		LanguageTypeROP:    true,
		LanguageTypeROR:    true,
		LanguageTypeROU:    true,
		LanguageTypeROW:    true,
		LanguageTypeRPN:    true,
		LanguageTypeRPT:    true,
		LanguageTypeRRI:    true,
		LanguageTypeRRO:    true,
		LanguageTypeRRT:    true,
		LanguageTypeRSB:    true,
		LanguageTypeRSI:    true,
		LanguageTypeRSL:    true,
		LanguageTypeRTC:    true,
		LanguageTypeRTH:    true,
		LanguageTypeRTM:    true,
		LanguageTypeRTW:    true,
		LanguageTypeRUB:    true,
		LanguageTypeRUC:    true,
		LanguageTypeRUE:    true,
		LanguageTypeRUF:    true,
		LanguageTypeRUG:    true,
		LanguageTypeRUH:    true,
		LanguageTypeRUI:    true,
		LanguageTypeRUK:    true,
		LanguageTypeRUO:    true,
		LanguageTypeRUP:    true,
		LanguageTypeRUQ:    true,
		LanguageTypeRUT:    true,
		LanguageTypeRUU:    true,
		LanguageTypeRUY:    true,
		LanguageTypeRUZ:    true,
		LanguageTypeRWA:    true,
		LanguageTypeRWK:    true,
		LanguageTypeRWM:    true,
		LanguageTypeRWO:    true,
		LanguageTypeRWR:    true,
		LanguageTypeRXD:    true,
		LanguageTypeRXW:    true,
		LanguageTypeRYN:    true,
		LanguageTypeRYS:    true,
		LanguageTypeRYU:    true,
		LanguageTypeSAA:    true,
		LanguageTypeSAB:    true,
		LanguageTypeSAC:    true,
		LanguageTypeSAD:    true,
		LanguageTypeSAE:    true,
		LanguageTypeSAF:    true,
		LanguageTypeSAH:    true,
		LanguageTypeSAI:    true,
		LanguageTypeSAJ:    true,
		LanguageTypeSAK:    true,
		LanguageTypeSAL:    true,
		LanguageTypeSAM:    true,
		LanguageTypeSAO:    true,
		LanguageTypeSAP:    true,
		LanguageTypeSAQ:    true,
		LanguageTypeSAR:    true,
		LanguageTypeSAS:    true,
		LanguageTypeSAT:    true,
		LanguageTypeSAU:    true,
		LanguageTypeSAV:    true,
		LanguageTypeSAW:    true,
		LanguageTypeSAX:    true,
		LanguageTypeSAY:    true,
		LanguageTypeSAZ:    true,
		LanguageTypeSBA:    true,
		LanguageTypeSBB:    true,
		LanguageTypeSBC:    true,
		LanguageTypeSBD:    true,
		LanguageTypeSBE:    true,
		LanguageTypeSBF:    true,
		LanguageTypeSBG:    true,
		LanguageTypeSBH:    true,
		LanguageTypeSBI:    true,
		LanguageTypeSBJ:    true,
		LanguageTypeSBK:    true,
		LanguageTypeSBL:    true,
		LanguageTypeSBM:    true,
		LanguageTypeSBN:    true,
		LanguageTypeSBO:    true,
		LanguageTypeSBP:    true,
		LanguageTypeSBQ:    true,
		LanguageTypeSBR:    true,
		LanguageTypeSBS:    true,
		LanguageTypeSBT:    true,
		LanguageTypeSBU:    true,
		LanguageTypeSBV:    true,
		LanguageTypeSBW:    true,
		LanguageTypeSBX:    true,
		LanguageTypeSBY:    true,
		LanguageTypeSBZ:    true,
		LanguageTypeSCA:    true,
		LanguageTypeSCB:    true,
		LanguageTypeSCE:    true,
		LanguageTypeSCF:    true,
		LanguageTypeSCG:    true,
		LanguageTypeSCH:    true,
		LanguageTypeSCI:    true,
		LanguageTypeSCK:    true,
		LanguageTypeSCL:    true,
		LanguageTypeSCN:    true,
		LanguageTypeSCO:    true,
		LanguageTypeSCP:    true,
		LanguageTypeSCQ:    true,
		LanguageTypeSCS:    true,
		LanguageTypeSCU:    true,
		LanguageTypeSCV:    true,
		LanguageTypeSCW:    true,
		LanguageTypeSCX:    true,
		LanguageTypeSDA:    true,
		LanguageTypeSDB:    true,
		LanguageTypeSDC:    true,
		LanguageTypeSDE:    true,
		LanguageTypeSDF:    true,
		LanguageTypeSDG:    true,
		LanguageTypeSDH:    true,
		LanguageTypeSDJ:    true,
		LanguageTypeSDK:    true,
		LanguageTypeSDL:    true,
		LanguageTypeSDM:    true,
		LanguageTypeSDN:    true,
		LanguageTypeSDO:    true,
		LanguageTypeSDP:    true,
		LanguageTypeSDR:    true,
		LanguageTypeSDS:    true,
		LanguageTypeSDT:    true,
		LanguageTypeSDU:    true,
		LanguageTypeSDV:    true,
		LanguageTypeSDX:    true,
		LanguageTypeSDZ:    true,
		LanguageTypeSEA:    true,
		LanguageTypeSEB:    true,
		LanguageTypeSEC:    true,
		LanguageTypeSED:    true,
		LanguageTypeSEE:    true,
		LanguageTypeSEF:    true,
		LanguageTypeSEG:    true,
		LanguageTypeSEH:    true,
		LanguageTypeSEI:    true,
		LanguageTypeSEJ:    true,
		LanguageTypeSEK:    true,
		LanguageTypeSEL:    true,
		LanguageTypeSEM:    true,
		LanguageTypeSEN:    true,
		LanguageTypeSEO:    true,
		LanguageTypeSEP:    true,
		LanguageTypeSEQ:    true,
		LanguageTypeSER:    true,
		LanguageTypeSES:    true,
		LanguageTypeSET:    true,
		LanguageTypeSEU:    true,
		LanguageTypeSEV:    true,
		LanguageTypeSEW:    true,
		LanguageTypeSEY:    true,
		LanguageTypeSEZ:    true,
		LanguageTypeSFB:    true,
		LanguageTypeSFE:    true,
		LanguageTypeSFM:    true,
		LanguageTypeSFS:    true,
		LanguageTypeSFW:    true,
		LanguageTypeSGA:    true,
		LanguageTypeSGB:    true,
		LanguageTypeSGC:    true,
		LanguageTypeSGD:    true,
		LanguageTypeSGE:    true,
		LanguageTypeSGG:    true,
		LanguageTypeSGH:    true,
		LanguageTypeSGI:    true,
		LanguageTypeSGJ:    true,
		LanguageTypeSGK:    true,
		LanguageTypeSGL:    true,
		LanguageTypeSGM:    true,
		LanguageTypeSGN:    true,
		LanguageTypeSGO:    true,
		LanguageTypeSGP:    true,
		LanguageTypeSGR:    true,
		LanguageTypeSGS:    true,
		LanguageTypeSGT:    true,
		LanguageTypeSGU:    true,
		LanguageTypeSGW:    true,
		LanguageTypeSGX:    true,
		LanguageTypeSGY:    true,
		LanguageTypeSGZ:    true,
		LanguageTypeSHA:    true,
		LanguageTypeSHB:    true,
		LanguageTypeSHC:    true,
		LanguageTypeSHD:    true,
		LanguageTypeSHE:    true,
		LanguageTypeSHG:    true,
		LanguageTypeSHH:    true,
		LanguageTypeSHI:    true,
		LanguageTypeSHJ:    true,
		LanguageTypeSHK:    true,
		LanguageTypeSHL:    true,
		LanguageTypeSHM:    true,
		LanguageTypeSHN:    true,
		LanguageTypeSHO:    true,
		LanguageTypeSHP:    true,
		LanguageTypeSHQ:    true,
		LanguageTypeSHR:    true,
		LanguageTypeSHS:    true,
		LanguageTypeSHT:    true,
		LanguageTypeSHU:    true,
		LanguageTypeSHV:    true,
		LanguageTypeSHW:    true,
		LanguageTypeSHX:    true,
		LanguageTypeSHY:    true,
		LanguageTypeSHZ:    true,
		LanguageTypeSIA:    true,
		LanguageTypeSIB:    true,
		LanguageTypeSID:    true,
		LanguageTypeSIE:    true,
		LanguageTypeSIF:    true,
		LanguageTypeSIG:    true,
		LanguageTypeSIH:    true,
		LanguageTypeSII:    true,
		LanguageTypeSIJ:    true,
		LanguageTypeSIK:    true,
		LanguageTypeSIL:    true,
		LanguageTypeSIM:    true,
		LanguageTypeSIO:    true,
		LanguageTypeSIP:    true,
		LanguageTypeSIQ:    true,
		LanguageTypeSIR:    true,
		LanguageTypeSIS:    true,
		LanguageTypeSIT:    true,
		LanguageTypeSIU:    true,
		LanguageTypeSIV:    true,
		LanguageTypeSIW:    true,
		LanguageTypeSIX:    true,
		LanguageTypeSIY:    true,
		LanguageTypeSIZ:    true,
		LanguageTypeSJA:    true,
		LanguageTypeSJB:    true,
		LanguageTypeSJD:    true,
		LanguageTypeSJE:    true,
		LanguageTypeSJG:    true,
		LanguageTypeSJK:    true,
		LanguageTypeSJL:    true,
		LanguageTypeSJM:    true,
		LanguageTypeSJN:    true,
		LanguageTypeSJO:    true,
		LanguageTypeSJP:    true,
		LanguageTypeSJR:    true,
		LanguageTypeSJS:    true,
		LanguageTypeSJT:    true,
		LanguageTypeSJU:    true,
		LanguageTypeSJW:    true,
		LanguageTypeSKA:    true,
		LanguageTypeSKB:    true,
		LanguageTypeSKC:    true,
		LanguageTypeSKD:    true,
		LanguageTypeSKE:    true,
		LanguageTypeSKF:    true,
		LanguageTypeSKG:    true,
		LanguageTypeSKH:    true,
		LanguageTypeSKI:    true,
		LanguageTypeSKJ:    true,
		LanguageTypeSKK:    true,
		LanguageTypeSKM:    true,
		LanguageTypeSKN:    true,
		LanguageTypeSKO:    true,
		LanguageTypeSKP:    true,
		LanguageTypeSKQ:    true,
		LanguageTypeSKR:    true,
		LanguageTypeSKS:    true,
		LanguageTypeSKT:    true,
		LanguageTypeSKU:    true,
		LanguageTypeSKV:    true,
		LanguageTypeSKW:    true,
		LanguageTypeSKX:    true,
		LanguageTypeSKY:    true,
		LanguageTypeSKZ:    true,
		LanguageTypeSLA:    true,
		LanguageTypeSLC:    true,
		LanguageTypeSLD:    true,
		LanguageTypeSLE:    true,
		LanguageTypeSLF:    true,
		LanguageTypeSLG:    true,
		LanguageTypeSLH:    true,
		LanguageTypeSLI:    true,
		LanguageTypeSLJ:    true,
		LanguageTypeSLL:    true,
		LanguageTypeSLM:    true,
		LanguageTypeSLN:    true,
		LanguageTypeSLP:    true,
		LanguageTypeSLQ:    true,
		LanguageTypeSLR:    true,
		LanguageTypeSLS:    true,
		LanguageTypeSLT:    true,
		LanguageTypeSLU:    true,
		LanguageTypeSLW:    true,
		LanguageTypeSLX:    true,
		LanguageTypeSLY:    true,
		LanguageTypeSLZ:    true,
		LanguageTypeSMA:    true,
		LanguageTypeSMB:    true,
		LanguageTypeSMC:    true,
		LanguageTypeSMD:    true,
		LanguageTypeSMF:    true,
		LanguageTypeSMG:    true,
		LanguageTypeSMH:    true,
		LanguageTypeSMI:    true,
		LanguageTypeSMJ:    true,
		LanguageTypeSMK:    true,
		LanguageTypeSML:    true,
		LanguageTypeSMM:    true,
		LanguageTypeSMN:    true,
		LanguageTypeSMP:    true,
		LanguageTypeSMQ:    true,
		LanguageTypeSMR:    true,
		LanguageTypeSMS:    true,
		LanguageTypeSMT:    true,
		LanguageTypeSMU:    true,
		LanguageTypeSMV:    true,
		LanguageTypeSMW:    true,
		LanguageTypeSMX:    true,
		LanguageTypeSMY:    true,
		LanguageTypeSMZ:    true,
		LanguageTypeSNB:    true,
		LanguageTypeSNC:    true,
		LanguageTypeSNE:    true,
		LanguageTypeSNF:    true,
		LanguageTypeSNG:    true,
		LanguageTypeSNH:    true,
		LanguageTypeSNI:    true,
		LanguageTypeSNJ:    true,
		LanguageTypeSNK:    true,
		LanguageTypeSNL:    true,
		LanguageTypeSNM:    true,
		LanguageTypeSNN:    true,
		LanguageTypeSNO:    true,
		LanguageTypeSNP:    true,
		LanguageTypeSNQ:    true,
		LanguageTypeSNR:    true,
		LanguageTypeSNS:    true,
		LanguageTypeSNU:    true,
		LanguageTypeSNV:    true,
		LanguageTypeSNW:    true,
		LanguageTypeSNX:    true,
		LanguageTypeSNY:    true,
		LanguageTypeSNZ:    true,
		LanguageTypeSOA:    true,
		LanguageTypeSOB:    true,
		LanguageTypeSOC:    true,
		LanguageTypeSOD:    true,
		LanguageTypeSOE:    true,
		LanguageTypeSOG:    true,
		LanguageTypeSOH:    true,
		LanguageTypeSOI:    true,
		LanguageTypeSOJ:    true,
		LanguageTypeSOK:    true,
		LanguageTypeSOL:    true,
		LanguageTypeSON:    true,
		LanguageTypeSOO:    true,
		LanguageTypeSOP:    true,
		LanguageTypeSOQ:    true,
		LanguageTypeSOR:    true,
		LanguageTypeSOS:    true,
		LanguageTypeSOU:    true,
		LanguageTypeSOV:    true,
		LanguageTypeSOW:    true,
		LanguageTypeSOX:    true,
		LanguageTypeSOY:    true,
		LanguageTypeSOZ:    true,
		LanguageTypeSPB:    true,
		LanguageTypeSPC:    true,
		LanguageTypeSPD:    true,
		LanguageTypeSPE:    true,
		LanguageTypeSPG:    true,
		LanguageTypeSPI:    true,
		LanguageTypeSPK:    true,
		LanguageTypeSPL:    true,
		LanguageTypeSPM:    true,
		LanguageTypeSPO:    true,
		LanguageTypeSPP:    true,
		LanguageTypeSPQ:    true,
		LanguageTypeSPR:    true,
		LanguageTypeSPS:    true,
		LanguageTypeSPT:    true,
		LanguageTypeSPU:    true,
		LanguageTypeSPV:    true,
		LanguageTypeSPX:    true,
		LanguageTypeSPY:    true,
		LanguageTypeSQA:    true,
		LanguageTypeSQH:    true,
		LanguageTypeSQJ:    true,
		LanguageTypeSQK:    true,
		LanguageTypeSQM:    true,
		LanguageTypeSQN:    true,
		LanguageTypeSQO:    true,
		LanguageTypeSQQ:    true,
		LanguageTypeSQR:    true,
		LanguageTypeSQS:    true,
		LanguageTypeSQT:    true,
		LanguageTypeSQU:    true,
		LanguageTypeSRA:    true,
		LanguageTypeSRB:    true,
		LanguageTypeSRC:    true,
		LanguageTypeSRE:    true,
		LanguageTypeSRF:    true,
		LanguageTypeSRG:    true,
		LanguageTypeSRH:    true,
		LanguageTypeSRI:    true,
		LanguageTypeSRK:    true,
		LanguageTypeSRL:    true,
		LanguageTypeSRM:    true,
		LanguageTypeSRN:    true,
		LanguageTypeSRO:    true,
		LanguageTypeSRQ:    true,
		LanguageTypeSRR:    true,
		LanguageTypeSRS:    true,
		LanguageTypeSRT:    true,
		LanguageTypeSRU:    true,
		LanguageTypeSRV:    true,
		LanguageTypeSRW:    true,
		LanguageTypeSRX:    true,
		LanguageTypeSRY:    true,
		LanguageTypeSRZ:    true,
		LanguageTypeSSA:    true,
		LanguageTypeSSB:    true,
		LanguageTypeSSC:    true,
		LanguageTypeSSD:    true,
		LanguageTypeSSE:    true,
		LanguageTypeSSF:    true,
		LanguageTypeSSG:    true,
		LanguageTypeSSH:    true,
		LanguageTypeSSI:    true,
		LanguageTypeSSJ:    true,
		LanguageTypeSSK:    true,
		LanguageTypeSSL:    true,
		LanguageTypeSSM:    true,
		LanguageTypeSSN:    true,
		LanguageTypeSSO:    true,
		LanguageTypeSSP:    true,
		LanguageTypeSSQ:    true,
		LanguageTypeSSR:    true,
		LanguageTypeSSS:    true,
		LanguageTypeSST:    true,
		LanguageTypeSSU:    true,
		LanguageTypeSSV:    true,
		LanguageTypeSSX:    true,
		LanguageTypeSSY:    true,
		LanguageTypeSSZ:    true,
		LanguageTypeSTA:    true,
		LanguageTypeSTB:    true,
		LanguageTypeSTD:    true,
		LanguageTypeSTE:    true,
		LanguageTypeSTF:    true,
		LanguageTypeSTG:    true,
		LanguageTypeSTH:    true,
		LanguageTypeSTI:    true,
		LanguageTypeSTJ:    true,
		LanguageTypeSTK:    true,
		LanguageTypeSTL:    true,
		LanguageTypeSTM:    true,
		LanguageTypeSTN:    true,
		LanguageTypeSTO:    true,
		LanguageTypeSTP:    true,
		LanguageTypeSTQ:    true,
		LanguageTypeSTR:    true,
		LanguageTypeSTS:    true,
		LanguageTypeSTT:    true,
		LanguageTypeSTU:    true,
		LanguageTypeSTV:    true,
		LanguageTypeSTW:    true,
		LanguageTypeSTY:    true,
		LanguageTypeSUA:    true,
		LanguageTypeSUB:    true,
		LanguageTypeSUC:    true,
		LanguageTypeSUE:    true,
		LanguageTypeSUG:    true,
		LanguageTypeSUI:    true,
		LanguageTypeSUJ:    true,
		LanguageTypeSUK:    true,
		LanguageTypeSUL:    true,
		LanguageTypeSUM:    true,
		LanguageTypeSUQ:    true,
		LanguageTypeSUR:    true,
		LanguageTypeSUS:    true,
		LanguageTypeSUT:    true,
		LanguageTypeSUV:    true,
		LanguageTypeSUW:    true,
		LanguageTypeSUX:    true,
		LanguageTypeSUY:    true,
		LanguageTypeSUZ:    true,
		LanguageTypeSVA:    true,
		LanguageTypeSVB:    true,
		LanguageTypeSVC:    true,
		LanguageTypeSVE:    true,
		LanguageTypeSVK:    true,
		LanguageTypeSVM:    true,
		LanguageTypeSVR:    true,
		LanguageTypeSVS:    true,
		LanguageTypeSVX:    true,
		LanguageTypeSWB:    true,
		LanguageTypeSWC:    true,
		LanguageTypeSWF:    true,
		LanguageTypeSWG:    true,
		LanguageTypeSWH:    true,
		LanguageTypeSWI:    true,
		LanguageTypeSWJ:    true,
		LanguageTypeSWK:    true,
		LanguageTypeSWL:    true,
		LanguageTypeSWM:    true,
		LanguageTypeSWN:    true,
		LanguageTypeSWO:    true,
		LanguageTypeSWP:    true,
		LanguageTypeSWQ:    true,
		LanguageTypeSWR:    true,
		LanguageTypeSWS:    true,
		LanguageTypeSWT:    true,
		LanguageTypeSWU:    true,
		LanguageTypeSWV:    true,
		LanguageTypeSWW:    true,
		LanguageTypeSWX:    true,
		LanguageTypeSWY:    true,
		LanguageTypeSXB:    true,
		LanguageTypeSXC:    true,
		LanguageTypeSXE:    true,
		LanguageTypeSXG:    true,
		LanguageTypeSXK:    true,
		LanguageTypeSXL:    true,
		LanguageTypeSXM:    true,
		LanguageTypeSXN:    true,
		LanguageTypeSXO:    true,
		LanguageTypeSXR:    true,
		LanguageTypeSXS:    true,
		LanguageTypeSXU:    true,
		LanguageTypeSXW:    true,
		LanguageTypeSYA:    true,
		LanguageTypeSYB:    true,
		LanguageTypeSYC:    true,
		LanguageTypeSYD:    true,
		LanguageTypeSYI:    true,
		LanguageTypeSYK:    true,
		LanguageTypeSYL:    true,
		LanguageTypeSYM:    true,
		LanguageTypeSYN:    true,
		LanguageTypeSYO:    true,
		LanguageTypeSYR:    true,
		LanguageTypeSYS:    true,
		LanguageTypeSYW:    true,
		LanguageTypeSYY:    true,
		LanguageTypeSZA:    true,
		LanguageTypeSZB:    true,
		LanguageTypeSZC:    true,
		LanguageTypeSZD:    true,
		LanguageTypeSZE:    true,
		LanguageTypeSZG:    true,
		LanguageTypeSZL:    true,
		LanguageTypeSZN:    true,
		LanguageTypeSZP:    true,
		LanguageTypeSZV:    true,
		LanguageTypeSZW:    true,
		LanguageTypeTAA:    true,
		LanguageTypeTAB:    true,
		LanguageTypeTAC:    true,
		LanguageTypeTAD:    true,
		LanguageTypeTAE:    true,
		LanguageTypeTAF:    true,
		LanguageTypeTAG:    true,
		LanguageTypeTAI:    true,
		LanguageTypeTAJ:    true,
		LanguageTypeTAK:    true,
		LanguageTypeTAL:    true,
		LanguageTypeTAN:    true,
		LanguageTypeTAO:    true,
		LanguageTypeTAP:    true,
		LanguageTypeTAQ:    true,
		LanguageTypeTAR:    true,
		LanguageTypeTAS:    true,
		LanguageTypeTAU:    true,
		LanguageTypeTAV:    true,
		LanguageTypeTAW:    true,
		LanguageTypeTAX:    true,
		LanguageTypeTAY:    true,
		LanguageTypeTAZ:    true,
		LanguageTypeTBA:    true,
		LanguageTypeTBB:    true,
		LanguageTypeTBC:    true,
		LanguageTypeTBD:    true,
		LanguageTypeTBE:    true,
		LanguageTypeTBF:    true,
		LanguageTypeTBG:    true,
		LanguageTypeTBH:    true,
		LanguageTypeTBI:    true,
		LanguageTypeTBJ:    true,
		LanguageTypeTBK:    true,
		LanguageTypeTBL:    true,
		LanguageTypeTBM:    true,
		LanguageTypeTBN:    true,
		LanguageTypeTBO:    true,
		LanguageTypeTBP:    true,
		LanguageTypeTBQ:    true,
		LanguageTypeTBR:    true,
		LanguageTypeTBS:    true,
		LanguageTypeTBT:    true,
		LanguageTypeTBU:    true,
		LanguageTypeTBV:    true,
		LanguageTypeTBW:    true,
		LanguageTypeTBX:    true,
		LanguageTypeTBY:    true,
		LanguageTypeTBZ:    true,
		LanguageTypeTCA:    true,
		LanguageTypeTCB:    true,
		LanguageTypeTCC:    true,
		LanguageTypeTCD:    true,
		LanguageTypeTCE:    true,
		LanguageTypeTCF:    true,
		LanguageTypeTCG:    true,
		LanguageTypeTCH:    true,
		LanguageTypeTCI:    true,
		LanguageTypeTCK:    true,
		LanguageTypeTCL:    true,
		LanguageTypeTCM:    true,
		LanguageTypeTCN:    true,
		LanguageTypeTCO:    true,
		LanguageTypeTCP:    true,
		LanguageTypeTCQ:    true,
		LanguageTypeTCS:    true,
		LanguageTypeTCT:    true,
		LanguageTypeTCU:    true,
		LanguageTypeTCW:    true,
		LanguageTypeTCX:    true,
		LanguageTypeTCY:    true,
		LanguageTypeTCZ:    true,
		LanguageTypeTDA:    true,
		LanguageTypeTDB:    true,
		LanguageTypeTDC:    true,
		LanguageTypeTDD:    true,
		LanguageTypeTDE:    true,
		LanguageTypeTDF:    true,
		LanguageTypeTDG:    true,
		LanguageTypeTDH:    true,
		LanguageTypeTDI:    true,
		LanguageTypeTDJ:    true,
		LanguageTypeTDK:    true,
		LanguageTypeTDL:    true,
		LanguageTypeTDN:    true,
		LanguageTypeTDO:    true,
		LanguageTypeTDQ:    true,
		LanguageTypeTDR:    true,
		LanguageTypeTDS:    true,
		LanguageTypeTDT:    true,
		LanguageTypeTDU:    true,
		LanguageTypeTDV:    true,
		LanguageTypeTDX:    true,
		LanguageTypeTDY:    true,
		LanguageTypeTEA:    true,
		LanguageTypeTEB:    true,
		LanguageTypeTEC:    true,
		LanguageTypeTED:    true,
		LanguageTypeTEE:    true,
		LanguageTypeTEF:    true,
		LanguageTypeTEG:    true,
		LanguageTypeTEH:    true,
		LanguageTypeTEI:    true,
		LanguageTypeTEK:    true,
		LanguageTypeTEM:    true,
		LanguageTypeTEN:    true,
		LanguageTypeTEO:    true,
		LanguageTypeTEP:    true,
		LanguageTypeTEQ:    true,
		LanguageTypeTER:    true,
		LanguageTypeTES:    true,
		LanguageTypeTET:    true,
		LanguageTypeTEU:    true,
		LanguageTypeTEV:    true,
		LanguageTypeTEW:    true,
		LanguageTypeTEX:    true,
		LanguageTypeTEY:    true,
		LanguageTypeTFI:    true,
		LanguageTypeTFN:    true,
		LanguageTypeTFO:    true,
		LanguageTypeTFR:    true,
		LanguageTypeTFT:    true,
		LanguageTypeTGA:    true,
		LanguageTypeTGB:    true,
		LanguageTypeTGC:    true,
		LanguageTypeTGD:    true,
		LanguageTypeTGE:    true,
		LanguageTypeTGF:    true,
		LanguageTypeTGG:    true,
		LanguageTypeTGH:    true,
		LanguageTypeTGI:    true,
		LanguageTypeTGJ:    true,
		LanguageTypeTGN:    true,
		LanguageTypeTGO:    true,
		LanguageTypeTGP:    true,
		LanguageTypeTGQ:    true,
		LanguageTypeTGR:    true,
		LanguageTypeTGS:    true,
		LanguageTypeTGT:    true,
		LanguageTypeTGU:    true,
		LanguageTypeTGV:    true,
		LanguageTypeTGW:    true,
		LanguageTypeTGX:    true,
		LanguageTypeTGY:    true,
		LanguageTypeTGZ:    true,
		LanguageTypeTHC:    true,
		LanguageTypeTHD:    true,
		LanguageTypeTHE:    true,
		LanguageTypeTHF:    true,
		LanguageTypeTHH:    true,
		LanguageTypeTHI:    true,
		LanguageTypeTHK:    true,
		LanguageTypeTHL:    true,
		LanguageTypeTHM:    true,
		LanguageTypeTHN:    true,
		LanguageTypeTHP:    true,
		LanguageTypeTHQ:    true,
		LanguageTypeTHR:    true,
		LanguageTypeTHS:    true,
		LanguageTypeTHT:    true,
		LanguageTypeTHU:    true,
		LanguageTypeTHV:    true,
		LanguageTypeTHW:    true,
		LanguageTypeTHX:    true,
		LanguageTypeTHY:    true,
		LanguageTypeTHZ:    true,
		LanguageTypeTIA:    true,
		LanguageTypeTIC:    true,
		LanguageTypeTID:    true,
		LanguageTypeTIE:    true,
		LanguageTypeTIF:    true,
		LanguageTypeTIG:    true,
		LanguageTypeTIH:    true,
		LanguageTypeTII:    true,
		LanguageTypeTIJ:    true,
		LanguageTypeTIK:    true,
		LanguageTypeTIL:    true,
		LanguageTypeTIM:    true,
		LanguageTypeTIN:    true,
		LanguageTypeTIO:    true,
		LanguageTypeTIP:    true,
		LanguageTypeTIQ:    true,
		LanguageTypeTIS:    true,
		LanguageTypeTIT:    true,
		LanguageTypeTIU:    true,
		LanguageTypeTIV:    true,
		LanguageTypeTIW:    true,
		LanguageTypeTIX:    true,
		LanguageTypeTIY:    true,
		LanguageTypeTIZ:    true,
		LanguageTypeTJA:    true,
		LanguageTypeTJG:    true,
		LanguageTypeTJI:    true,
		LanguageTypeTJL:    true,
		LanguageTypeTJM:    true,
		LanguageTypeTJN:    true,
		LanguageTypeTJO:    true,
		LanguageTypeTJS:    true,
		LanguageTypeTJU:    true,
		LanguageTypeTJW:    true,
		LanguageTypeTKA:    true,
		LanguageTypeTKB:    true,
		LanguageTypeTKD:    true,
		LanguageTypeTKE:    true,
		LanguageTypeTKF:    true,
		LanguageTypeTKG:    true,
		LanguageTypeTKK:    true,
		LanguageTypeTKL:    true,
		LanguageTypeTKM:    true,
		LanguageTypeTKN:    true,
		LanguageTypeTKP:    true,
		LanguageTypeTKQ:    true,
		LanguageTypeTKR:    true,
		LanguageTypeTKS:    true,
		LanguageTypeTKT:    true,
		LanguageTypeTKU:    true,
		LanguageTypeTKW:    true,
		LanguageTypeTKX:    true,
		LanguageTypeTKZ:    true,
		LanguageTypeTLA:    true,
		LanguageTypeTLB:    true,
		LanguageTypeTLC:    true,
		LanguageTypeTLD:    true,
		LanguageTypeTLF:    true,
		LanguageTypeTLG:    true,
		LanguageTypeTLH:    true,
		LanguageTypeTLI:    true,
		LanguageTypeTLJ:    true,
		LanguageTypeTLK:    true,
		LanguageTypeTLL:    true,
		LanguageTypeTLM:    true,
		LanguageTypeTLN:    true,
		LanguageTypeTLO:    true,
		LanguageTypeTLP:    true,
		LanguageTypeTLQ:    true,
		LanguageTypeTLR:    true,
		LanguageTypeTLS:    true,
		LanguageTypeTLT:    true,
		LanguageTypeTLU:    true,
		LanguageTypeTLV:    true,
		LanguageTypeTLW:    true,
		LanguageTypeTLX:    true,
		LanguageTypeTLY:    true,
		LanguageTypeTMA:    true,
		LanguageTypeTMB:    true,
		LanguageTypeTMC:    true,
		LanguageTypeTMD:    true,
		LanguageTypeTME:    true,
		LanguageTypeTMF:    true,
		LanguageTypeTMG:    true,
		LanguageTypeTMH:    true,
		LanguageTypeTMI:    true,
		LanguageTypeTMJ:    true,
		LanguageTypeTMK:    true,
		LanguageTypeTML:    true,
		LanguageTypeTMM:    true,
		LanguageTypeTMN:    true,
		LanguageTypeTMO:    true,
		LanguageTypeTMP:    true,
		LanguageTypeTMQ:    true,
		LanguageTypeTMR:    true,
		LanguageTypeTMS:    true,
		LanguageTypeTMT:    true,
		LanguageTypeTMU:    true,
		LanguageTypeTMV:    true,
		LanguageTypeTMW:    true,
		LanguageTypeTMY:    true,
		LanguageTypeTMZ:    true,
		LanguageTypeTNA:    true,
		LanguageTypeTNB:    true,
		LanguageTypeTNC:    true,
		LanguageTypeTND:    true,
		LanguageTypeTNE:    true,
		LanguageTypeTNF:    true,
		LanguageTypeTNG:    true,
		LanguageTypeTNH:    true,
		LanguageTypeTNI:    true,
		LanguageTypeTNK:    true,
		LanguageTypeTNL:    true,
		LanguageTypeTNM:    true,
		LanguageTypeTNN:    true,
		LanguageTypeTNO:    true,
		LanguageTypeTNP:    true,
		LanguageTypeTNQ:    true,
		LanguageTypeTNR:    true,
		LanguageTypeTNS:    true,
		LanguageTypeTNT:    true,
		LanguageTypeTNU:    true,
		LanguageTypeTNV:    true,
		LanguageTypeTNW:    true,
		LanguageTypeTNX:    true,
		LanguageTypeTNY:    true,
		LanguageTypeTNZ:    true,
		LanguageTypeTOB:    true,
		LanguageTypeTOC:    true,
		LanguageTypeTOD:    true,
		LanguageTypeTOE:    true,
		LanguageTypeTOF:    true,
		LanguageTypeTOG:    true,
		LanguageTypeTOH:    true,
		LanguageTypeTOI:    true,
		LanguageTypeTOJ:    true,
		LanguageTypeTOL:    true,
		LanguageTypeTOM:    true,
		LanguageTypeTOO:    true,
		LanguageTypeTOP:    true,
		LanguageTypeTOQ:    true,
		LanguageTypeTOR:    true,
		LanguageTypeTOS:    true,
		LanguageTypeTOU:    true,
		LanguageTypeTOV:    true,
		LanguageTypeTOW:    true,
		LanguageTypeTOX:    true,
		LanguageTypeTOY:    true,
		LanguageTypeTOZ:    true,
		LanguageTypeTPA:    true,
		LanguageTypeTPC:    true,
		LanguageTypeTPE:    true,
		LanguageTypeTPF:    true,
		LanguageTypeTPG:    true,
		LanguageTypeTPI:    true,
		LanguageTypeTPJ:    true,
		LanguageTypeTPK:    true,
		LanguageTypeTPL:    true,
		LanguageTypeTPM:    true,
		LanguageTypeTPN:    true,
		LanguageTypeTPO:    true,
		LanguageTypeTPP:    true,
		LanguageTypeTPQ:    true,
		LanguageTypeTPR:    true,
		LanguageTypeTPT:    true,
		LanguageTypeTPU:    true,
		LanguageTypeTPV:    true,
		LanguageTypeTPW:    true,
		LanguageTypeTPX:    true,
		LanguageTypeTPY:    true,
		LanguageTypeTPZ:    true,
		LanguageTypeTQB:    true,
		LanguageTypeTQL:    true,
		LanguageTypeTQM:    true,
		LanguageTypeTQN:    true,
		LanguageTypeTQO:    true,
		LanguageTypeTQP:    true,
		LanguageTypeTQQ:    true,
		LanguageTypeTQR:    true,
		LanguageTypeTQT:    true,
		LanguageTypeTQU:    true,
		LanguageTypeTQW:    true,
		LanguageTypeTRA:    true,
		LanguageTypeTRB:    true,
		LanguageTypeTRC:    true,
		LanguageTypeTRD:    true,
		LanguageTypeTRE:    true,
		LanguageTypeTRF:    true,
		LanguageTypeTRG:    true,
		LanguageTypeTRH:    true,
		LanguageTypeTRI:    true,
		LanguageTypeTRJ:    true,
		LanguageTypeTRK:    true,
		LanguageTypeTRL:    true,
		LanguageTypeTRM:    true,
		LanguageTypeTRN:    true,
		LanguageTypeTRO:    true,
		LanguageTypeTRP:    true,
		LanguageTypeTRQ:    true,
		LanguageTypeTRR:    true,
		LanguageTypeTRS:    true,
		LanguageTypeTRT:    true,
		LanguageTypeTRU:    true,
		LanguageTypeTRV:    true,
		LanguageTypeTRW:    true,
		LanguageTypeTRX:    true,
		LanguageTypeTRY:    true,
		LanguageTypeTRZ:    true,
		LanguageTypeTSA:    true,
		LanguageTypeTSB:    true,
		LanguageTypeTSC:    true,
		LanguageTypeTSD:    true,
		LanguageTypeTSE:    true,
		LanguageTypeTSF:    true,
		LanguageTypeTSG:    true,
		LanguageTypeTSH:    true,
		LanguageTypeTSI:    true,
		LanguageTypeTSJ:    true,
		LanguageTypeTSK:    true,
		LanguageTypeTSL:    true,
		LanguageTypeTSM:    true,
		LanguageTypeTSP:    true,
		LanguageTypeTSQ:    true,
		LanguageTypeTSR:    true,
		LanguageTypeTSS:    true,
		LanguageTypeTST:    true,
		LanguageTypeTSU:    true,
		LanguageTypeTSV:    true,
		LanguageTypeTSW:    true,
		LanguageTypeTSX:    true,
		LanguageTypeTSY:    true,
		LanguageTypeTSZ:    true,
		LanguageTypeTTA:    true,
		LanguageTypeTTB:    true,
		LanguageTypeTTC:    true,
		LanguageTypeTTD:    true,
		LanguageTypeTTE:    true,
		LanguageTypeTTF:    true,
		LanguageTypeTTG:    true,
		LanguageTypeTTH:    true,
		LanguageTypeTTI:    true,
		LanguageTypeTTJ:    true,
		LanguageTypeTTK:    true,
		LanguageTypeTTL:    true,
		LanguageTypeTTM:    true,
		LanguageTypeTTN:    true,
		LanguageTypeTTO:    true,
		LanguageTypeTTP:    true,
		LanguageTypeTTQ:    true,
		LanguageTypeTTR:    true,
		LanguageTypeTTS:    true,
		LanguageTypeTTT:    true,
		LanguageTypeTTU:    true,
		LanguageTypeTTV:    true,
		LanguageTypeTTW:    true,
		LanguageTypeTTY:    true,
		LanguageTypeTTZ:    true,
		LanguageTypeTUA:    true,
		LanguageTypeTUB:    true,
		LanguageTypeTUC:    true,
		LanguageTypeTUD:    true,
		LanguageTypeTUE:    true,
		LanguageTypeTUF:    true,
		LanguageTypeTUG:    true,
		LanguageTypeTUH:    true,
		LanguageTypeTUI:    true,
		LanguageTypeTUJ:    true,
		LanguageTypeTUL:    true,
		LanguageTypeTUM:    true,
		LanguageTypeTUN:    true,
		LanguageTypeTUO:    true,
		LanguageTypeTUP:    true,
		LanguageTypeTUQ:    true,
		LanguageTypeTUS:    true,
		LanguageTypeTUT:    true,
		LanguageTypeTUU:    true,
		LanguageTypeTUV:    true,
		LanguageTypeTUW:    true,
		LanguageTypeTUX:    true,
		LanguageTypeTUY:    true,
		LanguageTypeTUZ:    true,
		LanguageTypeTVA:    true,
		LanguageTypeTVD:    true,
		LanguageTypeTVE:    true,
		LanguageTypeTVK:    true,
		LanguageTypeTVL:    true,
		LanguageTypeTVM:    true,
		LanguageTypeTVN:    true,
		LanguageTypeTVO:    true,
		LanguageTypeTVS:    true,
		LanguageTypeTVT:    true,
		LanguageTypeTVU:    true,
		LanguageTypeTVW:    true,
		LanguageTypeTVY:    true,
		LanguageTypeTWA:    true,
		LanguageTypeTWB:    true,
		LanguageTypeTWC:    true,
		LanguageTypeTWD:    true,
		LanguageTypeTWE:    true,
		LanguageTypeTWF:    true,
		LanguageTypeTWG:    true,
		LanguageTypeTWH:    true,
		LanguageTypeTWL:    true,
		LanguageTypeTWM:    true,
		LanguageTypeTWN:    true,
		LanguageTypeTWO:    true,
		LanguageTypeTWP:    true,
		LanguageTypeTWQ:    true,
		LanguageTypeTWR:    true,
		LanguageTypeTWT:    true,
		LanguageTypeTWU:    true,
		LanguageTypeTWW:    true,
		LanguageTypeTWX:    true,
		LanguageTypeTWY:    true,
		LanguageTypeTXA:    true,
		LanguageTypeTXB:    true,
		LanguageTypeTXC:    true,
		LanguageTypeTXE:    true,
		LanguageTypeTXG:    true,
		LanguageTypeTXH:    true,
		LanguageTypeTXI:    true,
		LanguageTypeTXM:    true,
		LanguageTypeTXN:    true,
		LanguageTypeTXO:    true,
		LanguageTypeTXQ:    true,
		LanguageTypeTXR:    true,
		LanguageTypeTXS:    true,
		LanguageTypeTXT:    true,
		LanguageTypeTXU:    true,
		LanguageTypeTXX:    true,
		LanguageTypeTXY:    true,
		LanguageTypeTYA:    true,
		LanguageTypeTYE:    true,
		LanguageTypeTYH:    true,
		LanguageTypeTYI:    true,
		LanguageTypeTYJ:    true,
		LanguageTypeTYL:    true,
		LanguageTypeTYN:    true,
		LanguageTypeTYP:    true,
		LanguageTypeTYR:    true,
		LanguageTypeTYS:    true,
		LanguageTypeTYT:    true,
		LanguageTypeTYU:    true,
		LanguageTypeTYV:    true,
		LanguageTypeTYX:    true,
		LanguageTypeTYZ:    true,
		LanguageTypeTZA:    true,
		LanguageTypeTZH:    true,
		LanguageTypeTZJ:    true,
		LanguageTypeTZL:    true,
		LanguageTypeTZM:    true,
		LanguageTypeTZN:    true,
		LanguageTypeTZO:    true,
		LanguageTypeTZX:    true,
		LanguageTypeUAM:    true,
		LanguageTypeUAN:    true,
		LanguageTypeUAR:    true,
		LanguageTypeUBA:    true,
		LanguageTypeUBI:    true,
		LanguageTypeUBL:    true,
		LanguageTypeUBR:    true,
		LanguageTypeUBU:    true,
		LanguageTypeUBY:    true,
		LanguageTypeUDA:    true,
		LanguageTypeUDE:    true,
		LanguageTypeUDG:    true,
		LanguageTypeUDI:    true,
		LanguageTypeUDJ:    true,
		LanguageTypeUDL:    true,
		LanguageTypeUDM:    true,
		LanguageTypeUDU:    true,
		LanguageTypeUES:    true,
		LanguageTypeUFI:    true,
		LanguageTypeUGA:    true,
		LanguageTypeUGB:    true,
		LanguageTypeUGE:    true,
		LanguageTypeUGN:    true,
		LanguageTypeUGO:    true,
		LanguageTypeUGY:    true,
		LanguageTypeUHA:    true,
		LanguageTypeUHN:    true,
		LanguageTypeUIS:    true,
		LanguageTypeUIV:    true,
		LanguageTypeUJI:    true,
		LanguageTypeUKA:    true,
		LanguageTypeUKG:    true,
		LanguageTypeUKH:    true,
		LanguageTypeUKL:    true,
		LanguageTypeUKP:    true,
		LanguageTypeUKQ:    true,
		LanguageTypeUKS:    true,
		LanguageTypeUKU:    true,
		LanguageTypeUKW:    true,
		LanguageTypeUKY:    true,
		LanguageTypeULA:    true,
		LanguageTypeULB:    true,
		LanguageTypeULC:    true,
		LanguageTypeULE:    true,
		LanguageTypeULF:    true,
		LanguageTypeULI:    true,
		LanguageTypeULK:    true,
		LanguageTypeULL:    true,
		LanguageTypeULM:    true,
		LanguageTypeULN:    true,
		LanguageTypeULU:    true,
		LanguageTypeULW:    true,
		LanguageTypeUMA:    true,
		LanguageTypeUMB:    true,
		LanguageTypeUMC:    true,
		LanguageTypeUMD:    true,
		LanguageTypeUMG:    true,
		LanguageTypeUMI:    true,
		LanguageTypeUMM:    true,
		LanguageTypeUMN:    true,
		LanguageTypeUMO:    true,
		LanguageTypeUMP:    true,
		LanguageTypeUMR:    true,
		LanguageTypeUMS:    true,
		LanguageTypeUMU:    true,
		LanguageTypeUNA:    true,
		LanguageTypeUND:    true,
		LanguageTypeUNE:    true,
		LanguageTypeUNG:    true,
		LanguageTypeUNK:    true,
		LanguageTypeUNM:    true,
		LanguageTypeUNN:    true,
		LanguageTypeUNP:    true,
		LanguageTypeUNR:    true,
		LanguageTypeUNU:    true,
		LanguageTypeUNX:    true,
		LanguageTypeUNZ:    true,
		LanguageTypeUOK:    true,
		LanguageTypeUPI:    true,
		LanguageTypeUPV:    true,
		LanguageTypeURA:    true,
		LanguageTypeURB:    true,
		LanguageTypeURC:    true,
		LanguageTypeURE:    true,
		LanguageTypeURF:    true,
		LanguageTypeURG:    true,
		LanguageTypeURH:    true,
		LanguageTypeURI:    true,
		LanguageTypeURJ:    true,
		LanguageTypeURK:    true,
		LanguageTypeURL:    true,
		LanguageTypeURM:    true,
		LanguageTypeURN:    true,
		LanguageTypeURO:    true,
		LanguageTypeURP:    true,
		LanguageTypeURR:    true,
		LanguageTypeURT:    true,
		LanguageTypeURU:    true,
		LanguageTypeURV:    true,
		LanguageTypeURW:    true,
		LanguageTypeURX:    true,
		LanguageTypeURY:    true,
		LanguageTypeURZ:    true,
		LanguageTypeUSA:    true,
		LanguageTypeUSH:    true,
		LanguageTypeUSI:    true,
		LanguageTypeUSK:    true,
		LanguageTypeUSP:    true,
		LanguageTypeUSU:    true,
		LanguageTypeUTA:    true,
		LanguageTypeUTE:    true,
		LanguageTypeUTP:    true,
		LanguageTypeUTR:    true,
		LanguageTypeUTU:    true,
		LanguageTypeUUM:    true,
		LanguageTypeUUN:    true,
		LanguageTypeUUR:    true,
		LanguageTypeUUU:    true,
		LanguageTypeUVE:    true,
		LanguageTypeUVH:    true,
		LanguageTypeUVL:    true,
		LanguageTypeUWA:    true,
		LanguageTypeUYA:    true,
		LanguageTypeUZN:    true,
		LanguageTypeUZS:    true,
		LanguageTypeVAA:    true,
		LanguageTypeVAE:    true,
		LanguageTypeVAF:    true,
		LanguageTypeVAG:    true,
		LanguageTypeVAH:    true,
		LanguageTypeVAI:    true,
		LanguageTypeVAJ:    true,
		LanguageTypeVAL:    true,
		LanguageTypeVAM:    true,
		LanguageTypeVAN:    true,
		LanguageTypeVAO:    true,
		LanguageTypeVAP:    true,
		LanguageTypeVAR:    true,
		LanguageTypeVAS:    true,
		LanguageTypeVAU:    true,
		LanguageTypeVAV:    true,
		LanguageTypeVAY:    true,
		LanguageTypeVBB:    true,
		LanguageTypeVBK:    true,
		LanguageTypeVEC:    true,
		LanguageTypeVED:    true,
		LanguageTypeVEL:    true,
		LanguageTypeVEM:    true,
		LanguageTypeVEO:    true,
		LanguageTypeVEP:    true,
		LanguageTypeVER:    true,
		LanguageTypeVGR:    true,
		LanguageTypeVGT:    true,
		LanguageTypeVIC:    true,
		LanguageTypeVID:    true,
		LanguageTypeVIF:    true,
		LanguageTypeVIG:    true,
		LanguageTypeVIL:    true,
		LanguageTypeVIN:    true,
		LanguageTypeVIS:    true,
		LanguageTypeVIT:    true,
		LanguageTypeVIV:    true,
		LanguageTypeVKA:    true,
		LanguageTypeVKI:    true,
		LanguageTypeVKJ:    true,
		LanguageTypeVKK:    true,
		LanguageTypeVKL:    true,
		LanguageTypeVKM:    true,
		LanguageTypeVKO:    true,
		LanguageTypeVKP:    true,
		LanguageTypeVKT:    true,
		LanguageTypeVKU:    true,
		LanguageTypeVLP:    true,
		LanguageTypeVLS:    true,
		LanguageTypeVMA:    true,
		LanguageTypeVMB:    true,
		LanguageTypeVMC:    true,
		LanguageTypeVMD:    true,
		LanguageTypeVME:    true,
		LanguageTypeVMF:    true,
		LanguageTypeVMG:    true,
		LanguageTypeVMH:    true,
		LanguageTypeVMI:    true,
		LanguageTypeVMJ:    true,
		LanguageTypeVMK:    true,
		LanguageTypeVML:    true,
		LanguageTypeVMM:    true,
		LanguageTypeVMP:    true,
		LanguageTypeVMQ:    true,
		LanguageTypeVMR:    true,
		LanguageTypeVMS:    true,
		LanguageTypeVMU:    true,
		LanguageTypeVMV:    true,
		LanguageTypeVMW:    true,
		LanguageTypeVMX:    true,
		LanguageTypeVMY:    true,
		LanguageTypeVMZ:    true,
		LanguageTypeVNK:    true,
		LanguageTypeVNM:    true,
		LanguageTypeVNP:    true,
		LanguageTypeVOR:    true,
		LanguageTypeVOT:    true,
		LanguageTypeVRA:    true,
		LanguageTypeVRO:    true,
		LanguageTypeVRS:    true,
		LanguageTypeVRT:    true,
		LanguageTypeVSI:    true,
		LanguageTypeVSL:    true,
		LanguageTypeVSV:    true,
		LanguageTypeVTO:    true,
		LanguageTypeVUM:    true,
		LanguageTypeVUN:    true,
		LanguageTypeVUT:    true,
		LanguageTypeVWA:    true,
		LanguageTypeWAA:    true,
		LanguageTypeWAB:    true,
		LanguageTypeWAC:    true,
		LanguageTypeWAD:    true,
		LanguageTypeWAE:    true,
		LanguageTypeWAF:    true,
		LanguageTypeWAG:    true,
		LanguageTypeWAH:    true,
		LanguageTypeWAI:    true,
		LanguageTypeWAJ:    true,
		LanguageTypeWAK:    true,
		LanguageTypeWAL:    true,
		LanguageTypeWAM:    true,
		LanguageTypeWAN:    true,
		LanguageTypeWAO:    true,
		LanguageTypeWAP:    true,
		LanguageTypeWAQ:    true,
		LanguageTypeWAR:    true,
		LanguageTypeWAS:    true,
		LanguageTypeWAT:    true,
		LanguageTypeWAU:    true,
		LanguageTypeWAV:    true,
		LanguageTypeWAW:    true,
		LanguageTypeWAX:    true,
		LanguageTypeWAY:    true,
		LanguageTypeWAZ:    true,
		LanguageTypeWBA:    true,
		LanguageTypeWBB:    true,
		LanguageTypeWBE:    true,
		LanguageTypeWBF:    true,
		LanguageTypeWBH:    true,
		LanguageTypeWBI:    true,
		LanguageTypeWBJ:    true,
		LanguageTypeWBK:    true,
		LanguageTypeWBL:    true,
		LanguageTypeWBM:    true,
		LanguageTypeWBP:    true,
		LanguageTypeWBQ:    true,
		LanguageTypeWBR:    true,
		LanguageTypeWBT:    true,
		LanguageTypeWBV:    true,
		LanguageTypeWBW:    true,
		LanguageTypeWCA:    true,
		LanguageTypeWCI:    true,
		LanguageTypeWDD:    true,
		LanguageTypeWDG:    true,
		LanguageTypeWDJ:    true,
		LanguageTypeWDK:    true,
		LanguageTypeWDU:    true,
		LanguageTypeWDY:    true,
		LanguageTypeWEA:    true,
		LanguageTypeWEC:    true,
		LanguageTypeWED:    true,
		LanguageTypeWEG:    true,
		LanguageTypeWEH:    true,
		LanguageTypeWEI:    true,
		LanguageTypeWEM:    true,
		LanguageTypeWEN:    true,
		LanguageTypeWEO:    true,
		LanguageTypeWEP:    true,
		LanguageTypeWER:    true,
		LanguageTypeWES:    true,
		LanguageTypeWET:    true,
		LanguageTypeWEU:    true,
		LanguageTypeWEW:    true,
		LanguageTypeWFG:    true,
		LanguageTypeWGA:    true,
		LanguageTypeWGB:    true,
		LanguageTypeWGG:    true,
		LanguageTypeWGI:    true,
		LanguageTypeWGO:    true,
		LanguageTypeWGU:    true,
		LanguageTypeWGW:    true,
		LanguageTypeWGY:    true,
		LanguageTypeWHA:    true,
		LanguageTypeWHG:    true,
		LanguageTypeWHK:    true,
		LanguageTypeWHU:    true,
		LanguageTypeWIB:    true,
		LanguageTypeWIC:    true,
		LanguageTypeWIE:    true,
		LanguageTypeWIF:    true,
		LanguageTypeWIG:    true,
		LanguageTypeWIH:    true,
		LanguageTypeWII:    true,
		LanguageTypeWIJ:    true,
		LanguageTypeWIK:    true,
		LanguageTypeWIL:    true,
		LanguageTypeWIM:    true,
		LanguageTypeWIN:    true,
		LanguageTypeWIR:    true,
		LanguageTypeWIT:    true,
		LanguageTypeWIU:    true,
		LanguageTypeWIV:    true,
		LanguageTypeWIW:    true,
		LanguageTypeWIY:    true,
		LanguageTypeWJA:    true,
		LanguageTypeWJI:    true,
		LanguageTypeWKA:    true,
		LanguageTypeWKB:    true,
		LanguageTypeWKD:    true,
		LanguageTypeWKL:    true,
		LanguageTypeWKU:    true,
		LanguageTypeWKW:    true,
		LanguageTypeWKY:    true,
		LanguageTypeWLA:    true,
		LanguageTypeWLC:    true,
		LanguageTypeWLE:    true,
		LanguageTypeWLG:    true,
		LanguageTypeWLI:    true,
		LanguageTypeWLK:    true,
		LanguageTypeWLL:    true,
		LanguageTypeWLM:    true,
		LanguageTypeWLO:    true,
		LanguageTypeWLR:    true,
		LanguageTypeWLS:    true,
		LanguageTypeWLU:    true,
		LanguageTypeWLV:    true,
		LanguageTypeWLW:    true,
		LanguageTypeWLX:    true,
		LanguageTypeWLY:    true,
		LanguageTypeWMA:    true,
		LanguageTypeWMB:    true,
		LanguageTypeWMC:    true,
		LanguageTypeWMD:    true,
		LanguageTypeWME:    true,
		LanguageTypeWMH:    true,
		LanguageTypeWMI:    true,
		LanguageTypeWMM:    true,
		LanguageTypeWMN:    true,
		LanguageTypeWMO:    true,
		LanguageTypeWMS:    true,
		LanguageTypeWMT:    true,
		LanguageTypeWMW:    true,
		LanguageTypeWMX:    true,
		LanguageTypeWNB:    true,
		LanguageTypeWNC:    true,
		LanguageTypeWND:    true,
		LanguageTypeWNE:    true,
		LanguageTypeWNG:    true,
		LanguageTypeWNI:    true,
		LanguageTypeWNK:    true,
		LanguageTypeWNM:    true,
		LanguageTypeWNN:    true,
		LanguageTypeWNO:    true,
		LanguageTypeWNP:    true,
		LanguageTypeWNU:    true,
		LanguageTypeWNW:    true,
		LanguageTypeWNY:    true,
		LanguageTypeWOA:    true,
		LanguageTypeWOB:    true,
		LanguageTypeWOC:    true,
		LanguageTypeWOD:    true,
		LanguageTypeWOE:    true,
		LanguageTypeWOF:    true,
		LanguageTypeWOG:    true,
		LanguageTypeWOI:    true,
		LanguageTypeWOK:    true,
		LanguageTypeWOM:    true,
		LanguageTypeWON:    true,
		LanguageTypeWOO:    true,
		LanguageTypeWOR:    true,
		LanguageTypeWOS:    true,
		LanguageTypeWOW:    true,
		LanguageTypeWOY:    true,
		LanguageTypeWPC:    true,
		LanguageTypeWRA:    true,
		LanguageTypeWRB:    true,
		LanguageTypeWRD:    true,
		LanguageTypeWRG:    true,
		LanguageTypeWRH:    true,
		LanguageTypeWRI:    true,
		LanguageTypeWRK:    true,
		LanguageTypeWRL:    true,
		LanguageTypeWRM:    true,
		LanguageTypeWRN:    true,
		LanguageTypeWRO:    true,
		LanguageTypeWRP:    true,
		LanguageTypeWRR:    true,
		LanguageTypeWRS:    true,
		LanguageTypeWRU:    true,
		LanguageTypeWRV:    true,
		LanguageTypeWRW:    true,
		LanguageTypeWRX:    true,
		LanguageTypeWRY:    true,
		LanguageTypeWRZ:    true,
		LanguageTypeWSA:    true,
		LanguageTypeWSI:    true,
		LanguageTypeWSK:    true,
		LanguageTypeWSR:    true,
		LanguageTypeWSS:    true,
		LanguageTypeWSU:    true,
		LanguageTypeWSV:    true,
		LanguageTypeWTF:    true,
		LanguageTypeWTH:    true,
		LanguageTypeWTI:    true,
		LanguageTypeWTK:    true,
		LanguageTypeWTM:    true,
		LanguageTypeWTW:    true,
		LanguageTypeWUA:    true,
		LanguageTypeWUB:    true,
		LanguageTypeWUD:    true,
		LanguageTypeWUH:    true,
		LanguageTypeWUL:    true,
		LanguageTypeWUM:    true,
		LanguageTypeWUN:    true,
		LanguageTypeWUR:    true,
		LanguageTypeWUT:    true,
		LanguageTypeWUU:    true,
		LanguageTypeWUV:    true,
		LanguageTypeWUX:    true,
		LanguageTypeWUY:    true,
		LanguageTypeWWA:    true,
		LanguageTypeWWB:    true,
		LanguageTypeWWO:    true,
		LanguageTypeWWR:    true,
		LanguageTypeWWW:    true,
		LanguageTypeWXA:    true,
		LanguageTypeWXW:    true,
		LanguageTypeWYA:    true,
		LanguageTypeWYB:    true,
		LanguageTypeWYI:    true,
		LanguageTypeWYM:    true,
		LanguageTypeWYR:    true,
		LanguageTypeWYY:    true,
		LanguageTypeXAA:    true,
		LanguageTypeXAB:    true,
		LanguageTypeXAC:    true,
		LanguageTypeXAD:    true,
		LanguageTypeXAE:    true,
		LanguageTypeXAG:    true,
		LanguageTypeXAI:    true,
		LanguageTypeXAL:    true,
		LanguageTypeXAM:    true,
		LanguageTypeXAN:    true,
		LanguageTypeXAO:    true,
		LanguageTypeXAP:    true,
		LanguageTypeXAQ:    true,
		LanguageTypeXAR:    true,
		LanguageTypeXAS:    true,
		LanguageTypeXAT:    true,
		LanguageTypeXAU:    true,
		LanguageTypeXAV:    true,
		LanguageTypeXAW:    true,
		LanguageTypeXAY:    true,
		LanguageTypeXBA:    true,
		LanguageTypeXBB:    true,
		LanguageTypeXBC:    true,
		LanguageTypeXBD:    true,
		LanguageTypeXBE:    true,
		LanguageTypeXBG:    true,
		LanguageTypeXBI:    true,
		LanguageTypeXBJ:    true,
		LanguageTypeXBM:    true,
		LanguageTypeXBN:    true,
		LanguageTypeXBO:    true,
		LanguageTypeXBP:    true,
		LanguageTypeXBR:    true,
		LanguageTypeXBW:    true,
		LanguageTypeXBX:    true,
		LanguageTypeXBY:    true,
		LanguageTypeXCB:    true,
		LanguageTypeXCC:    true,
		LanguageTypeXCE:    true,
		LanguageTypeXCG:    true,
		LanguageTypeXCH:    true,
		LanguageTypeXCL:    true,
		LanguageTypeXCM:    true,
		LanguageTypeXCN:    true,
		LanguageTypeXCO:    true,
		LanguageTypeXCR:    true,
		LanguageTypeXCT:    true,
		LanguageTypeXCU:    true,
		LanguageTypeXCV:    true,
		LanguageTypeXCW:    true,
		LanguageTypeXCY:    true,
		LanguageTypeXDA:    true,
		LanguageTypeXDC:    true,
		LanguageTypeXDK:    true,
		LanguageTypeXDM:    true,
		LanguageTypeXDY:    true,
		LanguageTypeXEB:    true,
		LanguageTypeXED:    true,
		LanguageTypeXEG:    true,
		LanguageTypeXEL:    true,
		LanguageTypeXEM:    true,
		LanguageTypeXEP:    true,
		LanguageTypeXER:    true,
		LanguageTypeXES:    true,
		LanguageTypeXET:    true,
		LanguageTypeXEU:    true,
		LanguageTypeXFA:    true,
		LanguageTypeXGA:    true,
		LanguageTypeXGB:    true,
		LanguageTypeXGD:    true,
		LanguageTypeXGF:    true,
		LanguageTypeXGG:    true,
		LanguageTypeXGI:    true,
		LanguageTypeXGL:    true,
		LanguageTypeXGM:    true,
		LanguageTypeXGN:    true,
		LanguageTypeXGR:    true,
		LanguageTypeXGU:    true,
		LanguageTypeXGW:    true,
		LanguageTypeXHA:    true,
		LanguageTypeXHC:    true,
		LanguageTypeXHD:    true,
		LanguageTypeXHE:    true,
		LanguageTypeXHR:    true,
		LanguageTypeXHT:    true,
		LanguageTypeXHU:    true,
		LanguageTypeXHV:    true,
		LanguageTypeXIA:    true,
		LanguageTypeXIB:    true,
		LanguageTypeXII:    true,
		LanguageTypeXIL:    true,
		LanguageTypeXIN:    true,
		LanguageTypeXIP:    true,
		LanguageTypeXIR:    true,
		LanguageTypeXIV:    true,
		LanguageTypeXIY:    true,
		LanguageTypeXJB:    true,
		LanguageTypeXJT:    true,
		LanguageTypeXKA:    true,
		LanguageTypeXKB:    true,
		LanguageTypeXKC:    true,
		LanguageTypeXKD:    true,
		LanguageTypeXKE:    true,
		LanguageTypeXKF:    true,
		LanguageTypeXKG:    true,
		LanguageTypeXKH:    true,
		LanguageTypeXKI:    true,
		LanguageTypeXKJ:    true,
		LanguageTypeXKK:    true,
		LanguageTypeXKL:    true,
		LanguageTypeXKN:    true,
		LanguageTypeXKO:    true,
		LanguageTypeXKP:    true,
		LanguageTypeXKQ:    true,
		LanguageTypeXKR:    true,
		LanguageTypeXKS:    true,
		LanguageTypeXKT:    true,
		LanguageTypeXKU:    true,
		LanguageTypeXKV:    true,
		LanguageTypeXKW:    true,
		LanguageTypeXKX:    true,
		LanguageTypeXKY:    true,
		LanguageTypeXKZ:    true,
		LanguageTypeXLA:    true,
		LanguageTypeXLB:    true,
		LanguageTypeXLC:    true,
		LanguageTypeXLD:    true,
		LanguageTypeXLE:    true,
		LanguageTypeXLG:    true,
		LanguageTypeXLI:    true,
		LanguageTypeXLN:    true,
		LanguageTypeXLO:    true,
		LanguageTypeXLP:    true,
		LanguageTypeXLS:    true,
		LanguageTypeXLU:    true,
		LanguageTypeXLY:    true,
		LanguageTypeXMA:    true,
		LanguageTypeXMB:    true,
		LanguageTypeXMC:    true,
		LanguageTypeXMD:    true,
		LanguageTypeXME:    true,
		LanguageTypeXMF:    true,
		LanguageTypeXMG:    true,
		LanguageTypeXMH:    true,
		LanguageTypeXMJ:    true,
		LanguageTypeXMK:    true,
		LanguageTypeXML:    true,
		LanguageTypeXMM:    true,
		LanguageTypeXMN:    true,
		LanguageTypeXMO:    true,
		LanguageTypeXMP:    true,
		LanguageTypeXMQ:    true,
		LanguageTypeXMR:    true,
		LanguageTypeXMS:    true,
		LanguageTypeXMT:    true,
		LanguageTypeXMU:    true,
		LanguageTypeXMV:    true,
		LanguageTypeXMW:    true,
		LanguageTypeXMX:    true,
		LanguageTypeXMY:    true,
		LanguageTypeXMZ:    true,
		LanguageTypeXNA:    true,
		LanguageTypeXNB:    true,
		LanguageTypeXND:    true,
		LanguageTypeXNG:    true,
		LanguageTypeXNH:    true,
		LanguageTypeXNI:    true,
		LanguageTypeXNK:    true,
		LanguageTypeXNN:    true,
		LanguageTypeXNO:    true,
		LanguageTypeXNR:    true,
		LanguageTypeXNS:    true,
		LanguageTypeXNT:    true,
		LanguageTypeXNU:    true,
		LanguageTypeXNY:    true,
		LanguageTypeXNZ:    true,
		LanguageTypeXOC:    true,
		LanguageTypeXOD:    true,
		LanguageTypeXOG:    true,
		LanguageTypeXOI:    true,
		LanguageTypeXOK:    true,
		LanguageTypeXOM:    true,
		LanguageTypeXON:    true,
		LanguageTypeXOO:    true,
		LanguageTypeXOP:    true,
		LanguageTypeXOR:    true,
		LanguageTypeXOW:    true,
		LanguageTypeXPA:    true,
		LanguageTypeXPC:    true,
		LanguageTypeXPE:    true,
		LanguageTypeXPG:    true,
		LanguageTypeXPI:    true,
		LanguageTypeXPJ:    true,
		LanguageTypeXPK:    true,
		LanguageTypeXPM:    true,
		LanguageTypeXPN:    true,
		LanguageTypeXPO:    true,
		LanguageTypeXPP:    true,
		LanguageTypeXPQ:    true,
		LanguageTypeXPR:    true,
		LanguageTypeXPS:    true,
		LanguageTypeXPT:    true,
		LanguageTypeXPU:    true,
		LanguageTypeXPY:    true,
		LanguageTypeXQA:    true,
		LanguageTypeXQT:    true,
		LanguageTypeXRA:    true,
		LanguageTypeXRB:    true,
		LanguageTypeXRD:    true,
		LanguageTypeXRE:    true,
		LanguageTypeXRG:    true,
		LanguageTypeXRI:    true,
		LanguageTypeXRM:    true,
		LanguageTypeXRN:    true,
		LanguageTypeXRQ:    true,
		LanguageTypeXRR:    true,
		LanguageTypeXRT:    true,
		LanguageTypeXRU:    true,
		LanguageTypeXRW:    true,
		LanguageTypeXSA:    true,
		LanguageTypeXSB:    true,
		LanguageTypeXSC:    true,
		LanguageTypeXSD:    true,
		LanguageTypeXSE:    true,
		LanguageTypeXSH:    true,
		LanguageTypeXSI:    true,
		LanguageTypeXSJ:    true,
		LanguageTypeXSL:    true,
		LanguageTypeXSM:    true,
		LanguageTypeXSN:    true,
		LanguageTypeXSO:    true,
		LanguageTypeXSP:    true,
		LanguageTypeXSQ:    true,
		LanguageTypeXSR:    true,
		LanguageTypeXSS:    true,
		LanguageTypeXSU:    true,
		LanguageTypeXSV:    true,
		LanguageTypeXSY:    true,
		LanguageTypeXTA:    true,
		LanguageTypeXTB:    true,
		LanguageTypeXTC:    true,
		LanguageTypeXTD:    true,
		LanguageTypeXTE:    true,
		LanguageTypeXTG:    true,
		LanguageTypeXTH:    true,
		LanguageTypeXTI:    true,
		LanguageTypeXTJ:    true,
		LanguageTypeXTL:    true,
		LanguageTypeXTM:    true,
		LanguageTypeXTN:    true,
		LanguageTypeXTO:    true,
		LanguageTypeXTP:    true,
		LanguageTypeXTQ:    true,
		LanguageTypeXTR:    true,
		LanguageTypeXTS:    true,
		LanguageTypeXTT:    true,
		LanguageTypeXTU:    true,
		LanguageTypeXTV:    true,
		LanguageTypeXTW:    true,
		LanguageTypeXTY:    true,
		LanguageTypeXTZ:    true,
		LanguageTypeXUA:    true,
		LanguageTypeXUB:    true,
		LanguageTypeXUD:    true,
		LanguageTypeXUG:    true,
		LanguageTypeXUJ:    true,
		LanguageTypeXUL:    true,
		LanguageTypeXUM:    true,
		LanguageTypeXUN:    true,
		LanguageTypeXUO:    true,
		LanguageTypeXUP:    true,
		LanguageTypeXUR:    true,
		LanguageTypeXUT:    true,
		LanguageTypeXUU:    true,
		LanguageTypeXVE:    true,
		LanguageTypeXVI:    true,
		LanguageTypeXVN:    true,
		LanguageTypeXVO:    true,
		LanguageTypeXVS:    true,
		LanguageTypeXWA:    true,
		LanguageTypeXWC:    true,
		LanguageTypeXWD:    true,
		LanguageTypeXWE:    true,
		LanguageTypeXWG:    true,
		LanguageTypeXWJ:    true,
		LanguageTypeXWK:    true,
		LanguageTypeXWL:    true,
		LanguageTypeXWO:    true,
		LanguageTypeXWR:    true,
		LanguageTypeXWT:    true,
		LanguageTypeXWW:    true,
		LanguageTypeXXB:    true,
		LanguageTypeXXK:    true,
		LanguageTypeXXM:    true,
		LanguageTypeXXR:    true,
		LanguageTypeXXT:    true,
		LanguageTypeXYA:    true,
		LanguageTypeXYB:    true,
		LanguageTypeXYJ:    true,
		LanguageTypeXYK:    true,
		LanguageTypeXYL:    true,
		LanguageTypeXYT:    true,
		LanguageTypeXYY:    true,
		LanguageTypeXZH:    true,
		LanguageTypeXZM:    true,
		LanguageTypeXZP:    true,
		LanguageTypeYAA:    true,
		LanguageTypeYAB:    true,
		LanguageTypeYAC:    true,
		LanguageTypeYAD:    true,
		LanguageTypeYAE:    true,
		LanguageTypeYAF:    true,
		LanguageTypeYAG:    true,
		LanguageTypeYAH:    true,
		LanguageTypeYAI:    true,
		LanguageTypeYAJ:    true,
		LanguageTypeYAK:    true,
		LanguageTypeYAL:    true,
		LanguageTypeYAM:    true,
		LanguageTypeYAN:    true,
		LanguageTypeYAO:    true,
		LanguageTypeYAP:    true,
		LanguageTypeYAQ:    true,
		LanguageTypeYAR:    true,
		LanguageTypeYAS:    true,
		LanguageTypeYAT:    true,
		LanguageTypeYAU:    true,
		LanguageTypeYAV:    true,
		LanguageTypeYAW:    true,
		LanguageTypeYAX:    true,
		LanguageTypeYAY:    true,
		LanguageTypeYAZ:    true,
		LanguageTypeYBA:    true,
		LanguageTypeYBB:    true,
		LanguageTypeYBD:    true,
		LanguageTypeYBE:    true,
		LanguageTypeYBH:    true,
		LanguageTypeYBI:    true,
		LanguageTypeYBJ:    true,
		LanguageTypeYBK:    true,
		LanguageTypeYBL:    true,
		LanguageTypeYBM:    true,
		LanguageTypeYBN:    true,
		LanguageTypeYBO:    true,
		LanguageTypeYBX:    true,
		LanguageTypeYBY:    true,
		LanguageTypeYCH:    true,
		LanguageTypeYCL:    true,
		LanguageTypeYCN:    true,
		LanguageTypeYCP:    true,
		LanguageTypeYDA:    true,
		LanguageTypeYDD:    true,
		LanguageTypeYDE:    true,
		LanguageTypeYDG:    true,
		LanguageTypeYDK:    true,
		LanguageTypeYDS:    true,
		LanguageTypeYEA:    true,
		LanguageTypeYEC:    true,
		LanguageTypeYEE:    true,
		LanguageTypeYEI:    true,
		LanguageTypeYEJ:    true,
		LanguageTypeYEL:    true,
		LanguageTypeYEN:    true,
		LanguageTypeYER:    true,
		LanguageTypeYES:    true,
		LanguageTypeYET:    true,
		LanguageTypeYEU:    true,
		LanguageTypeYEV:    true,
		LanguageTypeYEY:    true,
		LanguageTypeYGA:    true,
		LanguageTypeYGI:    true,
		LanguageTypeYGL:    true,
		LanguageTypeYGM:    true,
		LanguageTypeYGP:    true,
		LanguageTypeYGR:    true,
		LanguageTypeYGU:    true,
		LanguageTypeYGW:    true,
		LanguageTypeYHA:    true,
		LanguageTypeYHD:    true,
		LanguageTypeYHL:    true,
		LanguageTypeYIA:    true,
		LanguageTypeYIF:    true,
		LanguageTypeYIG:    true,
		LanguageTypeYIH:    true,
		LanguageTypeYII:    true,
		LanguageTypeYIJ:    true,
		LanguageTypeYIK:    true,
		LanguageTypeYIL:    true,
		LanguageTypeYIM:    true,
		LanguageTypeYIN:    true,
		LanguageTypeYIP:    true,
		LanguageTypeYIQ:    true,
		LanguageTypeYIR:    true,
		LanguageTypeYIS:    true,
		LanguageTypeYIT:    true,
		LanguageTypeYIU:    true,
		LanguageTypeYIV:    true,
		LanguageTypeYIX:    true,
		LanguageTypeYIY:    true,
		LanguageTypeYIZ:    true,
		LanguageTypeYKA:    true,
		LanguageTypeYKG:    true,
		LanguageTypeYKI:    true,
		LanguageTypeYKK:    true,
		LanguageTypeYKL:    true,
		LanguageTypeYKM:    true,
		LanguageTypeYKN:    true,
		LanguageTypeYKO:    true,
		LanguageTypeYKR:    true,
		LanguageTypeYKT:    true,
		LanguageTypeYKU:    true,
		LanguageTypeYKY:    true,
		LanguageTypeYLA:    true,
		LanguageTypeYLB:    true,
		LanguageTypeYLE:    true,
		LanguageTypeYLG:    true,
		LanguageTypeYLI:    true,
		LanguageTypeYLL:    true,
		LanguageTypeYLM:    true,
		LanguageTypeYLN:    true,
		LanguageTypeYLO:    true,
		LanguageTypeYLR:    true,
		LanguageTypeYLU:    true,
		LanguageTypeYLY:    true,
		LanguageTypeYMA:    true,
		LanguageTypeYMB:    true,
		LanguageTypeYMC:    true,
		LanguageTypeYMD:    true,
		LanguageTypeYME:    true,
		LanguageTypeYMG:    true,
		LanguageTypeYMH:    true,
		LanguageTypeYMI:    true,
		LanguageTypeYMK:    true,
		LanguageTypeYML:    true,
		LanguageTypeYMM:    true,
		LanguageTypeYMN:    true,
		LanguageTypeYMO:    true,
		LanguageTypeYMP:    true,
		LanguageTypeYMQ:    true,
		LanguageTypeYMR:    true,
		LanguageTypeYMS:    true,
		LanguageTypeYMT:    true,
		LanguageTypeYMX:    true,
		LanguageTypeYMZ:    true,
		LanguageTypeYNA:    true,
		LanguageTypeYND:    true,
		LanguageTypeYNE:    true,
		LanguageTypeYNG:    true,
		LanguageTypeYNH:    true,
		LanguageTypeYNK:    true,
		LanguageTypeYNL:    true,
		LanguageTypeYNN:    true,
		LanguageTypeYNO:    true,
		LanguageTypeYNQ:    true,
		LanguageTypeYNS:    true,
		LanguageTypeYNU:    true,
		LanguageTypeYOB:    true,
		LanguageTypeYOG:    true,
		LanguageTypeYOI:    true,
		LanguageTypeYOK:    true,
		LanguageTypeYOL:    true,
		LanguageTypeYOM:    true,
		LanguageTypeYON:    true,
		LanguageTypeYOS:    true,
		LanguageTypeYOT:    true,
		LanguageTypeYOX:    true,
		LanguageTypeYOY:    true,
		LanguageTypeYPA:    true,
		LanguageTypeYPB:    true,
		LanguageTypeYPG:    true,
		LanguageTypeYPH:    true,
		LanguageTypeYPK:    true,
		LanguageTypeYPM:    true,
		LanguageTypeYPN:    true,
		LanguageTypeYPO:    true,
		LanguageTypeYPP:    true,
		LanguageTypeYPZ:    true,
		LanguageTypeYRA:    true,
		LanguageTypeYRB:    true,
		LanguageTypeYRE:    true,
		LanguageTypeYRI:    true,
		LanguageTypeYRK:    true,
		LanguageTypeYRL:    true,
		LanguageTypeYRM:    true,
		LanguageTypeYRN:    true,
		LanguageTypeYRS:    true,
		LanguageTypeYRW:    true,
		LanguageTypeYRY:    true,
		LanguageTypeYSC:    true,
		LanguageTypeYSD:    true,
		LanguageTypeYSG:    true,
		LanguageTypeYSL:    true,
		LanguageTypeYSN:    true,
		LanguageTypeYSO:    true,
		LanguageTypeYSP:    true,
		LanguageTypeYSR:    true,
		LanguageTypeYSS:    true,
		LanguageTypeYSY:    true,
		LanguageTypeYTA:    true,
		LanguageTypeYTL:    true,
		LanguageTypeYTP:    true,
		LanguageTypeYTW:    true,
		LanguageTypeYTY:    true,
		LanguageTypeYUA:    true,
		LanguageTypeYUB:    true,
		LanguageTypeYUC:    true,
		LanguageTypeYUD:    true,
		LanguageTypeYUE:    true,
		LanguageTypeYUF:    true,
		LanguageTypeYUG:    true,
		LanguageTypeYUI:    true,
		LanguageTypeYUJ:    true,
		LanguageTypeYUK:    true,
		LanguageTypeYUL:    true,
		LanguageTypeYUM:    true,
		LanguageTypeYUN:    true,
		LanguageTypeYUP:    true,
		LanguageTypeYUQ:    true,
		LanguageTypeYUR:    true,
		LanguageTypeYUT:    true,
		LanguageTypeYUU:    true,
		LanguageTypeYUW:    true,
		LanguageTypeYUX:    true,
		LanguageTypeYUY:    true,
		LanguageTypeYUZ:    true,
		LanguageTypeYVA:    true,
		LanguageTypeYVT:    true,
		LanguageTypeYWA:    true,
		LanguageTypeYWG:    true,
		LanguageTypeYWL:    true,
		LanguageTypeYWN:    true,
		LanguageTypeYWQ:    true,
		LanguageTypeYWR:    true,
		LanguageTypeYWT:    true,
		LanguageTypeYWU:    true,
		LanguageTypeYWW:    true,
		LanguageTypeYXA:    true,
		LanguageTypeYXG:    true,
		LanguageTypeYXL:    true,
		LanguageTypeYXM:    true,
		LanguageTypeYXU:    true,
		LanguageTypeYXY:    true,
		LanguageTypeYYR:    true,
		LanguageTypeYYU:    true,
		LanguageTypeYYZ:    true,
		LanguageTypeYZG:    true,
		LanguageTypeYZK:    true,
		LanguageTypeZAA:    true,
		LanguageTypeZAB:    true,
		LanguageTypeZAC:    true,
		LanguageTypeZAD:    true,
		LanguageTypeZAE:    true,
		LanguageTypeZAF:    true,
		LanguageTypeZAG:    true,
		LanguageTypeZAH:    true,
		LanguageTypeZAI:    true,
		LanguageTypeZAJ:    true,
		LanguageTypeZAK:    true,
		LanguageTypeZAL:    true,
		LanguageTypeZAM:    true,
		LanguageTypeZAO:    true,
		LanguageTypeZAP:    true,
		LanguageTypeZAQ:    true,
		LanguageTypeZAR:    true,
		LanguageTypeZAS:    true,
		LanguageTypeZAT:    true,
		LanguageTypeZAU:    true,
		LanguageTypeZAV:    true,
		LanguageTypeZAW:    true,
		LanguageTypeZAX:    true,
		LanguageTypeZAY:    true,
		LanguageTypeZAZ:    true,
		LanguageTypeZBC:    true,
		LanguageTypeZBE:    true,
		LanguageTypeZBL:    true,
		LanguageTypeZBT:    true,
		LanguageTypeZBW:    true,
		LanguageTypeZCA:    true,
		LanguageTypeZCH:    true,
		LanguageTypeZDJ:    true,
		LanguageTypeZEA:    true,
		LanguageTypeZEG:    true,
		LanguageTypeZEH:    true,
		LanguageTypeZEN:    true,
		LanguageTypeZGA:    true,
		LanguageTypeZGB:    true,
		LanguageTypeZGH:    true,
		LanguageTypeZGM:    true,
		LanguageTypeZGN:    true,
		LanguageTypeZGR:    true,
		LanguageTypeZHB:    true,
		LanguageTypeZHD:    true,
		LanguageTypeZHI:    true,
		LanguageTypeZHN:    true,
		LanguageTypeZHW:    true,
		LanguageTypeZHX:    true,
		LanguageTypeZIA:    true,
		LanguageTypeZIB:    true,
		LanguageTypeZIK:    true,
		LanguageTypeZIL:    true,
		LanguageTypeZIM:    true,
		LanguageTypeZIN:    true,
		LanguageTypeZIR:    true,
		LanguageTypeZIW:    true,
		LanguageTypeZIZ:    true,
		LanguageTypeZKA:    true,
		LanguageTypeZKB:    true,
		LanguageTypeZKD:    true,
		LanguageTypeZKG:    true,
		LanguageTypeZKH:    true,
		LanguageTypeZKK:    true,
		LanguageTypeZKN:    true,
		LanguageTypeZKO:    true,
		LanguageTypeZKP:    true,
		LanguageTypeZKR:    true,
		LanguageTypeZKT:    true,
		LanguageTypeZKU:    true,
		LanguageTypeZKV:    true,
		LanguageTypeZKZ:    true,
		LanguageTypeZLE:    true,
		LanguageTypeZLJ:    true,
		LanguageTypeZLM:    true,
		LanguageTypeZLN:    true,
		LanguageTypeZLQ:    true,
		LanguageTypeZLS:    true,
		LanguageTypeZLW:    true,
		LanguageTypeZMA:    true,
		LanguageTypeZMB:    true,
		LanguageTypeZMC:    true,
		LanguageTypeZMD:    true,
		LanguageTypeZME:    true,
		LanguageTypeZMF:    true,
		LanguageTypeZMG:    true,
		LanguageTypeZMH:    true,
		LanguageTypeZMI:    true,
		LanguageTypeZMJ:    true,
		LanguageTypeZMK:    true,
		LanguageTypeZML:    true,
		LanguageTypeZMM:    true,
		LanguageTypeZMN:    true,
		LanguageTypeZMO:    true,
		LanguageTypeZMP:    true,
		LanguageTypeZMQ:    true,
		LanguageTypeZMR:    true,
		LanguageTypeZMS:    true,
		LanguageTypeZMT:    true,
		LanguageTypeZMU:    true,
		LanguageTypeZMV:    true,
		LanguageTypeZMW:    true,
		LanguageTypeZMX:    true,
		LanguageTypeZMY:    true,
		LanguageTypeZMZ:    true,
		LanguageTypeZNA:    true,
		LanguageTypeZND:    true,
		LanguageTypeZNE:    true,
		LanguageTypeZNG:    true,
		LanguageTypeZNK:    true,
		LanguageTypeZNS:    true,
		LanguageTypeZOC:    true,
		LanguageTypeZOH:    true,
		LanguageTypeZOM:    true,
		LanguageTypeZOO:    true,
		LanguageTypeZOQ:    true,
		LanguageTypeZOR:    true,
		LanguageTypeZOS:    true,
		LanguageTypeZPA:    true,
		LanguageTypeZPB:    true,
		LanguageTypeZPC:    true,
		LanguageTypeZPD:    true,
		LanguageTypeZPE:    true,
		LanguageTypeZPF:    true,
		LanguageTypeZPG:    true,
		LanguageTypeZPH:    true,
		LanguageTypeZPI:    true,
		LanguageTypeZPJ:    true,
		LanguageTypeZPK:    true,
		LanguageTypeZPL:    true,
		LanguageTypeZPM:    true,
		LanguageTypeZPN:    true,
		LanguageTypeZPO:    true,
		LanguageTypeZPP:    true,
		LanguageTypeZPQ:    true,
		LanguageTypeZPR:    true,
		LanguageTypeZPS:    true,
		LanguageTypeZPT:    true,
		LanguageTypeZPU:    true,
		LanguageTypeZPV:    true,
		LanguageTypeZPW:    true,
		LanguageTypeZPX:    true,
		LanguageTypeZPY:    true,
		LanguageTypeZPZ:    true,
		LanguageTypeZQE:    true,
		LanguageTypeZRA:    true,
		LanguageTypeZRG:    true,
		LanguageTypeZRN:    true,
		LanguageTypeZRO:    true,
		LanguageTypeZRP:    true,
		LanguageTypeZRS:    true,
		LanguageTypeZSA:    true,
		LanguageTypeZSK:    true,
		LanguageTypeZSL:    true,
		LanguageTypeZSM:    true,
		LanguageTypeZSR:    true,
		LanguageTypeZSU:    true,
		LanguageTypeZTE:    true,
		LanguageTypeZTG:    true,
		LanguageTypeZTL:    true,
		LanguageTypeZTM:    true,
		LanguageTypeZTN:    true,
		LanguageTypeZTP:    true,
		LanguageTypeZTQ:    true,
		LanguageTypeZTS:    true,
		LanguageTypeZTT:    true,
		LanguageTypeZTU:    true,
		LanguageTypeZTX:    true,
		LanguageTypeZTY:    true,
		LanguageTypeZUA:    true,
		LanguageTypeZUH:    true,
		LanguageTypeZUM:    true,
		LanguageTypeZUN:    true,
		LanguageTypeZUY:    true,
		LanguageTypeZWA:    true,
		LanguageTypeZXX:    true,
		LanguageTypeZYB:    true,
		LanguageTypeZYG:    true,
		LanguageTypeZYJ:    true,
		LanguageTypeZYN:    true,
		LanguageTypeZYP:    true,
		LanguageTypeZZA:    true,
		LanguageTypeZZJ:    true,
	}
)

// Used to verify if a language is valid or not. Don't need lock because we have many
// readers and no writers. Rob Pike sad that it's ok
// (https://groups.google.com/forum/#!msg/golang-nuts/HpLWnGTp-n8/hyUYmnWJqiQJ)
func LanguageTypeExists(languageType string) bool {
	// Normalize input
	languageType = strings.ToLower(languageType)
	languageType = strings.TrimSpace(languageType)

	_, ok := languageTypes[LanguageType(languageType)]
	return ok
}

// List of possible region types, that are a subcategory of the language
const (
	RegionTypeAA   RegionType = "AA"     // Private use
	RegionTypeAC   RegionType = "AC"     // Ascension Island
	RegionTypeAD   RegionType = "AD"     // Andorra
	RegionTypeAE   RegionType = "AE"     // United Arab Emirates
	RegionTypeAF   RegionType = "AF"     // Afghanistan
	RegionTypeAG   RegionType = "AG"     // Antigua and Barbuda
	RegionTypeAI   RegionType = "AI"     // Anguilla
	RegionTypeAL   RegionType = "AL"     // Albania
	RegionTypeAM   RegionType = "AM"     // Armenia
	RegionTypeAN   RegionType = "AN"     // Netherlands Antilles
	RegionTypeAO   RegionType = "AO"     // Angola
	RegionTypeAQ   RegionType = "AQ"     // Antarctica
	RegionTypeAR   RegionType = "AR"     // Argentina
	RegionTypeAS   RegionType = "AS"     // American Samoa
	RegionTypeAT   RegionType = "AT"     // Austria
	RegionTypeAU   RegionType = "AU"     // Australia
	RegionTypeAW   RegionType = "AW"     // Aruba
	RegionTypeAX   RegionType = "AX"     // Åland Islands
	RegionTypeAZ   RegionType = "AZ"     // Azerbaijan
	RegionTypeBA   RegionType = "BA"     // Bosnia and Herzegovina
	RegionTypeBB   RegionType = "BB"     // Barbados
	RegionTypeBD   RegionType = "BD"     // Bangladesh
	RegionTypeBE   RegionType = "BE"     // Belgium
	RegionTypeBF   RegionType = "BF"     // Burkina Faso
	RegionTypeBG   RegionType = "BG"     // Bulgaria
	RegionTypeBH   RegionType = "BH"     // Bahrain
	RegionTypeBI   RegionType = "BI"     // Burundi
	RegionTypeBJ   RegionType = "BJ"     // Benin
	RegionTypeBL   RegionType = "BL"     // Saint Barthélemy
	RegionTypeBM   RegionType = "BM"     // Bermuda
	RegionTypeBN   RegionType = "BN"     // Brunei Darussalam
	RegionTypeBO   RegionType = "BO"     // Bolivia
	RegionTypeBQ   RegionType = "BQ"     // Bonaire, Sint Eustatius and Saba
	RegionTypeBR   RegionType = "BR"     // Brazil
	RegionTypeBS   RegionType = "BS"     // Bahamas
	RegionTypeBT   RegionType = "BT"     // Bhutan
	RegionTypeBU   RegionType = "BU"     // Burma
	RegionTypeBV   RegionType = "BV"     // Bouvet Island
	RegionTypeBW   RegionType = "BW"     // Botswana
	RegionTypeBY   RegionType = "BY"     // Belarus
	RegionTypeBZ   RegionType = "BZ"     // Belize
	RegionTypeCA   RegionType = "CA"     // Canada
	RegionTypeCC   RegionType = "CC"     // Cocos (Keeling) Islands
	RegionTypeCD   RegionType = "CD"     // The Democratic Republic of the Congo
	RegionTypeCF   RegionType = "CF"     // Central African Republic
	RegionTypeCG   RegionType = "CG"     // Congo
	RegionTypeCH   RegionType = "CH"     // Switzerland
	RegionTypeCI   RegionType = "CI"     // Côte d'Ivoire
	RegionTypeCK   RegionType = "CK"     // Cook Islands
	RegionTypeCL   RegionType = "CL"     // Chile
	RegionTypeCM   RegionType = "CM"     // Cameroon
	RegionTypeCN   RegionType = "CN"     // China
	RegionTypeCO   RegionType = "CO"     // Colombia
	RegionTypeCP   RegionType = "CP"     // Clipperton Island
	RegionTypeCR   RegionType = "CR"     // Costa Rica
	RegionTypeCS   RegionType = "CS"     // Serbia and Montenegro
	RegionTypeCU   RegionType = "CU"     // Cuba
	RegionTypeCV   RegionType = "CV"     // Cape Verde
	RegionTypeCW   RegionType = "CW"     // Curaçao
	RegionTypeCX   RegionType = "CX"     // Christmas Island
	RegionTypeCY   RegionType = "CY"     // Cyprus
	RegionTypeCZ   RegionType = "CZ"     // Czech Republic
	RegionTypeDD   RegionType = "DD"     // German Democratic Republic
	RegionTypeDE   RegionType = "DE"     // Germany
	RegionTypeDG   RegionType = "DG"     // Diego Garcia
	RegionTypeDJ   RegionType = "DJ"     // Djibouti
	RegionTypeDK   RegionType = "DK"     // Denmark
	RegionTypeDM   RegionType = "DM"     // Dominica
	RegionTypeDO   RegionType = "DO"     // Dominican Republic
	RegionTypeDZ   RegionType = "DZ"     // Algeria
	RegionTypeEA   RegionType = "EA"     // Ceuta, Melilla
	RegionTypeEC   RegionType = "EC"     // Ecuador
	RegionTypeEE   RegionType = "EE"     // Estonia
	RegionTypeEG   RegionType = "EG"     // Egypt
	RegionTypeEH   RegionType = "EH"     // Western Sahara
	RegionTypeER   RegionType = "ER"     // Eritrea
	RegionTypeES   RegionType = "ES"     // Spain
	RegionTypeET   RegionType = "ET"     // Ethiopia
	RegionTypeEU   RegionType = "EU"     // European Union
	RegionTypeFI   RegionType = "FI"     // Finland
	RegionTypeFJ   RegionType = "FJ"     // Fiji
	RegionTypeFK   RegionType = "FK"     // Falkland Islands (Malvinas)
	RegionTypeFM   RegionType = "FM"     // Federated States of Micronesia
	RegionTypeFO   RegionType = "FO"     // Faroe Islands
	RegionTypeFR   RegionType = "FR"     // France
	RegionTypeFX   RegionType = "FX"     // Metropolitan France
	RegionTypeGA   RegionType = "GA"     // Gabon
	RegionTypeGB   RegionType = "GB"     // United Kingdom
	RegionTypeGD   RegionType = "GD"     // Grenada
	RegionTypeGE   RegionType = "GE"     // Georgia
	RegionTypeGF   RegionType = "GF"     // French Guiana
	RegionTypeGG   RegionType = "GG"     // Guernsey
	RegionTypeGH   RegionType = "GH"     // Ghana
	RegionTypeGI   RegionType = "GI"     // Gibraltar
	RegionTypeGL   RegionType = "GL"     // Greenland
	RegionTypeGM   RegionType = "GM"     // Gambia
	RegionTypeGN   RegionType = "GN"     // Guinea
	RegionTypeGP   RegionType = "GP"     // Guadeloupe
	RegionTypeGQ   RegionType = "GQ"     // Equatorial Guinea
	RegionTypeGR   RegionType = "GR"     // Greece
	RegionTypeGS   RegionType = "GS"     // South Georgia and the South Sandwich Islands
	RegionTypeGT   RegionType = "GT"     // Guatemala
	RegionTypeGU   RegionType = "GU"     // Guam
	RegionTypeGW   RegionType = "GW"     // Guinea-Bissau
	RegionTypeGY   RegionType = "GY"     // Guyana
	RegionTypeHK   RegionType = "HK"     // Hong Kong
	RegionTypeHM   RegionType = "HM"     // Heard Island and McDonald Islands
	RegionTypeHN   RegionType = "HN"     // Honduras
	RegionTypeHR   RegionType = "HR"     // Croatia
	RegionTypeHT   RegionType = "HT"     // Haiti
	RegionTypeHU   RegionType = "HU"     // Hungary
	RegionTypeIC   RegionType = "IC"     // Canary Islands
	RegionTypeID   RegionType = "ID"     // Indonesia
	RegionTypeIE   RegionType = "IE"     // Ireland
	RegionTypeIL   RegionType = "IL"     // Israel
	RegionTypeIM   RegionType = "IM"     // Isle of Man
	RegionTypeIN   RegionType = "IN"     // India
	RegionTypeIO   RegionType = "IO"     // British Indian Ocean Territory
	RegionTypeIQ   RegionType = "IQ"     // Iraq
	RegionTypeIR   RegionType = "IR"     // Islamic Republic of Iran
	RegionTypeIS   RegionType = "IS"     // Iceland
	RegionTypeIT   RegionType = "IT"     // Italy
	RegionTypeJE   RegionType = "JE"     // Jersey
	RegionTypeJM   RegionType = "JM"     // Jamaica
	RegionTypeJO   RegionType = "JO"     // Jordan
	RegionTypeJP   RegionType = "JP"     // Japan
	RegionTypeKE   RegionType = "KE"     // Kenya
	RegionTypeKG   RegionType = "KG"     // Kyrgyzstan
	RegionTypeKH   RegionType = "KH"     // Cambodia
	RegionTypeKI   RegionType = "KI"     // Kiribati
	RegionTypeKM   RegionType = "KM"     // Comoros
	RegionTypeKN   RegionType = "KN"     // Saint Kitts and Nevis
	RegionTypeKP   RegionType = "KP"     // Democratic People's Republic of Korea
	RegionTypeKR   RegionType = "KR"     // Republic of Korea
	RegionTypeKW   RegionType = "KW"     // Kuwait
	RegionTypeKY   RegionType = "KY"     // Cayman Islands
	RegionTypeKZ   RegionType = "KZ"     // Kazakhstan
	RegionTypeLA   RegionType = "LA"     // Lao People's Democratic Republic
	RegionTypeLB   RegionType = "LB"     // Lebanon
	RegionTypeLC   RegionType = "LC"     // Saint Lucia
	RegionTypeLI   RegionType = "LI"     // Liechtenstein
	RegionTypeLK   RegionType = "LK"     // Sri Lanka
	RegionTypeLR   RegionType = "LR"     // Liberia
	RegionTypeLS   RegionType = "LS"     // Lesotho
	RegionTypeLT   RegionType = "LT"     // Lithuania
	RegionTypeLU   RegionType = "LU"     // Luxembourg
	RegionTypeLV   RegionType = "LV"     // Latvia
	RegionTypeLY   RegionType = "LY"     // Libya
	RegionTypeMA   RegionType = "MA"     // Morocco
	RegionTypeMC   RegionType = "MC"     // Monaco
	RegionTypeMD   RegionType = "MD"     // Moldova
	RegionTypeME   RegionType = "ME"     // Montenegro
	RegionTypeMF   RegionType = "MF"     // Saint Martin (French part)
	RegionTypeMG   RegionType = "MG"     // Madagascar
	RegionTypeMH   RegionType = "MH"     // Marshall Islands
	RegionTypeMK   RegionType = "MK"     // The Former Yugoslav Republic of Macedonia
	RegionTypeML   RegionType = "ML"     // Mali
	RegionTypeMM   RegionType = "MM"     // Myanmar
	RegionTypeMN   RegionType = "MN"     // Mongolia
	RegionTypeMO   RegionType = "MO"     // Macao
	RegionTypeMP   RegionType = "MP"     // Northern Mariana Islands
	RegionTypeMQ   RegionType = "MQ"     // Martinique
	RegionTypeMR   RegionType = "MR"     // Mauritania
	RegionTypeMS   RegionType = "MS"     // Montserrat
	RegionTypeMT   RegionType = "MT"     // Malta
	RegionTypeMU   RegionType = "MU"     // Mauritius
	RegionTypeMV   RegionType = "MV"     // Maldives
	RegionTypeMW   RegionType = "MW"     // Malawi
	RegionTypeMX   RegionType = "MX"     // Mexico
	RegionTypeMY   RegionType = "MY"     // Malaysia
	RegionTypeMZ   RegionType = "MZ"     // Mozambique
	RegionTypeNA   RegionType = "NA"     // Namibia
	RegionTypeNC   RegionType = "NC"     // New Caledonia
	RegionTypeNE   RegionType = "NE"     // Niger
	RegionTypeNF   RegionType = "NF"     // Norfolk Island
	RegionTypeNG   RegionType = "NG"     // Nigeria
	RegionTypeNI   RegionType = "NI"     // Nicaragua
	RegionTypeNL   RegionType = "NL"     // Netherlands
	RegionTypeNO   RegionType = "NO"     // Norway
	RegionTypeNP   RegionType = "NP"     // Nepal
	RegionTypeNR   RegionType = "NR"     // Nauru
	RegionTypeNT   RegionType = "NT"     // Neutral Zone
	RegionTypeNU   RegionType = "NU"     // Niue
	RegionTypeNZ   RegionType = "NZ"     // New Zealand
	RegionTypeOM   RegionType = "OM"     // Oman
	RegionTypePA   RegionType = "PA"     // Panama
	RegionTypePE   RegionType = "PE"     // Peru
	RegionTypePF   RegionType = "PF"     // French Polynesia
	RegionTypePG   RegionType = "PG"     // Papua New Guinea
	RegionTypePH   RegionType = "PH"     // Philippines
	RegionTypePK   RegionType = "PK"     // Pakistan
	RegionTypePL   RegionType = "PL"     // Poland
	RegionTypePM   RegionType = "PM"     // Saint Pierre and Miquelon
	RegionTypePN   RegionType = "PN"     // Pitcairn
	RegionTypePR   RegionType = "PR"     // Puerto Rico
	RegionTypePS   RegionType = "PS"     // State of Palestine
	RegionTypePT   RegionType = "PT"     // Portugal
	RegionTypePW   RegionType = "PW"     // Palau
	RegionTypePY   RegionType = "PY"     // Paraguay
	RegionTypeQA   RegionType = "QA"     // Qatar
	RegionTypeQMQZ RegionType = "QM..QZ" // Private use
	RegionTypeRE   RegionType = "RE"     // Réunion
	RegionTypeRO   RegionType = "RO"     // Romania
	RegionTypeRS   RegionType = "RS"     // Serbia
	RegionTypeRU   RegionType = "RU"     // Russian Federation
	RegionTypeRW   RegionType = "RW"     // Rwanda
	RegionTypeSA   RegionType = "SA"     // Saudi Arabia
	RegionTypeSB   RegionType = "SB"     // Solomon Islands
	RegionTypeSC   RegionType = "SC"     // Seychelles
	RegionTypeSD   RegionType = "SD"     // Sudan
	RegionTypeSE   RegionType = "SE"     // Sweden
	RegionTypeSG   RegionType = "SG"     // Singapore
	RegionTypeSH   RegionType = "SH"     // Saint Helena, Ascension and Tristan da Cunha
	RegionTypeSI   RegionType = "SI"     // Slovenia
	RegionTypeSJ   RegionType = "SJ"     // Svalbard and Jan Mayen
	RegionTypeSK   RegionType = "SK"     // Slovakia
	RegionTypeSL   RegionType = "SL"     // Sierra Leone
	RegionTypeSM   RegionType = "SM"     // San Marino
	RegionTypeSN   RegionType = "SN"     // Senegal
	RegionTypeSO   RegionType = "SO"     // Somalia
	RegionTypeSR   RegionType = "SR"     // Suriname
	RegionTypeSS   RegionType = "SS"     // South Sudan
	RegionTypeST   RegionType = "ST"     // Sao Tome and Principe
	RegionTypeSU   RegionType = "SU"     // Union of Soviet Socialist Republics
	RegionTypeSV   RegionType = "SV"     // El Salvador
	RegionTypeSX   RegionType = "SX"     // Sint Maarten (Dutch part)
	RegionTypeSY   RegionType = "SY"     // Syrian Arab Republic
	RegionTypeSZ   RegionType = "SZ"     // Swaziland
	RegionTypeTA   RegionType = "TA"     // Tristan da Cunha
	RegionTypeTC   RegionType = "TC"     // Turks and Caicos Islands
	RegionTypeTD   RegionType = "TD"     // Chad
	RegionTypeTF   RegionType = "TF"     // French Southern Territories
	RegionTypeTG   RegionType = "TG"     // Togo
	RegionTypeTH   RegionType = "TH"     // Thailand
	RegionTypeTJ   RegionType = "TJ"     // Tajikistan
	RegionTypeTK   RegionType = "TK"     // Tokelau
	RegionTypeTL   RegionType = "TL"     // Timor-Leste
	RegionTypeTM   RegionType = "TM"     // Turkmenistan
	RegionTypeTN   RegionType = "TN"     // Tunisia
	RegionTypeTO   RegionType = "TO"     // Tonga
	RegionTypeTP   RegionType = "TP"     // East Timor
	RegionTypeTR   RegionType = "TR"     // Turkey
	RegionTypeTT   RegionType = "TT"     // Trinidad and Tobago
	RegionTypeTV   RegionType = "TV"     // Tuvalu
	RegionTypeTW   RegionType = "TW"     // Taiwan, Province of China
	RegionTypeTZ   RegionType = "TZ"     // United Republic of Tanzania
	RegionTypeUA   RegionType = "UA"     // Ukraine
	RegionTypeUG   RegionType = "UG"     // Uganda
	RegionTypeUM   RegionType = "UM"     // United States Minor Outlying Islands
	RegionTypeUS   RegionType = "US"     // United States
	RegionTypeUY   RegionType = "UY"     // Uruguay
	RegionTypeUZ   RegionType = "UZ"     // Uzbekistan
	RegionTypeVA   RegionType = "VA"     // Holy See (Vatican City State)
	RegionTypeVC   RegionType = "VC"     // Saint Vincent and the Grenadines
	RegionTypeVE   RegionType = "VE"     // Venezuela
	RegionTypeVG   RegionType = "VG"     // British Virgin Islands
	RegionTypeVI   RegionType = "VI"     // U.S. Virgin Islands
	RegionTypeVN   RegionType = "VN"     // Viet Nam
	RegionTypeVU   RegionType = "VU"     // Vanuatu
	RegionTypeWF   RegionType = "WF"     // Wallis and Futuna
	RegionTypeWS   RegionType = "WS"     // Samoa
	RegionTypeXAXZ RegionType = "XA..XZ" // Private use
	RegionTypeYD   RegionType = "YD"     // Democratic Yemen
	RegionTypeYE   RegionType = "YE"     // Yemen
	RegionTypeYT   RegionType = "YT"     // Mayotte
	RegionTypeYU   RegionType = "YU"     // Yugoslavia
	RegionTypeZA   RegionType = "ZA"     // South Africa
	RegionTypeZM   RegionType = "ZM"     // Zambia
	RegionTypeZR   RegionType = "ZR"     // Zaire
	RegionTypeZW   RegionType = "ZW"     // Zimbabwe
	RegionTypeZZ   RegionType = "ZZ"     // Private use
	RegionType001  RegionType = "001"    // World
	RegionType002  RegionType = "002"    // Africa
	RegionType003  RegionType = "003"    // North America
	RegionType005  RegionType = "005"    // South America
	RegionType009  RegionType = "009"    // Oceania
	RegionType011  RegionType = "011"    // Western Africa
	RegionType013  RegionType = "013"    // Central America
	RegionType014  RegionType = "014"    // Eastern Africa
	RegionType015  RegionType = "015"    // Northern Africa
	RegionType017  RegionType = "017"    // Middle Africa
	RegionType018  RegionType = "018"    // Southern Africa
	RegionType019  RegionType = "019"    // Americas
	RegionType021  RegionType = "021"    // Northern America
	RegionType029  RegionType = "029"    // Caribbean
	RegionType030  RegionType = "030"    // Eastern Asia
	RegionType034  RegionType = "034"    // Southern Asia
	RegionType035  RegionType = "035"    // South-Eastern Asia
	RegionType039  RegionType = "039"    // Southern Europe
	RegionType053  RegionType = "053"    // Australia and New Zealand
	RegionType054  RegionType = "054"    // Melanesia
	RegionType057  RegionType = "057"    // Micronesia
	RegionType061  RegionType = "061"    // Polynesia
	RegionType142  RegionType = "142"    // Asia
	RegionType143  RegionType = "143"    // Central Asia
	RegionType145  RegionType = "145"    // Western Asia
	RegionType150  RegionType = "150"    // Europe
	RegionType151  RegionType = "151"    // Eastern Europe
	RegionType154  RegionType = "154"    // Northern Europe
	RegionType155  RegionType = "155"    // Western Europe
	RegionType419  RegionType = "419"    // Latin America and the Caribbean
)

// Used to determinate from a given language, what is the specific dialect (region) that we will
// assume
type RegionType string

// Structure used to identify if a region exists or not
var (
	regionTypes map[RegionType]bool = map[RegionType]bool{
		RegionTypeAA:   true,
		RegionTypeAC:   true,
		RegionTypeAD:   true,
		RegionTypeAE:   true,
		RegionTypeAF:   true,
		RegionTypeAG:   true,
		RegionTypeAI:   true,
		RegionTypeAL:   true,
		RegionTypeAM:   true,
		RegionTypeAN:   true,
		RegionTypeAO:   true,
		RegionTypeAQ:   true,
		RegionTypeAR:   true,
		RegionTypeAS:   true,
		RegionTypeAT:   true,
		RegionTypeAU:   true,
		RegionTypeAW:   true,
		RegionTypeAX:   true,
		RegionTypeAZ:   true,
		RegionTypeBA:   true,
		RegionTypeBB:   true,
		RegionTypeBD:   true,
		RegionTypeBE:   true,
		RegionTypeBF:   true,
		RegionTypeBG:   true,
		RegionTypeBH:   true,
		RegionTypeBI:   true,
		RegionTypeBJ:   true,
		RegionTypeBL:   true,
		RegionTypeBM:   true,
		RegionTypeBN:   true,
		RegionTypeBO:   true,
		RegionTypeBQ:   true,
		RegionTypeBR:   true,
		RegionTypeBS:   true,
		RegionTypeBT:   true,
		RegionTypeBU:   true,
		RegionTypeBV:   true,
		RegionTypeBW:   true,
		RegionTypeBY:   true,
		RegionTypeBZ:   true,
		RegionTypeCA:   true,
		RegionTypeCC:   true,
		RegionTypeCD:   true,
		RegionTypeCF:   true,
		RegionTypeCG:   true,
		RegionTypeCH:   true,
		RegionTypeCI:   true,
		RegionTypeCK:   true,
		RegionTypeCL:   true,
		RegionTypeCM:   true,
		RegionTypeCN:   true,
		RegionTypeCO:   true,
		RegionTypeCP:   true,
		RegionTypeCR:   true,
		RegionTypeCS:   true,
		RegionTypeCU:   true,
		RegionTypeCV:   true,
		RegionTypeCW:   true,
		RegionTypeCX:   true,
		RegionTypeCY:   true,
		RegionTypeCZ:   true,
		RegionTypeDD:   true,
		RegionTypeDE:   true,
		RegionTypeDG:   true,
		RegionTypeDJ:   true,
		RegionTypeDK:   true,
		RegionTypeDM:   true,
		RegionTypeDO:   true,
		RegionTypeDZ:   true,
		RegionTypeEA:   true,
		RegionTypeEC:   true,
		RegionTypeEE:   true,
		RegionTypeEG:   true,
		RegionTypeEH:   true,
		RegionTypeER:   true,
		RegionTypeES:   true,
		RegionTypeET:   true,
		RegionTypeEU:   true,
		RegionTypeFI:   true,
		RegionTypeFJ:   true,
		RegionTypeFK:   true,
		RegionTypeFM:   true,
		RegionTypeFO:   true,
		RegionTypeFR:   true,
		RegionTypeFX:   true,
		RegionTypeGA:   true,
		RegionTypeGB:   true,
		RegionTypeGD:   true,
		RegionTypeGE:   true,
		RegionTypeGF:   true,
		RegionTypeGG:   true,
		RegionTypeGH:   true,
		RegionTypeGI:   true,
		RegionTypeGL:   true,
		RegionTypeGM:   true,
		RegionTypeGN:   true,
		RegionTypeGP:   true,
		RegionTypeGQ:   true,
		RegionTypeGR:   true,
		RegionTypeGS:   true,
		RegionTypeGT:   true,
		RegionTypeGU:   true,
		RegionTypeGW:   true,
		RegionTypeGY:   true,
		RegionTypeHK:   true,
		RegionTypeHM:   true,
		RegionTypeHN:   true,
		RegionTypeHR:   true,
		RegionTypeHT:   true,
		RegionTypeHU:   true,
		RegionTypeIC:   true,
		RegionTypeID:   true,
		RegionTypeIE:   true,
		RegionTypeIL:   true,
		RegionTypeIM:   true,
		RegionTypeIN:   true,
		RegionTypeIO:   true,
		RegionTypeIQ:   true,
		RegionTypeIR:   true,
		RegionTypeIS:   true,
		RegionTypeIT:   true,
		RegionTypeJE:   true,
		RegionTypeJM:   true,
		RegionTypeJO:   true,
		RegionTypeJP:   true,
		RegionTypeKE:   true,
		RegionTypeKG:   true,
		RegionTypeKH:   true,
		RegionTypeKI:   true,
		RegionTypeKM:   true,
		RegionTypeKN:   true,
		RegionTypeKP:   true,
		RegionTypeKR:   true,
		RegionTypeKW:   true,
		RegionTypeKY:   true,
		RegionTypeKZ:   true,
		RegionTypeLA:   true,
		RegionTypeLB:   true,
		RegionTypeLC:   true,
		RegionTypeLI:   true,
		RegionTypeLK:   true,
		RegionTypeLR:   true,
		RegionTypeLS:   true,
		RegionTypeLT:   true,
		RegionTypeLU:   true,
		RegionTypeLV:   true,
		RegionTypeLY:   true,
		RegionTypeMA:   true,
		RegionTypeMC:   true,
		RegionTypeMD:   true,
		RegionTypeME:   true,
		RegionTypeMF:   true,
		RegionTypeMG:   true,
		RegionTypeMH:   true,
		RegionTypeMK:   true,
		RegionTypeML:   true,
		RegionTypeMM:   true,
		RegionTypeMN:   true,
		RegionTypeMO:   true,
		RegionTypeMP:   true,
		RegionTypeMQ:   true,
		RegionTypeMR:   true,
		RegionTypeMS:   true,
		RegionTypeMT:   true,
		RegionTypeMU:   true,
		RegionTypeMV:   true,
		RegionTypeMW:   true,
		RegionTypeMX:   true,
		RegionTypeMY:   true,
		RegionTypeMZ:   true,
		RegionTypeNA:   true,
		RegionTypeNC:   true,
		RegionTypeNE:   true,
		RegionTypeNF:   true,
		RegionTypeNG:   true,
		RegionTypeNI:   true,
		RegionTypeNL:   true,
		RegionTypeNO:   true,
		RegionTypeNP:   true,
		RegionTypeNR:   true,
		RegionTypeNT:   true,
		RegionTypeNU:   true,
		RegionTypeNZ:   true,
		RegionTypeOM:   true,
		RegionTypePA:   true,
		RegionTypePE:   true,
		RegionTypePF:   true,
		RegionTypePG:   true,
		RegionTypePH:   true,
		RegionTypePK:   true,
		RegionTypePL:   true,
		RegionTypePM:   true,
		RegionTypePN:   true,
		RegionTypePR:   true,
		RegionTypePS:   true,
		RegionTypePT:   true,
		RegionTypePW:   true,
		RegionTypePY:   true,
		RegionTypeQA:   true,
		RegionTypeQMQZ: true,
		RegionTypeRE:   true,
		RegionTypeRO:   true,
		RegionTypeRS:   true,
		RegionTypeRU:   true,
		RegionTypeRW:   true,
		RegionTypeSA:   true,
		RegionTypeSB:   true,
		RegionTypeSC:   true,
		RegionTypeSD:   true,
		RegionTypeSE:   true,
		RegionTypeSG:   true,
		RegionTypeSH:   true,
		RegionTypeSI:   true,
		RegionTypeSJ:   true,
		RegionTypeSK:   true,
		RegionTypeSL:   true,
		RegionTypeSM:   true,
		RegionTypeSN:   true,
		RegionTypeSO:   true,
		RegionTypeSR:   true,
		RegionTypeSS:   true,
		RegionTypeST:   true,
		RegionTypeSU:   true,
		RegionTypeSV:   true,
		RegionTypeSX:   true,
		RegionTypeSY:   true,
		RegionTypeSZ:   true,
		RegionTypeTA:   true,
		RegionTypeTC:   true,
		RegionTypeTD:   true,
		RegionTypeTF:   true,
		RegionTypeTG:   true,
		RegionTypeTH:   true,
		RegionTypeTJ:   true,
		RegionTypeTK:   true,
		RegionTypeTL:   true,
		RegionTypeTM:   true,
		RegionTypeTN:   true,
		RegionTypeTO:   true,
		RegionTypeTP:   true,
		RegionTypeTR:   true,
		RegionTypeTT:   true,
		RegionTypeTV:   true,
		RegionTypeTW:   true,
		RegionTypeTZ:   true,
		RegionTypeUA:   true,
		RegionTypeUG:   true,
		RegionTypeUM:   true,
		RegionTypeUS:   true,
		RegionTypeUY:   true,
		RegionTypeUZ:   true,
		RegionTypeVA:   true,
		RegionTypeVC:   true,
		RegionTypeVE:   true,
		RegionTypeVG:   true,
		RegionTypeVI:   true,
		RegionTypeVN:   true,
		RegionTypeVU:   true,
		RegionTypeWF:   true,
		RegionTypeWS:   true,
		RegionTypeXAXZ: true,
		RegionTypeYD:   true,
		RegionTypeYE:   true,
		RegionTypeYT:   true,
		RegionTypeYU:   true,
		RegionTypeZA:   true,
		RegionTypeZM:   true,
		RegionTypeZR:   true,
		RegionTypeZW:   true,
		RegionTypeZZ:   true,
		RegionType001:  true,
		RegionType002:  true,
		RegionType003:  true,
		RegionType005:  true,
		RegionType009:  true,
		RegionType011:  true,
		RegionType013:  true,
		RegionType014:  true,
		RegionType015:  true,
		RegionType017:  true,
		RegionType018:  true,
		RegionType019:  true,
		RegionType021:  true,
		RegionType029:  true,
		RegionType030:  true,
		RegionType034:  true,
		RegionType035:  true,
		RegionType039:  true,
		RegionType053:  true,
		RegionType054:  true,
		RegionType057:  true,
		RegionType061:  true,
		RegionType142:  true,
		RegionType143:  true,
		RegionType145:  true,
		RegionType150:  true,
		RegionType151:  true,
		RegionType154:  true,
		RegionType155:  true,
		RegionType419:  true,
	}
)

// Used to verify if a region is valid or not. Don't need lock because we have many
// readers and no writers. Rob Pike sad that it's ok
// (https://groups.google.com/forum/#!msg/golang-nuts/HpLWnGTp-n8/hyUYmnWJqiQJ)
func RegionTypeExists(regionType string) bool {
	// Normalize input
	regionType = strings.ToUpper(regionType)
	regionType = strings.TrimSpace(regionType)

	_, ok := regionTypes[RegionType(regionType)]
	return ok
}

// Useful function to check if a language with region or not is valid
func LanguageIsValid(language string) bool {
	if LanguageTypeExists(language) {
		return true
	}

	languageParts := strings.Split(language, "-")
	if len(languageParts) != 2 {
		return false
	}

	return LanguageTypeExists(languageParts[0]) && RegionTypeExists(languageParts[1])
}
