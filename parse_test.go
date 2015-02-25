// parse_test
package main

import (
	"testing"
)

func TestParseDetail(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e)
		}
	}()
	crash := new(Crash)
	parseDetail(crash,
		`{mem_info=HeapMax:48,DvmTotal:18503,DvmFree:3915,Pss:75081,Private:61732,Shared:6232}{detail=java.lang.ArrayIndexOutOfBoundsException: length=12; index=-2<br>at java.util.ArrayList.get(ArrayList.java:306)<br>at android.support.v4.app.FragmentManagerImpl.getFragment(FragmentManager.java:579)<br>at com.baidu.mapframework.app.fpstack.BaseTask.onActivityResult(BaseTask.java:398)<br>at android.app.Activity.dispatchActivityResult(Activity.java:5385)<br>at android.app.ActivityThread.deliverResults(ActivityThread.java:3145)<br>at android.app.ActivityThread.handleSendResult(ActivityThread.java:3192)<br>at android.app.ActivityThread.access$1100(ActivityThread.java:133)<br>at android.app.ActivityThread$H.handleMessage(ActivityThread.java:1251)<br>at android.os.Handler.dispatchMessage(Handler.java:99)<br>at android.os.Looper.loop(Looper.java:137)<br>at android.app.ActivityThread.main(ActivityThread.java:4813)<br>at java.lang.reflect.Method.invokeNative(Native Method)<br>at java.lang.reflect.Method.invoke(Method.java:511)<br>at com.android.internal.os.ZygoteInit$MethodAndArgsCaller.run(ZygoteInit.java:792)<br>at com.android.internal.os.ZygoteInit.main(ZygoteInit.java:559)<br>at dalvik.system.NativeStart.main(Native Method)}{reason=java.lang.RuntimeException: Failure delivering result ResultInfo{who=null, request=-1, result=-1, data=Intent { act=android.intent.action.DIAL dat=tel:xxxxxxxxxxx }} to activity {com.baidu.BaiduMap/com.baidu.baidumaps.MapsActivity}: java.lang.ArrayIndexOutOfBoundsException: length=12; index=-2}{active_thread=36}{net=9}{locx=12953557}{locy=4836374}{cpu_abi2=armeabi}{coms_info=map.android.baidu.ishare_1.1.9-map.android.baidu.weekend_1.0.8-map.android.baidu.hotel_1.6.6-map.android.baidu.scenery_3.2.3-map.android.baidu.maplab_1.0.4-map.android.baidu.violation_1.0.2-map.android.baidu.cater_1.3.6-map.android.baidu.aoi_1.1.2-map.android.baidu.cater_1.3.3-map.android.baidu.caterbooking_1.1.2-map.android.baidu.diagnose_1.0.2-map.android.baidu.grouponcomponent_1.1.2-map.android.baidu.hotel_1.6.3-map.android.baidu.ishare_1.1.8-map.android.baidu.maplab_1.0.3-map.android.baidu.movie_1.4.3-map.android.baidu.pano_1.1.3-map.android.baidu.qrcode_1.0.7-map.android.baidu.rentcar_1.3.2-map.android.baidu.scenery_3.2.2-map.android.baidu.takeout_1.6.3-map.android.baidu.taxi_3.1.2-map.android.baidu.violation_1.0.0-map.android.baidu.voice_1.2.2-map.android.baidu.websdk_1.2.13-map.android.baidu.weekend_1.0.6-}{pages=#0:com.component.android.main|com.baidu.baidumaps.MapsActivity@4195b710|com.baidu.baidumaps.base.MapFramePage#1:com.component.android.main|com.baidu.baidumaps.MapsActivity@4195b710|com.baidu.baidumaps.poi.page.PoiSearchPage#2:map.android.baidu.cater_1.3.6_1101248456|com.baidu.baidumaps.MapsActivity@4195b710|map.android.baidu.cater.page.CateringPoiListPage}{cpu_abi=armeabi-v7a}`)
}
