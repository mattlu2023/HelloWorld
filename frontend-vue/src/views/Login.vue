<template>
  <div class="login-container">
    <div class="login-background"></div>
    <div class="login-content">
      <div class="login-card glass-effect">
        <div class="login-brand">
          <div class="brand-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="9" y1="9" x2="15" y2="9"></line>
              <line x1="12" y1="12" x2="12" y2="15"></line>
            </svg>
          </div>
          <h1 class="brand-title">广告 BI 系统</h1>
          <p class="brand-subtitle">数据驱动的营销决策平台</p>
        </div>
        
        <el-form :model="loginForm" :rules="rules" ref="formRef" class="login-form">
          <el-form-item prop="username">
            <div class="input-wrapper">
              <div class="input-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </div>
              <el-input
                v-model="loginForm.username"
                placeholder="用户名"
                size="large"
                class="custom-input"
                @focus="handleInputFocus('username')"
                @blur="handleInputBlur('username')"
              />
            </div>
          </el-form-item>
          
          <el-form-item prop="password">
            <div class="input-wrapper">
              <div class="input-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </div>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password
                class="custom-input"
                @focus="handleInputFocus('password')"
                @blur="handleInputBlur('password')"
                @keyup.enter="handleLogin"
              />
            </div>
          </el-form-item>
          
          <el-form-item class="form-actions">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleLogin"
              class="login-btn"
            >
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-footer">
          <span>默认账户: admin / admin123</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const focusedField = ref('')

const loginForm = reactive({
  username: 'admin',
  password: 'admin123'
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleInputFocus = (field) => {
  focusedField.value = field
}

const handleInputBlur = () => {
  focusedField.value = ''
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await login(loginForm)
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user_id', res.data.user_id)
        localStorage.setItem('username', res.data.username)
        
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  overflow: hidden;
}

.login-background {
  position: absolute;
  inset: 0;
  background: 
    radial-gradient(ellipse at 20% 30%, rgba(6, 182, 212, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 70%, rgba(15, 23, 42, 0.1) 0%, transparent 50%),
    linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
}

.login-content {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
  padding: var(--spacing-6);
  animation: fadeInUp 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-card {
  border-radius: var(--radius-2xl);
  padding: var(--spacing-10);
  box-shadow: var(--shadow-xl);
}

.login-brand {
  text-align: center;
  margin-bottom: var(--spacing-10);
}

.brand-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto var(--spacing-6);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border-radius: var(--radius-xl);
  color: white;
  box-shadow: var(--shadow-lg);
}

.brand-title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  margin: 0 0 var(--spacing-2);
  letter-spacing: -0.02em;
}

.brand-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.login-form {
  margin-bottom: var(--spacing-8);
}

.input-wrapper {
  position: relative;
  transition: all var(--transition-normal);
}

.input-wrapper:focus-within {
  transform: translateY(-2px);
}

.input-icon {
  position: absolute;
  left: var(--spacing-4);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-muted);
  z-index: 1;
  transition: color var(--transition-fast);
}

.input-wrapper:focus-within .input-icon {
  color: var(--color-accent);
}

.custom-input {
  padding-left: calc(var(--spacing-4) + 28px);
}

.custom-input :deep(.el-input__wrapper) {
  border-radius: var(--radius-lg);
  border-width: 2px;
  border-color: var(--color-border);
  transition: all var(--transition-normal);
}

.custom-input:focus-within :deep(.el-input__wrapper) {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(6, 182, 212, 0.1);
}

.form-actions {
  margin-bottom: 0;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  letter-spacing: 0.05em;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border: none;
  border-radius: var(--radius-lg);
  transition: all var(--transition-normal);
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-lg);
}

.login-btn:active:not(:disabled) {
  transform: translateY(0);
}

.login-footer {
  text-align: center;
  padding-top: var(--spacing-6);
  border-top: 1px solid var(--color-border-light);
}

.login-footer span {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}
</style>