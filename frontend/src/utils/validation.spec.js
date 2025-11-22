import { describe, it, expect } from 'vitest'
import {
  validateEmail,
  validateDisplayName,
  validatePassword,
  validateConfirmPassword,
  validateInitialPoints,
  validateBidderRegistrationForm,
  hasErrors
} from './validation'

describe('validation utilities', () => {
  describe('validateEmail', () => {
    it('returns error when email is empty', () => {
      expect(validateEmail('')).toBe('メールアドレスを入力してください')
      expect(validateEmail('   ')).toBe('メールアドレスを入力してください')
      expect(validateEmail(null)).toBe('メールアドレスを入力してください')
    })

    it('returns error when email format is invalid', () => {
      expect(validateEmail('invalid-email')).toBe('メールアドレスの形式が正しくありません')
      expect(validateEmail('test@')).toBe('メールアドレスの形式が正しくありません')
      expect(validateEmail('@example.com')).toBe('メールアドレスの形式が正しくありません')
      expect(validateEmail('test@example')).toBe('メールアドレスの形式が正しくありません')
    })

    it('returns error when email exceeds max length', () => {
      const longEmail = 'a'.repeat(246) + '@example.com' // 256 characters
      expect(validateEmail(longEmail)).toBe('メールアドレスは255文字以内で入力してください')
    })

    it('returns null for valid email', () => {
      expect(validateEmail('test@example.com')).toBeNull()
      expect(validateEmail('user.name@example.co.jp')).toBeNull()
      expect(validateEmail('test+tag@example.com')).toBeNull()
    })
  })

  describe('validateDisplayName', () => {
    it('returns null when display name is empty (optional field)', () => {
      expect(validateDisplayName('')).toBeNull()
      expect(validateDisplayName('   ')).toBeNull()
      expect(validateDisplayName(null)).toBeNull()
    })

    it('returns error when display name exceeds max length', () => {
      const longName = 'a'.repeat(101)
      expect(validateDisplayName(longName)).toBe('表示名は100文字以内で入力してください')
    })

    it('returns null for valid display name', () => {
      expect(validateDisplayName('入札者01')).toBeNull()
      expect(validateDisplayName('Test Bidder')).toBeNull()
      expect(validateDisplayName('a'.repeat(100))).toBeNull()
    })
  })

  describe('validatePassword', () => {
    it('returns error when password is empty', () => {
      expect(validatePassword('')).toBe('パスワードを入力してください')
      expect(validatePassword('   ')).toBe('パスワードを入力してください')
      expect(validatePassword(null)).toBe('パスワードを入力してください')
    })

    it('returns error when password is too short', () => {
      expect(validatePassword('short')).toBe('パスワードは8文字以上で入力してください')
      expect(validatePassword('1234567')).toBe('パスワードは8文字以上で入力してください')
    })

    it('returns null for valid password', () => {
      expect(validatePassword('password123')).toBeNull()
      expect(validatePassword('12345678')).toBeNull()
      expect(validatePassword('a'.repeat(8))).toBeNull()
      expect(validatePassword('a'.repeat(100))).toBeNull()
    })
  })

  describe('validateConfirmPassword', () => {
    it('returns error when confirmation password is empty', () => {
      expect(validateConfirmPassword('password123', '')).toBe('確認用パスワードを入力してください')
      expect(validateConfirmPassword('password123', '   ')).toBe('確認用パスワードを入力してください')
      expect(validateConfirmPassword('password123', null)).toBe('確認用パスワードを入力してください')
    })

    it('returns error when passwords do not match', () => {
      expect(validateConfirmPassword('password123', 'password456')).toBe('パスワードが一致しません')
      expect(validateConfirmPassword('password', 'Password')).toBe('パスワードが一致しません')
    })

    it('returns null when passwords match', () => {
      expect(validateConfirmPassword('password123', 'password123')).toBeNull()
      expect(validateConfirmPassword('12345678', '12345678')).toBeNull()
    })
  })

  describe('validateInitialPoints', () => {
    it('returns null when points is empty (optional field)', () => {
      expect(validateInitialPoints('')).toBeNull()
      expect(validateInitialPoints(null)).toBeNull()
      expect(validateInitialPoints(undefined)).toBeNull()
    })

    it('returns error when points is not a number', () => {
      expect(validateInitialPoints('abc')).toBe('数値を入力してください')
      expect(validateInitialPoints('12abc')).toBe('数値を入力してください')
    })

    it('returns error when points is not an integer', () => {
      expect(validateInitialPoints('100.5')).toBe('整数で入力してください')
      expect(validateInitialPoints(100.5)).toBe('整数で入力してください')
      expect(validateInitialPoints('99.99')).toBe('整数で入力してください')
    })

    it('returns error when points is negative', () => {
      expect(validateInitialPoints('-100')).toBe('0以上の整数を入力してください')
      expect(validateInitialPoints(-1)).toBe('0以上の整数を入力してください')
    })

    it('returns null for valid points', () => {
      expect(validateInitialPoints('0')).toBeNull()
      expect(validateInitialPoints(0)).toBeNull()
      expect(validateInitialPoints('100')).toBeNull()
      expect(validateInitialPoints(100)).toBeNull()
      expect(validateInitialPoints('1000000')).toBeNull()
      expect(validateInitialPoints(1000000)).toBeNull()
    })
  })

  describe('validateBidderRegistrationForm', () => {
    it('returns all errors when all fields are empty', () => {
      const errors = validateBidderRegistrationForm({
        email: '',
        password: '',
        confirmPassword: '',
        display_name: '',
        initial_points: ''
      })

      expect(errors.email).toBe('メールアドレスを入力してください')
      expect(errors.password).toBe('パスワードを入力してください')
      expect(errors.confirmPassword).toBe('確認用パスワードを入力してください')
      expect(errors.display_name).toBeUndefined() // Optional field
      expect(errors.initial_points).toBeUndefined() // Optional field
    })

    it('returns no errors for valid form data with all fields', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: '入札者01',
        initial_points: '1000'
      })

      expect(Object.keys(errors).length).toBe(0)
    })

    it('returns no errors for valid form data with only required fields', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: '',
        initial_points: ''
      })

      expect(Object.keys(errors).length).toBe(0)
    })

    it('returns email error when email format is invalid', () => {
      const errors = validateBidderRegistrationForm({
        email: 'invalid-email',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: '',
        initial_points: ''
      })

      expect(errors.email).toBe('メールアドレスの形式が正しくありません')
    })

    it('returns password errors when passwords do not match', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'different123',
        display_name: '',
        initial_points: ''
      })

      expect(errors.confirmPassword).toBe('パスワードが一致しません')
    })

    it('returns display_name error when exceeds max length', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: 'a'.repeat(101),
        initial_points: ''
      })

      expect(errors.display_name).toBe('表示名は100文字以内で入力してください')
    })

    it('returns initial_points error when negative', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: '',
        initial_points: '-100'
      })

      expect(errors.initial_points).toBe('0以上の整数を入力してください')
    })

    it('returns multiple errors when multiple fields are invalid', () => {
      const errors = validateBidderRegistrationForm({
        email: 'invalid-email',
        password: 'short',
        confirmPassword: 'different',
        display_name: 'a'.repeat(101),
        initial_points: '-100'
      })

      expect(errors.email).toBeDefined()
      expect(errors.password).toBeDefined()
      expect(errors.confirmPassword).toBeDefined()
      expect(errors.display_name).toBeDefined()
      expect(errors.initial_points).toBeDefined()
    })

    it('allows zero as valid initial points', () => {
      const errors = validateBidderRegistrationForm({
        email: 'bidder@example.com',
        password: 'password123',
        confirmPassword: 'password123',
        display_name: '',
        initial_points: '0'
      })

      expect(errors.initial_points).toBeUndefined()
      expect(Object.keys(errors).length).toBe(0)
    })
  })

  describe('hasErrors', () => {
    it('returns false for empty object', () => {
      expect(hasErrors({})).toBe(false)
    })

    it('returns true for object with errors', () => {
      expect(hasErrors({ email: 'Error message' })).toBe(true)
      expect(hasErrors({ 
        email: 'Error', 
        password: 'Error' 
      })).toBe(true)
    })

    it('returns true even with one error', () => {
      expect(hasErrors({ 
        email: 'Error message',
        password: null,
        display_name: undefined
      })).toBe(true)
    })
  })
})
