package iio

//cgo: pkg-config
//#include <iio.h>

import(
  "C"
  "fmt"
)

/** @brief Retrieve a pointer to the iio_context structure
 * @param dev A pointer to an iio_device structure
 * @return A pointer to an iio_context structure */
func DeviceGetContext(dev *Device)(*Context,error){
  var err error
  var ctx *Context
  ctx = C.iio_device_get_context(dev.handle)
  if(ctx == nil){
    err = fmt.printf("Error getting context from device\n")
  }
  return ctx,err
}

/** @brief Retrieve the device ID (e.g. <b><i>iio:device0</i></b>)
 * @param dev A pointer to an iio_device structure
 * @return A pointer to a static NULL-terminated string */
func DeviceGetId(dev *Device)(string, error){
  var err error
  var str string
  var strLen int

  res := C.iio_device_get_name(dev.handle)
  if res == nil{
    err = fmt.Errorf("Device id not found") 
    return "",err    
  }
  strLen = C.strlen(res)
  str = unsafe.String(res,strLen)
  return str, nil
}


/** @brief Retrieve the device label (e.g. <b><i>lo_pll0_rx_adf4351</i></b>)
 * @param dev A pointer to an iio_device structure
 * @return A pointer to a static NULL-terminated string
 *
 * <b>NOTE:</b> if the device has no label, NULL is returned. */
func DeviceGetLabel(dev *Device)(string, error){
  var err error
  var str string
  var strLen int

  res := C.iio_device_get_label(dev.handle)
  if res == nil{
    err = fmt.Errorf("Device id not found") 
    return "",err    
  }
  strLen = C.strlen(res)
  str = unsafe.String(res,strLen)
  return str, nil
}


/** @brief Enumerate the channels of the given device
 * @param dev A pointer to an iio_device structure
 * @return The number of channels found */
func DeviceGetChannelsCount(dev *Device)(int){
  return int(C.iio_device_get_channels_count(dev.handle))
}






