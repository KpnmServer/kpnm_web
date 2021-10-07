
'use strict';

const emailreg = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/;
const namereg  = /^[A-Za-z_-][0-9A-Za-z_-]{1,31}$/;
const pwdreg  = /^[A-Za-z][0-9A-Za-z_+\-*/!@#$%^&()~\[\]{}|=,.<>;:'"]{7,127}$/;

function shakeBox(box){
	box.removeClass("animation-shake-time3");
	setTimeout(()=>{box.addClass("animation-shake-time3");}, 50);
}

function bindEnterClick(base, btn){
	btn.keydown(function(event){ if(event.keyCode == 13){event.preventDefault();} });
	base.keyup(function(event){ if(event.keyCode == 13){btn.click();} });
}

var getting_captimg = false;
function updateCaptcha(){
	if(getting_captimg){
		return;
	}
	getting_captimg = true;
	$.ajax({
		url: "/user/api/captcha/image",
		type: "GET",
		success: function(res){
			if(res.status === "ok" ){
				$('#verify-captcha-img').prop("src", res.data);
				return;
			}
			console.error("get chaptcha error:", res);
		},
		complete: function(){
			getting_captimg = false;
		}
	})
}

var _register_data = {
	email: "",
	emailtk: "",
	name: "",
	pwd: ""
}

function submitRegister(){
	const errorbox = $('#verify-error');
	$.ajax({
		url: "/user/api/register",
		type: "POST",
		data: {
			emailtk: _register_data["emailtk"],
			name: _register_data["name"],
			pwd: _register_data["pwd"]
		},
		success: function(res){
			if(res.status === "ok"){
				alert("register success");
				errorbox.hide();
				errorbox.text("");
				window.location.replace("/user/login");
				return;
			}
			if(res.status === "error"){
				errorbox.text(res.errorMessage);
				errorbox.show();
				shakeBox(errorbox);
				updateCaptcha();
				return;
			}
			console.log("error res:", res);
		}
	});
}

function changeToEmailVerify(){
	updateCaptcha();
	$('#verify-email-addr').text(_register_data["email"]);
	$('#register-info').hide();
	$('#verify-email').show();
	let vci = $("#verify-captcha-input");
	vci.val('');
	vci.focus();
}

function changeToUserSetting(){
	_register_data = {
		emailtk: "",
		name: "",
		pwd: ""
	}
	$('#verify-email').hide();
	$('#register-info').show();
}

$(document).ready(function(){
	function checkEmailVal(){
		const email = $('#register-email-input').val();
		const errbox = $('#register-email-error');
		if(!emailreg.test(email)){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}
	function checkUsernameVal(){
		const username = $('#register-username-input').val();
		const errbox = $('#register-username-error');
		if(!namereg.test(username)){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}
	function checkPasswordVal(){
		const password = $('#register-password-input').val();
		const errbox = $('#register-password-error');
		if(!pwdreg.test(password)){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}
	function checkPwdagainVal(){
		if(!checkPasswordVal()){
			return false;
		}
		const password = $('#register-password-input').val();
		const pwdagain = $('#register-pwdagain-input').val();
		const errbox = $('#register-pwdagain-error');
		if(password !== pwdagain){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}

	$('#register-email-input').blur(checkEmailVal);
	$('#register-username-input').blur(checkUsernameVal);
	$('#register-password-input').blur(checkPasswordVal);
	$('#register-pwdagain-input').blur(checkPwdagainVal);
	$('#register-next').click(function(){
		let okmail = checkEmailVal(), okname = checkUsernameVal(), okpwd = checkPasswordVal() && checkPwdagainVal();
		if(!okmail || !okname || !okpwd){
			return
		}
		const errorbox = $('#register-error');
		const email = $('#register-email-input').val();
		const name = $('#register-username-input').val();
		const pwd = $('#register-password-input').val();
		$.ajax({
			type: "POST",
			url: "/user/api/regcheck",
			data: {
				email: email,
				name: name,
				pwd: pwd
			},
			success: function(res){
				if(res.status === "ok"){
					errorbox.hide();
					_register_data["emailtk"] = "";
					_register_data["name"] = name;
					_register_data["pwd"] = pwd;
					_register_data["email"] = email;
					changeToEmailVerify();
					return;
				}
				if(res.status === "error"){
					errorbox.text(res.errorMessage);
					errorbox.show();
					shakeBox(errorbox);
				}
			}
		})
	});
	bindEnterClick($('#register-info'), $('#register-next'));

	var sending_email = 0;
	function sendEmailCode(){
		const errorbox = $('#verify-captcha-error');
		let sendint = sending_email - new Date().getTime();
		if(sendint > 0){
			errorbox.text(Math.floor(sendint / 1000) + "秒后才能再次发送邮件");
			errorbox.show();
			shakeBox(errorbox);
			return
		}
		sending_email = new Date().getTime() + 60 * 1000;
		const captcode = $("#verify-captcha-input").val();
		$.ajax({
			type: "POST",
			url: "/user/api/verify/email/send",
			data: {
				email: _register_data["email"],
				capt: captcode
			},
			success: function(res){
				if(res.status === "ok"){
					errorbox.hide();
					errorbox.text("");
					alert("我们已将验证码发送至 '" + _register_data["email"] + "', 请查收");
					let veci = $("#verify-email-captcha-input");
					veci.val('');
					veci.focus();
					return;
				}
				if(res.status === "error"){
					$("#verify-captcha-input").val('');
					$("#verify-captcha-input").focus();
					if(res.error == "CaptcodeError"){
						errorbox.text("验证码错误");
						errorbox.show();
						shakeBox(errorbox);
						sending_email = 0;
					}else{
						errorbox.text(res.errorMessage);
						errorbox.show();
						shakeBox(errorbox);
					}
					updateCaptcha();
					return;
				}
				console.log("error res:", res);
			}
		});
	}
	function verifyEmail(){
		const errorbox = $('#verify-email-captcha-error');
		const emailcode = $("#verify-email-captcha-input").val();
		const captcode = $("#verify-captcha-input").val();
		$.ajax({
			type: "POST",
			url: "/user/api/verify/email",
			data: {
				code: emailcode,
				capt: captcode
			},
			success: function(res){
				if(res.status === "ok"){
					errorbox.hide();
					errorbox.text("");
					_register_data["emailtk"] = res.token;
					submitRegister();
					return;
				}
				if(res.status === "error"){
					let vci = $("#verify-captcha-input");
					vci.val('');
					vci.focus();
					if(res.error == "CaptcodeError"){
						let errbox = $('#verify-captcha-error');
						errbox.text("验证码错误");
						errbox.show();
						shakeBox(errbox);
						sending_email = 0;
					}else{
						errorbox.text(res.errorMessage);
						errorbox.show();
						shakeBox(errorbox);
					}
					updateCaptcha();
					return;
				}
				console.log("error res:", res);
			}
		});
	}
	$('#back-user-setting').click(changeToUserSetting);
	$('#verify-email-send').click(sendEmailCode);
	$('#register-submit').click(verifyEmail);
	bindEnterClick($('#verify-captcha'), $('#verify-email-send'));
	bindEnterClick($('#verify-email-captcha'), $('#register-submit'));
	$('#verify-captcha-img').click(updateCaptcha);
});
