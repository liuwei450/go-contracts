package router

import (
	"encoding/json"
	"net/http"
	"go-contracts/service"
	"go-contracts/util"
)

const (
	// ERC20相关API路由
	ERC20_ALLOWANCE     = "/api/erc20/allowance"
	ERC20_APPROVE       = "/api/erc20/approve"
	ERC20_TRANSFER      = "/api/erc20/transfer"
	ERC20_TRANSFER_FROM = "/api/erc20/transfer_from"
	ERC20_BALANCE       = "/api/erc20/balance"
	ERC20_TOTAL_SUPPLY  = "/api/erc20/total_supply"
	ERC20_TOKEN_INFO    = "/api/erc20/token_info"
)

// ERC20Allowance 处理ERC20授权查询请求
func (h Routes) ERC20Allowance(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20AllowanceParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20Allowance(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20Approve 处理ERC20授权请求
func (h Routes) ERC20Approve(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20ApproveParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20Approve(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20Transfer 处理ERC20转账请求
func (h Routes) ERC20Transfer(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20TransferParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20Transfer(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20TransferFrom 处理ERC20授权转账请求
func (h Routes) ERC20TransferFrom(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20TransferFromParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20TransferFrom(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20Balance 处理ERC20余额查询请求
func (h Routes) ERC20Balance(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20BalanceParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20Balance(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20TotalSupply 处理ERC20总供应量查询请求
func (h Routes) ERC20TotalSupply(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20ContractParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20TotalSupply(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result.String())
}

// ERC20TokenInfo 处理ERC20代币信息查询请求
func (h Routes) ERC20TokenInfo(w http.ResponseWriter, r *http.Request) {
	var params service.ERC20ContractParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlerError(w, err)
		return
	}

	result, err := h.svc.ERC20TokenInfo(r.Context(), params)
	if err != nil {
		handlerError(w, err)
		return
	}

	handlerSuccess(w, result)
}

// 通用响应处理函数
func handlerSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func handlerError(w http.ResponseWriter, err error) {
	util.Log.Error("API error", "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 500,
		"msg":  err.Error(),
	})
}