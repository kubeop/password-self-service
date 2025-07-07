document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('resetPasswordForm');
    const sendCodeBtn = document.getElementById('sendCodeBtn');
    const usernameInput = document.getElementById('username');
    const codeInput = document.getElementById('code');
    const newPasswordInput = document.getElementById('newPassword');
    const confirmPasswordInput = document.getElementById('confirmPassword');
    const usernameError = document.getElementById('usernameError');
    const codeError = document.getElementById('codeError');
    const passwordError = document.getElementById('passwordError');
    const confirmError = document.getElementById('confirmError');
    const successMessage = document.getElementById('successMessage');

    // 验证码倒计时功能
    let countdown = 0;
    let countdownTimer = null;

    // 表单验证函数
    const validateForm = () => {
        let isValid = true;
        const username = usernameInput.value.trim();
        const code = codeInput.value.trim();
        const newPassword = newPasswordInput.value;
        const confirmPassword = confirmPasswordInput.value;

        // 用户名验证
        if (!username) {
            usernameError.textContent = '请输入用户名';
            isValid = false;
        } else if (username.length < 2 || username.length > 20) {
            usernameError.textContent = '用户名长度必须在2-20个字符之间';
            isValid = false;
        } else {
            usernameError.textContent = '';
        }

        // 验证码验证
        if (!code) {
            codeError.textContent = '请输入验证码';
            isValid = false;
        } else {
            codeError.textContent = '';
        }

        // 密码验证
        if (!newPassword) {
            passwordError.textContent = '请输入新密码';
            isValid = false;
        } else if (newPassword.length < 8) {
            passwordError.textContent = '密码长度不能少于8个字符';
            isValid = false;
        } else {
            passwordError.textContent = '';
        }

        // 确认密码验证
        if (!confirmPassword) {
            confirmError.textContent = '请确认新密码';
            isValid = false;
        } else if (confirmPassword !== newPassword) {
            confirmError.textContent = '两次输入的密码不一致';
            isValid = false;
        } else {
            confirmError.textContent = '';
        }

        return isValid;
    };

    // 发送验证码
    sendCodeBtn.addEventListener('click', async () => {
        const username = usernameInput.value.trim();

        // 简单验证用户名
        if (!username) {
            usernameError.textContent = '请输入用户名';
            return;
        } else if (username.length < 2 || username.length > 20) {
            usernameError.textContent = '用户名长度必须在2-20个字符之间';
            return;
        } else {
            usernameError.textContent = '';
        }

        // 如果正在倒计时，不重复发送
        if (countdown > 0) return;

        try {
            // 调用发送验证码API
            const response = await fetch('/api/v1/captcha/send', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    category: 'password'
                })
            });

            const data = await response.json();

            if (data.code === 0) {
                // 发送成功，开始倒计时
                countdown = 60;
                sendCodeBtn.disabled = true;
                sendCodeBtn.textContent = `重新发送(${countdown}s)`;

                countdownTimer = setInterval(() => {
                    countdown--;
                    sendCodeBtn.textContent = `重新发送(${countdown}s)`;

                    if (countdown <= 0) {
                        clearInterval(countdownTimer);
                        sendCodeBtn.disabled = false;
                        sendCodeBtn.textContent = '获取验证码';
                    }
                }, 1000);

                codeError.textContent = '';
                successMessage.textContent = '验证码已发送，请查收';
            } else {
                codeError.textContent = data.message || '获取验证码失败，请重试';
            }
        } catch (error) {
            console.error('发送验证码失败:', error);
            codeError.textContent = '网络错误，请稍后重试';
        }
    });

    // 表单提交
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        try {
            const response = await fetch('/api/v1/reset-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: usernameInput.value.trim(),
                    code: codeInput.value.trim(),
                    newPassword: newPasswordInput.value
                })
            });

            const data = await response.json();

            if (data.code === 0) {
                // 重置成功
                successMessage.textContent = '密码重置成功，请使用新密码登录';
                form.reset();
                // 3秒后跳转到首页
                setTimeout(() => {
                    window.location.href = 'index.html';
                }, 3000);
            } else {
                // 显示错误信息
                if (data.message.includes('验证码')) {
                    codeError.textContent = data.message;
                } else if (data.message.includes('用户名')) {
                    usernameError.textContent = data.message;
                } else {
                    passwordError.textContent = data.message;
                }
            }
        } catch (error) {
            console.error('重置密码失败:', error);
            passwordError.textContent = '网络错误，请稍后重试';
        }
    });
});