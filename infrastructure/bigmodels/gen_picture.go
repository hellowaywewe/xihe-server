package bigmodels

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
)

func (s *service) GenPicture(user domain.Account, desc string) (string, error) {
	r := new(singlePicture)

	err := s.genPicture(user, desc, s.singlePictures, r)
	if err != nil {
		return "", err
	}

	return r.picture()
}

func (s *service) GenPictures(user domain.Account, desc string) ([]string, error) {
	r := new(multiplePictures)

	err := s.genPicture(user, desc, s.multiplePictures, r)
	if err != nil {
		return nil, err
	}

	return r.picture()
}

func (s *service) genPicture(
	user domain.Account, desc string,
	ec chan string, result interface{},
) error {
	select {
	case e := <-ec:
		err := s.sendReqToGenPicture(user, e, desc, result)
		ec <- e

		return err

	default:
		return errors.New("busy")
	}
}

type pictureGenerateOpt struct {
	Desc string `json:"input_text"`
	User string `json:"user_name"`
}

type singlePicture struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Output string `json:"output_image_url"`
}

func (p *singlePicture) picture() (string, error) {
	if p.Status == -1 {
		return "", errors.New(p.Msg)
	}

	return p.Output, nil
}

type multiplePictures struct {
	Status int      `json:"status"`
	Msg    string   `json:"msg"`
	Output []string `json:"output_image_url"`
}

func (p *multiplePictures) picture() ([]string, error) {
	if p.Status == -1 {
		return nil, errors.New(p.Msg)
	}

	return p.Output, nil
}

func (s *service) sendReqToGenPicture(
	user domain.Account, endpoint, desc string, r interface{},
) (err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	opt := pictureGenerateOpt{
		Desc: desc,
		User: user.Account(),
	}

	body, err := utils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	_, err = s.hc.ForwardTo(req, r)

	return err
}
