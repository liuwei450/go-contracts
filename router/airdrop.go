package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"go-contracts/service"
	"go-contracts/util"
)

// AirdropBnb 处理BNB空投请求
func (h Routes) AirdropBnb(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 2. 解析请求参数
	var params service.AirdropParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		util.Log.Error("解析请求参数失败", "error", err)
		h.handleError(w, http.StatusBadRequest, "无效的请求参数格式: %v", err)
		return
	}
	defer r.Body.Close()

	// 3. 创建上下文
	ctx := r.Context()

	// 4. 调用服务层的AirdropBnb方法
	if err := h.svc.AirdropBnb(ctx, params); err != nil {
		util.Log.Error("BNB空投失败", "error", err)
		h.handleError(w, http.StatusInternalServerError, "执行空投失败: %v", err)
		return
	}

	// 5. 返回成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "BNB空投请求已提交成功",
		"data":    nil,
	})
}

// handleError 统一处理错误响应
func (h Routes) handleError(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    statusCode,
		"message": fmt.Sprintf(format, args...),
		"data":    nil,
	})
}

// AirdropERC20 处理ERC20空投请求
func (h Routes) AirdropERC20(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 2. 解析请求参数
	var params service.AirdropParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		util.Log.Error("解析请求参数失败", "error", err)
	h.handleError(w, http.StatusBadRequest, "无效的请求参数格式: %v", err)
	return
	}
	defer r.Body.Close()

	// 3. 创建上下文
	ctx := r.Context()

	// 4. 调用服务层的AirdropERC20方法
	if err := h.svc.AirdropERC20(ctx, params); err != nil {
		util.Log.Error("ERC20空投失败", "error", err)
	h.handleError(w, http.StatusInternalServerError, "执行空投失败: %v", err)
	return
	}

	// 5. 返回成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "ERC20空投请求已提交成功",
		"data":    nil,
	})
}
