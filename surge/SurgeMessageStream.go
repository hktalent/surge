package surge

import (
	"encoding/binary"
	"io"

	log "github.com/sirupsen/logrus"
)

// SessionWrite writes to session
func SessionWrite(Session *Session, Data []byte, ID byte) (err error) {

	//Package identifier to know what we are sending
	packID := make([]byte, 1)
	packID[0] = ID

	//Create buffer of 4 bytes to put the size of the package
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(len(Data)))

	//append pack and buff
	buff = append(packID, buff...)

	//Write data
	buff = append(buff, Data...)
	_, err = Session.session.Write(buff)
	if err != nil {
		log.Fatal(err)
	}

	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator[Session.FileHash] += len(Data)
	bandwidthAccumulatorMapLock.Unlock()

	return err
}

//SessionRead reads from session
func SessionRead(Session *Session) (data []byte, ID byte, err error) {
	headerBuffer := make([]byte, 5) //int32 size of header + 1 for packid

	// the header of 4 bytes + 1 for packid
	_, err = io.ReadFull(Session.reader, headerBuffer)
	if err != nil {
		if err.Error() == "session closed" {
			log.Println(err)
			return nil, 0x0, err
		}
		log.Println(err)
		return nil, 0x0, err
	}

	//Get the packid
	packID := headerBuffer[0]
	log.Println(packID)

	//Get the size from the bytes
	sizeBytes := append(headerBuffer[:0], headerBuffer[1:]...)

	size := binary.LittleEndian.Uint32(sizeBytes)

	data = make([]byte, size)

	// read the full message, or return an error
	_, err = io.ReadFull(Session.reader, data[:int(size)])
	if err != nil {
		log.Println(err)
		return nil, 0x0, err
	}

	//Write add to download
	bandwidthAccumulatorMapLock.Lock()
	downloadBandwidthAccumulator[Session.FileHash] += int(size)
	bandwidthAccumulatorMapLock.Unlock()

	return data, packID, err
}
