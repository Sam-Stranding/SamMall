package redis

//type INews interface {
//	GetLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string, captcha string) (*do.LarkSmsCodeResp, error)
//}
//
//type News struct {
//	redis *redis.Client
//}
//
//func NewNews(adaptor adaptor.IAdaptor) *News {
//	return &News{
//		redis: adaptor.GetRedis(),
//	}
//}
//
//func (n *News) GetLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string, captcha string) (*do.LarkSmsCodeResp, error) {
//	url := fmt.Sprintf("%s/open-apis/im/v1/messages?receive_id_type=open_id", larkHost)
//	headers := map[string]string{
//		"Authorization": "Bearer " + TenantAccessToken,
//		"Context-Type":  headerCT,
//	}
//	body := map[string]string{
//		"receive_id": UserOpenID,
//		"msg_type":   "text",
//		"content":    fmt.Sprintf(`{"text":"%s"}`, captcha),
//	}
//	resp := &do.LarkSmsCodeResp{}
//	err := http.Post(ctx, url, headers, body, resp)
//	if err != nil {
//		logger.Error("GetLarkSmsCode error", zap.Error(err))
//		return nil, err
//	}
//	return resp, nil
//}
