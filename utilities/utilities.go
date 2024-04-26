package NamiUtilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	math "math/rand"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func OverwriteImplantConfig(ipaddr string, port int, arch int, debug bool, name string, implantName string) {
	f, err := os.Create(filepath.FromSlash("implant/config.implant"))
	if err != nil {
		fmt.Println(err)
	}
	key, err := os.ReadFile(filepath.FromSlash("_keys/public.pem"))
	if err != nil {
		fmt.Println(err)
	}
	formattedString := fmt.Sprintf("{\"IP_ADDRESS\": \"%s\", \"PORT\": \"%v\", \"ARCH\": \"%v\", \"DEBUG\": \"%t\", \"NAME\": \"%s\", \"REGNAME\": \"%s\", \"KEY\": \"%x\"}", ipaddr, port, arch, debug, implantName, name, key)
	f.WriteString(formattedString)
}

func GenerateRandomImplantName() string {
	math.Seed(time.Now().Unix())
	implantName := implantNames[math.Intn(len(implantNames))]
	return implantName
}

func GenerateRSAKeys() {
	if _, err := os.Stat("_keys"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("_keys", os.ModePerm)
	}
	reader := rand.Reader
	bitSize := 4096

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		fmt.Println("\n\t[-] Unable to generate RSA key!")
	}
	publicKey := key.PublicKey

	pubFile, err := os.Create(filepath.FromSlash("_keys/public.pem"))
	if err != nil {
		fmt.Println("\n\t[-] Unable to create public RSA key file!")
	}
	defer pubFile.Close()

	privFile, err := os.Create(filepath.FromSlash("_keys/private.pem"))
	if err != nil {
		fmt.Println("\n\t[-] Unable to create private RSA key file!")
	}
	defer privFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	asn1Bytes, err := asn1.Marshal(publicKey)
	if err != nil {
		fmt.Println("\n\t[-] Unable to encode public key to ASN.1!")
	}
	var publicPem = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	err = pem.Encode(privFile, privateKey)
	if err != nil {
		fmt.Println("\n\t[-] Unable to write RSA private key to file!")
	}

	err = pem.Encode(pubFile, publicPem)
	if err != nil {
		fmt.Println("\n\t[-] Unable to write RSA public PEM to file!")
	}
}

func DecryptOAEP(cipherText string, privKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("Nami")

	rng := rand.Reader
	plaintext, _ := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	return string(plaintext)
}

func DecryptAES(passphrase, ciphertext string) []byte {
	salt, _ := hex.DecodeString(ciphertext[0:24])
	iv, _ := hex.DecodeString(ciphertext[24:48])
	data, _ := hex.DecodeString(ciphertext[48:])
	key := pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return data
}

var implantNames = []string{"Abdullah", "Absalom", "Acilia", "Adele", "Aggie68", "Agotogi", "Agsilly", "Agyo", "AhhoDesunenIX", "AhhoZurako", "Ahiru", "Aisa", "Akehende", "AkudaiKanzaburo", "Akumai", "Aladine", "Albion", "Ally", "Alpacaman", "Alvida", "Amadob", "Amazon", "AnZengaiina", "And", "Andre", "Anjo", "Ankoro", "Antonio", "Aphelandra", "Aramaki", "AremoGanmi", "Arlong", "Arrow", "Arthur", "AruyutayanV", "Asahija", "Asashichi", "AshuraDoji", "Aswa", "Atagoyama", "Atlas", "Atmos", "Attach", "AvaloPizarro", "AyeséMar", "Azuki", "Babanuki", "Babe", "Baby5", "Bacura", "Baggaley", "Bakezo", "BanDedessinée", "Banchi", "Banchina", "Bankuro", "Banshee", "Banzaburo", "Bao", "BaoHuang", "Bariete", "Barrel", "Barry", "BartholomewKuma", "Bartolomeo", "Bas", "BasilHawkins", "Basilisk", "Bastille", "Batchee", "Batman", "Baxcon", "BeerVI", "BeloBetty", "Belladonna", "Bellamy", "Bellett", "Bell-mere", "BennBeckman", "Bentham", "Bepo", "Bian", "BigPan", "Bimine", "Bingo", "Bishamon", "Biyo", "Bizarre", "BlackMaria", "Blackback", "Blamenco", "Blenheim", "Blondie", "BlueFan", "BlueGilly", "Bluejam", "Blueno", "Blumarine", "BoaHancock", "BoaMarigold", "BoaSandersonia", "Bobbin", "BobbyFunk", "Bobomba", "Bogard", "Bomba", "Bomba", "Bongo", "BonkPunch", "Boo", "Boodle", "Boogie", "Borsalino", "BourbonJr", "Braham", "Brahm", "Brannew", "Brew", "Briscola", "Brocca", "Brogy", "Brontosaurus", "Brook", "Bubblegum", "Buche", "Buchi", "BuckinghamStussy", "Buffalo", "Buggy", "Buhichuck", "BuildingSnake", "Bunbuku", "Bungo", "BunnyJoe", "Bushon", "Busshiri", "Byron", "Cabaji", "CaesarClown", "Caimanlady", "Camie", "Camel", "Cancer", "Candre", "Cands", "Candy", "CaponeBege", "CaponePez", "Capote", "Caribou", "Carmel", "Carne", "Carrot", "Catacombo", "CatarinaDevon", "Cavendish", "CBGallant", "Cerberus", "Cezar", "Chabo", "ChadrosHigelyges", "Chaka", "Chao", "Chap", "Chappe", "Charlos", "CharlotteAkimeg", "CharlotteAllmeg", "CharlotteAmande", "CharlotteAnana", "CharlotteAngel", "CharlotteAnglais", "CharlotteBasans", "CharlotteBasskarte", "CharlotteBavarois", "CharlotteBrownie", "CharlotteBroyé", "CharlotteBrûlée", "CharlotteCabaletta", "CharlotteCadenza", "CharlotteChiboust", "CharlotteChiffon", "CharlotteCinnamon", "CharlotteCitron", "CharlotteCompo", "CharlotteCompote", "CharlotteCornstarch", "CharlotteCounter", "CharlotteCracker", "CharlotteCustard", "CharlotteDacquoise", "CharlotteDaifuku", "CharlotteDe-Chat", "CharlotteDolce", "CharlotteDosmarche", "CharlotteDragée", "CharlotteEffilée", "CharlotteFlampe", "CharlotteFuyumeg", "CharlotteGala", "CharlotteGalette", "CharlotteHachée", "CharlotteHarumeg", "CharlotteHigh-Fat", "CharlotteJoconde", "CharlotteJoscarpone", "CharlotteKanten", "CharlotteKatakuri", "CharlotteKato", "CharlotteLaurin", "CharlotteLinlin", "CharlotteLola", "CharlotteMaple", "CharlotteMarble", "CharlotteMarnier", "CharlotteMascarpone", "CharlotteMash", "CharlotteMelise", "CharlotteMobile", "CharlotteMondée", "CharlotteMont-dOr", "CharlotteMontb", "CharlotteMoscato", "CharlotteMozart", "CharlotteMyukuru", "CharlotteNewgo", "CharlotteNewichi", "CharlotteNewji", "CharlotteNewsan", "CharlotteNewshi", "CharlotteNoisette", "CharlotteNormande", "CharlotteNougat", "CharlotteNusstorte", "CharlotteNutmeg", "CharlotteOpera", "CharlotteOven", "CharlottePanna", "CharlottePerospero", "CharlottePoire", "CharlottePraline", "CharlottePrim", "CharlottePudding", "CharlotteRaisin", "CharlotteSaint-Marc", "CharlotteSmoothie", "CharlotteSnack", "CharlotteTablet", "CharlotteWafers", "CharlotteWiro", "CharlotteYuen", "CharlotteZuccotto", "Chess", "Chesskippa", "Chessmarimo", "Chew", "Chichilisia", "Chicken", "Chicory", "Chimney", "Chinjao", "Chihaya", "Chiya", "Cho", "Chocho", "ChocoPolice", "Chocolat", "Choi", "Chome", "Chouchou", "Chuchun", "Chuji", "Clione", "Clover", "Coburn", "Cocoa", "Cocox", "Colscon", "Columbus", "Comil", "Compo", "Concelot", "Conis", "Corgi", "Coribou", "Cornelia", "Cosette", "Cosmo", "Cosmos", "Cotton", "Couran", "Cowboy", "Crocodile", "Crocus", "Curiel", "CurlyDadan", "Cyrano", "DR", "Dachoman", "DaddyDee", "Dagama", "Daidalos", "Daifugo", "Daigin", "Daikoku", "Daikon", "Daisy", "Dalmatian", "Dalton", "Damask", "Daruma", "DazBonez", "DecalvanBrothers", "Delacuaji", "Dellinger", "DemaroBlack", "Den", "Denjiro", "DevilDias", "Diamante", "Didit", "Diesel", "DiezBarrels", "DirtBoss", "DiscJ", "Disco", "Dive", "DobbyIbadonbo", "Doberman", "Dobon", "DocQ", "Dogra", "Dogya", "DohaIttankaII", "Doll", "Doma", "Domino", "Domo-kun", "Donannoyo", "Donovan", "Donquino", "DonquixoteDoflamingo", "DonquixoteHoming", "DonquixoteMjosgard", "DonquixoteRosinante", "Doran", "Doringo", "Dorry", "Dosun", "Dotaku", "DraculeMihawk", "DragonNumberThirteen", "Draw", "Drip", "Drophy", "DrugPeclo", "DuFeld", "DuckyBree", "Duval", "Eddy", "Edison", "EdwardNewgate", "EdwardWeevil", "Eiri", "Egana", "EggplantSoldier", "Eikon", "ElizabelloII", "Elmy", "Emma", "EmporioIvankov", "Enel", "Enishida", "Epoida", "Époni", "Erik", "Erio", "Esta", "EthanbaronVNusjuro", "EustassKid", "Farafra", "Farul", "Faust", "FenBock", "FigarlandGarling", "Finamore", "Fillonce", "Fishbonen", "FisherTiger", "Flapper", "Flare", "ForestBoss", "Forliewbs", "Fossa", "Fourtricks", "Foxy", "Francois", "Franky", "Fuga", "Fugar", "FugetsuOmusubi", "Fujin", "Fukaboshi", "Fukumi", "Fukurokuju", "Fukurou", "Fullbody", "Funkfreed", "Furrari", "Fuza", "Gab", "Gaburu", "GaikotsuYukichi", "Gaimon", "Gal", "Galaxy", "Galdino", "Gally", "GamaPyonnosuke", "Gambia", "Gambo", "GanFall", "Gancho", "Ganryu", "Ganryu", "Gatherine", "Gatz", "Gazelleman", "GeckoMoria", "Gedatsu", "Gem", "Genbo", "Genrin", "Genzo", "GeorgeBlack", "GeorgeMach", "Gerd", "Gerotini", "Giberson", "Gig", "Gimlet", "Gin", "Gina", "Ginko", "Ginnosuke", "Ginrummy", "GiroChintaro", "Giolla", "Gion", "Giovanni", "Gismonda", "Gladius", "Gloriosa", "Glove", "Gode", "GoingMerry", "Goki", "GolDRoger", "Goldberg", "GoldfishPrincess", "Gomorrah", "Gonbe", "Goo", "Gorilla", "Gorishiro", "Goro", "Gorobe", "Gotti", "Grabar", "Gram", "GreatMichael", "Guernica", "Gyaro", "Gyoro", "Gyoru", "Gyro", "Hack", "Hack", "Hajrudin", "Hakowan", "Hakugan", "HamBurger", "Hamburg", "Hamlet", "Hammond", "Han", "Hana", "Hangan", "Hanger", "Hanji", "Hanji", "Hannyabal", "Hanzo", "HappaYamao", "Happygun", "Hara", "Haredas", "Harisenbon", "HaritsuKendiyo", "Haruta", "Hasami", "Hatcha", "Hatchan", "Hattori", "Heat", "Helmeppo", "Heppoko", "Hera", "Hera", "Heracles", "Herb", "Hewitt", "Hibari", "Hidayu", "Higuma", "Hihimaru", "Hikoichi", "Hildon", "Hina", "Hinokizu", "Hip", "HippoGentleman", "Hiramera", "Hiriluk", "Hiroshi", "Hitsugisukan", "Ho", "Hocha", "Hocker", "HodyJones", "Hoe", "Hogback", "Hoichael", "Holedem", "Holy", "Home", "Hongo", "Honner", "Hop", "Hotei", "Hotori", "House", "Hublot", "Humphrey", "Hustle", "Hyogoro", "Hyota", "Hyoutauros", "Hyouzou", "Ibiributsu", "Ibusu", "Iceburg", "Ichika", "Ideaman", "Ideo", "Igaram", "IkarosMuch", "Ikkaku", "Ikkaku", "Imu", "Inazuma", "Inbi", "Indigo", "Inhel", "Inoichibannosuke", "Inuarashi", "Inukai", "Inuppe", "Ipponmatsu", "Ipponume", "Isa", "IshigoShitemanna", "Ishilly", "IslandEater", "Islewan", "Issho", "Isuka", "Itomimizu", "IvanX", "Izou", "Iwatobi", "Jabra", "Jack", "Jack-in-the-Box", "Jacksonbanner", "Jaguar", "JaguarDSaul", "Jaki", "Jalmack", "Jango", "Jarul", "JaygarciaSaturn", "JeanAngo", "JeanBart", "JeanGoen", "Jeep", "Jeet", "Jerry", "Jero", "JesusBurgess", "JewWall", "JewelryBonney", "Jibuemon", "JigokuBenten", "Jigoro", "Jigra", "Jinbe", "Jiron", "Jizo", "Jo", "Jobo", "John", "JohnGiant", "Johnny", "Jorge", "Joseph", "Joshu", "Jorul", "JoyBoy", "Jozu", "Judy", "Jujiro", "Juki", "Julius", "Junan", "Kabu", "Kadar", "Kagero", "Kagikko", "Kaidou", "Kairen", "Kairiken", "KairoKureyo", "Kaku", "Kaku", "Kakukaku", "Kakunoshin", "Kalgara", "Kalifa", "Kamakiri", "Kamekichi", "Kamijiro", "Kaneshiro", "Kanezenny", "Kappa", "Karasu", "Karma", "Karoo", "Kasa", "Kasagoba", "Kashii", "Kashigami", "Katsuzo", "Kawamatsu", "Kaya", "Kazekage", "Kazenbo", "Kebi", "Kechatch", "Keith", "Kentauros", "KellyFunk", "Kerville", "Kibagaeru", "Kibin", "Kiev", "Kikipatsu", "Kiku", "Kikunojo", "Kikyo", "Killer", "Kimel", "Kinemon", "Kinbo", "Kinderella", "King", "Kingbaum", "Kinga", "Kingdew", "Kinoko", "Kirinkodanuki", "Kirintauros", "Kisegawa", "Kitton", "Kiwi", "Koala", "Kobe", "Koby", "Koda", "Koito", "Kojuro", "Kokoro", "Komachiyo", "Komane", "Konbu", "Kong", "Kop", "Kotatsu", "Kotetsu", "Kotori", "Koyama", "Koza", "Koze", "KozukiHiyori", "KozukiMomonosuke", "KozukiOden", "KozukiSukiyaki", "KozukiToki", "Krieg", "Kujaku", "Kukai", "Kumadori", "KumadoriYamanbako", "Kumae", "Kumagoro", "KumaguchiIchiro", "Kumashi", "Kuni", "Kunyun", "Kureha", "Kuro", "Kurokoma", "Kuromarimo", "Kuroobi", "Kurosawa", "Kurotsuru", "KurozumiHigurashi", "KurozumiKanjuro", "KurozumiOrochi", "KurozumiSemimaru", "KurozumiTama", "Kuzan", "Kyros", "Kyuin", "Kyuji", "Kyukyu", "Laboon", "Lacroix", "Lacuba", "LadyTree", "Laffitte", "LaoG", "Laskey", "Lassoo", "Lemoncheese", "Leo", "Leonero", "Lilith", "Lily", "Limejuice", "Lindbergh", "Lines", "Lionbuta", "LipDoughty", "LittleOarsJr", "Loki", "Lola", "Lonz", "LordoftheCoast", "LouisArnote", "LuckyRoux", "Lulis", "Machvise", "Macro", "Macro", "Madilloman", "Magellan", "Magra", "Maha", "Maidy", "Maki", "Makino", "Makko", "Manboshi", "Mani", "Manjaro", "Mansherry", "Marco", "MarcusMars", "Margarita", "Marguerite", "Mari", "MariaNapole", "Marianne", "Marie", "Marilyn", "Marin", "Mario", "MarshallDTeach", "Marumieta", "Mash", "Mashikaku", "Masira", "MaskedDeuce", "Massui", "Master", "MasteroftheWaters", "Matryosaka", "Matryoseka", "Matryosoka", "Matryosuka", "Matsuge", "Maujii", "Maynard", "MAXMarx", "Mayushika", "McGuy", "McKinley", "Meadows", "Mecha-Shark", "Megalo", "Mero", "Merry", "Michael", "Mihar", "Mikazuki", "MikioItoo", "Mikita", "Milky", "MilletPine", "Milo", "Minatomo", "Minatomo", "MinisteroftheLeft", "MinisteroftheRight", "Minochihuahua", "Minokoala", "Minorhinoceros", "Minoruba", "Minotaurus", "Minozebra", "Misery", "MissCatherina", "MissFathersDay", "MissFriday", "MissMonday", "MissMothersDay", "MissSaturday", "MissThursday", "MissTuesday", "Miyagi", "Mizerka", "Mizuira", "MizutaMadaisuki", "MizutaMawaritosuki", "Moai", "MocDonald", "Mocha", "Mochi", "Moda", "Mohji", "Momonga", "Momoo", "Monda", "Monet", "Monjii", "Monjiro", "MonkeyDDragon", "MonkeyDGarp", "MonkeyDLuffy", "Monster", "MontBlancCricket", "MontBlancNoland", "Moodie", "MoonIsaacJr", "Moqueca", "Morgan", "Morgans", "Morley", "Mornin", "Mororon", "Mosh", "Motobaro", "Motzel", "Mounblutain", "MountainGod", "MountainRicky", "Mouseman", "Mousse", "Moyle", "Mozambia", "Mozu", "Mr6", "Mr7", "Mr7", "Mr9", "Mr10", "Mr11", "Mr12", "Mr13", "MrBeans", "MrLove", "MrMellow", "MrMomora", "MrSacrifice", "MrShimizu", "Muchana", "Mugren", "MukkashimiTower", "Mummy", "MummyMee", "Muret", "Nako", "Nami", "Namur", "Nangi", "Napoleon", "Nashi", "Natto", "NazuKetagari", "NefertariCobra", "NefertariDLili", "NefertariTiti", "NefertariVivi", "Neggy", "NegikumaMaria", "Nekomamushi", "Nekozaemon", "Neptune", "Nerine", "Nero", "NeronaImu", "Nezumi", "NicoOlvia", "NicoRobin", "Nigeratta", "Nika", "Nika", "Nin", "Ninjin", "Ninth", "Nitro", "Nnke-kun", "NobleCroc", "Nokotti", "Nojiko", "Nola", "NoraGitsune", "Noriko", "Nosgarl", "Nubon", "NugireYainu", "Nure-Onna", "Nyasha", "Oars", "Octopako", "Ochoku", "Ohm", "Oide", "Oimo", "Okame", "Okome", "Oli", "Oliva", "Omasa", "Onigumo", "Onimaru", "Oran", "Orlumbus", "Ossamondo", "Otohime", "OutlookIII", "Packy", "Pagaya", "PageOne", "Palms", "Pandaman", "Pandawoman", "Pandora", "PankutaDakeyan", "Pansy", "Papaneel", "Papas", "Pappag", "Pascia", "Patty", "Paulie", "Pavlik", "Peachbeard", "Pearl", "PeepleyLulu", "Pedro", "Pekkori", "Pekoms", "Pell", "Pellini", "Penguin", "Peppoko", "Perona", "Peterman", "Petermoo", "Pica", "Pickles", "Pierre", "Piiman", "Pike", "Pinky", "Pisaro", "Pluming", "Poker", "Pomp", "Poppoko", "Poppy", "Porche", "Porchemy", "Poro", "PortgasDAce", "PortgasDRouge", "Potaufeu", "Potsun", "Pound", "PrinceGrus", "Prometheus", "PuddingPudding", "Pudos", "Puppu", "Pururu", "PX-1", "PX-4", "PX-5", "PX-7", "Pythagoras", "Queen", "QueenMamaChanter", "Quincy", "Rabbitman", "Rabiyan", "Raccoon", "Raideen", "Raijin", "Raizo", "Raki", "Rakuda", "Rakuyo", "Ramba", "Ramen", "Rampo", "Ran", "Randolph", "Rangram", "Ratel", "Rebecca", "Reck", "Recycle-Wan", "Reforte", "Reuder", "Richie", "Rika", "RikuDoldoIII", "Rindo", "Rint", "Ripper", "Ririka", "RiskyBrothers", "RiskyBrothers", "Rivers", "Road", "RobLucci", "Robson", "Roche", "RocheTomson", "Rock", "RocksDXebec", "Rockstar", "Roddy", "Roji", "Rokkaku", "Rokki", "RollingLogan", "RoronoaArashi", "RoronoaPinzoro", "RoronoaZoro", "Roshio", "Ross", "Rosward", "Rowing", "Roxanne", "RugBear", "Run", "Rush", "Russian", "Ryuboshi", "Ryunosuke", "S-Bear", "S-Hawk", "S-Shark", "S-Snake", "Saber", "Sabo", "Sadi", "Sai", "Saikoro", "Sakazuki", "Saki", "Saldeath", "Salome", "Sam", "Samosa", "SamuraiBatts", "Sancrin", "Sandayu", "Sanji", "SanjuanWolf", "Sanka", "Sapi", "Sarahebi", "Sarfunkel", "SarieNantokanette", "Sarquiss", "Saru", "Sarutobi", "Sasaki", "Satori", "Sauce", "Scarlett", "Schollzo", "ScopperGaban", "Scotch", "Scotch", "ScratchmenApoo", "Seaboar", "SeagullGunsNozdon", "Seamars", "Seira", "Seki", "Sengoku", "Sennorikyuru", "SenorPink", "Sentomaru", "Serizawa", "Seto", "Shachi", "Shaka", "Shakuyaku", "Shalria", "Sham", "Shanba", "ShandiaChief", "Shanks", "Sharinguru", "Sheepshead", "Shelly", "ShepherdJuPeter", "Shiki", "ShimotsukiFuriko", "ShimotsukiKoushirou", "ShimotsukiKozaburo", "ShimotsukiKuina", "ShimotsukiRyuma", "ShimotsukiUshimaru", "ShimotsukiYasuie", "ShinDetamaruka", "ShinJaiya", "Shine", "Shinobu", "Shinosuke", "Shion", "Shioyaki", "Shirahoshi", "Shiryu", "Shishilian", "Shoujou", "Shu", "Shura", "Shyarly", "Sicily", "SilverAxe", "SilversRayleigh", "Sind", "Skull", "Sleepy", "Smiley", "Smoker", "Smooge", "Snakeman", "Sodom", "Solitaire", "Some", "Sonieh", "Sora", "Sora", "Soro", "Spacey", "Spandam", "Spandine", "Spartan", "Spector", "Speed", "SpeedJiru", "Spencer", "Sphinx", "Splash", "Splatter", "Spoil", "Squard", "Stainless", "Stalker", "Stansen", "StealthBlack", "Stefan", "Sterry", "Stevie", "Stomp", "Stool", "Strawberry", "Streusen", "Stroganoff", "Stronger", "Stussy", "Su", "Sugamichi", "Sugar", "Suke", "SukoshibaKanishitoru", "Suleiman", "Sunbell", "Surume", "SweetPea", "TBone", "Tabuhachiro", "Tacos", "Takao", "Take", "Tama", "Tamachibi", "Tamago", "Tamagon", "Tamanegi", "TankLepanto", "Tansui", "Tarara", "Tararan", "Taro", "Taroimo", "Tashigi", "Tate", "TeaIV", "TegataRingana", "Tenjo-Sagari", "Tensei", "Terracotta", "Terry", "TerryGilteo", "Tera", "Teru", "Tetsu", "ThalassaLucas", "Thatch", "Tibany", "Tilestone", "Togare", "Tokijiro", "Tokikake", "Toko", "Tom", "TomatoGang", "Tonjit", "TonyTonyChopper", "TopmanWarcury", "Torasaburo", "Tori", "Toto", "ToyamaTsujigiro", "TrafalgarDWaterLaw", "TrafalgarLami", "Trebol", "Tristan", "Tritobu", "TsugaruUmi", "Tsukimi", "Tsunagoro", "Tsunokkov", "Tsuru", "TsurueMonnosuke", "Tsurujo", "Turco", "Tyrannosaurus", "Ubau", "Ucy", "Uhho", "Uholisia", "UK", "Ukkari", "Ukon", "Ulti", "Ultraking", "Umit", "UnforgivableMask", "Uni", "Unigaro", "Urashima", "Urouge", "Usagihebi", "UsaguchiHideo", "Usakkov", "Ushiano", "Usopp", "Uta", "Uwattsura", "Uzu", "UzukiTempura", "VanAugur", "VanderDecken", "VanderDeckenIX", "VascoShot", "Vegapunk", "Vergo", "VeryGood", "VictoriaCindry", "VictoriaShirutonDoruyanaika", "VinsmokeIchiji", "VinsmokeJudge", "VinsmokeNiji", "VinsmokeReiju", "VinsmokeSora", "VinsmokeYonji", "Viola", "Vista", "Vitan", "Vito", "Wadatsumi", "Wakasa", "WallZombie", "Wallace", "Wallem", "Wanda", "Wany", "Wanyudo", "Wanze", "Wapol", "Warashi", "Warazane", "WarunoFurishiro", "Wellington", "Wheel", "WhiteyBay", "Whos-Who", "Wicca", "WillieGallon", "Wire", "WoopSlap", "Wyper", "XDrake", "Yama", "Yamakaji", "Yamato", "Yame", "Yamenahare", "Yamon", "Yarisugi", "Yasopp", "Yatappe", "Yazaemon", "Yokan", "Yokozuna", "Yomo", "Yonka", "YonkaTwo", "York", "Yorki", "Yosaku", "Yoshimoto", "Yotsubane", "Yu", "Yui", "Yuki", "Yurikah", "Zabo", "Zadie", "Zala", "Zambai", "Zanki", "Zeff", "Zeo", "Zepo", "Zeus", "Zodia", "Zotto", "Zucca", "Zunesha"}
