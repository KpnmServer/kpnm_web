
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
				$('#login-captcha-img').prop("src", res.data);
				return;
			}
			console.error("get chaptcha error:", res);
		},
		complete: function(){
			getting_captimg = false;
		}
	})
}

$(document).ready(function(){
	function checkUsernameVal(){
		const username = $('#login-username-input').val();
		const errbox = $('#login-username-error');
		if(!emailreg.test(username) && !namereg.test(username)){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}
	function checkPasswordVal(){
		const password = $('#login-password-input').val();
		const errbox = $('#login-password-error');
		if(!pwdreg.test(password)){
			errbox.show();
			shakeBox(errbox);
			return false;
		}
		errbox.hide();
		return true;
	}

	$('#login-username-input').blur(checkUsernameVal);
	$('#login-password-input').blur(checkPasswordVal);
	$('#login-submit').click(function(){
		let okname = checkUsernameVal(), okpwd = checkPasswordVal();
		if(!okname || !okpwd){
			return
		}
		const errorbox = $('#login-error');
		const user = $('#login-username-input').val();
		const pwd = $('#login-password-input').val();
		const captcode = $('#login-captcha-input').val();
		const capterrbox = $("#login-captcha-error");
		$.ajax({
			type: "POST",
			url: "/user/api/login",
			data: {
				user: user,
				pwd: pwd,
				capt: captcode
			},
			success: function(res){
				capterrbox.hide();
				errorbox.hide();
				if(res.status === "ok"){
					window.location.replace('/user');
					return;
				}
				if(res.status === "error"){
					let lci = $("#login-captcha-input");
					lci.val('');
					lci.focus();
					updateCaptcha();
					switch(res.error){
					case "CaptcodeError":{
						capterrbox.text("验证码错误");
						capterrbox.show();
						shakeBox(capterrbox);
						break;
					}
					case "PasswordError":{
						let lpi = $("#login-password-input");
						lpi.val('');
						lpi.focus();
						let errbox = $("#login-password-error");
						errbox.text("密码错误");
						errbox.show();
						shakeBox(errbox);
						break;
					}
					default:
						errorbox.text(res.errorMessage);
						errorbox.show();
						shakeBox(errorbox);
					}
				}
			}
		})
	});
	bindEnterClick($('#login'), $('#login-submit'));
	$('#login-captcha-img').click(updateCaptcha);
	updateCaptcha();
});
