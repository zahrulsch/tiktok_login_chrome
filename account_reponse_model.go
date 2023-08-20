package tiktokloginchrome

type AccountResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Seller struct {
			Status       int    `json:"status"`
			SellerType   int    `json:"seller_type"`
			BusinessType int    `json:"business_type"`
			BaseGeoIDL0  string `json:"base_geo_id_l0"`
			ShopName     string `json:"shop_name"`
			SellerID     string `json:"seller_id"`
			ShopCode     string `json:"shop_code"`
			SellerName   string `json:"seller_name"`
			RegionCode   string `json:"region_code"`
			Logo         struct {
				Height  int      `json:"height"`
				Width   int      `json:"width"`
				URLList []string `json:"url_list"`
			} `json:"logo"`
		} `json:"seller"`
		LogoAuditStatus int `json:"logo_audit_status"`
		AuditLogo       struct {
			Height  int      `json:"height"`
			Width   int      `json:"width"`
			URLList []string `json:"url_list"`
		} `json:"audit_logo"`
		IsLogoOverFrequency        bool   `json:"is_logo_over_frequency"`
		ShopNameAuditStatus        int    `json:"shop_name_audit_status"`
		AuditShopName              string `json:"audit_shop_name"`
		IsShopNameOverFrequency    bool   `json:"is_shop_name_over_frequency"`
		ShopNameLastNNaturalMonths int    `json:"shop_name_last_n_natural_months"`
		ShopNameFreq               int    `json:"shop_name_freq"`
		ShopNameFreqClose          bool   `json:"shop_name_freq_close"`
	} `json:"data"`
}
