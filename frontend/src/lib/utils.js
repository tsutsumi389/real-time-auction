import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * Tailwind CSS クラスをマージするユーティリティ関数
 * @param {...import('clsx').ClassValue} inputs
 * @returns {string}
 */
export function cn(...inputs) {
  return twMerge(clsx(inputs))
}
