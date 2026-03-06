<template>
  <div class="drift-app">
    <!-- 大海背景 -->
    <div class="ocean-bg">
      <div class="wave wave1"></div>
      <div class="wave wave2"></div>
      <div class="wave wave3"></div>
    </div>

    <div class="app-inner">
    <!-- 顶部 -->
    <header class="app-header">
      <span class="app-title">漂流瓶</span>
      <button class="text-btn" @click="goPage('myBottles')">我的瓶子</button>
    </header>

    <!-- 首页 -->
    <div v-if="page === 'home'" class="page">
      <div class="bottle-float">🍾</div>
      <p class="slogan">将心事装进瓶子，漂向远方</p>
      <div class="btn-group">
        <button class="btn btn-orange" @click="goPage('throw')">扔瓶子</button>
        <button class="btn btn-white" @click="pickABottle">捞瓶子</button>
      </div>
    </div>

    <!-- 扔瓶子 -->
    <div v-else-if="page === 'throw'" class="page">
      <button class="back-btn" @click="goPage('home')">← 返回</button>
      <h2 class="page-title">写下你的心事</h2>
      <div class="form-group">
        <label class="form-label">你的昵称</label>
        <input v-model="throwForm.nickname" class="input" placeholder="匿名旅人" maxlength="20" />
      </div>
      <div class="form-group">
        <label class="form-label">瓶子里装什么？</label>
        <textarea
          v-model="throwForm.content"
          class="textarea"
          placeholder="写下想说的话，最多500字..."
          maxlength="500"
          rows="6"
        ></textarea>
        <div class="word-count">{{ throwForm.content.length }}/500</div>
      </div>
      <button class="btn btn-orange btn-full" :disabled="throwing" @click="handleThrow">
        {{ throwing ? '扔出去中...' : '扔出去！' }}
      </button>
      <div v-if="throwSuccess" class="success-tip">🍾 瓶子已漂向大海...</div>
    </div>

    <!-- 捞到的瓶子 -->
    <div v-else-if="page === 'pickedBottle'" class="page">
      <button class="back-btn" @click="goPage('home')">← 返回大海</button>
      <h2 class="page-title">你捞到了一个瓶子</h2>
      <div v-if="pickedBottle" class="bottle-card">
        <div class="card-meta">来自：{{ pickedBottle.nickname || '匿名旅人' }}</div>
        <div class="card-content">{{ pickedBottle.content }}</div>
      </div>
      <div v-if="!showReplyForm" class="btn-group mt16">
        <button class="btn btn-blue" @click="showReplyForm = true">回复TA</button>
        <button class="btn btn-white" @click="pickABottle">再捞一个</button>
      </div>
      <div v-if="showReplyForm" class="reply-area">
        <div class="form-group">
          <label class="form-label">你的昵称</label>
          <input v-model="replyForm.nickname" class="input" placeholder="匿名旅人" maxlength="20" />
        </div>
        <div class="form-group">
          <label class="form-label">回复内容</label>
          <textarea
            v-model="replyForm.content"
            class="textarea"
            placeholder="写下你的回复..."
            maxlength="500"
            rows="4"
          ></textarea>
          <div class="word-count">{{ replyForm.content.length }}/500</div>
        </div>
        <button class="btn btn-blue btn-full" :disabled="replying" @click="handleReply">
          {{ replying ? '发送中...' : '发送回复' }}
        </button>
        <div v-if="replySuccess" class="success-tip">💌 回复已发出~</div>
      </div>
    </div>

    <!-- 我的瓶子 -->
    <div v-else-if="page === 'myBottles'" class="page">
      <button class="back-btn" @click="goPage('home')">← 返回大海</button>
      <h2 class="page-title">我的瓶子</h2>
      <div v-if="myBottleList.length === 0 && !loadingBottles" class="empty-state">
        <div class="empty-icon">🌊</div>
        <p>还没有扔过瓶子，去大海里扔一个吧！</p>
      </div>
      <div
        v-for="bottle in myBottleList"
        :key="bottle.ID"
        class="bottle-card clickable"
        @click="viewDetail(bottle)"
      >
        <div class="card-meta">
          <span>{{ formatDate(bottle.CreatedAt) }}</span>
          <span :class="['status-tag', statusClass(bottle.status)]">{{ statusLabel(bottle.status) }}</span>
        </div>
        <div class="card-content">{{ bottle.content }}</div>
      </div>
      <button v-if="myBottleTotal > myBottleList.length" class="text-btn load-more" @click="loadMoreBottles">
        加载更多
      </button>
    </div>

    <!-- 瓶子详情 -->
    <div v-else-if="page === 'bottleDetail'" class="page">
      <button class="back-btn" @click="goPage('myBottles')">← 返回我的瓶子</button>
      <h2 class="page-title">瓶子详情</h2>
      <div v-if="detailBottle" class="bottle-card">
        <div class="card-meta">投出时间：{{ formatDate(detailBottle.CreatedAt) }}</div>
        <div class="card-content">{{ detailBottle.content }}</div>
      </div>
      <div class="replies-section">
        <h3 class="replies-title">收到的回复 ({{ detailReplies.length }})</h3>
        <div v-if="detailReplies.length === 0" class="empty-state small">
          <p>还没有人回复你</p>
        </div>
        <div v-for="reply in detailReplies" :key="reply.ID" class="reply-card">
          <div class="card-meta">{{ reply.nickname || '匿名旅人' }} · {{ formatDate(reply.CreatedAt) }}</div>
          <div class="card-content">{{ reply.content }}</div>
        </div>
      </div>
    </div>

    <!-- 昵称初始化弹窗 -->
    <div v-if="showNicknameModal" class="modal-mask">
      <div class="modal">
        <h3 class="modal-title">欢迎来到漂流瓶</h3>
        <p class="modal-desc">给自己取一个昵称吧（后续可在扔瓶子时修改）</p>
        <input v-model="nicknameInput" class="input" placeholder="匿名旅人" maxlength="20" />
        <button class="btn btn-blue btn-full mt16" @click="saveNickname">开始漂流</button>
      </div>
    </div>

    <!-- Toast 提示 -->
    <div v-if="toastMsg" class="toast">{{ toastMsg }}</div>

    <!-- 底部广告位 -->
    <div v-if="showAd" class="ad-container">
      <div v-if="loadingAd" class="ad-content">
        <div class="ad-loading">广告加载中...</div>
      </div>
      <div v-else-if="currentAd" class="ad-content">
        <div class="ad-badge">广告</div>
        <div class="ad-main">
          <div class="ad-title">{{ currentAd.title }}</div>
          <div class="ad-desc">{{ currentAd.description }}</div>
          <a href="javascript:void(0)" @click="handleAdClick" class="ad-link">查看详情</a>
          <div v-if="adList.length > 1" class="ad-indicators">
            <span 
              v-for="(ad, index) in adList" 
              :key="ad.id"
              class="ad-indicator"
              :class="{ active: index === currentAdIndex }"
              @click="switchAd(index)"
            ></span>
          </div>
        </div>
        <div class="ad-actions">
          <button class="ad-refresh" @click="refreshAd" title="刷新广告">
            ↻
          </button>
          <button class="ad-close" @click="closeAd" title="关闭广告">×</button>
        </div>
      </div>
    </div>
    </div><!-- /app-inner -->
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { throwBottle, pickBottle, replyBottle, getMyBottles, getBottleDetail } from '@/api/driftBottle'
import { getSuitableAd, recordAdView, recordAdClick } from '@/api/advertisement'

// 设备ID
const getDeviceId = () => {
  let id = localStorage.getItem('drift_device_id')
  if (!id) {
    id = 'dev_' + Math.random().toString(36).slice(2) + Date.now().toString(36)
    localStorage.setItem('drift_device_id', id)
  }
  return id
}
const deviceId = getDeviceId()

// Toast
const toastMsg = ref('')
let toastTimer = null
const showToast = (msg) => {
  toastMsg.value = msg
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toastMsg.value = '' }, 2500)
}

// 昵称
const showNicknameModal = ref(false)
const nicknameInput = ref('')
const savedNickname = ref(localStorage.getItem('drift_nickname') || '')

const saveNickname = () => {
  const name = nicknameInput.value.trim() || '匿名旅人'
  savedNickname.value = name
  localStorage.setItem('drift_nickname', name)
  showNicknameModal.value = false
}

// 页面
const page = ref('home')
const goPage = (p) => { page.value = p }

// 广告位功能
const showAd = ref(true)
const currentAd = ref(null)
const loadingAd = ref(false)
const adList = ref([])
const currentAdIndex = ref(0)
const adTimer = ref(null)


// 关闭广告
const closeAd = () => {
  showAd.value = false
  stopAdRotation()
  // 可以设置一段时间后再显示，比如24小时
  localStorage.setItem('drift_ad_closed', Date.now().toString())
}

// 检查广告是否在24小时内被关闭过
const isAdClosed = () => {
  const closedTime = localStorage.getItem('drift_ad_closed')
  if (!closedTime) return false
  return Date.now() - parseInt(closedTime) < 24 * 60 * 60 * 1000
}

// 获取广告数据
const fetchAdData = async () => {
  if (loadingAd.value) return
  
  loadingAd.value = true
  try {
    const res = await getSuitableAd({ deviceId })
    if (res.code === 0 && res.data) {
      let ads = Array.isArray(res.data) ? res.data : [res.data]
      
      adList.value = ads
      
      // 有forcePopup广告则强制显示，否则检查是否被关闭过
      const hasForceAd = ads.some(ad => ad.forcePopup === true)
      showAd.value = hasForceAd ? true : !isAdClosed()
      
      if (adList.value.length > 0) {
        currentAd.value = adList.value[0]
        startAdRotation()
        
        // 记录广告展示
        if (currentAd.value && currentAd.value.id) {
          recordAdView({
            adId: currentAd.value.id,
            deviceId,
            position: 'bottom'
          }).catch(() => {})
        }
      }
    } else {
      // 如果API没有返回广告，使用默认广告
      const defaultAds = getDefaultAds()
      adList.value = defaultAds
      currentAd.value = defaultAds[0]
      startAdRotation()
    }
  } catch (error) {
    console.error('获取广告失败:', error)
    // 使用默认广告
    const defaultAds = getDefaultAds()
    adList.value = defaultAds
    currentAd.value = defaultAds[0]
    startAdRotation()
  } finally {
    loadingAd.value = false
  }
}

// 广告轮播
const startAdRotation = () => {
  if (adList.value.length <= 1) return
  
  clearInterval(adTimer.value)
  adTimer.value = setInterval(() => {
    currentAdIndex.value = (currentAdIndex.value + 1) % adList.value.length
    currentAd.value = adList.value[currentAdIndex.value]
    
    // 记录新广告展示
    if (currentAd.value && currentAd.value.id) {
      recordAdView({
        adId: currentAd.value.id,
        deviceId,
        position: 'bottom'
      }).catch(() => {})
    }
  }, 10000) // 每10秒切换一次广告
}

// 停止广告轮播
const stopAdRotation = () => {
  clearInterval(adTimer.value)
  adTimer.value = null
}

// 手动切换广告
const switchAd = (index) => {
  if (index >= 0 && index < adList.value.length) {
    currentAdIndex.value = index
    currentAd.value = adList.value[index]
    
    // 重置轮播计时器
    if (adTimer.value) {
      clearInterval(adTimer.value)
      startAdRotation()
    }
    
    // 记录广告展示
    if (currentAd.value && currentAd.value.id) {
      recordAdView({
        adId: currentAd.value.id,
        deviceId,
        position: 'bottom'
      }).catch(() => {})
    }
  }
}

// 获取默认广告列表（当API失败时使用）
const getDefaultAds = () => {
  return [
    {
      id: 1,
      title: '探索更多有趣应用',
      description: '发现更多创意H5应用，体验不一样的互动乐趣',
      link: 'https://example.com/more-apps'
    },
    {
      id: 2,
      title: 'AI创意工具推荐',
      description: '使用AI工具提升创作效率，激发无限创意',
      link: 'https://example.com/ai-tools'
    },
    {
      id: 3,
      title: '开发者学习资源',
      description: '免费学习前端开发，掌握最新技术栈',
      link: 'https://example.com/learn-dev'
    },
    {
      id: 4,
      title: '漂流瓶技巧分享',
      description: '学习如何写出吸引人的漂流瓶内容',
      link: 'https://example.com/drift-tips'
    }
  ]
}

// 广告链接点击处理
const handleAdClick = () => {
  if (!currentAd.value) return
  
  // 记录广告点击
  if (currentAd.value.id) {
    recordAdClick({
      adId: currentAd.value.id,
      deviceId,
      position: 'bottom'
    }).catch(() => {}) // 静默失败
  }
  
  // 在新标签页打开链接
  window.open(currentAd.value.link, '_blank')
}

// 刷新广告
const refreshAd = async () => {
  stopAdRotation()
  await fetchAdData()
}

// 日期格式化
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 状态
const statusLabel = (s) => ({ 1: '漂流中', 2: '已被捞起', 3: '已回复' }[s] || '未知')
const statusClass = (s) => ({ 1: 'tag-blue', 2: 'tag-orange', 3: 'tag-green' }[s] || '')

// 扔瓶子
const throwing = ref(false)
const throwSuccess = ref(false)
const throwForm = ref({ nickname: savedNickname.value, content: '' })

const handleThrow = async () => {
  if (!throwForm.value.content.trim()) { showToast('瓶子里没有内容哦'); return }
  throwing.value = true
  throwSuccess.value = false
  try {
    const res = await throwBottle({ nickname: throwForm.value.nickname || '匿名旅人', content: throwForm.value.content, deviceId })
    if (res.code === 0) {
      throwSuccess.value = true
      throwForm.value.content = ''
      setTimeout(() => { throwSuccess.value = false; goPage('home') }, 2000)
    } else {
      showToast(res.msg || '扔瓶子失败')
    }
  } finally {
    throwing.value = false
  }
}

// 捞瓶子
const pickedBottle = ref(null)
const pickABottle = async () => {
  const res = await pickBottle({ deviceId })
  if (res.code === 0) {
    pickedBottle.value = res.data
    showReplyForm.value = false
    replySuccess.value = false
    replyForm.value.content = ''
    goPage('pickedBottle')
  } else {
    showToast(res.msg || '暂时没有可捞的瓶子')
  }
}

// 回复瓶子
const replying = ref(false)
const replySuccess = ref(false)
const showReplyForm = ref(false)
const replyForm = ref({ nickname: savedNickname.value, content: '' })

const handleReply = async () => {
  if (!replyForm.value.content.trim()) { showToast('回复内容不能为空'); return }
  replying.value = true
  replySuccess.value = false
  try {
    const res = await replyBottle({ bottleId: pickedBottle.value.ID, nickname: replyForm.value.nickname || '匿名旅人', content: replyForm.value.content, deviceId })
    if (res.code === 0) {
      replySuccess.value = true
      replyForm.value.content = ''
      setTimeout(() => { replySuccess.value = false; goPage('home') }, 2000)
    } else {
      showToast(res.msg || '回复失败')
    }
  } finally {
    replying.value = false
  }
}

// 我的瓶子
const myBottleList = ref([])
const myBottleTotal = ref(0)
const myBottlePage = ref(1)
const loadingBottles = ref(false)

const loadMyBottles = async (reset = true) => {
  if (reset) { myBottlePage.value = 1; myBottleList.value = [] }
  loadingBottles.value = true
  try {
    const res = await getMyBottles({ deviceId, page: myBottlePage.value, pageSize: 10 })
    if (res.code === 0) {
      myBottleList.value = reset ? res.data.list : [...myBottleList.value, ...res.data.list]
      myBottleTotal.value = res.data.total
    }
  } finally {
    loadingBottles.value = false
  }
}

const loadMoreBottles = async () => { myBottlePage.value++; await loadMyBottles(false) }

watch(page, (val) => { if (val === 'myBottles') loadMyBottles() })

// 瓶子详情
const detailBottle = ref(null)
const detailReplies = ref([])

const viewDetail = async (bottle) => {
  detailBottle.value = bottle
  detailReplies.value = []
  goPage('bottleDetail')
  const res = await getBottleDetail({ id: bottle.ID })
  if (res.code === 0) {
    detailBottle.value = res.data.bottle
    detailReplies.value = res.data.replies || []
  }
}

onMounted(() => {
  if (!savedNickname.value) showNicknameModal.value = true
  throwForm.value.nickname = savedNickname.value
  replyForm.value.nickname = savedNickname.value
  
  // 初始化广告
  fetchAdData()
})

// 组件卸载时清理
import { onUnmounted } from 'vue'
onUnmounted(() => {
  stopAdRotation()
  clearTimeout(toastTimer)
})
</script>

<style scoped>
/* ===== 基础 ===== */
.drift-app {
  min-height: 100vh;
  background: linear-gradient(180deg, #d4e8f0 0%, #a8d5e8 40%, #4ba8d1 100%);
  position: relative;
  overflow-x: hidden;
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
  -webkit-font-smoothing: antialiased;
  display: flex;
  justify-content: center;
}

.app-inner {
  width: 100%;
  max-width: 480px;
  min-height: 100vh;
  position: relative;
  display: flex;
  flex-direction: column;
}

/* ===== 海浪 ===== */
.ocean-bg {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 160px;
  pointer-events: none;
  z-index: 0;
}

.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 200%;
  height: 100px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 50% 50% 0 0;
  animation: wave-move 8s linear infinite;
}
.wave2 { background: rgba(255,255,255,0.1); animation-duration: 10s; animation-delay: -2s; height: 80px; }
.wave3 { background: rgba(255,255,255,0.08); animation-duration: 12s; animation-delay: -4s; height: 60px; }

@keyframes wave-move {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}

/* ===== Header ===== */
.app-header {
  position: relative;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
}
.app-title { font-size: 20px; font-weight: 700; color: #fff; text-shadow: 0 2px 8px rgba(0,0,0,0.2); }

/* ===== Page ===== */
.page {
  position: relative;
  z-index: 10;
  padding: 16px 20px 120px;
  flex: 1;
}

.back-btn {
  background: none;
  border: none;
  color: rgba(255,255,255,0.85);
  font-size: 14px;
  cursor: pointer;
  padding: 0;
  margin-bottom: 16px;
}
.back-btn:hover { color: #fff; }

.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 20px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.15);
}

/* ===== 首页 ===== */
.bottle-float {
  font-size: 80px;
  text-align: center;
  margin: 30px 0 10px;
  animation: float 4s ease-in-out infinite;
}
@keyframes float {
  0%,100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}
.slogan { text-align: center; color: rgba(255,255,255,0.9); font-size: 15px; margin-bottom: 36px; }

/* ===== 按钮 ===== */
.btn-group { display: flex; gap: 14px; justify-content: center; }
.btn {
  height: 48px;
  padding: 0 32px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: transform 0.15s, opacity 0.15s;
}
.btn:active { transform: scale(0.95); }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-orange { background: #e8a85c; color: #fff; box-shadow: 0 4px 16px rgba(232,168,92,0.45); }
.btn-white  { background: rgba(255,255,255,0.92); color: #1b6b8f; box-shadow: 0 4px 16px rgba(0,0,0,0.1); }
.btn-blue   { background: #4ba8d1; color: #fff; box-shadow: 0 4px 16px rgba(75,168,209,0.45); }
.btn-full   { width: 100%; margin-top: 16px; }

.text-btn { background: none; border: none; color: rgba(255,255,255,0.85); font-size: 14px; cursor: pointer; }
.text-btn:hover { color: #fff; }

/* ===== 表单 ===== */
.form-group { margin-bottom: 16px; }
.form-label { display: block; color: rgba(255,255,255,0.9); font-size: 14px; margin-bottom: 8px; }
.input, .textarea {
  width: 100%;
  background: rgba(255,255,255,0.92);
  border: none;
  border-radius: 12px;
  padding: 12px 14px;
  font-size: 15px;
  color: #303133;
  outline: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  font-family: inherit;
}
.textarea { resize: none; }
.word-count { text-align: right; color: rgba(255,255,255,0.7); font-size: 12px; margin-top: 4px; }

/* ===== 卡片 ===== */
.bottle-card {
  background: rgba(255,255,255,0.92);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}
.bottle-card.clickable { cursor: pointer; transition: transform 0.2s; }
.bottle-card.clickable:hover { transform: translateY(-2px); }

.card-meta { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: #909399; margin-bottom: 8px; }
.card-content { font-size: 15px; color: #303133; line-height: 1.7; white-space: pre-wrap; word-break: break-all; }

/* ===== 状态标签 ===== */
.status-tag { font-size: 11px; padding: 2px 8px; border-radius: 10px; }
.tag-blue   { background: #e8f4fd; color: #4ba8d1; }
.tag-orange { background: #fdf6ec; color: #e6a23c; }
.tag-green  { background: #f0f9eb; color: #67c23a; }

/* ===== 回复 ===== */
.reply-area { margin-top: 20px; }
.mt16 { margin-top: 16px; }

.replies-section { margin-top: 20px; }
.replies-title { font-size: 16px; color: rgba(255,255,255,0.9); margin-bottom: 12px; }
.reply-card {
  background: rgba(255,255,255,0.85);
  border-radius: 12px;
  padding: 12px 14px;
  margin-bottom: 10px;
}

/* ===== 空状态 ===== */
.empty-state { text-align: center; color: rgba(255,255,255,0.8); padding: 40px 0; }
.empty-state.small { padding: 16px 0; }
.empty-icon { font-size: 48px; margin-bottom: 12px; }
.load-more { display: block; margin: 12px auto 0; }

/* ===== 成功提示 ===== */
.success-tip { text-align: center; color: #fff; font-size: 16px; margin-top: 16px; }

/* ===== 昵称弹窗 ===== */
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.modal {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  width: 100%;
  max-width: 320px;
}
.modal-title { font-size: 18px; font-weight: 700; color: #303133; margin-bottom: 8px; }
.modal-desc  { font-size: 14px; color: #606266; margin-bottom: 16px; }
.modal .input { background: #f5f7fa; box-shadow: none; border: 1px solid #e4e7ed; }
.modal .btn-blue { box-shadow: none; }

/* ===== Toast ===== */
.toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0,0,0,0.65);
  color: #fff;
  font-size: 14px;
  padding: 10px 20px;
  border-radius: 20px;
  z-index: 200;
  white-space: nowrap;
}

/* ===== 底部广告位 ===== */
.ad-container {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 480px;
  background: rgba(255, 255, 255, 0.95);
  border-top: 1px solid #e0e0e0;
  padding: 12px 16px;
  z-index: 150;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.ad-content {
  display: flex;
  align-items: center;
  position: relative;
}

.ad-badge {
  background: #ff6b6b;
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-right: 12px;
  flex-shrink: 0;
}

.ad-main {
  flex: 1;
  min-width: 0;
}

.ad-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ad-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ad-link {
  display: inline-block;
  font-size: 12px;
  color: #4ba8d1;
  text-decoration: none;
  margin-top: 4px;
  font-weight: 500;
}

.ad-link:hover {
  text-decoration: underline;
}

.ad-close {
  background: none;
  border: none;
  color: #999;
  font-size: 20px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  margin-left: 8px;
  border-radius: 50%;
}

.ad-close:hover {
  background: #f5f5f5;
  color: #666;
}

.ad-loading {
  font-size: 13px;
  color: #999;
  text-align: center;
  width: 100%;
  padding: 8px 0;
}

/* 广告指示器 */
.ad-indicators {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-top: 8px;
}

.ad-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #ddd;
  cursor: pointer;
  transition: background 0.3s, transform 0.3s;
}

.ad-indicator:hover {
  background: #bbb;
  transform: scale(1.2);
}

.ad-indicator.active {
  background: #4ba8d1;
  transform: scale(1.2);
}

/* 广告操作按钮 */
.ad-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-left: 8px;
}

.ad-refresh {
  background: none;
  border: none;
  color: #999;
  font-size: 16px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.3s;
}

.ad-refresh:hover {
  background: #f5f5f5;
  color: #4ba8d1;
  transform: rotate(90deg);
}


</style>
