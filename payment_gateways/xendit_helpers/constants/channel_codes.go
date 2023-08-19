package constants

const (
    ChannelTypeBank = iota + 1
    ChannelTypeEWallet
)

type ChannelCodes struct {
    ChannelCode string
    ChannelType string
    Name        string
}

type ChannelType int
type Channels map[ChannelType][]ChannelCodes

func GetPayoutChannel() Channels {
    payoutChannel := Channels{
        ChannelTypeBank:    getBankChannel(),
        ChannelTypeEWallet: getEWalletChannel(),
    }

    return payoutChannel
}

func getBankChannel() []ChannelCodes {
    bankChannel := []ChannelCodes{
        {
            ChannelCode: "ACEH",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Aceh",
        },
        {
            ChannelCode: "ACEH_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Aceh Syariah (UUS)",
        },
        {
            ChannelCode: "AGRONIAGA",
            ChannelType: "Bank",
            Name:        "Bank Raya Indonesia (formerly called Bank BRI Agroniaga)",
        },
        {
            ChannelCode: "ALADIN",
            ChannelType: "Bank",
            Name:        "Bank Aladin Syariah (formerly Bank Maybank Syariah Indonesia)",
        },
        {
            ChannelCode: "ALLO",
            ChannelType: "Bank",
            Name:        "Allo Bank Indonesia (formerly Bank Harda Internasional)",
        },
        {
            ChannelCode: "AMAR",
            ChannelType: "Bank",
            Name:        "Bank Amar Indonesia (formerly called Anglomas International Bank)",
        },
        {
            ChannelCode: "ANZ",
            ChannelType: "Bank",
            Name:        "Bank ANZ Indonesia",
        },
        {
            ChannelCode: "ARTHA",
            ChannelType: "Bank",
            Name:        "Bank Artha Graha International",
        },
        {
            ChannelCode: "BALI",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Bali",
        },
        {
            ChannelCode: "BAML",
            ChannelType: "Bank",
            Name:        "Bank of America Merill-Lynch",
        },
        {
            ChannelCode: "BANTEN",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Banten (formerly called Bank Pundi Indonesia)",
        },
        {
            ChannelCode: "BCA",
            ChannelType: "Bank",
            Name:        "Bank Central Asia (BCA)",
        },
        {
            ChannelCode: "BCA_DIGITAL",
            ChannelType: "Bank",
            Name:        "Bank Central Asia Digital (BluBCA)",
        },
        {
            ChannelCode: "BCA_SYR",
            ChannelType: "Bank",
            Name:        "Bank Central Asia (BCA) Syariah",
        },
        {
            ChannelCode: "BENGKULU",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Bengkulu",
        },
        {
            ChannelCode: "BISNIS_INTERNASIONAL",
            ChannelType: "Bank",
            Name:        "Bank Bisnis Internasional",
        },
        {
            ChannelCode: "BJB",
            ChannelType: "Bank",
            Name:        "Bank BJB",
        },
        {
            ChannelCode: "BJB_SYR",
            ChannelType: "Bank",
            Name:        "Bank BJB Syariah",
        },
        {
            ChannelCode: "BNC",
            ChannelType: "Bank",
            Name:        "Bank Neo Commerce (formerly Bank Yudha Bhakti)",
        },
        {
            ChannelCode: "BNI",
            ChannelType: "Bank",
            Name:        "Bank Negara Indonesia (BNI)",
        },
        {
            ChannelCode: "BNP_PARIBAS",
            ChannelType: "Bank",
            Name:        "Bank BNP Paribas",
        },
        {
            ChannelCode: "BOC",
            ChannelType: "Bank",
            Name:        "Bank of China (BOC)",
        },
        {
            ChannelCode: "BRI",
            ChannelType: "Bank",
            Name:        "Bank Rakyat Indonesia (BRI)",
        },
        {
            ChannelCode: "BSI",
            ChannelType: "Bank",
            Name:        "Bank Syariah Indonesia (BSI)",
        },
        {
            ChannelCode: "BTN",
            ChannelType: "Bank",
            Name:        "Bank Tabungan Negara (BTN)",
        },
        {
            ChannelCode: "BTN_UUS",
            ChannelType: "Bank",
            Name:        "Bank Tabungan Negara Syariah (BTN UUS)",
        },
        {
            ChannelCode: "BTPN_SYARIAH",
            ChannelType: "Bank",
            Name:        "BTPN Syariah (formerly called BTPN UUS and Bank Sahabat Purba Danarta)",
        },
        {
            ChannelCode: "BUKOPIN",
            ChannelType: "Bank",
            Name:        "Bank Bukopin",
        },
        {
            ChannelCode: "BUKOPIN_SYR",
            ChannelType: "Bank",
            Name:        "Bank Syariah Bukopin",
        },
        {
            ChannelCode: "BUMI_ARTA",
            ChannelType: "Bank",
            Name:        "Bank Bumi Arta",
        },
        {
            ChannelCode: "CAPITAL",
            ChannelType: "Bank",
            Name:        "Bank Capital Indonesia",
        },
        {
            ChannelCode: "CCB",
            ChannelType: "Bank",
            Name:        "China Construction Bank Indonesia (formerly called Bank Antar Daerah and Bank Windu Kentjana International)",
        },
        {
            ChannelCode: "CHINATRUST",
            ChannelType: "Bank",
            Name:        "Bank Chinatrust Indonesia",
        },
        {
            ChannelCode: "CIMB",
            ChannelType: "Bank",
            Name:        "Bank CIMB Niaga",
        },
        {
            ChannelCode: "CIMB_UUS",
            ChannelType: "Bank",
            Name:        "Bank CIMB Niaga Syariah (UUS)",
        },
        {
            ChannelCode: "CITIBANK",
            ChannelType: "Bank",
            Name:        "Citibank",
        },
        {
            ChannelCode: "COMMONWEALTH",
            ChannelType: "Bank",
            Name:        "Bank Commonwealth",
        },
        {
            ChannelCode: "DAERAH_ISTIMEWA",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Daerah Istimewa Yogyakarta (DIY)",
        },
        {
            ChannelCode: "DAERAH_ISTIMEWA_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Daerah Istimewa Yogyakarta Syariah (DIY UUS)",
        },
        {
            ChannelCode: "DANAMON",
            ChannelType: "Bank",
            Name:        "Bank Danamon",
        },
        {
            ChannelCode: "DANAMON_UUS",
            ChannelType: "Bank",
            Name:        "Bank Danamon Syariah (UUS)",
        },
        {
            ChannelCode: "DBS",
            ChannelType: "Bank",
            Name:        "Bank DBS Indonesia",
        },
        {
            ChannelCode: "DEUTSCHE",
            ChannelType: "Bank",
            Name:        "Deutsche Bank",
        },
        {
            ChannelCode: "DINAR_INDONESIA",
            ChannelType: "Bank",
            Name:        "Bank Dinar Indonesia",
        },
        {
            ChannelCode: "DKI",
            ChannelType: "Bank",
            Name:        "Bank DKI",
        },
        {
            ChannelCode: "DKI_UUS",
            ChannelType: "Bank",
            Name:        "Bank DKI Syariah (UUS)",
        },
        {
            ChannelCode: "FAMA",
            ChannelType: "Bank",
            Name:        "Bank Fama International",
        },
        {
            ChannelCode: "GANESHA",
            ChannelType: "Bank",
            Name:        "Bank Ganesha",
        },
        {
            ChannelCode: "HANA",
            ChannelType: "Bank",
            Name:        "Bank Hana",
        },
        {
            ChannelCode: "HSBC",
            ChannelType: "Bank",
            Name:        "HSBC Indonesia (formerly called Bank Ekonomi Raharja)",
        },
        {
            ChannelCode: "HSBC_UUS",
            ChannelType: "Bank",
            Name:        "Hongkong and Shanghai Bank Corporation Syariah (HSBC UUS)",
        },
        {
            ChannelCode: "IBK",
            ChannelType: "Bank",
            Name:        "Bank IBK Indonesia (formerly Bank Agris)",
        },
        {
            ChannelCode: "ICBC",
            ChannelType: "Bank",
            Name:        "Bank ICBC Indonesia",
        },
        {
            ChannelCode: "INA_PERDANA",
            ChannelType: "Bank",
            Name:        "Bank Ina Perdania",
        },
        {
            ChannelCode: "INDEX_SELINDO",
            ChannelType: "Bank",
            Name:        "Bank Index Selindo",
        },
        {
            ChannelCode: "INDIA",
            ChannelType: "Bank",
            Name:        "Bank of India Indonesia",
        },
        {
            ChannelCode: "JAGO",
            ChannelType: "Bank",
            Name:        "Bank Jago (formerly Bank Artos Indonesia)",
        },
        {
            ChannelCode: "JAMBI",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jambi",
        },
        {
            ChannelCode: "JAMBI_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jambi Syariah (UUS)",
        },
        {
            ChannelCode: "JASA_JAKARTA",
            ChannelType: "Bank",
            Name:        "Bank Jasa Jakarta",
        },
        {
            ChannelCode: "JAWA_TENGAH",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jawa Tengah",
        },
        {
            ChannelCode: "JAWA_TENGAH_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jawa Tengah Syariah (UUS)",
        },
        {
            ChannelCode: "JAWA_TIMUR",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jawa Timur",
        },
        {
            ChannelCode: "JAWA_TIMUR_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Jawa Timur Syariah (UUS)",
        },
        {
            ChannelCode: "JPMORGAN",
            ChannelType: "Bank",
            Name:        "JP Morgan Chase Bank",
        },
        {
            ChannelCode: "JTRUST",
            ChannelType: "Bank",
            Name:        "Bank JTrust Indonesia (formerly called Bank Mutiara)",
        },
        {
            ChannelCode: "KALIMANTAN_BARAT",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Barat",
        },
        {
            ChannelCode: "KALIMANTAN_BARAT_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Barat Syariah (UUS)",
        },
        {
            ChannelCode: "KALIMANTAN_SELATAN",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Selatan",
        },
        {
            ChannelCode: "KALIMANTAN_SELATAN_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Selatan Syariah (UUS)",
        },
        {
            ChannelCode: "KALIMANTAN_TENGAH",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Tengah",
        },
        {
            ChannelCode: "KALIMANTAN_TIMUR",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Timur",
        },
        {
            ChannelCode: "KALIMANTAN_TIMUR_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Kalimantan Timur Syariah (UUS)",
        },
        {
            ChannelCode: "LAMPUNG",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Lampung",
        },
        {
            ChannelCode: "MALUKU",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Maluku",
        },
        {
            ChannelCode: "MANDIRI",
            ChannelType: "Bank",
            Name:        "Bank Mandiri",
        },
        {
            ChannelCode: "MANDIRI_TASPEN",
            ChannelType: "Bank",
            Name:        "Mandiri Taspen Pos (formerly called Bank Sinar Harapan Bali)",
        },
        {
            ChannelCode: "MASPION",
            ChannelType: "Bank",
            Name:        "Bank Maspion Indonesia",
        },
        {
            ChannelCode: "MAYAPADA",
            ChannelType: "Bank",
            Name:        "Bank Mayapada International",
        },
        {
            ChannelCode: "MAYBANK",
            ChannelType: "Bank",
            Name:        "Bank Maybank",
        },
        {
            ChannelCode: "MAYBANK_SYR",
            ChannelType: "Bank",
            Name:        "Bank Maybank Syariah Indonesia",
        },
        {
            ChannelCode: "MAYORA",
            ChannelType: "Bank",
            Name:        "Bank Mayora",
        },
        {
            ChannelCode: "MEGA",
            ChannelType: "Bank",
            Name:        "Bank Mega",
        },
        {
            ChannelCode: "MEGA_SYR",
            ChannelType: "Bank",
            Name:        "Bank Syariah Mega",
        },
        {
            ChannelCode: "MESTIKA_DHARMA",
            ChannelType: "Bank",
            Name:        "Bank Mestika Dharma",
        },
        {
            ChannelCode: "MIZUHO",
            ChannelType: "Bank",
            Name:        "Bank Mizuho Indonesia",
        },
        {
            ChannelCode: "MNC_INTERNASIONAL",
            ChannelType: "Bank",
            Name:        "Bank MNC Internasional",
        },
        {
            ChannelCode: "MUAMALAT",
            ChannelType: "Bank",
            Name:        "Bank Muamalat Indonesia",
        },
        {
            ChannelCode: "MULTI_ARTA_SENTOSA",
            ChannelType: "Bank",
            Name:        "Bank Multi Arta Sentosa",
        },
        {
            ChannelCode: "NATIONALNOBU",
            ChannelType: "Bank",
            Name:        "Bank Nationalnobu",
        },
        {
            ChannelCode: "NUSA_TENGGARA_BARAT",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Nusa Tenggara Barat",
        },
        {
            ChannelCode: "NUSA_TENGGARA_TIMUR",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Nusa Tenggara Timur",
        },
        {
            ChannelCode: "OCBC",
            ChannelType: "Bank",
            Name:        "Bank OCBC NISP",
        },
        {
            ChannelCode: "OCBC_UUS",
            ChannelType: "Bank",
            Name:        "Bank OCBC NISP Syariah (UUS)",
        },
        {
            ChannelCode: "OKE",
            ChannelType: "Bank",
            Name:        "Bank Oke Indonesia (formerly called Bank Andara)",
        },
        {
            ChannelCode: "PANIN",
            ChannelType: "Bank",
            Name:        "Bank Panin",
        },
        {
            ChannelCode: "PANIN_SYR",
            ChannelType: "Bank",
            Name:        "Bank Panin Syariah",
        },
        {
            ChannelCode: "PAPUA",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Papua",
        },
        {
            ChannelCode: "PERMATA",
            ChannelType: "Bank",
            Name:        "Bank Permata",
        },
        {
            ChannelCode: "PERMATA_UUS",
            ChannelType: "Bank",
            Name:        "Bank Permata Syariah (UUS)",
        },
        {
            ChannelCode: "PRIMA_MASTER",
            ChannelType: "Bank",
            Name:        "Prima Master Bank",
        },
        {
            ChannelCode: "QNB_INDONESIA",
            ChannelType: "Bank",
            Name:        "Bank QNB Indonesia (formerly called Bank QNB Kesawan)",
        },
        {
            ChannelCode: "RABOBANK",
            ChannelType: "Bank",
            Name:        "Bank Rabobank International Indonesia",
        },
        {
            ChannelCode: "RESONA",
            ChannelType: "Bank",
            Name:        "Bank Resona Perdania",
        },
        {
            ChannelCode: "RIAU_DAN_KEPRI",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Riau Dan Kepri",
        },
        {
            ChannelCode: "SAHABAT_SAMPOERNA",
            ChannelType: "Bank",
            Name:        "Bank Sahabat Sampoerna",
        },
        {
            ChannelCode: "SBI_INDONESIA",
            ChannelType: "Bank",
            Name:        "Bank SBI Indonesia",
        },
        {
            ChannelCode: "SEABANK",
            ChannelType: "Bank",
            Name:        "Bank Seabank Indonesia (formerly Bank Kesejahteraan Ekonomi)",
        },
        {
            ChannelCode: "SHINHAN",
            ChannelType: "Bank",
            Name:        "Bank Shinhan Indonesia (formerly called Bank Metro Express)",
        },
        {
            ChannelCode: "SINARMAS",
            ChannelType: "Bank",
            Name:        "Bank Sinarmas",
        },
        {
            ChannelCode: "SINARMAS_UUS",
            ChannelType: "Bank",
            Name:        "Bank Sinarmas Syariah (UUS)",
        },
        {
            ChannelCode: "STANDARD_CHARTERED",
            ChannelType: "Bank",
            Name:        "Standard Charted Bank",
        },
        {
            ChannelCode: "SULAWESI",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sulawesi Tengah",
        },
        {
            ChannelCode: "SULAWESI_TENGGARA",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sulawesi Tenggara",
        },
        {
            ChannelCode: "SULSELBAR",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sulselbar",
        },
        {
            ChannelCode: "SULSELBAR_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sulselbar Syariah (UUS)",
        },
        {
            ChannelCode: "SULUT",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sulut",
        },
        {
            ChannelCode: "SUMATERA_BARAT",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumatera Barat",
        },
        {
            ChannelCode: "SUMATERA_BARAT_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumatera Barat Syariah (UUS)",
        },
        {
            ChannelCode: "SUMSEL_DAN_BABEL",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumsel Dan Babel",
        },
        {
            ChannelCode: "SUMSEL_DAN_BABEL_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumsel Dan Babel Syariah (UUS)",
        },
        {
            ChannelCode: "SUMUT",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumut",
        },
        {
            ChannelCode: "SUMUT_UUS",
            ChannelType: "Bank",
            Name:        "Bank Pembangunan Daerah Sumut Syariah (UUS)",
        },
        {
            ChannelCode: "TABUNGAN_PENSIUNAN_NASIONAL",
            ChannelType: "Bank",
            Name:        "Bank Tabungan Pensiunan Nasional (BTPN)",
        },
        {
            ChannelCode: "TOKYO",
            ChannelType: "Bank",
            Name:        "Bank of Tokyo Mitsubishi UFJ (MUFJ)",
        },
        {
            ChannelCode: "UOB",
            ChannelType: "Bank",
            Name:        "Bank UOB Indonesia",
        },
        {
            ChannelCode: "VICTORIA_INTERNASIONAL",
            ChannelType: "Bank",
            Name:        "Bank Victoria Internasional",
        },
        {
            ChannelCode: "VICTORIA_SYR",
            ChannelType: "Bank",
            Name:        "Bank Victoria Syariah",
        },
        {
            ChannelCode: "WOORI",
            ChannelType: "Bank",
            Name:        "Bank Woori Indonesia",
        },
    }

    return bankChannel
}

func getEWalletChannel() []ChannelCodes {
    ewalletChannel := []ChannelCodes{
        {
            ChannelCode: "DANA",
            ChannelType: "E-Wallet",
            Name:        "DANA",
        },
        {
            ChannelCode: "GOPAY",
            ChannelType: "E-Wallet",
            Name:        "GoPay",
        },
        {
            ChannelCode: "LINKAJA",
            ChannelType: "E-Wallet",
            Name:        "LinkAja",
        },
        {
            ChannelCode: "OVO",
            ChannelType: "E-Wallet",
            Name:        "OVO",
        },
        {
            ChannelCode: "SHOPEEPAY",
            ChannelType: "E-Wallet",
            Name:        "ShopeePay",
        },
    }

    return ewalletChannel
}
