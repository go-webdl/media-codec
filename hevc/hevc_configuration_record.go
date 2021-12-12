package hevc

import (
	"encoding/binary"
	"io"
)

// 8.3.3.1 HEVC decoder configuration record

// This subclause specifies the decoder configuration information for ISO/IEC
// 23008-2 video content. This record contains the size of the length field used
// in each sample to indicate the length of its contained NAL units as well as
// the parameter sets, if stored in the sample entry. This record is externally
// framed (its size shall be supplied by the structure that contains it).
//
// This record contains a version field. This version of the specification
// defines version 1 of this record. Incompatible changes to the record will be
// indicated by a change of version number. Readers shall not attempt to decode
// this record or the streams to which it applies if the version number is
// unrecognized.
//
// Compatible extensions to this record will extend it and will not change the
// configuration version code. Readers should be prepared to ignore unrecognised
// data beyond the definition of the data they understand.
//
// The values for general_profile_space, general_tier_flag, general_profile_idc,
// general_profile_compatibility_flags, general_constraint_indicator_flags,
// general_level_idc, min_spatial_segmentation_idc, chroma_format_idc,
// bit_depth_luma_minus8 and bit_depth_chroma_minus8 shall be valid for all
// parameter sets that are activated when the stream described by this record is
// decoded (referred to as “all the parameter sets” in the following sentences
// in this paragraph). Specifically, the following restrictions apply.
//
//   — The value of general_profile_space in all the parameter sets shall be
//     identical.
//   — The tier indication general_tier_flag shall indicate a tier equal to or
//     greater than the highest tier indicated in all the parameter sets.
//   - The profile indication general_profile_idc shall indicate a profile to
//     which the stream associated with this configuration record conforms.
//
// If the sequence parameter sets are marked with different profiles, then the
// stream may need examination to determine which profile, if any, the entire
// stream conforms to. If the entire stream is not examined, or the examination
// reveals that there is no profile to which the entire stream conforms, then
// the entire stream shall be split into two or more sub-streams with separate
// configuration records in which these rules can be met.
//
//   - Each bit in general_profile_compatibility_flags may only be set if all
//     the parameter sets set that bit.
//   - Each bit in general_constraint_indicator_flags may only be set if all the
//     parameter sets set that bit.
//   - The level indication general_level_idc shall indicate a level of
//     capability equal to or greater than the highest level indicated for the
//     highest tier in all the parameter sets.
//   - The min_spatial_segmentation_idc indication shall indicate a level of
//     spatial segmentation equal to or less than the lowest level of spatial
//     segmentation indicated in all the parameter sets.
//   - The value of chroma_format_idc in all the parameter sets shall be
//     identical.
//   - The value of bit_depth_luma_minus8 in all the parameter sets shall be
//     identical.
//   - The value of bit_depth_chroma_minus8 in all the parameter sets shall be
//     identical.
//
// Explicit indication can be provided in the HEVC Decoder Configuration Record
// about the chroma format and bit depth as well as other important format
// information used by the HEVC video elementary stream. Each type of such
// information shall be identical in all parameter sets, if present, in a single
// HEVC configuration record. If two sequences differ in any type of such
// information, two different HEVC sample entries shall be used. If the two
// sequences differ in color space indications in their VUI information, then
// two different HEVC sample entries are also required.
//
// There is a set of arrays to carry initialization NAL units. The NAL unit
// types are restricted to indicate SPS, PPS, VPS, prefix SEI, and suffix SEI
// NAL units only. NAL unit types that are reserved in ISO/IEC 23008-2 and in
// this specification may acquire a definition in future, and readers should
// ignore arrays with reserved or unpermitted values of NAL unit type.
//
// > NOTE This “tolerant” behaviour is designed so that errors are not raised,
// allowing the possibility of backwards-compatible extensions to these arrays
// in future specifications.
//
// It is recommended that the arrays be in the order VPS, SPS, PPS, prefix SEI,
// suffix SEI.
//
// When general_non_packed_constraint_flag (bit 3 of the 6-byte
// general_constraint_indicator_flags) is equal to 0 and some of the samples
// referring to this sample entry represent frame-packed content and any of the
// default display windows specified by the active SPSs for the samples
// referring to this sample entry covers more than one constituent frame of the
// frame-packed content, the techniques described in ISO/IEC 14496-12:2015, 8.15
// (“Post-decoder requirements on media”) using the scheme type “stvi” shall be
// used. In this case, the stereo_scheme in the Stereo Video Box should be set
// to 1, to indicate that the frame packing scheme used in HEVC is the same as
// in AVC.
type HEVCDecoderConfigurationRecord struct {
	ConfigurationVersion             uint8
	GeneralProfileSpace              uint8
	GeneralTierFlag                  bool
	GenertalProfileIndicator         uint8
	GeneralProfileCompatibilityFlags uint32
	GeneralConstraintIndicatorFlags  uint64
	GeneralLevelIndicator            uint8
	MinSpatialSegmentationIndicator  uint16
	ParallelismType                  uint8
	ChromaFormatIndicator            uint8
	BitDepthLumaMinus8               uint8
	BitDepthChromaMinus8             uint8
	AvgFrameRate                     uint16
	ConstantFrameRate                uint8
	NumTemporalLayers                uint8
	TemporalIDNested                 uint8
	LengthSizeMinusOne               uint8
	NaluArrays                       []NaluArray
}

type NaluArray struct {
	ArrayCompleteness bool
	NALUnitType       NaluType
	NALUs             [][]byte
}

func (b *HEVCDecoderConfigurationRecord) RecordSize() (size uint32) {
	// unsigned int(8) configurationVersion = 1;
	// unsigned int(2) general_profile_space;
	// unsigned int(1) general_tier_flag;
	// unsigned int(5) general_profile_idc;
	// unsigned int(32) general_profile_compatibility_flags;
	// unsigned int(48) general_constraint_indicator_flags;
	// unsigned int(8) general_level_idc;
	// bit(4) reserved = '1111'b;
	// unsigned int(12) min_spatial_segmentation_idc;
	// bit(6) reserved = '111111'b;
	// unsigned int(2) parallelismType;
	// bit(6) reserved = '111111'b;
	// unsigned int(2) chroma_format_idc;
	// bit(5) reserved = '11111'b;
	// unsigned int(3) bit_depth_luma_minus8;
	// bit(5) reserved = '11111'b;
	// unsigned int(3) bit_depth_chroma_minus8;
	// unsigned int(16) avgFrameRate;
	// unsigned int(2) constantFrameRate;
	// unsigned int(3) numTemporalLayers;
	// unsigned int(1) temporalIdNested;
	// unsigned int(2) lengthSizeMinusOne;
	// unsigned int(8) numOfArrays;
	size += 23
	// unsigned int(1) array_completeness;
	// bit(1) reserved = 0;
	// unsigned int(6) NAL_unit_type;
	// unsigned int(16) numNalus;
	size += 3 * uint32(len(b.NaluArrays))
	var naluCount uint32
	for _, entry := range b.NaluArrays {
		naluCount += uint32(len(entry.NALUs))
		for _, nalu := range entry.NALUs {
			size += uint32(len(nalu)) // bit(8*nalUnitLength) nalUnit;
		}
	}
	size += 2 * naluCount // unsigned int(16) nalUnitLength;
	return
}

func (b *HEVCDecoderConfigurationRecord) RecordRead(r io.Reader) (err error) {
	var tmp [23]uint8
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	b.ConfigurationVersion = tmp[0]
	b.GeneralProfileSpace = tmp[1] >> 6
	b.GeneralTierFlag = ((tmp[1] >> 5) & 0b1) > 0
	b.GenertalProfileIndicator = tmp[1] & 0b11111
	b.GeneralProfileCompatibilityFlags = uint32(tmp[2])<<24 | uint32(tmp[3])<<16 | uint32(tmp[4])<<8 | uint32(tmp[5])
	b.GeneralConstraintIndicatorFlags = uint64(tmp[6])<<40 | uint64(tmp[7])<<32 | uint64(tmp[8])<<24 | uint64(tmp[9])<<16 | uint64(tmp[10])<<8 | uint64(tmp[11])
	b.GeneralLevelIndicator = tmp[12]
	b.MinSpatialSegmentationIndicator = uint16(tmp[13]&0b1111)<<8 | uint16(tmp[14])
	b.ParallelismType = tmp[15] & 0b11
	b.ChromaFormatIndicator = tmp[16] & 0b11
	b.BitDepthLumaMinus8 = tmp[17] & 0b111
	b.BitDepthChromaMinus8 = tmp[18] & 0b111
	b.AvgFrameRate = uint16(tmp[19])<<8 | uint16(tmp[20])
	b.ConstantFrameRate = tmp[21] >> 6
	b.NumTemporalLayers = (tmp[21] >> 3) & 0b111
	b.TemporalIDNested = (tmp[21] >> 2) & 0b1
	b.LengthSizeMinusOne = tmp[21] & 0b11
	entryCount := tmp[22]
	b.NaluArrays = make([]NaluArray, entryCount)
	for i := uint8(0); i < entryCount; i++ {
		if err = binary.Read(r, binary.BigEndian, tmp[:3]); err != nil {
			return
		}
		b.NaluArrays[i].ArrayCompleteness = (tmp[0] >> 7) > 0
		b.NaluArrays[i].NALUnitType = NaluType(tmp[0] & 0b111111)
		naluCount := uint16(tmp[1]&0b1111)<<8 | uint16(tmp[2])
		b.NaluArrays[i].NALUs = make([][]byte, naluCount)
		for j := uint16(0); j < naluCount; j++ {
			var naluLength uint16
			if err = binary.Read(r, binary.BigEndian, &naluLength); err != nil {
				return
			}
			b.NaluArrays[i].NALUs[j] = make([]byte, naluLength)
			if _, err = io.ReadFull(r, b.NaluArrays[i].NALUs[j]); err != nil {
				return
			}
		}
	}
	return
}

func (b *HEVCDecoderConfigurationRecord) RecordWrite(w io.Writer) (err error) {
	var tmp uint8
	if err = binary.Write(w, binary.BigEndian, b.ConfigurationVersion); err != nil {
		return
	}
	tmp = (b.GeneralProfileSpace << 6) | (b.GenertalProfileIndicator & 0b11111)
	if b.GeneralTierFlag {
		tmp |= 0b10000
	}
	if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.GeneralProfileCompatibilityFlags); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint16(b.GeneralConstraintIndicatorFlags>>32)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(b.GeneralConstraintIndicatorFlags)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.GeneralLevelIndicator); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.MinSpatialSegmentationIndicator|(0b1111<<12)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.ParallelismType|0b111111); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.ChromaFormatIndicator|0b111111); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.BitDepthLumaMinus8|0b11111); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.BitDepthChromaMinus8|0b11111); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.AvgFrameRate); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, (b.ConstantFrameRate<<6)|(b.NumTemporalLayers&0b111)<<3|(b.TemporalIDNested&0b1)<<2|(b.LengthSizeMinusOne&0b11)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint8(len(b.NaluArrays))); err != nil {
		return
	}
	for _, entry := range b.NaluArrays {
		var tmp uint8
		tmp |= uint8(entry.NALUnitType) & 0b00111111
		if entry.ArrayCompleteness {
			tmp |= 0b10000000
		}
		if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint16(len(entry.NALUs))); err != nil {
			return
		}
		for _, nalu := range entry.NALUs {
			if err = binary.Write(w, binary.BigEndian, uint16(len(nalu))); err != nil {
				return
			}
			if err = binary.Write(w, binary.BigEndian, nalu); err != nil {
				return
			}
		}
	}
	return
}

// CreateHEVCDecoderConfigurationRecord - extract information from vps, sps, pps and fill HEVCDecoderConfigurationRecord with that
func CreateHEVCDecoderConfigurationRecord(vpsNalus, spsNalus, ppsNalus [][]byte, vpsComplete, spsComplete, ppsComplete bool) (HEVCDecoderConfigurationRecord, error) {
	sps, err := ParseSPSNALUnit(spsNalus[0])
	if err != nil {
		return HEVCDecoderConfigurationRecord{}, err
	}
	var naluArrays []NaluArray
	naluArrays = append(naluArrays, NaluArray{vpsComplete, NALU_VPS, vpsNalus})
	naluArrays = append(naluArrays, NaluArray{spsComplete, NALU_SPS, spsNalus})
	naluArrays = append(naluArrays, NaluArray{ppsComplete, NALU_PPS, ppsNalus})
	ptf := sps.ProfileTierLevel
	return HEVCDecoderConfigurationRecord{
		ConfigurationVersion:             1,
		GeneralProfileSpace:              ptf.GeneralProfileSpace,
		GeneralTierFlag:                  ptf.GeneralTierFlag,
		GenertalProfileIndicator:         ptf.GeneralProfileIndicator,
		GeneralProfileCompatibilityFlags: ptf.GeneralProfileCompatibilityFlags,
		GeneralConstraintIndicatorFlags:  ptf.GeneralConstraintIndicatorFlags,
		GeneralLevelIndicator:            ptf.GeneralLevelIndicator,
		MinSpatialSegmentationIndicator:  0, // Set as default value
		ParallelismType:                  0, // Set as default value
		ChromaFormatIndicator:            sps.ChromaFormatIndicator,
		BitDepthLumaMinus8:               sps.BitDepthLumaMinus8,
		BitDepthChromaMinus8:             sps.BitDepthChromaMinus8,
		AvgFrameRate:                     0,          // Set as default value
		ConstantFrameRate:                0,          // Set as default value
		NumTemporalLayers:                0,          // Set as default value
		TemporalIDNested:                 0,          // Set as default value
		LengthSizeMinusOne:               3,          // only support 4-byte length
		NaluArrays:                       naluArrays, // VPS, SPS, PPS nalus with complete flag
	}, nil
}
