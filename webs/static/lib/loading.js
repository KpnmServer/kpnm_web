
'use strict';

(function(){
	const mount = document.createElement('div');
	document.getElementsByTagName('body')[0].appendChild(mount);
	mount.style.display = 'none';
	mount.style.width = '100%';
	mount.style.height = '100%';
	mount.style.backgroundColor = '#0009';
	mount.style.position = 'fixed';
	mount.style.top = '0';
	mount.style.left = '0';
	const loadingbox = document.createElement('div');
	mount.appendChild(loadingbox);
	loadingbox.style.width = '16rem';
	loadingbox.style.position = 'absolute';
	loadingbox.style.left = loadingbox.style.top = 'calc(calc(100% - ' + loadingbox.style.width + ') / 2)';
	const loadingcv = document.createElement('canvas');
	loadingbox.appendChild(loadingcv);
	const DIAMETER = 1000;
	loadingcv.width = DIAMETER;
	loadingcv.height = DIAMETER;
	loadingcv.style.height = loadingcv.style.width = loadingbox.style.width;
	const loadingtext = document.createElement('div');
	loadingbox.appendChild(loadingtext);
	loadingtext.style.width = '100%';
	loadingtext.style.marginTop = '1rem';
	loadingtext.style.textAlign = 'center';
	loadingtext.style.color = '#fff';
	loadingtext.style.fontSize = '1.3rem';

	var loadingtexts = [];
	const cvctx = loadingcv.getContext('2d');
	const [centerx, centery] = [DIAMETER / 2, DIAMETER / 2];
	var drawintervalid = undefined;
	var timeouts = [];
	function whenTimeout(n){
		for(let i in timeouts){
			if(timeouts[i] === n){
				timeouts.splice(i, 1);
				break;
			}
		}
	}
	function setDTimeout(call, time){
		var n = window.setTimeout(()=>{
			whenTimeout(n);
			call();
		}, time);
		timeouts.push(n);
		return n;
	}
	function stopDrawer(){
		if(drawintervalid !== undefined){
			clearInterval(drawintervalid);
			drawintervalid = undefined;
		}
	}
	function stopAllTimeouts(){
		let ts = timeouts;
		timeouts = [];
		ts.forEach((i)=>{
			clearTimeout(i);
		});
	}
	function stopAll(){
		stopAllTimeouts();
		stopDrawer();
	}
	var is_closing = 0;
	var linelength = 0;
	var degrees1 = 0;
	function drawCycle(){
		cvctx.strokeStyle = '#fff';
		cvctx.lineWidth = DIAMETER / 40;
		let r = (DIAMETER - cvctx.lineWidth) / 2;
		cvctx.beginPath();
		cvctx.moveTo(centerx + r, centery);
		cvctx.arc(centerx, centery, r, 0, Math.PI * 2, false);
		cvctx.closePath();
		cvctx.stroke();
		cvctx.lineWidth = DIAMETER / 200;
		cvctx.beginPath();
		cvctx.moveTo(centerx - 60, centery);
		cvctx.lineTo(centerx + 60, centery);
		cvctx.moveTo(centerx, centery - 60);
		cvctx.lineTo(centerx, centery + 60);
		cvctx.closePath();
		cvctx.stroke();
	}
	function drawDiameter(d, a=0){
		let r = d / 2;
		if(r < 100){
			return;
		}
		cvctx.save();
		cvctx.strokeStyle = '#fff';
		cvctx.lineWidth = DIAMETER / 40;
		cvctx.translate(centerx, centery);
		cvctx.rotate((a % 360) * Math.PI / 180);
		cvctx.beginPath();
		cvctx.moveTo(-r, 0);
		cvctx.lineTo(-100, 0);
		cvctx.moveTo(r, 0);
		cvctx.lineTo(100, 0);
		cvctx.closePath();
		cvctx.stroke();
		cvctx.restore();
	}
	const D_FRAME_LEN = 30;
	const D_FRAME = DIAMETER / D_FRAME_LEN;
	const DIAMETER_S = Math.sqrt(DIAMETER);
	var good_len = [];
	for(let i = 0;i < DIAMETER;i += D_FRAME){
		good_len.push(Math.sqrt(i) * DIAMETER_S);
	}
	good_len.push(DIAMETER);
	function beginAni(){
		mount.style.display = 'block';
		is_closing = 0;
		if(drawintervalid !== undefined){
			return;
		}
		stopAll();
		let i = 0;
		for(let j = Math.round(Math.pow(linelength / DIAMETER_S, 2) / D_FRAME); j < D_FRAME_LEN ;j++){
			let d = good_len[j];
			setDTimeout(()=>{
				cvctx.clearRect(0, 0, loadingcv.width, loadingcv.height);
				drawDiameter(d);
				linelength = d;
				drawCycle();
			}, (i++) * 30);
		}
		setDTimeout(()=>{
			degrees1 = 0;
			linelength = DIAMETER;
			mainAni();
		}, i * 30);
	}
	function afterAni(){
		is_closing = 1;
		if(drawintervalid === undefined){
			stopAllTimeouts();
			afterAni0();
		}
	}
	function afterAni0(){
		is_closing = 0;
		let i = 0;
		for(let j = Math.round(Math.pow(linelength / DIAMETER_S, 2) / D_FRAME); j >= 0 ;j--){
			let d = good_len[j];
			setDTimeout(()=>{
				cvctx.clearRect(0, 0, loadingcv.width, loadingcv.height);
				drawDiameter(d);
				linelength = d;
				drawCycle();
			}, (i++) * 30);
		}
		setDTimeout(()=>{
			mount.style.display = 'none';
			linelength = 0;
		}, i * 30);
	}
	function mainAni(){
		if(is_closing === 1){
			afterAni0();
			return;
		}
		drawintervalid = setInterval(()=>{
			if(is_closing >= 0){
				if(is_closing === 1 && (90 - Math.abs(90 - degrees1 % 180) <= 4)){
					degrees1 = 0;
					stopDrawer();
					afterAni0();
					return;
				}
				degrees1 = (degrees1 + 4) % 360;
			}
			cvctx.clearRect(0, 0, loadingcv.width, loadingcv.height);
			drawDiameter(linelength, degrees1);
			drawCycle();
		}, 30);
	}
	function startAnimation(){
		beginAni();
	}
	function closeAnimation(){
		afterAni();
	}
	var loading_count = 0;
	function beginLoad(){
		if((loading_count++) == 0){
			startAnimation();
			loadingtext.innerText = 'Loading...';
		}
	}
	function endLoad(){
		if((--loading_count) == 0){
			closeAnimation();
		}
	}
	var loading = async function(data){
		beginLoad();
		var re = undefined;
		try{
			if(data.reslove){
				re = await new Promise((rs, rj)=>{data.call(rs, rj)});
			}else if(!data.sync){
				re = await data.call();
			}else{
				re = data.call();
			}
			if(typeof data.success === "function"){
				data.success(re);
			}
		}finally{
			endLoad();
			if(typeof data.complete === "function"){
				data.complete();
			}
		}
	}
	loading.isloading = function(){
		return drawintervalid !== undefined || timeouts.length;
	};
	window.loading = loading;
})();
