package models

type ApiResponse struct {
	Response Response
}

type Response struct {
	ResultCd  string       // 결과코드
	ResultMsg string       // 결과메세지
	NumOfRows string       // 한 페이지 결과 수
	Items     []FlowerItem // 꽃 항목들
}

type FlowerItem struct {
	SaleDate   string // 경매일자
	FlowerGubn string // 화훼부류명
	PumName    string // 품목명
	GoodName   string // 품종명
	LvNm       string // 등급명
	MaxAmt     string // 최고가
	MinAmt     string // 최저가
	AvgAmt     string // 평균가
	TotAmt     string // 총금액
	TotQty     string // 총물량
}
