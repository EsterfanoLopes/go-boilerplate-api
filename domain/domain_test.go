package domain_test

import (
	"go-boilerplate/domain"
	"go-boilerplate/test"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIsValidWeekday(t *testing.T) {
	testCases := []struct {
		name     string
		weekday  string
		expected bool
	}{
		{
			name:     "valid weekday",
			weekday:  "MONDAY",
			expected: true,
		},
		{
			name:     "invalid weekday",
			weekday:  "XXX",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.IsValidWeekday(tc.weekday)
			if result != tc.expected {
				t.Errorf("unexpected weekday validation result %t", result)
				return
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	testCases := []struct {
		name          string
		address       *domain.Address
		expectedError string
	}{
		{
			name: "valid address",
			address: &domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "195",
				Complement:   "Next to the most expansive bakery in the world",
				ZipCode:      "01415001",
			},
		},
		{
			name: "missing zip code",
			address: &domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "195",
			},
			expectedError: "zipCode: cannot be blank.",
		},
		{
			name: "missing state",
			address: &domain.Address{
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "195",
				Complement:   "Next to the most expansive bakery in the world",
				ZipCode:      "01415001",
			},
			expectedError: "state: cannot be blank.",
		},
		{
			name: "missing city",
			address: &domain.Address{
				State:        "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "195",
				Complement:   "Next to the most expansive bakery in the world",
				ZipCode:      "01415001",
			},
			expectedError: "city: cannot be blank.",
		},
		{
			name: "missing street",
			address: &domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				StreetNumber: "195",
				ZipCode:      "01415001",
			},
			expectedError: "street: cannot be blank.",
		},
		{
			name: "missing street number",
			address: &domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				ZipCode:      "01415001",
			},
			expectedError: "streetNumber: cannot be blank.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.address.Validate()
			test.AssertError(t, err, tc.expectedError)
		})
	}
}

func TestIsFromPortals(t *testing.T) {
	testCases := []struct {
		Origin   domain.Origin
		expected bool
	}{
		{
			Origin:   domain.OriginNone,
			expected: false,
		},
		{
			Origin:   domain.Sms,
			expected: false,
		},
		{
			Origin:   domain.EmailMarketing,
			expected: false,
		},
		{
			Origin:   domain.ExternalAdvertising,
			expected: false,
		},
		{
			Origin:   domain.ActiveOffer,
			expected: false,
		},
		{
			Origin:   domain.Telephone,
			expected: false,
		},
		{
			Origin:   domain.Recommendation,
			expected: false,
		},
		{
			Origin:   domain.AdvertiserSite,
			expected: false,
		},
		{
			Origin:   domain.Other,
			expected: false,
		},
		{
			Origin:   domain.VivaReal,
			expected: true,
		},
		{
			Origin:   domain.Zap,
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Origin.String(), func(t *testing.T) {
			result := tc.Origin.IsFromPortals()
			if result != tc.expected {
				t.Errorf("unexpected is from portals result %t", result)
				return
			}
		})
	}
}

func TestParseAccessToken(t *testing.T) {
	testCases := []struct {
		name         string
		accountID    string
		advertiserID string
		JWTSecret    string
		expHours     int
	}{
		{
			name:         "valid token",
			accountID:    "4bdcacda-549c-44e3-b763-f2f7d4c73252",
			advertiserID: "ab5c0ff0-de16-b6d5-d0a5-3bd3e25b6b46",
			JWTSecret:    "test-secret",
			expHours:     1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := domain.GenerateAccessToken(tc.accountID, tc.advertiserID, tc.JWTSecret, tc.expHours)
			if err != nil {
				t.Errorf("unexpected error generating document token %s", err)
				return
			}

			accountID, advertiserID, err := domain.ParseAccessToken(token, tc.JWTSecret)
			if err != nil {
				t.Errorf("unexpected error parsing document token %s", err)
				return
			}
			if accountID != tc.accountID {
				t.Errorf("unexpected account id %s", accountID)
				return
			}
			if advertiserID != tc.advertiserID {
				t.Errorf("unexpected account id %s", accountID)
				return
			}
		})
	}
}

func TestAddressAnonymize(t *testing.T) {
	testCases := []struct {
		name     string
		address  domain.Address
		expected domain.Address
	}{
		{
			name: "full address",
			address: domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "195",
				Complement:   "CS-1",
				ZipCode:      "01415001",
			},
			expected: domain.Address{
				State:        "Sao Paulo",
				City:         "Sao Paulo",
				Neighborhood: "Bela Vista",
				Street:       "Bela Cintra",
				StreetNumber: "***",
				Complement:   "****",
				ZipCode:      "01415001",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.address.Anonymize()
			if diff := cmp.Diff(result, tc.expected); diff != "" {
				t.Errorf("unexpected anonymous address %s", diff)
				return
			}

		})
	}
}

func TestValidateTenantInfo(t *testing.T) {
	testCases := []struct {
		name       string
		tenantInfo domain.TenantInfo
		expected   string
	}{
		{
			name: "valid tenant info",
			tenantInfo: domain.TenantInfo{
				Adults:   1,
				LiveWith: domain.Alone,
			},
		},
		{
			name: "no adults",
			tenantInfo: domain.TenantInfo{
				Children: 0,
				LiveWith: domain.Alone,
				Pets:     0,
			},
			expected: "adults: cannot be blank.",
		},
		{
			name: "missing live with type",
			tenantInfo: domain.TenantInfo{
				Adults: 1,
			},
			expected: "liveWith: cannot be blank.",
		},
		{
			name: "negative numbers",
			tenantInfo: domain.TenantInfo{
				Adults:   -1,
				Children: -1,
				Pets:     -1,
				LiveWith: domain.Alone,
			},
			expected: "adults: must be no less than 1; children: must be no less than 0; pets: must be no less than 0.",
		},
		{
			name: "pets without description",
			tenantInfo: domain.TenantInfo{
				Adults:   1,
				LiveWith: domain.Alone,
				Pets:     1,
			},
			expected: "petsDescription: cannot be blank.",
		},
		{
			name: "pets description without pets",
			tenantInfo: domain.TenantInfo{
				Adults:          1,
				LiveWith:        domain.Alone,
				PetsDescription: "cats n the cradle",
			},
			expected: "pets: cannot be blank.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.tenantInfo.Validate()
			test.AssertError(t, err, tc.expected)
		})
	}
}

func TestHost(t *testing.T) {
	testCases := []struct {
		listingOrigin domain.ListingOrigin
		expected      string
	}{
		{
			listingOrigin: domain.PortalZap,
			expected:      "www.zapimoveis.com.br",
		},
		{
			listingOrigin: domain.PortalVivaReal,
			expected:      "www.vivareal.com.br",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.listingOrigin.String(), func(t *testing.T) {
			result := tc.listingOrigin.Host()
			if result != tc.expected {
				t.Errorf("unexpected host %s", tc.expected)
				return
			}
		})
	}
}

func TestFredoURL(t *testing.T) {
	testCases := []struct {
		name          string
		listingOrigin domain.ListingOrigin
		expected      string
	}{
		{
			name:          "fredo url for portal viva real",
			listingOrigin: domain.PortalVivaReal,
			expected:      "https://qa-negociacao.vivareal.com.br",
		},
		{
			name:          "fredo url for portal zap",
			listingOrigin: domain.PortalZap,
			expected:      "https://qa-negociacao.zapimoveis.com.br",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.listingOrigin.FredoURL()
			if result != tc.expected {
				t.Errorf("unexpected fredo url %s", result)
			}
		})
	}
}

func TestPortalURL(t *testing.T) {
	testCases := []struct {
		name          string
		listingOrigin domain.ListingOrigin
		expected      string
	}{
		{
			name:          "portal url for viva real",
			listingOrigin: domain.PortalVivaReal,
			expected:      "https://www.vivareal.com.br",
		},
		{
			name:          "portal url for zap",
			listingOrigin: domain.PortalZap,
			expected:      "https://www.zapimoveis.com.br",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.listingOrigin.PortalURL()
			if result != tc.expected {
				t.Errorf("unexpected portal url %s", result)
			}
		})
	}
}

func TestTemperatureFromScore(t *testing.T) {
	testCases := []struct {
		name     string
		score    int
		expected domain.Temperature
	}{
		{
			name:     "hot temperature",
			score:    70,
			expected: domain.TemperatureHot,
		},
		{
			name:     "warm temperature start",
			score:    50,
			expected: domain.TemperatureWarm,
		},
		{
			name:     "warm temperature end",
			score:    69,
			expected: domain.TemperatureWarm,
		},
		{
			name:     "cold temperature",
			score:    49,
			expected: domain.TemperatureCold,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.TemperatureFromScore(tc.score)
			if diff := cmp.Diff(result, tc.expected); diff != "" {
				t.Errorf("unexpected temperature %s", diff)
				return
			}
		})
	}
}

func TestTemperatureScoreRange(t *testing.T) {
	testCases := []struct {
		temperature   domain.Temperature
		expectedStart int
		expectedEnd   int
	}{
		{
			temperature:   domain.TemperatureHot,
			expectedStart: 70,
		},
		{
			temperature:   domain.TemperatureWarm,
			expectedStart: 50,
			expectedEnd:   69,
		},
		{
			temperature: domain.TemperatureCold,
			expectedEnd: 49,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.temperature.String(), func(t *testing.T) {
			start, end := tc.temperature.ScoreRange()
			if start != tc.expectedStart {
				t.Errorf("unexpected temperature score start %d", start)
			}

			if end != tc.expectedEnd {
				t.Errorf("unexpected temperature score end %d", start)
			}
		})
	}
}
