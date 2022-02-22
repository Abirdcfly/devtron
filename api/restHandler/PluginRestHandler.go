/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package restHandler

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PluginRestHandler interface {
	SavePlugin(w http.ResponseWriter, r *http.Request)
	UpdatePlugin(w http.ResponseWriter, r *http.Request)
	FindByPlugin(w http.ResponseWriter, r *http.Request)
	DeletePlugin(w http.ResponseWriter, r *http.Request)
}

type PluginRestHandlerImpl struct {
	logger     *zap.SugaredLogger
	repository repository.PluginRepository
}

type plugin struct {
	Id                   int    `json:"Id"`
	Name                 string `json:"Name"`
	Description          string `json:"Description"`
	Body                 string `json:"Body"`
	StepTemplateLanguage string `json:"StepTemplateLanguage"`
	StepTemplate         string `json:"StepTemplate"`
}

type pluginInputs struct {
	Id          int    `json:"Id"`
	Name        string `json:"Name"`
	Value       string `json:"Value"`
	Description string `json:"Description"`
}

func NewPluginRestHandlerImpl(logger *zap.SugaredLogger, repository repository.PluginRepository) *PluginRestHandlerImpl {
	pluginRestHandler := &PluginRestHandlerImpl{
		logger:     logger,
		repository: repository,
	}
	return pluginRestHandler
}

func (handler PluginRestHandlerImpl) SavePlugin(w http.ResponseWriter, r *http.Request) {
	//for checking
	decoder := json.NewDecoder(r.Body)
	println(decoder)
	var bean plugin
	err := decoder.Decode(&bean)
	if err != nil {
		common.WriteJsonResp(w, err, "Plugin Id couldn't be parsed from input", http.StatusBadRequest)
	}
	test := &repository.Plugin{
		Id:                   bean.Id,
		Name:                 bean.Name,
		Description:          bean.Description,
		Body:                 bean.Body,
		StepTemplateLanguage: bean.StepTemplateLanguage,
		StepTemplate:         bean.StepTemplate,
	}
	err = handler.repository.Save(test)
	if err != nil {
		common.WriteJsonResp(w, err, "Plugin couldn't be saved", http.StatusInternalServerError)
	}
	common.WriteJsonResp(w, err, "Stored Successfully", http.StatusOK)
}

func (handler PluginRestHandlerImpl) UpdatePlugin(w http.ResponseWriter, r *http.Request) {
	//for checking
	decoder := json.NewDecoder(r.Body)
	println(decoder)
	var bean plugin
	err := decoder.Decode(&bean)
	if err != nil {
		handler.logger.Errorw("decode err", "err", err)
		common.WriteJsonResp(w, err, "Plugin Id couldn't be parsed from input", http.StatusBadRequest)
	}

	test := &repository.Plugin{
		Id:                   bean.Id,
		Name:                 bean.Name,
		Description:          bean.Description,
		Body:                 bean.Body,
		StepTemplateLanguage: bean.StepTemplateLanguage,
		StepTemplate:         bean.StepTemplate,
	}
	err = handler.repository.Update(test)
	if err != nil {
		common.WriteJsonResp(w, err, "Plugin couldn't be updated", http.StatusInternalServerError)
	}
	common.WriteJsonResp(w, err, "Update Successful", http.StatusOK)
}

func (handler PluginRestHandlerImpl) FindByPlugin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	/* #nosec */
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		handler.logger.Errorw("decode err", "err", err)
		common.WriteJsonResp(w, err, "Plugin Id couldn't be parsed from input", http.StatusBadRequest)
	}
	values, err := handler.repository.FindByAppId(id)
	if err != nil {
		common.WriteJsonResp(w, err, "Plugin not found", http.StatusInternalServerError)
	}
	common.WriteJsonResp(w, err, values, http.StatusOK)
}

func (handler PluginRestHandlerImpl) DeletePlugin(w http.ResponseWriter, r *http.Request) {
	//for checking
	vars := mux.Vars(r)
	/* #nosec */
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		handler.logger.Errorw("decode err", "err", err)
		common.WriteJsonResp(w, err, "Plugin Id couldn't be parsed from input", http.StatusBadRequest)
	}

	err = handler.repository.Delete(id)
	if err != nil {
		common.WriteJsonResp(w, err, "Plugin couldn't be deleted", http.StatusInternalServerError)
	}
	common.WriteJsonResp(w, err, "Delete Successful", http.StatusOK)
}