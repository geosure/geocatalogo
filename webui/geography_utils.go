package webui

import "strings"

// CountryEmoji returns the emoji flag for a country code
func CountryEmoji(code string) string {
	code = strings.ToLower(code)

	// Map country codes to emoji flags
	emojiMap := map[string]string{
		// ===== CLANKR CONTINENT - AI PROVIDER "COUNTRIES" =====
		"anthropic": "🟣", // Anthropic (Claude) - Purple (brand color)
		"openai":    "🟢", // OpenAI (GPT) - Green
		"google":    "🔵", // Google (Gemini) - Blue
		"gro":       "🟠", // GRO (internal bots) - Orange
		"aws":       "🟡", // AWS (infrastructure) - Yellow/Gold
		"inference": "🧠", // Inference providers - Brain

		// ===== VERBS CONTINENT - APP CATEGORIES =====
		"user-facing": "👥", // User-facing apps
		"internal":    "🔧", // Internal tools
		"services":    "⚙️", // API services

		// ===== TEAM CONTINENT - TEAM ROLES =====
		"leadership":       "👔", // Leadership team
		"engineering":      "⚙️", // Engineering team
		"data-engineering": "📊", // Data engineering
		"research":         "🔬", // Research contributors
		"advisors":         "🎓", // Advisors & mentors

		// North America
		"us": "🇺🇸",
		"ca": "🇨🇦",
		"mx": "🇲🇽",

		// South America
		"ar": "🇦🇷",
		"br": "🇧🇷",
		"cl": "🇨🇱",
		"co": "🇨🇴",
		"pe": "🇵🇪",
		"ve": "🇻🇪",
		"ec": "🇪🇨",
		"bo": "🇧🇴",
		"py": "🇵🇾",
		"uy": "🇺🇾",

		// Europe
		"gb": "🇬🇧",
		"uk": "🇬🇧",
		"de": "🇩🇪",
		"fr": "🇫🇷",
		"it": "🇮🇹",
		"es": "🇪🇸",
		"pt": "🇵🇹",
		"nl": "🇳🇱",
		"be": "🇧🇪",
		"ch": "🇨🇭",
		"at": "🇦🇹",
		"se": "🇸🇪",
		"no": "🇳🇴",
		"dk": "🇩🇰",
		"fi": "🇫🇮",
		"ie": "🇮🇪",
		"pl": "🇵🇱",
		"cz": "🇨🇿",
		"hu": "🇭🇺",
		"ro": "🇷🇴",
		"gr": "🇬🇷",
		"bg": "🇧🇬",
		"hr": "🇭🇷",
		"sk": "🇸🇰",
		"si": "🇸🇮",
		"ee": "🇪🇪",
		"lv": "🇱🇻",
		"lt": "🇱🇹",
		"ru": "🇷🇺",
		"ua": "🇺🇦",
		"tr": "🇹🇷",
		"eu": "🇪🇺",

		// Asia
		"cn": "🇨🇳",
		"jp": "🇯🇵",
		"kr": "🇰🇷",
		"in": "🇮🇳",
		"id": "🇮🇩",
		"th": "🇹🇭",
		"vn": "🇻🇳",
		"ph": "🇵🇭",
		"my": "🇲🇾",
		"sg": "🇸🇬",
		"bd": "🇧🇩",
		"pk": "🇵🇰",
		"mm": "🇲🇲",
		"kh": "🇰🇭",
		"la": "🇱🇦",
		"mn": "🇲🇳",
		"np": "🇳🇵",
		"lk": "🇱🇰",
		"af": "🇦🇫",
		"iq": "🇮🇶",
		"ir": "🇮🇷",
		"sa": "🇸🇦",
		"ae": "🇦🇪",
		"il": "🇮🇱",
		"jo": "🇯🇴",
		"lb": "🇱🇧",
		"sy": "🇸🇾",
		"ye": "🇾🇪",
		"om": "🇴🇲",
		"kw": "🇰🇼",
		"qa": "🇶🇦",
		"bh": "🇧🇭",

		// Africa
		"za": "🇿🇦",
		"ng": "🇳🇬",
		"eg": "🇪🇬",
		"ke": "🇰🇪",
		"et": "🇪🇹",
		"gh": "🇬🇭",
		"tz": "🇹🇿",
		"ug": "🇺🇬",
		"dz": "🇩🇿",
		"ma": "🇲🇦",
		"ao": "🇦🇴",
		"mz": "🇲🇿",
		"mg": "🇲🇬",
		"cm": "🇨🇲",
		"ci": "🇨🇮",
		"ne": "🇳🇪",
		"bf": "🇧🇫",
		"ml": "🇲🇱",
		"mw": "🇲🇼",
		"zm": "🇿🇲",
		"sn": "🇸🇳",
		"so": "🇸🇴",
		"td": "🇹🇩",
		"gn": "🇬🇳",
		"rw": "🇷🇼",
		"bj": "🇧🇯",
		"tn": "🇹🇳",
		"bi": "🇧🇮",
		"ss": "🇸🇸",
		"tg": "🇹🇬",
		"sl": "🇸🇱",
		"ly": "🇱🇾",
		"lr": "🇱🇷",
		"mr": "🇲🇷",
		"cf": "🇨🇫",
		"ga": "🇬🇦",
		"gw": "🇬🇼",
		"gq": "🇬🇶",
		"mu": "🇲🇺",
		"sz": "🇸🇿",
		"dj": "🇩🇯",
		"km": "🇰🇲",
		"cv": "🇨🇻",
		"st": "🇸🇹",
		"sc": "🇸🇨",

		// Oceania
		"au": "🇦🇺",
		"nz": "🇳🇿",
		"fj": "🇫🇯",
		"pg": "🇵🇬",
		"nc": "🇳🇨",
		"pf": "🇵🇫",
		"ws": "🇼🇸",
		"to": "🇹🇴",
		"vu": "🇻🇺",
		"sb": "🇸🇧",
	}

	if emoji, ok := emojiMap[code]; ok {
		return emoji
	}

	// Default: use generic pin emoji
	return "📍"
}

// ContinentEmoji returns the emoji for a continent
func ContinentEmoji(continent string) string {
	continent = strings.ToLower(continent)

	emojiMap := map[string]string{
		"north-america": "🌎",
		"south-america": "🌎",
		"europe":        "🌍",
		"africa":        "🌍",
		"asia":          "🌏",
		"oceania":       "🌏",
		"global":        "🌐",
		"clankr":        "🤖", // Clankr - AI Agent Continent
		"verbs":         "🎯", // Verbs - Application Continent
		"team":          "👥", // Team - The People Behind GRO
		"infra":         "🏗️", // Infra - AWS Infrastructure (Planet Infra)
	}

	if emoji, ok := emojiMap[continent]; ok {
		return emoji
	}

	return "🌐"
}

// CountryName returns the full country name for a code
func CountryName(code string) string {
	code = strings.ToLower(code)

	nameMap := map[string]string{
		// ===== CLANKR CONTINENT - AI PROVIDER "COUNTRIES" =====
		"anthropic": "Anthropic",
		"openai":    "OpenAI",
		"google":    "Google",
		"gro":       "GRO",
		"aws":       "AWS",
		"inference": "Inference",

		// ===== VERBS CONTINENT - APP CATEGORIES =====
		"user-facing": "User-Facing Apps",
		"internal":    "Internal Tools",
		"services":    "API Services",

		// ===== TEAM CONTINENT - TEAM ROLES =====
		"leadership":       "Leadership",
		"engineering":      "Engineering",
		"data-engineering": "Data Engineering",
		"research":         "Research",
		"advisors":         "Advisors",

		// North America
		"us": "United States",
		"ca": "Canada",
		"mx": "Mexico",

		// South America
		"ar": "Argentina",
		"br": "Brazil",
		"cl": "Chile",
		"co": "Colombia",
		"pe": "Peru",
		"ve": "Venezuela",
		"ec": "Ecuador",
		"bo": "Bolivia",
		"py": "Paraguay",
		"uy": "Uruguay",

		// Europe
		"gb": "United Kingdom",
		"uk": "United Kingdom",
		"de": "Germany",
		"fr": "France",
		"it": "Italy",
		"es": "Spain",
		"pt": "Portugal",
		"nl": "Netherlands",
		"be": "Belgium",
		"ch": "Switzerland",
		"at": "Austria",
		"se": "Sweden",
		"no": "Norway",
		"dk": "Denmark",
		"fi": "Finland",
		"ie": "Ireland",
		"pl": "Poland",
		"cz": "Czech Republic",
		"hu": "Hungary",
		"ro": "Romania",
		"gr": "Greece",
		"bg": "Bulgaria",
		"hr": "Croatia",
		"sk": "Slovakia",
		"si": "Slovenia",
		"ee": "Estonia",
		"lv": "Latvia",
		"lt": "Lithuania",
		"ru": "Russia",
		"ua": "Ukraine",
		"tr": "Turkey",
		"eu": "European Union",

		// Asia
		"cn": "China",
		"jp": "Japan",
		"kr": "South Korea",
		"in": "India",
		"id": "Indonesia",
		"th": "Thailand",
		"vn": "Vietnam",
		"ph": "Philippines",
		"my": "Malaysia",
		"sg": "Singapore",
		"bd": "Bangladesh",
		"pk": "Pakistan",
		"mm": "Myanmar",
		"kh": "Cambodia",
		"la": "Laos",
		"mn": "Mongolia",
		"np": "Nepal",
		"lk": "Sri Lanka",
		"af": "Afghanistan",
		"iq": "Iraq",
		"ir": "Iran",
		"sa": "Saudi Arabia",
		"ae": "United Arab Emirates",
		"il": "Israel",
		"jo": "Jordan",
		"lb": "Lebanon",
		"sy": "Syria",
		"ye": "Yemen",
		"om": "Oman",
		"kw": "Kuwait",
		"qa": "Qatar",
		"bh": "Bahrain",

		// Africa
		"za": "South Africa",
		"ng": "Nigeria",
		"eg": "Egypt",
		"ke": "Kenya",
		"et": "Ethiopia",
		"gh": "Ghana",
		"tz": "Tanzania",
		"ug": "Uganda",
		"dz": "Algeria",
		"ma": "Morocco",
		"ao": "Angola",
		"mz": "Mozambique",
		"mg": "Madagascar",
		"cm": "Cameroon",
		"ci": "Côte d'Ivoire",
		"ne": "Niger",
		"bf": "Burkina Faso",
		"ml": "Mali",
		"mw": "Malawi",
		"zm": "Zambia",
		"sn": "Senegal",
		"so": "Somalia",
		"td": "Chad",
		"gn": "Guinea",
		"rw": "Rwanda",
		"bj": "Benin",
		"tn": "Tunisia",
		"bi": "Burundi",
		"ss": "South Sudan",
		"tg": "Togo",
		"sl": "Sierra Leone",
		"ly": "Libya",
		"lr": "Liberia",

		// Oceania
		"au": "Australia",
		"nz": "New Zealand",
		"fj": "Fiji",
		"pg": "Papua New Guinea",
	}

	if name, ok := nameMap[code]; ok {
		return name
	}

	return strings.ToUpper(code)
}
