
'use strict';

(async function(){
	await loadJsAsync('/static/lib/loading.js');

	function formatSize(size){
		if(size < 1000){
			return size + " Byte";
		}
		size /= 1024;
		if(size < 1000){
			return size.toFixed(2) + " KB";
		}
		size /= 1024;
		if(size < 1000){
			return size.toFixed(2) + " MB";
		}
		size /= 1024;
		if(size < 1000){
			return size.toFixed(2) + " GB";
		}
		size /= 1024;
		return size.toFixed(2) + " TB";
	}

	function formatCpuTime(time){
		time = Math.floor(time * 1000);
		const millisecond = Math.floor(time % 1000); time /= 1000;
		const second = Math.floor(time % 60); time /= 60;
		const minute = Math.floor(time % 60); time /= 60;
		const hour = Math.floor(time % 24); time /= 24;
		const day = Math.floor(time);
		var string = second + '.' + millisecond;
		if(minute){
			string = minute + ':' + string;
		}
		if(hour){
			string = hour + ':' + string;
		}
		if(day){
			string = day + ':' + string;
		}
		return string;
	}

	var svrid = "";
	var flushtimeout = 1000;

	var SVRSTATUS = {
		status: 'UNKNOWN',
		ticks: 0,
		java_version: '',
		os: '',
		max_mem: 0,
		total_mem: 0,
		used_mem: 0,
		cpu_num: 0,
		cpu_load: 0,
		cpu_time: 0,
	};


	function drawCycle(ctx, x, y, border, background="#fff"){
		ctx.fillStyle = border;
		ctx.beginPath();
		ctx.moveTo(x + 10, y);
		ctx.arc(x, y, 10, 0, Math.PI * 2, false);
		ctx.fill();
		ctx.fillStyle = background;
		ctx.beginPath();
		ctx.moveTo(x + 5, y);
		ctx.arc(x, y, 5, 0, Math.PI * 2, false);
		ctx.fill();
	}

	const drawMemoryStatus = (function(){
		const DATA_MAX = 14;
		var datas = [];
		for(let i = 0;i < DATA_MAX;i++){
			datas.push([0, 0]);
		}

		function drawBorder(ctx){
			ctx.font = "26px serif";
			ctx.textBaseline = "middle";
			ctx.textAlign = "right";
			ctx.fillStyle = "#000";
			ctx.lineWidth = 2;
			ctx.strokeStyle = "#999";
			ctx.beginPath();
			for(let i = 10; i >= 0 ;i -= 1){
				ctx.fillText((100 - i * 10) + "%", 60, 50 + (i * 50));
				ctx.moveTo(70, 50 + (i * 50));
				ctx.lineTo(ctx.canvas.width - 20, 50 + (i * 50));
			}
			ctx.stroke();

			ctx.strokeStyle = "#000";
			ctx.lineWidth = 4;
			ctx.beginPath();
			ctx.moveTo(70, 30);
			ctx.lineTo(70, ctx.canvas.height - 50);
			ctx.lineTo(ctx.canvas.width - 20, ctx.canvas.height - 50);
			ctx.stroke();
		}

		function getMemPos(ctx, dt){
			return ctx.canvas.height - (ctx.canvas.height - 100) * dt - 50;
		}
		function drawDatas(ctx){
			ctx.fillStyle = "#0f07";
			ctx.strokeStyle = "#0f0";
			ctx.lineWidth = 8;
			ctx.beginPath();
			let [x, y] = [0, 0];
			[x, y] = [70, getMemPos(ctx, datas[0][0])];
			ctx.moveTo(x, y);
			for(let i = 1; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getMemPos(ctx, datas[i][0])];
				ctx.lineTo(x, y);
			}
			ctx.stroke();
			ctx.lineTo(x, ctx.canvas.height - 50);
			ctx.lineTo(70, ctx.canvas.height - 50);
			ctx.fill();
			for(let i = 0; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getMemPos(ctx, datas[i][0])];
				drawCycle(ctx, x, y, ctx.strokeStyle);
			}

			ctx.fillStyle = "#03f7";
			ctx.strokeStyle = "#03f";
			ctx.lineWidth = 8;
			ctx.beginPath();
			[x, y] = [70, getMemPos(ctx, datas[0][1])];
			ctx.moveTo(x, y);
			for(let i = 1; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getMemPos(ctx, datas[i][1])];
				ctx.lineTo(x, y);
			}
			ctx.stroke();
			ctx.lineTo(x, ctx.canvas.height - 50);
			ctx.lineTo(70, ctx.canvas.height - 50);
			ctx.fill();
			for(let i = 0; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getMemPos(ctx, datas[i][1])];
				drawCycle(ctx, x, y, ctx.strokeStyle);
			}
		}

		return function(ctx){
			if(SVRSTATUS.max_mem > 0){
				datas.push([SVRSTATUS.total_mem / SVRSTATUS.max_mem, SVRSTATUS.used_mem / SVRSTATUS.max_mem]);
			}
			if(datas.length > DATA_MAX){
				datas = datas.slice(datas.length - DATA_MAX);
			}
			ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
			drawBorder(ctx);
			drawDatas(ctx);
		}
	})();

	const drawCpuStatus = (function(){
		const DATA_MAX = 14;
		var datas = [];
		for(let i = 0;i < DATA_MAX;i++){
			datas.push(0);
		}

		function drawBorder(ctx){
			ctx.font = "26px serif";
			ctx.textBaseline = "middle";
			ctx.textAlign = "right";
			ctx.fillStyle = "#000";
			ctx.lineWidth = 2;
			ctx.strokeStyle = "#999";
			ctx.beginPath();
			for(let i = 10; i >= 0 ;i -= 1){
				ctx.fillText(((100 - i * 10) * SVRSTATUS.cpu_num) + "%", 60, 50 + (i * 50));
				ctx.moveTo(70, 50 + (i * 50));
				ctx.lineTo(ctx.canvas.width - 20, 50 + (i * 50));
			}
			ctx.stroke();

			ctx.strokeStyle = "#000";
			ctx.lineWidth = 4;
			ctx.beginPath();
			ctx.moveTo(70, 30);
			ctx.lineTo(70, ctx.canvas.height - 50);
			ctx.lineTo(ctx.canvas.width - 20, ctx.canvas.height - 50);
			ctx.stroke();
		}

		function getCpuPos(ctx, dt){
			if(SVRSTATUS.cpu_num === 0){
				return ctx.canvas.height - 50;
			}
			return ctx.canvas.height - (ctx.canvas.height - 100) * (dt / SVRSTATUS.cpu_num) - 50;
		}
		function drawDatas(ctx){
			ctx.fillStyle = "#f307";
			ctx.strokeStyle = "#f30";
			ctx.lineWidth = 8;
			ctx.beginPath();
			let [x, y] = [70, getCpuPos(ctx, datas[0])];
			ctx.moveTo(x, y);
			for(let i = 1; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getCpuPos(ctx, datas[i])];
				ctx.lineTo(x, y);
			}
			ctx.stroke();
			ctx.lineTo(x, ctx.canvas.height - 50);
			ctx.lineTo(70, ctx.canvas.height - 50);
			ctx.fill();
			for(let i = 0; i < datas.length ;i++){
				[x, y] = [70 + i * 70, getCpuPos(ctx, datas[i])];
				drawCycle(ctx, x, y, ctx.strokeStyle);
			}
		}

		return function(ctx){
			datas.push(SVRSTATUS.cpu_load);
			if(datas.length > DATA_MAX){
				datas = datas.slice(datas.length - DATA_MAX);
			}
			ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
			drawBorder(ctx);
			drawDatas(ctx);
		}
	})();

	function setup($){
		const memory_canvas_ctx = document.getElementById('images-memory').getContext('2d');
		const cpu_canvas_ctx = document.getElementById('images-cpu').getContext('2d');
		drawMemoryStatus(memory_canvas_ctx);
		drawCpuStatus(cpu_canvas_ctx);

		function onStatusUpdate(updateinfo=false, updatestatus=true){
			if(SVRSTATUS.status === "UNKNOWN"){
				$("#status-ticks").hide();
				$("#status-total_mem").hide();
				$("#status-used_mem").hide();
				$("#status-cpu_load").hide();
				$("#status-cpu_time").hide();
				$("#status-cpu_num").hide();
				$("#status-java_version").hide();
				$("#status-os").hide();
				$("#status-max_mem").hide();
				return;
			}
			if(updateinfo){
				$("#status-cpu_num-v").text(SVRSTATUS.cpu_num);
				$("#status-java_version-v").text(SVRSTATUS.java_version);
				$("#status-os-v").text(SVRSTATUS.os);
				$("#status-max_mem-v").text(formatSize(SVRSTATUS.max_mem));
				$("#status-cpu_num").show();
				$("#status-java_version").show();
				$("#status-os").show();
				$("#status-max_mem").show();
			}
			if(updatestatus){
				$("#status-ticks-v").text(SVRSTATUS.ticks);
				$("#status-total_mem-v").text(formatSize(SVRSTATUS.total_mem));
				$("#status-used_mem-v").text(formatSize(SVRSTATUS.used_mem));
				$("#status-cpu_load-v").text((SVRSTATUS.cpu_load * 100.0).toFixed(2));
				$("#status-cpu_time-v").text(formatCpuTime(SVRSTATUS.cpu_time));
				$("#status-status").show();
				$("#status-ticks").show();
				$("#status-total_mem").show();
				$("#status-used_mem").show();
				$("#status-cpu_load").show();
				$("#status-cpu_time").show();
			}
			if((updateinfo || updatestatus) && SVRSTATUS.status !== "STOPPED"){
				drawMemoryStatus(memory_canvas_ctx);
				drawCpuStatus(cpu_canvas_ctx);
			}
		}

		((window.WebSocket)?(function(){
			var socket = null;
			var lastopen = 0;
			function openSocket(){
				if(socket !== null){
					socket.close();
				}
				lastopen = new Date().getTime();
				socket = new window.WebSocket(
					(window.location.protocol === "http:" ?"ws://" :"wss://") + window.location.host + "/zcs/api/svrstatusws?svr=" + SERVER_NAME);
				socket.onopen = function(event){
					socket.send(JSON.stringify({status: "init"}));
				}
				socket.onmessage = function(event){
					const res = JSON.parse(event.data);
					let [updateinfo, updatestatus] = [false, false];
					if(res.status !== undefined){
						if(SVRSTATUS.status === "UNKNOWN" && res.status !== "UNKNOWN"){
							[updateinfo, updatestatus] = [true, true];
						}
						SVRSTATUS.status = res.status;
						$("#status-status-v").text(SVRSTATUS.status);
					}
					if(res.ticks !== undefined){
						SVRSTATUS.ticks = res.ticks;
						SVRSTATUS.total_mem = res.total_mem;
						SVRSTATUS.used_mem = res.used_mem;
						SVRSTATUS.cpu_load = res.cpu_load;
						SVRSTATUS.cpu_time = res.cpu_time;
						updatestatus = true;
					}
					if(res.os !== undefined){
						SVRSTATUS.cpu_num = res.cpu_num;
						SVRSTATUS.java_version = res.java_version;
						SVRSTATUS.os = res.os;
						SVRSTATUS.max_mem = res.max_mem;
						updateinfo = true;
					}
					onStatusUpdate(updateinfo, updatestatus);
				}
				socket.onclose = function(event){
					socket = null;
					setTimeout(openSocket, (new Date().getTime() - lastopen < 10000) ?3000 :500);
				}
			}
			openSocket();
		}):(function(){
			function flushInfo(){
				$.ajax({
					url: "/zcs/api/svrinfo?svr=" + SERVER_NAME,
					type: "GET",
					success: function(res){
						svrid = res.id;
						flushtimeout = Math.max(res.interval * 1000, 500);
						SVRSTATUS.cpu_num = res.cpu_num;
						SVRSTATUS.java_version = res.java_version;
						SVRSTATUS.os = res.os;
						SVRSTATUS.max_mem = res.max_mem;
						onStatusUpdate(true, false);
					}
				});
			}

			var flushing_status = false;
			function flushStatus(){
				if(flushing_status){
					return;
				}
				flushing_status = true;
				$.ajax({
					url: "/zcs/api/svrstatus?svr=" + SERVER_NAME,
					type: "GET",
					success: function(res){
						const rid = res.id;
						if(rid != svrid){
							flushInfo();
						}
						const errstr = res.errstr;
						if(errstr){
							$("#status-errstr-v").text(res.errstr);
							$("#status-errstr").show();
							return;
						}
						SVRSTATUS.status = res.status;
						SVRSTATUS.ticks = res.ticks;
						SVRSTATUS.total_mem = res.total_mem;
						SVRSTATUS.used_mem = res.used_mem;
						SVRSTATUS.cpu_load = res.cpu_load;
						SVRSTATUS.cpu_time = res.cpu_time;
						$("#status-status-v").text(SVRSTATUS.status);
						onStatusUpdate();
					},
					complete: function(){
						flushing_status = false;
					}
				});
			}

			$.ajax({
				url: "/zcs/api/svrinfo?svr=" + SERVER_NAME,
				type: "GET",
				success: function(res){
					svrid = res.id;
					flushtimeout = Math.max(res.interval * 1000, 500);
					SVRSTATUS.cpu_num = res.cpu_num;
					SVRSTATUS.java_version = res.java_version;
					SVRSTATUS.os = res.os;
					SVRSTATUS.max_mem = res.max_mem;
					$("#status-cpu_num-v").text(res.cpu_num);
					$("#status-java_version-v").text(res.java_version);
					$("#status-os-v").text(res.os);
					$("#status-max_mem-v").text(formatSize(res.max_mem));
					$("#status-cpu_num").show();
					$("#status-java_version").show();
					$("#status-os").show();
					$("#status-max_mem").show();
				}
			});
			function timeFlushStatus(){
				flushStatus();
				setTimeout(timeFlushStatus, flushtimeout);
			}
			timeFlushStatus();
		}))();
	}

	loading({
		reslove: true,
		call: async function(rs){
			await loadJsAsync('/static/lib/jquery-3.6.0.min.js');
			$(document).ready(function(){
				setup($);
				rs();
			});
		}
	});
})();

