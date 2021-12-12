package avc

import (
	"encoding/binary"
	"io"
)

// 5.3.3.1 AVC decoder configuration record

// This record contains the size of the length field used in each sample to
// indicate the length of its contained NAL units as well as the initial
// parameter sets. This record is externally framed (its size shall be supplied
// by the structure that contains it).
//
// This record contains a version field. This version of the specification
// defines version 1 of this record. Incompatible changes to the record will be
// indicated by a change of version number. Readers shall not attempt to decode
// this record or the streams to which it applies if the version number is
// unrecognized.
//
// Compatible extensions to this record will extend it and will not change the
// configuration version code. Readers should be prepared to ignore unrecognized
// data beyond the definition of the data they understand (e.g. after the
// parameter sets in this specification).
//
// When used to provide the configuration of
//
//   — a parameter set elementary stream, and
//   — a video elementary stream used in conjunction with a parameter set
//     elementary stream,
//
// the configuration record shall contain no sequence or picture parameter sets
// (numOfSequenceParameterSets and numOfPictureParameterSets shall both have the
// value 0).
//
// When used to provide the configuration of a video elementary stream used
// without a parameter set elementary stream, the configuration record may or
// may not contain sequence or picture parameter sets
// (numOfSequenceParameterSets or numOfPictureParameterSets may or may not have
// the value 0).
//
// The values for AVCProfileIndication, AVCLevelIndication, and the flags that
// indicate profile compatibility shall be valid for all parameter sets of the
// stream described by this record. The level indication shall indicate a level
// of capability equal to or greater than the highest level indicated in the
// included parameter sets; each profile compatibility flag may only be set if
// all the included parameter sets set that flag. The profile indication shall
// indicate a profile to which the entire stream associated with this
// configuration record conforms. If the sequence parameter sets are marked with
// different profiles, and the relevant profile compatibility flags are all
// zero, then the stream may need examination to determine which profile, if
// any, the entire stream conforms to. If the entire stream is not examined, or
// the examination reveals that there is no profile to which the entire stream
// conforms, then the stream shall be split into two or more sub-streams with
// separate configuration records in which these rules can be met.
//
// Explicit indication can be provided in the AVC Decoder Configuration Record
// about the chroma format and bit depth used by the AVC video elementary
// stream. The parameter ʹchroma_format_idcʹ present in the sequence parameter
// set in AVC specifies the chroma sampling relative to the luma sampling.
// Similarly the parameters ʹbit_depth_luma_minus8ʹ and
// ʹbit_depth_chroma_minus8ʹ in the sequence parameter set specify the bit depth
// of the samples of the luma and chroma arrays. The values of
// chroma_format_idc, bit_depth_luma_minus8ʹ and ʹbit_depth_chroma_minus8ʹ shall
// be identical in all sequence parameter sets in a single AVC configuration
// record. If two sequences differ in any of these values, two different AVC
// configuration records will be needed. If the two sequences differ in color
// space indications in their VUI information, then two different configuration
// records are also required.
//
// The array of sequence parameter sets, and the array of picture parameter
// sets, may contain SEI messages of a “declarative” nature, that is, those that
// provide information about the stream as a whole. An example of such an SEI is
// a user-data SEI. Such SEIs may also be placed in a parameter set elementary
// stream. NAL unit types that are reserved in ISO/IEC 14496-10 and in this
// specification may acquire a definition in future, and readers should ignore
// NAL units with reserved values of NAL unit type when they are present in
// these arrays.
//
// NOTE 1 This “tolerant” behaviour is designed so that errors are not raised,
// allowing the possibility of backwards-compatible extensions to these arrays
// in future specifications.
//
// When Sequence Parameter Set Extension NAL units occur in this record in
// profiles other than those indicated for the array specific to such NAL units
// (profile_idc not equal to any of 100, 110, 122, 144), they should be placed
// in the Sequence Parameter Set Array.
//
// NOTE 2 The profile identified by profile_idc value 144 is deprecated in
// ISO/IEC 14496-10.
type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion uint8

	// contains the profile code as defined in ISO/IEC 14496-10.
	AVCProfileIndication uint8

	// is a byte defined exactly the same as the byte that occurs between the
	// profile_IDC and level_IDC in a sequence parameter set (SPS), as defined
	// in ISO/IEC 14496-10.
	ProfileCompatibility uint8

	// contains the level code as defined in ISO/IEC 14496-10.
	AVCLevelIndication uint8

	// indicates the length in bytes of the NALUnitLength field in an AVC video
	// sample or AVC parameter set sample of the associated stream minus one.
	// For example, a size of one byte is indicated with a value of 0. The value
	// of this field shall be one of 0, 1, or 3 corresponding to a length
	// encoded with 1, 2, or 4 bytes, respectively.
	LengthSizeMinusOne uint8

	// SPSs that are used as the initial set of SPSs for decoding the AVC
	// elementary stream.
	SequenceParameterSets []AVCSequenceParameterSet

	// picture parameter sets (PPSs) that are used as the initial set of PPSs
	// for decoding the AVC elementary stream.
	PictureParameterSets []AVCPictureParameterSet

	// contains the chroma_format indicator as defined by the chroma_format_idc
	// parameter in ISO/IEC 14496-10.
	ChromaFormat uint8

	// indicates the bit depth of the samples in the Luma arrays. For example, a
	// bit depth of 8 is indicated with a value of zero (BitDepth = 8 +
	// bit_depth_luma_minus8). The value of this field shall be in the range of
	// 0 to 4, inclusive.
	BitDepthLumaMinus8 uint8

	// indicates the bit depth of the samples in the Chroma arrays. For example,
	// a bit depth of 8 is indicated with a value of zero (BitDepth = 8 +
	// bit_depth_chroma_minus8). The value of this field shall be in the range
	// of 0 to 4, inclusive.
	BitDepthChromaMinus8 uint8

	// Sequence Parameter Set Extensions that are used for decoding the AVC
	// elementary stream.
	SequenceParameterSetExts []AVCSequenceParameterSetExt
}

type AVCSequenceParameterSet struct {
	// contains a SPS NAL unit, as specified in ISO/IEC 14496-10. SPSs shall
	// occur in order of ascending parameter set identifier with gaps being
	// allowed.
	NALUnit []byte
}

type AVCPictureParameterSet struct {
	// contains a PPS NAL unit, as specified in ISO/IEC 14496-10. PPSs shall
	// occur in order of ascending parameter set identifier with gaps being
	// allowed.
	NALUnit []byte
}

type AVCSequenceParameterSetExt struct {
	// contains a SPS Extension NAL unit, as specified in ISO/IEC 14496-10.
	NALUnit []byte
}

func (b *AVCDecoderConfigurationRecord) RecordSize() (size uint32) {
	// unsigned int(8) configurationVersion = 1;
	// unsigned int(8) AVCProfileIndication;
	// unsigned int(8) profile_compatibility;
	// unsigned int(8) AVCLevelIndication;
	// bit(6) reserved = '111111'b;
	// unsigned int(2) lengthSizeMinusOne;
	// bit(3) reserved = '111'b;
	// unsigned int(5) numOfSequenceParameterSets;
	size += 6
	// for (i=0; i< numOfSequenceParameterSets; i++) {
	//     unsigned int(16) sequenceParameterSetLength ;
	//     bit(8*sequenceParameterSetLength) sequenceParameterSetNALUnit;
	// }
	for _, sps := range b.SequenceParameterSets {
		size += 2 + uint32(len(sps.NALUnit))
	}
	// unsigned int(8) numOfPictureParameterSets;
	size += 1
	// for (i=0; i< numOfPictureParameterSets; i++) {
	//     unsigned int(16) pictureParameterSetLength;
	//     bit(8*pictureParameterSetLength) pictureParameterSetNALUnit;
	// }
	for _, pps := range b.PictureParameterSets {
		size += 2 + uint32(len(pps.NALUnit))
	}
	if b.AVCProfileIndication == 100 || b.AVCProfileIndication == 110 || b.AVCProfileIndication == 122 || b.AVCProfileIndication == 144 {
		// bit(6) reserved = '111111'b;
		// unsigned int(2) chroma_format;
		// bit(5) reserved = '11111'b;
		// unsigned int(3) bit_depth_luma_minus8;
		// bit(5) reserved = '11111'b;
		// unsigned int(3) bit_depth_chroma_minus8;
		// unsigned int(8) numOfSequenceParameterSetExt;
		size += 4
		// for (i=0; i< numOfSequenceParameterSetExt; i++) {
		//     unsigned int(16) sequenceParameterSetExtLength;
		//     bit(8*sequenceParameterSetExtLength) sequenceParameterSetExtNALUnit;
		// }
		for _, spse := range b.SequenceParameterSetExts {
			size += 2 + uint32(len(spse.NALUnit))
		}
	}
	return
}

func (b *AVCDecoderConfigurationRecord) RecordRead(r io.Reader) (err error) {
	var tmp [6]uint8
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	b.ConfigurationVersion = tmp[0]
	b.AVCProfileIndication = tmp[1]
	b.ProfileCompatibility = tmp[2]
	b.AVCLevelIndication = tmp[3]
	b.LengthSizeMinusOne = tmp[4] & 0b11
	numOfSequenceParameterSets := tmp[5] & 0b11111
	b.SequenceParameterSets = make([]AVCSequenceParameterSet, numOfSequenceParameterSets)
	for i := uint8(0); i < numOfSequenceParameterSets; i++ {
		var sequenceParameterSetLength uint16
		if err = binary.Read(r, binary.BigEndian, &sequenceParameterSetLength); err != nil {
			return
		}
		b.SequenceParameterSets[i].NALUnit = make([]byte, sequenceParameterSetLength)
		if _, err = io.ReadFull(r, b.SequenceParameterSets[i].NALUnit); err != nil {
			return
		}
	}
	var numOfPictureParameterSets uint8
	if err = binary.Read(r, binary.BigEndian, &numOfPictureParameterSets); err != nil {
		return
	}
	b.PictureParameterSets = make([]AVCPictureParameterSet, numOfPictureParameterSets)
	for i := uint8(0); i < numOfPictureParameterSets; i++ {
		var pictureParameterSetLength uint16
		if err = binary.Read(r, binary.BigEndian, &pictureParameterSetLength); err != nil {
			return
		}
		b.PictureParameterSets[i].NALUnit = make([]byte, pictureParameterSetLength)
		if _, err = io.ReadFull(r, b.PictureParameterSets[i].NALUnit); err != nil {
			return
		}
	}
	if b.AVCProfileIndication == 100 || b.AVCProfileIndication == 110 || b.AVCProfileIndication == 122 || b.AVCProfileIndication == 144 {
		if err = binary.Read(r, binary.BigEndian, tmp[:4]); err != nil {
			return
		}
		b.ChromaFormat = tmp[0] & 0b11
		b.BitDepthLumaMinus8 = tmp[1] & 0b111
		b.BitDepthChromaMinus8 = tmp[2] & 0b111
		numOfSequenceParameterSetExt := tmp[3]
		b.SequenceParameterSetExts = make([]AVCSequenceParameterSetExt, numOfSequenceParameterSetExt)
		for i := uint8(0); i < numOfSequenceParameterSetExt; i++ {
			var sequenceParameterSetExtLength uint16
			if err = binary.Read(r, binary.BigEndian, &sequenceParameterSetExtLength); err != nil {
				return
			}
			b.SequenceParameterSetExts[i].NALUnit = make([]byte, sequenceParameterSetExtLength)
			if _, err = io.ReadFull(r, b.SequenceParameterSetExts[i].NALUnit); err != nil {
				return
			}
		}
	}
	return
}

func (b *AVCDecoderConfigurationRecord) RecordWrite(w io.Writer) (err error) {
	if err = binary.Write(w, binary.BigEndian, b.ConfigurationVersion); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.AVCProfileIndication); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.ProfileCompatibility); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.AVCLevelIndication); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.LengthSizeMinusOne|0b11111100); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint8(len(b.SequenceParameterSets))|0b11100000); err != nil {
		return
	}
	for i := 0; i < len(b.SequenceParameterSets); i++ {
		if err = binary.Write(w, binary.BigEndian, uint16(len(b.SequenceParameterSets[i].NALUnit))); err != nil {
			return
		}
		if _, err = w.Write(b.SequenceParameterSets[i].NALUnit); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, uint8(len(b.PictureParameterSets))); err != nil {
		return
	}
	for i := 0; i < len(b.PictureParameterSets); i++ {
		if err = binary.Write(w, binary.BigEndian, uint16(len(b.PictureParameterSets[i].NALUnit))); err != nil {
			return
		}
		if _, err = w.Write(b.PictureParameterSets[i].NALUnit); err != nil {
			return
		}
	}
	if b.AVCProfileIndication == 100 || b.AVCProfileIndication == 110 || b.AVCProfileIndication == 122 || b.AVCProfileIndication == 144 {
		if err = binary.Write(w, binary.BigEndian, b.ChromaFormat|0b11111100); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.BitDepthLumaMinus8|0b11111000); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.BitDepthChromaMinus8|0b11111000); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint8(len(b.SequenceParameterSetExts))); err != nil {
			return
		}
		for i := 0; i < len(b.SequenceParameterSetExts); i++ {
			if err = binary.Write(w, binary.BigEndian, uint16(len(b.SequenceParameterSetExts[i].NALUnit))); err != nil {
				return
			}
			if _, err = w.Write(b.SequenceParameterSetExts[i].NALUnit); err != nil {
				return
			}
		}
	}
	return
}
