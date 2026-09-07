package main

// Field specs and seed defaults for the three Firestore-backed pages. The
// defaults are the original hand-written copy from when these pages were
// static HTML, so nothing goes blank until the client actually edits a
// field - see loadPageContent.

var homeFields = []fieldSpec{
	{"badge", "Odznak nad nadpisem", "line", "Hero"},
	{"h1", "Nadpis (jméno)", "line", "Hero"},
	{"h2", "Podnadpis", "line", "Hero"},
	{"body", "Úvodní text", "text", "Hero"},

	{"card1_title", "Box 1 - nadpis", "line", "Hero – tři boxy"},
	{"card1_body", "Box 1 - text", "text", "Hero – tři boxy"},
	{"card2_title", "Box 2 - nadpis", "line", "Hero – tři boxy"},
	{"card2_body", "Box 2 - text", "text", "Hero – tři boxy"},
	{"card3_title", "Box 3 - nadpis", "line", "Hero – tři boxy"},
	{"card3_body", "Box 3 - text", "text", "Hero – tři boxy"},

	{"services_heading", "Nadpis sekce", "line", "S čím vám mohu pomoci"},
	{"services_subheading", "Podnadpis sekce", "line", "S čím vám mohu pomoci"},

	{"tile1_title", "Dlaždice 1 - nadpis", "line", "Dlaždice 1"},
	{"tile1_body", "Dlaždice 1 - text", "text", "Dlaždice 1"},
	{"tile1_highlight", "Dlaždice 1 - zvýrazněná věta", "line", "Dlaždice 1"},
	{"tile1_bullets", "Dlaždice 1 - seznam (řádek = položka)", "list", "Dlaždice 1"},

	{"tile2_title", "Dlaždice 2 - nadpis", "line", "Dlaždice 2"},
	{"tile2_body", "Dlaždice 2 - text", "text", "Dlaždice 2"},
	{"tile2_bullets", "Dlaždice 2 - seznam (řádek = položka)", "list", "Dlaždice 2"},

	{"tile3_title", "Dlaždice 3 - nadpis", "line", "Dlaždice 3"},
	{"tile3_body", "Dlaždice 3 - text", "text", "Dlaždice 3"},
	{"tile3_bullets", "Dlaždice 3 - seznam (řádek = položka)", "list", "Dlaždice 3"},

	{"activities_heading", "Nadpis", "line", "Klíčové aktivity"},
	{"activities_badges", "Štítky (řádek = jeden štítek)", "list", "Klíčové aktivity"},

	{"process_heading", "Nadpis sekce", "line", "Jak spolupráce probíhá"},
	{"process_subheading", "Podnadpis sekce", "text", "Jak spolupráce probíhá"},

	{"step1_title", "Krok 1 - nadpis", "line", "Kroky spolupráce"},
	{"step1_body", "Krok 1 - text", "text", "Kroky spolupráce"},
	{"step2_title", "Krok 2 - nadpis", "line", "Kroky spolupráce"},
	{"step2_body", "Krok 2 - text", "text", "Kroky spolupráce"},
	{"step3_title", "Krok 3 - nadpis", "line", "Kroky spolupráce"},
	{"step3_body", "Krok 3 - text", "text", "Kroky spolupráce"},
	{"step4_title", "Krok 4 - nadpis", "line", "Kroky spolupráce"},
	{"step4_body", "Krok 4 - text", "text", "Kroky spolupráce"},
	{"step5_title", "Krok 5 - nadpis", "line", "Kroky spolupráce"},
	{"step5_body", "Krok 5 - text", "text", "Kroky spolupráce"},

	{"contact_heading", "Nadpis", "line", "Kontakt"},
	{"contact_intro", "Text vedle e-mailu", "text", "Kontakt"},
}

var homeDefaults = map[string]string{
	"badge": "Kybernetická bezpečnost • NIS2 • AI Act",
	"h1":    "Michaela Streckerová",
	"h2":    "Pomáhám organizacím nastavit bezpečnost tak, aby fungovala v praxi a nejen na papíře.",
	"body":  "Nastavíme společně bezpečnost tak, aby fungovala v každodenním provozu, splňovala zákonné požadavky a přirozeně zapadla do toho, jak vaše firma opravdu funguje. Zaměříme se na praktické kroky vycházející z NIS2 a AI Actu, které zapadnou do vašeho běžného provozu a budou srozumitelné a použitelné pro každého.",

	"card1_title": "Praktické kroky místo „papír pro papír“",
	"card1_body":  "Procesy a pravidla, které lidé opravdu používají.",
	"card2_title": "NIS2 & ZoKB bez zbytečné zátěže",
	"card2_body":  "Co udělat hned a co může počkat – jasné priority.",
	"card3_title": "AI Act srozumitelně",
	"card3_body":  "Bez strašení, s pravidly a kontrolními body.",

	"services_heading":    "S čím vám mohu pomoci",
	"services_subheading": "Tři oblasti, které nejčastěji řeším s klienty.",

	"tile1_title":     "AI ve firmě & požadavky AI Act",
	"tile1_body":      "AI může firmě hodně pomoci, ale jen ve chvíli, kdy je zavedena chytře a bezpečně. Společně nastavíme, kde vám AI dává smysl, jak ji bezpečně používat a co po vás chce AI Act.",
	"tile1_highlight": "Prakticky, bez strašení a bez zbytečných komplikací.",
	"tile1_bullets": "kde má pro vás AI reálný přínos\n" +
		"jak ji bezpečně a zodpovědně integrovat do provozu\n" +
		"co po vás AI Act skutečně požaduje\n" +
		"jasné a jednoduché postupy, pravidla a kontrolní body\n" +
		"příprava dokumentace i doporučení pro vedení",

	"tile2_title": "Procesy, role & bezpečnost podle NIS2 a ZoKB",
	"tile2_body":  "Když víte, co máte dělat a kdo za co odpovídá, všechno najednou funguje. Pomohu vám převést NIS2 a požadavky Zákona o kybernetické bezpečnosti do praxe tak, aby zapadly do vašeho běžného fungování, a ne aby ho komplikovaly.",
	"tile2_bullets": "nastavení rolí a odpovědností (CISO, vedení, IT, projektové role, role dle ZoKB)\n" +
		"procesy, které lidé opravdu používají\n" +
		"řízení rizik a rozhodování\n" +
		"praktické směrnice a postupy\n" +
		"přehled, co je nutné udělat hned a co může bez problémů počkat",

	"tile3_title": "Školení & bezpečnostní kultura",
	"tile3_body":  "Bezpečnost stojí hlavně na lidech. Když vědí proč a jak, začnou jednat správně automaticky. Připravím školení a interní komunikaci tak, aby si zaměstnanci odnesli hlavně to „proč“.",
	"tile3_bullets": "školení na míru firmě a prostředí (on-site i on-line)\n" +
		"krátké, srozumitelné formáty (prezentace, videa, microlearning)\n" +
		"interní články a kampaně\n" +
		"návody, které jsou jasné a použitelné",

	"activities_heading": "Klíčové aktivity",
	"activities_badges": "AI ve firmě & požadavky AI Actu\n" +
		"Procesy, role & bezpečnost podle NIS2 a ZoKB\n" +
		"Školení & bezpečnostní kultura",

	"process_heading":    "Jak spolupráce probíhá",
	"process_subheading": "Vždy mi záleží na tom, aby spolupráce byla jednoduchá a praktická. Žádné složité procesy, žádné nekonečné schůzky.",

	"step1_title": "1) Úvodní konzultace",
	"step1_body":  "Krátce si vyjasníme, co potřebujete vyřešit, jaké máte cíle a v jakém stavu jsou vaše procesy a dokumentace. Konzultace může být on-line, telefonicky či osobně.",
	"step2_title": "2) Rychlá analýza současného stavu",
	"step2_body":  "Podívám se na vaše prostředí, role, dokumenty a procesy. Nejde o audit, spíše praktický přehled, co funguje, co je v souladu a co je potřeba doplnit.",
	"step3_title": "3) Návrh konkrétních kroků",
	"step3_body":  "Připravím přehledné doporučení: co udělat hned, co počká a co může být dlouhodobý plán – tak, aby kroky zapadly do běžného provozu.",
	"step4_title": "4) Společné zavedení změn",
	"step4_body":  "Nastavíme role, procesy, pravidla pro AI, směrnice i školení. Krok za krokem, aby změny byly funkční a pochopitelné pro všechny.",
	"step5_title": "5) Předání podkladů a podpora",
	"step5_body":  "Dostanete jasné a přehledné materiály, postupy a dokumentaci, kterou můžete dál používat. Po dohodě i dlouhodobá podpora, školení nebo průběžné revize.",

	"contact_heading": "Kontakt",
	"contact_intro":   "Máte otázku, zájem o konzultaci nebo řešíte konkrétní požadavek? Napište mi pár vět níže a domluvíme další kroky.",
}

var aboutFields = []fieldSpec{
	{"heading", "Nadpis stránky", "line", "Úvod"},
	{"intro1", "Odstavec 1", "text", "Úvod"},
	{"intro2", "Odstavec 2", "text", "Úvod"},
	{"intro3", "Odstavec 3", "text", "Úvod"},

	{"today_heading", "Nadpis", "line", "Co dělám dnes"},
	{"today_body", "Text", "text", "Co dělám dnes"},
	{"today_bullets", "Seznam (řádek = položka)", "list", "Co dělám dnes"},

	{"sidebar_name", "Jméno pod fotkou", "line", "Postranní panel"},
	{"sidebar_role", "Popisek pod jménem", "line", "Postranní panel"},
	{"diff_heading", "Nadpis - Co mě odlišuje", "line", "Postranní panel"},
	{"diff_bullets", "Co mě odlišuje - seznam", "list", "Postranní panel"},
	{"edu_heading", "Nadpis - Vzdělání a certifikace", "line", "Postranní panel"},
	{"edu_bullets", "Vzdělání a certifikace - seznam", "list", "Postranní panel"},

	{"howwork_heading", "Nadpis - Jak pracuji", "line", "Jak a proč pracuji"},
	{"howwork_body", "Text - Jak pracuji", "text", "Jak a proč pracuji"},
	{"why_heading", "Nadpis - Proč to dělám", "line", "Jak a proč pracuji"},
	{"why_body", "Text - Proč to dělám", "text", "Jak a proč pracuji"},

	{"references_heading", "Nadpis", "line", "Reference"},
	{"reference_quote", "Text reference", "text", "Reference"},
	{"reference_name", "Jméno", "line", "Reference"},
	{"reference_role", "Pozice/role", "line", "Reference"},
}

var aboutDefaults = map[string]string{
	"heading": "O mně",
	"intro1":  "Jmenuji se Michaela Streckerová a několik let se pohybuji na pomezí kybernetické bezpečnosti, IT projektů a komunikace.",
	"intro2":  "Mám za sebou cestu od administrativy a podpory ISO procesů až po roli Manažerky kybernetické bezpečnosti ve veřejném sektoru.",
	"intro3":  "V praxi spojuji dva světy, které jsou často absolutně protichůdné: technické požadavky a lidsky pochopitelné vysvětlení toho, proč jsou důležité.",

	"today_heading": "Co dělám dnes",
	"today_body":    "V Technologické agentuře ČR působím jako manažer kybernetické bezpečnosti. Současně se jako freelancer věnuji bezpečnostním konzultacím a pomáhám klientům nastavovat procesy, compliance, dokumentaci a bezpečnostní požadavky tak, aby to dávalo smysl nejen bezpečnostně, ale i prakticky.",
	"today_bullets": "tvořím bezpečnostní politiky a strategie\n" +
		"analýzy rizik\n" +
		"navrhuji a nastavuji opatření\n" +
		"řeším incidenty\n" +
		"pomáhám při zavádění AI či implementaci AI Act\n" +
		"zajišťuji školení, přednášky a informovanost zaměstnanců\n" +
		"hlídám plnění legislativy a standardů (ISO 27001, AI Act, ZoKB)",

	"sidebar_name": "Michaela Streckerová",
	"sidebar_role": "Kybernetická bezpečnost • NIS2 • AI Act",
	"diff_heading": "Co mě odlišuje",
	"diff_bullets": "Umím přeložit „tech-speak“ do normální řeči.\n" +
		"Kombinuji analytické myšlení s empatií a praxí z projektového řízení.\n" +
		"Rozumím veřejné správě i soukromým firmám, ISO, řízení IT i veřejným zakázkám.",
	"edu_heading": "Vzdělání a certifikace",
	"edu_bullets": "Ing. – Management, VŠE\n" +
		"Manažer kybernetické bezpečnosti – TAYLLORCOX (2024)\n" +
		"Bezpečnostní prověrka NBÚ – „Důvěrné“\n" +
		"Zkušenosti s ISO 27001, riziky a auditní přípravou\n" +
		"Opakované každoroční certifikace NÚKIB",

	"howwork_heading": "Jak pracuji",
	"howwork_body":    "Můj přístup stojí na jednoduché myšlence: bezpečnost má pomáhat, ne překážet. Nesnažím se vytvářet zbytečnou byrokracii ve stylu „papír pro papír“, ale prostředí, kde lidé vědí proč je to důležité a jak jim to může pomoci.",
	"why_heading":     "Proč to dělám",
	"why_body":        "Kybernetická bezpečnost je často komplikované a nepopulární téma. Baví mě z něj dělat něco pochopitelného a užitečného. Nejvíc mě těší moment, kdy vidím, že díky pár správným krokům se organizace dokáže chránit lépe a lidé se v ní neztrácí.",

	"references_heading": "Reference",
	"reference_quote": "Spolupráce s Míšou byla prostě jízda – v tom nejlepším slova smyslu. Každá jsme přicházely z jiného světa, ona z IT a kyberbezpečnosti, já z práva, ale velmi rychle jsme našly společnou řeč. Má vzácný talent převést technické „IT řeči“ do srozumitelné podoby a současně přesně říct, co je potřeba dotáhnout, aby věci fungovaly.\n\n" +
		"Na naší spolupráci si cením hlavně jejího přístupu: je profesionální, efektivní a přirozeně pragmatická. Nehraje si na zbytečné procesy ani nedramatizuje detaily, které nemají reálný dopad. Soustředí se na výsledek a smysluplná řešení – a přitom zůstává otevřená debatě i právním doporučením, když dávají dobrý smysl.\n\n" +
		"Pro mě osobně to byla zkušenost, díky které jsem se naučila lépe propojovat obory a dívat se na problémy v širších souvislostech. Pokud chcete spolupracovat s někým, kdo spojuje odbornost, umění věci dotahovat a zdravý rozum, Michaela je přesně partner, kterého hledáte.",
	"reference_name": "Jana Syručková",
	"reference_role": "Advokátka",
}

func servicesBlockFields(n, section string) []fieldSpec {
	return []fieldSpec{
		{"s" + n + "_title", "Nadpis", "line", section},
		{"s" + n + "_intro", "Úvodní text", "text", section},
		{"s" + n + "_solve_heading", "Nadpis - co vyřešíme/nastavíme/připravím", "line", section},
		{"s" + n + "_solve_bullets", "Seznam (řádek = položka)", "list", section},
		{"s" + n + "_get_heading", "Nadpis - co získáte", "line", section},
		{"s" + n + "_get_bullets", "Seznam (řádek = položka)", "list", section},
	}
}

var servicesFields = append(
	append(
		append(
			[]fieldSpec{{"heading", "Nadpis stránky", "line", "Nadpis"}},
			servicesBlockFields("1", "Služba 1 – AI ve firmě")...,
		),
		servicesBlockFields("2", "Služba 2 – NIS2 a ZoKB")...,
	),
	servicesBlockFields("3", "Služba 3 – Školení")...,
)

var servicesDefaults = map[string]string{
	"heading": "S čím vám mohu pomoci",

	"s1_title":         "AI ve firmě & požadavky AI Act",
	"s1_intro":         "AI může být pro firmu velkou výhodou, pokud je zavedená chytře a bezpečně. Může ale být také velkou brzdou, pokud ji implementujeme jen proto, že je to zrovna trend, a nevíme jak to udělat správně. Pomohu vám najít konkrétní situace, kde má AI skutečný přínos, nastavit pravidla jejího používání a zároveň splnit to, co vám ukládá AI Act.",
	"s1_solve_heading": "Co společně vyřešíme",
	"s1_solve_bullets": "identifikace oblastí, kde má AI reálný dopad na efektivitu a úspory\n" +
		"bezpečné a zodpovědné zavedení AI do firemních procesů\n" +
		"posouzení rizikovosti podle AI Actu (minimal, limited, high-risk systems)\n" +
		"tvorba jednoduchých pravidel a kontrolních bodů\n" +
		"nastavení interních postupů pro zaměstnance i vedení\n" +
		"úprava nebo vytvoření bezpečnostní politiky\n" +
		"jasná dokumentace a podklady k auditu",
	"s1_get_heading": "Co získáte",
	"s1_get_bullets": "přehled, kde dává AI smysl a kde ne\n" +
		"bezpečné a kontrolované použití AI ve firmě\n" +
		"splnění legislativních požadavků ČR a AI Actu\n" +
		"srozumitelné materiály a doporučení, která lze hned zavést\n" +
		"větší efektivitu a méně manuální práce tam, kde to opravdu funguje",

	"s2_title":         "Procesy, role & bezpečnost podle NIS2 a ZoKB",
	"s2_intro":         "Požadavky NIS2 a nového Zákona o kybernetické bezpečnosti nemusí být složité. Pomohu vám zjistit, co se vás opravdu týká, jak to máte splnit, a převést povinnosti do takových kroků, které zapadnou do vašeho provozu bez zbytečných komplikací.",
	"s2_solve_heading": "Co společně nastavíme",
	"s2_solve_bullets": "role a odpovědnosti (CISO, vedení, IT, projektové role, role dle ZoKB)\n" +
		"funkční procesy, které lidé opravdu používají\n" +
		"řízení rizik a jasná rozhodovací schémata\n" +
		"incident management a reporting povinností\n" +
		"bezpečnost dodavatelů a třetích stran\n" +
		"srozumitelné směrnice a interní postupy\n" +
		"přehledné priority – co udělat hned a co může počkat",
	"s2_get_heading": "Co získáte",
	"s2_get_bullets": "praktické a funkční nastavení kybernetické bezpečnosti\n" +
		"splnění požadavků NIS2 a ZoKB bez zbytečné zátěže\n" +
		"jasné povinnosti pro všechny zúčastněné role\n" +
		"méně chaosu a lepší koordinaci mezi IT, vedením a uživateli\n" +
		"procesy, které firmu chrání, ale nebrzdí její běžné fungování",

	"s3_title":         "Školení & bezpečnostní kultura",
	"s3_intro":         "Bezpečnost začíná u lidí. A když lidé rozumějí tomu „proč“ a vidí v tom smysl, jednají bezpečně úplně jiným způsobem. Připravím školení a komunikaci, která je srozumitelná, krátká a použitelná v každodenní praxi.",
	"s3_solve_heading": "Co pro vás připravím",
	"s3_solve_bullets": "školení na míru prostředí vaší firmy (on-site i online)\n" +
		"krátké formáty – prezentace, videa, microlearning\n" +
		"interní články, newslettery a bezpečnostní kampaně\n" +
		"používání AI ve firmě včetně bezpečnostních zásad\n" +
		"postupy a návody, které zaměstnanci skutečně využijí",
	"s3_get_heading": "Co získáte",
	"s3_get_bullets": "zaměstnance, kteří skutečně chápou rizika a správné chování\n" +
		"splnění požadavku na proškolenost zaměstnanců i vedení\n" +
		"snížení chyb způsobených nepozorností nebo neznalostí\n" +
		"posílení bezpečnostní kultury bez zbytečné zátěže\n" +
		"dlouhodobé materiály, které můžete používat opakovaně",
}
