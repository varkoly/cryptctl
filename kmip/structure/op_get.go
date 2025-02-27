// cryptctl - Copyright (c) 2017 SUSE Linux GmbH, Germany
// This source code is licensed under GPL version 3 that can be found in LICENSE file.
package structure

import (
	"cryptctl/kmip/ttlv"
	"errors"
	"fmt"
)

// KMIP request message 420078
type SGetRequest struct {
	SRequestHeader    SRequestHeader    // IBatchCount is assumed to be 1 in serialisation operations
	SRequestBatchItem SRequestBatchItem // payload is SRequestPayloadGet
}

func (getReq *SGetRequest) SerialiseToTTLV() ttlv.Item {
	getReq.SRequestHeader.IBatchCount.Value = 1
	ret := ttlv.NewStructure(TagRequestMessage, getReq.SRequestHeader.SerialiseToTTLV(), getReq.SRequestBatchItem.SerialiseToTTLV())
	return ret
}
func (getReq *SGetRequest) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagRequestMessage, TagRequestHeader, &getReq.SRequestHeader); err != nil {
		return err
	}
	if val := getReq.SRequestHeader.IBatchCount.Value; val != 1 {
		return fmt.Errorf("SGetRequest.DeserialiseFromTTLV: was expecting exactly 1 item, but received %d instead.", val)
	}
	getReq.SRequestBatchItem = SRequestBatchItem{SRequestPayload: &SRequestPayloadGet{}}
	if err := DecodeStructItem(in, TagRequestMessage, TagBatchItem, &getReq.SRequestBatchItem); err != nil {
		return err
	}
	if getReq.SRequestBatchItem.EOperation.Value != ValOperationGet {
		return errors.New("SGetRequest.DeserialiseFromTTLV: input is not a get request")
	}
	return nil
}

// 420079 - request payload from a get request
type SRequestPayloadGet struct {
	TUniqueID ttlv.Text // 420094
}

func (getPayload *SRequestPayloadGet) SerialiseToTTLV() ttlv.Item {
	getPayload.TUniqueID.Tag = TagUniqueID
	return ttlv.NewStructure(TagRequestPayload, &getPayload.TUniqueID)
}
func (getPayload *SRequestPayloadGet) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagRequestPayload, TagUniqueID, &getPayload.TUniqueID); err != nil {
		return err
	}
	return nil
}

// KMIP response message 42007b
type SGetResponse struct {
	SResponseHeader    SResponseHeader    // IBatchCount is assumed to be 1 in serialisation operations
	SResponseBatchItem SResponseBatchItem // payload is SResponsePayloadGet
}

func (getResp *SGetResponse) SerialiseToTTLV() ttlv.Item {
	getResp.SResponseHeader.IBatchCount.Value = 1
	ret := ttlv.NewStructure(TagResponseMessage, getResp.SResponseHeader.SerialiseToTTLV(), getResp.SResponseBatchItem.SerialiseToTTLV())
	return ret
}
func (getResp *SGetResponse) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagResponseMessage, TagResponseHeader, &getResp.SResponseHeader); err != nil {
		return err
	}
	if val := getResp.SResponseHeader.IBatchCount.Value; val != 1 {
		return fmt.Errorf("SGetResponse.DeserialiseFromTTLV: was expecting exactly 1 item, but received %d instead.", val)
	}
	getResp.SResponseBatchItem = SResponseBatchItem{SResponsePayload: &SResponsePayloadGet{}}
	if err := DecodeStructItem(in, TagResponseMessage, TagBatchItem, &getResp.SResponseBatchItem); err != nil {
		return err
	}
	if getResp.SResponseBatchItem.EOperation.Value != ValOperationGet {
		return errors.New("SGetResponse.DeserialiseFromTTLV: input is not a get response")
	}
	return nil
}

// 42007c - response payload from a get response
type SResponsePayloadGet struct {
	EObjectType   ttlv.Enumeration // 420057
	TUniqueID     ttlv.Text        // 420094
	SSymmetricKey SSymmetricKey    // 42008f
}

func (getPayload *SResponsePayloadGet) SerialiseToTTLV() ttlv.Item {
	getPayload.EObjectType.Tag = TagObjectType
	getPayload.TUniqueID.Tag = TagUniqueID
	return ttlv.NewStructure(TagResponsePayload, &getPayload.EObjectType, &getPayload.TUniqueID, getPayload.SSymmetricKey.SerialiseToTTLV())
}
func (getPayload *SResponsePayloadGet) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagResponsePayload, TagObjectType, &getPayload.EObjectType); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagResponsePayload, TagUniqueID, &getPayload.TUniqueID); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagResponsePayload, TagSymmetricKey, &getPayload.SSymmetricKey); err != nil {
		return err
	}
	return nil
}

// 42008f
type SSymmetricKey struct {
	SKeyBlock SKeyBlock
}

func (symKey *SSymmetricKey) SerialiseToTTLV() ttlv.Item {
	return ttlv.NewStructure(TagSymmetricKey, symKey.SKeyBlock.SerialiseToTTLV())
}
func (symKey *SSymmetricKey) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagSymmetricKey, TagKeyBlock, &symKey.SKeyBlock); err != nil {
		return err
	}
	return nil
}

// 420040
type SKeyBlock struct {
	EFormatType      ttlv.Enumeration // 420042
	SKeyValue        SKeyValue
	ECryptoAlgorithm ttlv.Enumeration // 420028
	ECryptoLen       ttlv.Integer     // 42002a
}

func (block *SKeyBlock) SerialiseToTTLV() ttlv.Item {
	block.EFormatType.Tag = TagFormatType
	block.ECryptoAlgorithm.Tag = TagCryptoAlgorithm
	block.ECryptoLen.Tag = TagCryptoLen
	return ttlv.NewStructure(TagKeyBlock, &block.EFormatType, block.SKeyValue.SerialiseToTTLV(), &block.ECryptoAlgorithm, &block.ECryptoLen)
}
func (block *SKeyBlock) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagKeyBlock, TagFormatType, &block.EFormatType); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagKeyBlock, TagKeyValue, &block.SKeyValue); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagKeyBlock, TagCryptoAlgorithm, &block.ECryptoAlgorithm); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagKeyBlock, TagCryptoLen, &block.ECryptoLen); err != nil {
		return err
	}
	return nil
}

// 420045 - this is value of an encryption key, not to be confused with a key-value pair.
type SKeyValue struct {
	BKeyMaterial ttlv.Bytes // 420043
}

func (key *SKeyValue) SerialiseToTTLV() ttlv.Item {
	key.BKeyMaterial.Tag = TagKeyMaterial
	return ttlv.NewStructure(TagKeyValue, &key.BKeyMaterial)
}

func (key *SKeyValue) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagKeyValue, TagKeyMaterial, &key.BKeyMaterial); err != nil {
		return err
	}
	return nil
}
