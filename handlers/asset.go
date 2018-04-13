package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/utils"
)

func setAssetMux() {
	mux.PostFunc("/assets", commonHandler(postAsset))
	mux.GetFunc("/assets/:filename", commonHandler(getAsset))
}

func postAsset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:     "MultipartForm parse error. (Create asset item)",
			ErrorName: models.ERROR_NAME_INVALID_JSON,
		}
		respondErr(w, r, http.StatusBadRequest, pd)
		return
	}

	file, header, err := r.FormFile("asset")
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "asset",
					Reason: "asset is required, but it's empty.",
				},
			},
		}
		respondErr(w, r, http.StatusBadRequest, pd)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	asset, pd := services.PostAsset(r.Context(), contentType, file)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", asset)
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	filename := bone.GetValue(r, "filename")
	assetID := utils.GetFileNameWithoutExt(filename)
	ifModifiedSince := r.Header.Get("If-Modified-Since")

	bytes, asset, pd := services.GetAsset(r.Context(), assetID, ifModifiedSince)
	if pd != nil {
		respondErr(w, r, pd.Status, pd)
		return
	}

	setLastModified(w, asset.Modified)
	// w.Header().Set("Cache-Control", "max-age:86400, public")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Set("Content-Type", http.DetectContentType(bytes))
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
