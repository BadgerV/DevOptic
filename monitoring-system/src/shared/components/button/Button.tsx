import React, { forwardRef, useId } from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  label?: string;
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
  loading?: boolean;
  variant?: "default" | "outlined" | "filled" | "text";
  size?: "small" | "medium" | "large";
  fullWidth?: boolean;
  className?: string;
  iconClassName?: string;
  startIconClassName?: string;
  endIconClassName?: string;
  loadingIndicator?: React.ReactNode;
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      label,
      startIcon,
      endIcon,
      loading = false,
      variant = "default",
      size = "medium",
      fullWidth = false,
      className = "",
      iconClassName = "",
      startIconClassName = "",
      endIconClassName = "",
      loadingIndicator,
      disabled,
      ...props
    },
    ref
  ) => {
    const btnId = useId();

    const baseStyles = {
      button: {
        display: "inline-flex",
        alignItems: "center",
        justifyContent: "center",
        gap: "8px",
        borderRadius: "6px",
        border: "1px solid transparent",
        fontWeight: 500,
        cursor: disabled || loading ? "not-allowed" : "pointer",
        opacity: disabled || loading ? 0.6 : 1,
        transition: "all 0.2s ease-in-out",
        width: fullWidth ? "100%" : "auto",
        minWidth: "fit-content",
      },
      sizes: {
        small: {
          padding: "6px 12px",
          fontSize: "14px",
        },
        medium: {
          padding: "8px 16px",
          fontSize: "16px",
        },
        large: {
          padding: "12px 20px",
          fontSize: "18px",
        },
      },
      variants: {
        default: {
          backgroundColor: "#3b82f6",
          color: "#fff",
          border: "1px solid #3b82f6",
        },
        outlined: {
          backgroundColor: "transparent",
          color: "#3b82f6",
          border: "1px solid #3b82f6",
        },
        filled: {
          backgroundColor: "#1f2937",
          color: "#fff",
          border: "1px solid transparent",
        },
        text: {
          backgroundColor: "transparent",
          color: "#3b82f6",
          border: "none",
        },
      },
      icon: {
        display: "flex",
        alignItems: "center",
        fontSize: size === "small" ? "16px" : size === "large" ? "20px" : "18px",
      },
    };

    const combinedStyles = {
      ...baseStyles.button,
      ...baseStyles.sizes[size],
      ...baseStyles.variants[variant],
    };

    return (
      <button
        ref={ref}
        id={btnId}
        style={combinedStyles}
        className={`button-base ${className}`}
        disabled={disabled || loading}
        {...props}
      >
        {loading && (
          <span
            className={`button-loading ${iconClassName}`}
            style={baseStyles.icon}
          >
            {loadingIndicator ?? (
              <span
                style={{
                  width: "16px",
                  height: "16px",
                  border: "2px solid currentColor",
                  borderTopColor: "transparent",
                  borderRadius: "50%",
                  animation: "spin 0.6s linear infinite",
                }}
              ></span>
            )}
          </span>
        )}

        {!loading && startIcon && (
          <span
            className={`button-start-icon ${iconClassName} ${startIconClassName}`}
            style={baseStyles.icon}
          >
            {startIcon}
          </span>
        )}

        {label && <span className="button-label">{label}</span>}

        {!loading && endIcon && (
          <span
            className={`button-end-icon ${iconClassName} ${endIconClassName}`}
            style={baseStyles.icon}
          >
            {endIcon}
          </span>
        )}
      </button>
    );
  }
);

export default Button;
