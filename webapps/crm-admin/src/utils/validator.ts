/** 定义一个验证接口 */
export interface Verifier {
  verify(validator: Validator): void;
}

/** 通用验证器类 */
export class Validator {
  /** 收集错误信息 */
  public errors: Record<string, string>;

  constructor(errors?: Record<string, string>) {
    this.errors = errors ?? {};
  }

  /** ========= 静态方法 ========= */

  /** 执行校验 */
  static doVerify(obj: Verifier): Validator {
    const v = new Validator();
    obj.verify(v);
    return v;
  }

  /** 校验是否是手机号 */
  static isPhone(input?: string | null): boolean {
    if (!input) return false;
    const phoneRx = /^1[3-9]\d{9}$/;
    return phoneRx.test(input);
  }

  /** ========= 实例方法 ========= */

  /** 添加一个错误 */
  addError(field: string, message: string): void {
    const msg = this.errors[field];
    if (!msg) {
      this.errors[field] = message;
    }
  }

  /** 是否校验通过 */
  isValid(): boolean {
    return Object.keys(this.errors).length === 0;
  }

  /**
   * 根据 field 获取错误
   * 不传 field 时返回第一个错误
   */
  errorBy(field?: string): string {
    if (Object.keys(this.errors).length === 0) return "";

    if (field && field.trim()) {
      return this.errors[field] ?? "";
    }

    for (const key in this.errors) {
      const val = this.errors[key];
      if (val?.trim()) {
        return val;
      }
    }

    return "";
  }

  /** 是否是其中之一 */
  oneOf<T>(val: T, options: T[]): boolean {
    return options.includes(val);
  }

  /**
   * 当 expr 为 false 时添加错误
   * expr 表示“是否满足条件”
   */
  check(expr: boolean, field: string, message: string): void {
    if (!expr) {
      this.addError(field, message);
    }
  }

  /** 正则匹配 */
  matches(input: string | null | undefined, rx: RegExp): boolean {
    if (!input) return false;
    return rx.test(input);
  }

  /** 必填校验（trim 后为空则报错） */
  require(
    field: string,
    input: string | null | undefined,
    message?: string
  ): boolean {
    if (!input || input.trim().length === 0) {
      this.addError(field, message ?? `${field} 不能为空`);
      return false;
    }
    return true;
  }

  /** 最小长度 */
  minLen(
    size: number,
    field: string,
    input: string | null | undefined,
    message?: string
  ): boolean {
    if (!input || input.length < size) {
      const msg = message ?? `${field} 长度必须 >= ${size}`;
      this.addError(field, msg);
      return false;
    }
    return true;
  }

  /** 最大长度 */
  maxLen(
    size: number,
    field: string,
    input: string | null | undefined,
    message?: string
  ): boolean {
    if (!input || input.length > size) {
      const msg = message ?? `${field} 长度必须 <= ${size}`;
      this.addError(field, msg);
      return false;
    }
    return true;
  }
}
