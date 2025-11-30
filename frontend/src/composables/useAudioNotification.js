/**
 * Audio Notification Composable
 * Web Audio APIを使用してプログラム生成音による通知を提供
 */
import { ref } from 'vue'

export function useAudioNotification() {
  const audioContext = ref(null)
  const isEnabled = ref(false)

  /**
   * AudioContextを初期化（ユーザー操作時に呼び出す必要あり）
   */
  function initAudio() {
    if (!audioContext.value) {
      audioContext.value = new (window.AudioContext || window.webkitAudioContext)()
    }
    
    // suspended状態の場合はresume
    if (audioContext.value.state === 'suspended') {
      audioContext.value.resume()
    }
    
    isEnabled.value = true
    return audioContext.value
  }

  /**
   * 基本的なビープ音を生成
   * @param {number} frequency - 周波数(Hz)
   * @param {number} duration - 持続時間(秒)
   * @param {string} type - 波形タイプ ('sine', 'square', 'sawtooth', 'triangle')
   * @param {number} volume - 音量 (0-1)
   */
  function playTone(frequency, duration, type = 'sine', volume = 0.3) {
    if (!audioContext.value || !isEnabled.value) {
      return
    }

    const ctx = audioContext.value
    const oscillator = ctx.createOscillator()
    const gainNode = ctx.createGain()

    oscillator.connect(gainNode)
    gainNode.connect(ctx.destination)

    oscillator.type = type
    oscillator.frequency.setValueAtTime(frequency, ctx.currentTime)

    // フェードイン・アウトでクリックノイズを防ぐ
    gainNode.gain.setValueAtTime(0, ctx.currentTime)
    gainNode.gain.linearRampToValueAtTime(volume, ctx.currentTime + 0.01)
    gainNode.gain.linearRampToValueAtTime(0, ctx.currentTime + duration)

    oscillator.start(ctx.currentTime)
    oscillator.stop(ctx.currentTime + duration)
  }

  /**
   * 価格開示時の通知音: 高めのビープ2回
   */
  function playPriceOpenedSound() {
    if (!audioContext.value || !isEnabled.value) {
      return
    }

    // 高めの2回ビープ (880Hz = A5)
    playTone(880, 0.12, 'sine', 0.25)
    setTimeout(() => {
      playTone(880, 0.12, 'sine', 0.25)
    }, 150)
  }

  /**
   * 入札成功時の通知音: 短い上昇音
   */
  function playBidSuccessSound() {
    if (!audioContext.value || !isEnabled.value) {
      return
    }

    const ctx = audioContext.value
    const oscillator = ctx.createOscillator()
    const gainNode = ctx.createGain()

    oscillator.connect(gainNode)
    gainNode.connect(ctx.destination)

    oscillator.type = 'sine'
    
    // 上昇音 (C5 → E5 → G5)
    oscillator.frequency.setValueAtTime(523, ctx.currentTime)        // C5
    oscillator.frequency.setValueAtTime(659, ctx.currentTime + 0.08)  // E5
    oscillator.frequency.setValueAtTime(784, ctx.currentTime + 0.16)  // G5

    gainNode.gain.setValueAtTime(0, ctx.currentTime)
    gainNode.gain.linearRampToValueAtTime(0.3, ctx.currentTime + 0.02)
    gainNode.gain.linearRampToValueAtTime(0.2, ctx.currentTime + 0.16)
    gainNode.gain.linearRampToValueAtTime(0, ctx.currentTime + 0.25)

    oscillator.start(ctx.currentTime)
    oscillator.stop(ctx.currentTime + 0.25)
  }

  /**
   * 他者入札時の通知音: 低めの単音
   */
  function playOtherBidSound() {
    if (!audioContext.value || !isEnabled.value) {
      return
    }

    // 低めの単音 (440Hz = A4)
    playTone(440, 0.15, 'sine', 0.2)
  }

  /**
   * 落札時の祝福音: 上昇アルペジオ
   */
  function playWinSound() {
    if (!audioContext.value || !isEnabled.value) {
      return
    }

    // C major アルペジオ
    const notes = [523, 659, 784, 1047] // C5, E5, G5, C6
    notes.forEach((freq, index) => {
      setTimeout(() => {
        playTone(freq, 0.2, 'sine', 0.25)
      }, index * 100)
    })
  }

  /**
   * オーディオを無効化
   */
  function disableAudio() {
    isEnabled.value = false
  }

  /**
   * オーディオを有効化
   */
  function enableAudio() {
    if (audioContext.value) {
      if (audioContext.value.state === 'suspended') {
        audioContext.value.resume()
      }
      isEnabled.value = true
    } else {
      initAudio()
    }
  }

  /**
   * クリーンアップ
   */
  function cleanup() {
    if (audioContext.value) {
      audioContext.value.close()
      audioContext.value = null
    }
    isEnabled.value = false
  }

  return {
    isEnabled,
    initAudio,
    playPriceOpenedSound,
    playBidSuccessSound,
    playOtherBidSound,
    playWinSound,
    enableAudio,
    disableAudio,
    cleanup,
  }
}
