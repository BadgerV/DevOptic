import React, { useState, forwardRef, useId } from "react";

// Main reusable Input component
const Input = forwardRef(
  (
    {
      type = "text",
      label,
      placeholder,
      value,
      onChange,
      onFocus,
      onBlur,
      disabled = false,
      required = false,
      error,
      helperText,
      startIcon,
      endIcon,
      className = "",
      inputClassName = "",
      labelClassName = "",
      containerClassName = "",
      inputWrapperClassName = "",
      iconClassName = "",
      startIconClassName = "",
      endIconClassName = "",
      helperTextClassName = "",
      errorClassName = "",
      variant = "default",
      size = "medium",
      floatingLabel = false, // New prop for floating label behavior
      ...props
    }: any,
    ref
  ) => {
    const [isFocused, setIsFocused] = useState(false);
    const inputId = useId();
    const errorId = useId();
    const helperTextId = useId();

    // Determine if label should be in "floated" position
    const isLabelFloated =
      floatingLabel &&
      (isFocused || value || type === "date" || type === "time");

    const handleFocus = (e: any) => {
      setIsFocused(true);
      onFocus?.(e);
    };

    const handleBlur = (e: any) => {
      setIsFocused(false);
      onBlur?.(e);
    };

    const baseStyles = {
      container: {
        display: "flex",
        flexDirection: "column",
        gap: "4px",
        width: "100%",
      },
      inputWrapper: {
        position: "relative",
        display: "flex",
        alignItems: "center",
      },
      floatingInputWrapper: {
        position: "relative",
        display: "flex",
        alignItems: "center",
        marginTop: "8px", // Space for floating label
      },
      input: {
        width: "100%",
        border: "1px solid #d1d5db",
        borderRadius: "6px",
        padding:
          size === "small"
            ? "8px 12px"
            : size === "large"
            ? "14px 16px"
            : "10px 14px",
        fontSize:
          size === "small" ? "14px" : size === "large" ? "18px" : "16px",
        lineHeight: "1.5",
        transition: "all 0.2s ease-in-out",
        outline: "none",
        backgroundColor: disabled ? "#f9fafb" : "#ffffff",
        color: disabled ? "#9ca3af" : "#111827",
        cursor: disabled ? "not-allowed" : "text",
        paddingLeft: startIcon
          ? size === "small"
            ? "40px"
            : size === "large"
            ? "48px"
            : "44px"
          : undefined,
        paddingRight: endIcon
          ? size === "small"
            ? "40px"
            : size === "large"
            ? "48px"
            : "44px"
          : undefined,

        // Focus state styles
        "&:focus": {
          border: "1px solid #9ca3af", // grey border
          backgroundColor: "#f3f4f6", // light grey background while focused
        },

        // Optional hover
        "&:hover": {
          borderColor: "#9ca3af",
        },
      },
      label: {
        fontSize:
          size === "small" ? "13px" : size === "large" ? "16px" : "14px",
        fontWeight: "500",
        color: "#374151",
        marginBottom: "2px",
      },
      floatingLabel: {
        position: "absolute",
        left: startIcon
          ? size === "small"
            ? "40px"
            : size === "large"
            ? "48px"
            : "44px"
          : size === "small"
          ? "12px"
          : size === "large"
          ? "16px"
          : "14px",
        top: "50%",
        transform: "translateY(-50%)",
        fontSize:
          size === "small" ? "14px" : size === "large" ? "18px" : "16px",
        fontWeight: "400",
        color: "#9ca3af",
        backgroundColor: "#ffffff",
        padding: "0 4px",
        pointerEvents: "none",
        transition: "all 0.2s ease-in-out",
        zIndex: 1,
      },
      floatingLabelActive: {
        top: "0px",
        transform: "translateY(-50%)",
        fontSize:
          size === "small" ? "12px" : size === "large" ? "14px" : "13px",
        fontWeight: "500",
        color: isFocused ? "#3b82f6" : "#374151",
        left: size === "small" ? "8px" : size === "large" ? "12px" : "10px",
      },
      icon: {
        position: "absolute",
        top: "50%",
        transform: "translateY(-50%)",
        color: "#6b7280",
        fontSize:
          size === "small" ? "16px" : size === "large" ? "20px" : "18px",
        pointerEvents: "none",
      },
      startIcon: {
        left: size === "small" ? "12px" : size === "large" ? "16px" : "14px",
      },
      endIcon: {
        right: size === "small" ? "12px" : size === "large" ? "16px" : "14px",
      },
      helperText: {
        fontSize: "12px",
        color: "#6b7280",
        marginTop: "2px",
      },
      errorText: {
        fontSize: "12px",
        color: "#ef4444",
        marginTop: "2px",
      },
    };

    // Dynamic styles based on state and variant
    let dynamicInputStyles: any = { ...baseStyles.input };

    if (isFocused && !disabled) {
      dynamicInputStyles.borderColor = "#3b82f6";
      dynamicInputStyles.boxShadow = "0 0 0 3px rgba(59, 130, 246, 0.1)";
    }

    if (error) {
      dynamicInputStyles.borderColor = "#ef4444";
      if (isFocused) {
        dynamicInputStyles.boxShadow = "0 0 0 3px rgba(239, 68, 68, 0.1)";
      }
    }

    // Variant-specific styles
    if (variant === "outlined") {
      dynamicInputStyles.backgroundColor = "transparent";
      dynamicInputStyles.borderWidth = "2px";
    } else if (variant === "filled") {
      dynamicInputStyles.backgroundColor = "#f3f4f6";
      dynamicInputStyles.border = "1px solid transparent";
      dynamicInputStyles.borderBottom = "2px solid #d1d5db";
      dynamicInputStyles.borderRadius = "6px 6px 0 0";
    } else if (variant === "underline") {
      dynamicInputStyles.backgroundColor = "transparent";
      dynamicInputStyles.border = "none";
      dynamicInputStyles.borderBottom = "2px solid #d1d5db";
      dynamicInputStyles.borderRadius = "0";
      dynamicInputStyles.padding =
        size === "small" ? "6px 0" : size === "large" ? "12px 0" : "8px 0";
    }

    // Floating label styles for different variants
    let floatingLabelStyles = { ...baseStyles.floatingLabel };
    if (isLabelFloated) {
      floatingLabelStyles = {
        ...floatingLabelStyles,
        ...baseStyles.floatingLabelActive,
      };
    }

    // Adjust floating label background for different variants
    if (variant === "filled") {
      floatingLabelStyles.backgroundColor = isLabelFloated
        ? "#ffffff"
        : "transparent";
    } else if (variant === "underline") {
      floatingLabelStyles.backgroundColor = "transparent";
      floatingLabelStyles.padding = "0";
    }

    return (
      <div
        style={{ ...baseStyles.container } as React.CSSProperties}
        className={`input-container ${containerClassName}`}
      >
        {label && !floatingLabel && (
          <label
            htmlFor={inputId}
            style={{ ...baseStyles.label }}
            className={`input-label ${labelClassName}`}
          >
            {label}
            {required && (
              <span style={{ color: "#ef4444", marginLeft: "2px" }}>*</span>
            )}
          </label>
        )}

        <div
          style={
            (floatingLabel
              ? baseStyles.floatingInputWrapper
              : baseStyles.inputWrapper) as React.CSSProperties
          }
          className={`input-wrapper ${inputWrapperClassName}`}
        >
          {floatingLabel && label && (
            <label
              htmlFor={inputId}
              style={{ floatingLabelStyles } as React.CSSProperties}
              className={`input-floating-label ${labelClassName}`}
            >
              {label}
              {required && (
                <span style={{ color: "#ef4444", marginLeft: "2px" }}>*</span>
              )}
            </label>
          )}

          {startIcon && (
            <span
              style={
                {
                  ...baseStyles.icon,
                  ...baseStyles.startIcon,
                } as React.CSSProperties
              }
              className={`input-icon input-start-icon ${iconClassName} ${startIconClassName}`}
            >
              {startIcon}
            </span>
          )}

          <input
            ref={ref}
            id={inputId}
            type={type}
            placeholder={floatingLabel ? "" : placeholder}
            value={value}
            onChange={onChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            disabled={disabled}
            required={required}
            style={{ ...dynamicInputStyles }}
            className={`input-field ${className} ${inputClassName}`}
            aria-describedby={
              error ? errorId : helperText ? helperTextId : undefined
            }
            aria-invalid={error ? "true" : "false"}
            {...props}
          />

          {endIcon && (
            <span
              style={
                {
                  ...baseStyles.icon,
                  ...baseStyles.endIcon,
                } as React.CSSProperties
              }
              className={`input-icon input-end-icon ${iconClassName} ${endIconClassName}`}
            >
              {endIcon}
            </span>
          )}
        </div>

        {error && (
          <span
            id={errorId}
            style={baseStyles.errorText}
            className={`input-error ${errorClassName}`}
            role="alert"
          >
            {error}
          </span>
        )}

        {helperText && !error && (
          <span
            id={helperTextId}
            style={baseStyles.helperText}
            className={`input-helper-text ${helperTextClassName}`}
          >
            {helperText}
          </span>
        )}
      </div>
    );
  }
);

// Input.displayName = 'Input';

// // Demo component showing different usage scenarios
// const InputDemo = () => {
//   const [formData, setFormData] = useState({
//     email: '',
//     password: '',
//     search: '',
//     phone: '',
//     amount: '',
//     message: ''
//   });

//   const [errors, setErrors] = useState({});

//   const handleInputChange = (field) => (e) => {
//     setFormData(prev => ({
//       ...prev,
//       [field]: e.target.value
//     }));

//     // Clear error when user starts typing
//     if (errors[field]) {
//       setErrors(prev => ({
//         ...prev,
//         [field]: ''
//       }));
//     }
//   };

//   const validateEmail = () => {
//     if (!formData.email.includes('@')) {
//       setErrors(prev => ({ ...prev, email: 'Please enter a valid email address' }));
//     }
//   };

//   return (
//     <div style={{ padding: '40px', maxWidth: '800px', margin: '0 auto', fontFamily: 'system-ui, sans-serif' }}>
//       <h1 style={{ marginBottom: '40px', color: '#111827', fontSize: '28px', fontWeight: 'bold' }}>
//         Reusable Input Component Demo
//       </h1>

//       <div style={{ display: 'grid', gap: '32px' }}>
//         {/* Basic Text Input */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Basic Text Inputs
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Email Address"
//               type="email"
//               placeholder="Enter your email"
//               value={formData.email}
//               onChange={handleInputChange('email')}
//               onBlur={validateEmail}
//               error={errors.email}
//               required
//               helperText="We'll never share your email with anyone else"
//             />

//             <Input
//               label="Password"
//               type="password"
//               placeholder="Enter your password"
//               value={formData.password}
//               onChange={handleInputChange('password')}
//               required
//               helperText="Must be at least 8 characters long"
//             />
//           </div>
//         </section>

//         {/* Floating Labels Demo */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Floating Labels
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Email Address"
//               type="email"
//               floatingLabel
//               value={formData.email}
//               onChange={handleInputChange('email')}
//               onBlur={validateEmail}
//               error={errors.email}
//               required
//               helperText="Enter your email address"
//             />

//             <Input
//               label="Full Name"
//               floatingLabel
//               value={formData.search}
//               onChange={handleInputChange('search')}
//               helperText="First and last name"
//             />

//             <Input
//               label="Password"
//               type="password"
//               floatingLabel
//               value={formData.password}
//               onChange={handleInputChange('password')}
//               required
//             />

//             <Input
//               label="Phone Number"
//               type="tel"
//               floatingLabel
//               startIcon="📞"
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//               helperText="Include country code"
//             />
//           </div>
//         </section>

//         {/* Floating Labels with Different Variants */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Floating Labels - Different Variants
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Default with Floating"
//               variant="default"
//               floatingLabel
//               value={formData.email}
//               onChange={handleInputChange('email')}
//             />

//             <Input
//               label="Outlined with Floating"
//               variant="outlined"
//               floatingLabel
//               value={formData.password}
//               onChange={handleInputChange('password')}
//             />

//             <Input
//               label="Filled with Floating"
//               variant="filled"
//               floatingLabel
//               value={formData.search}
//               onChange={handleInputChange('search')}
//             />

//             <Input
//               label="Underline with Floating"
//               variant="underline"
//               floatingLabel
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//             />
//           </div>
//         </section>
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Different Sizes
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Small Input"
//               size="small"
//               placeholder="Small size input"
//               value={formData.search}
//               onChange={handleInputChange('search')}
//             />

//             <Input
//               label="Medium Input (Default)"
//               size="medium"
//               placeholder="Medium size input"
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//             />

//             <Input
//               label="Large Input"
//               size="large"
//               placeholder="Large size input"
//               value={formData.amount}
//               onChange={handleInputChange('amount')}
//             />
//           </div>
//         </section>

//         {/* With Icons */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             With Icons
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Search"
//               placeholder="Search..."
//               startIcon="🔍"
//               value={formData.search}
//               onChange={handleInputChange('search')}
//             />

//             <Input
//               label="Phone Number"
//               type="tel"
//               placeholder="+1 (555) 123-4567"
//               startIcon="📞"
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//             />

//             <Input
//               label="Amount"
//               type="number"
//               placeholder="0.00"
//               startIcon="💰"
//               endIcon="USD"
//               value={formData.amount}
//               onChange={handleInputChange('amount')}
//             />
//           </div>
//         </section>

//         {/* Different Variants */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Different Variants
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Default Variant"
//               variant="default"
//               placeholder="Default styling"
//               value={formData.email}
//               onChange={handleInputChange('email')}
//             />

//             <Input
//               label="Outlined Variant"
//               variant="outlined"
//               placeholder="Outlined styling"
//               value={formData.password}
//               onChange={handleInputChange('password')}
//             />

//             <Input
//               label="Filled Variant"
//               variant="filled"
//               placeholder="Filled styling"
//               value={formData.search}
//               onChange={handleInputChange('search')}
//             />

//             <Input
//               label="Underline Variant"
//               variant="underline"
//               placeholder="Underline styling"
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//             />
//           </div>
//         </section>

//         {/* States */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Different States
//           </h2>
//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Normal State"
//               placeholder="This is a normal input"
//             />

//             <Input
//               label="Disabled State"
//               placeholder="This input is disabled"
//               disabled
//               value="Can't edit this"
//             />

//             <Input
//               label="Error State"
//               placeholder="This has an error"
//               error="This field is required"
//               value=""
//             />

//             <Input
//               label="With Helper Text"
//               placeholder="Input with helper text"
//               helperText="This is some helpful information"
//               value={formData.message}
//               onChange={handleInputChange('message')}
//             />
//           </div>
//         </section>

//         {/* Fully Customizable with CSS Classes */}
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Fully Overridable with CSS Classes
//           </h2>
//           <style>{`
//             /* Complete style override examples */
//             .modern-input-container {
//               background: linear-gradient(145deg, #f0f0f0, #cacaca);
//               padding: 20px;
//               border-radius: 15px;
//               box-shadow: 5px 5px 10px #bebebe, -5px -5px 10px #ffffff;
//             }

//             .modern-input-field {
//               background: linear-gradient(145deg, #ffffff, #e6e6e6) !important;
//               border: none !important;
//               border-radius: 25px !important;
//               padding: 15px 20px !important;
//               box-shadow: inset 2px 2px 5px #d1d1d1, inset -2px -2px 5px #ffffff !important;
//               font-size: 16px !important;
//               transition: all 0.3s ease !important;
//             }

//             .modern-input-field:focus {
//               box-shadow: inset 4px 4px 8px #c1c1c1, inset -4px -4px 8px #ffffff !important;
//               transform: scale(1.02);
//             }

//             .modern-input-label {
//               color: #666 !important;
//               font-weight: 600 !important;
//               font-size: 14px !important;
//               text-transform: uppercase !important;
//               letter-spacing: 1px !important;
//             }

//             .glassmorphism-container {
//               background: rgba(255, 255, 255, 0.25);
//               backdrop-filter: blur(10px);
//               border: 1px solid rgba(255, 255, 255, 0.18);
//               border-radius: 20px;
//               padding: 25px;
//               box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
//             }

//             .glassmorphism-field {
//               background: rgba(255, 255, 255, 0.1) !important;
//               border: 1px solid rgba(255, 255, 255, 0.3) !important;
//               border-radius: 15px !important;
//               backdrop-filter: blur(5px) !important;
//               color: #333 !important;
//               font-size: 16px !important;
//               padding: 12px 16px !important;
//             }

//             .glassmorphism-field:focus {
//               background: rgba(255, 255, 255, 0.2) !important;
//               border-color: rgba(255, 255, 255, 0.5) !important;
//               box-shadow: 0 0 20px rgba(255, 255, 255, 0.3) !important;
//             }

//             .glassmorphism-label {
//               color: #444 !important;
//               font-weight: 500 !important;
//               text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8) !important;
//             }

//             .retro-container {
//               background: #222;
//               border: 2px solid #00ff41;
//               padding: 20px;
//               font-family: 'Courier New', monospace;
//               position: relative;
//             }

//             .retro-container::before {
//               content: '';
//               position: absolute;
//               top: 0;
//               left: 0;
//               right: 0;
//               bottom: 0;
//               background: repeating-linear-gradient(
//                 0deg,
//                 transparent,
//                 transparent 2px,
//                 rgba(0, 255, 65, 0.03) 2px,
//                 rgba(0, 255, 65, 0.03) 4px
//               );
//               pointer-events: none;
//             }

//             .retro-field {
//               background: #000 !important;
//               border: 1px solid #00ff41 !important;
//               color: #00ff41 !important;
//               font-family: 'Courier New', monospace !important;
//               font-size: 16px !important;
//               padding: 10px !important;
//               border-radius: 0 !important;
//               box-shadow: 0 0 10px rgba(0, 255, 65, 0.3) !important;
//             }

//             .retro-field:focus {
//               box-shadow: 0 0 20px rgba(0, 255, 65, 0.6) !important;
//               text-shadow: 0 0 5px #00ff41 !important;
//             }

//             .retro-label {
//               color: #00ff41 !important;
//               font-family: 'Courier New', monospace !important;
//               text-transform: uppercase !important;
//               font-size: 12px !important;
//               letter-spacing: 2px !important;
//             }

//             .custom-icons {
//               color: #ff6b6b !important;
//               font-size: 20px !important;
//             }

//             .custom-error {
//               background: #ffe6e6 !important;
//               color: #d32f2f !important;
//               padding: 8px 12px !important;
//               border-radius: 8px !important;
//               border-left: 4px solid #d32f2f !important;
//               font-weight: 500 !important;
//             }

//             .custom-helper {
//               background: #e3f2fd !important;
//               color: #1976d2 !important;
//               padding: 6px 10px !important;
//               border-radius: 6px !important;
//               font-size: 13px !important;
//             }
//           `}</style>

//           <div style={{ display: 'grid', gap: '30px' }}>
//             <Input
//               label="Neumorphism Style"
//               placeholder="Modern neumorphic design"
//               containerClassName="modern-input-container"
//               inputClassName="modern-input-field"
//               labelClassName="modern-input-label"
//               value={formData.email}
//               onChange={handleInputChange('email')}
//             />

//             <Input
//               label="Glassmorphism Effect"
//               floatingLabel
//               containerClassName="glassmorphism-container"
//               inputClassName="glassmorphism-field"
//               labelClassName="glassmorphism-label"
//               value={formData.search}
//               onChange={handleInputChange('search')}
//               helperText="Frosted glass effect"
//               helperTextClassName="custom-helper"
//             />

//             <Input
//               label="Retro Terminal"
//               placeholder="Enter command..."
//               containerClassName="retro-container"
//               inputClassName="retro-field"
//               labelClassName="retro-label"
//               startIcon=">"
//               startIconClassName="custom-icons"
//               value={formData.phone}
//               onChange={handleInputChange('phone')}
//               error={formData.phone.length > 0 && formData.phone.length < 3 ? "Command too short" : ""}
//               errorClassName="custom-error"
//             />
//           </div>
//         </section>
//         <section>
//           <h2 style={{ marginBottom: '16px', fontSize: '20px', fontWeight: '600', color: '#374151' }}>
//             Custom Styling Example
//           </h2>
//           <style>{`
//             .custom-input {
//               background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
//               color: white;
//               border: none !important;
//               border-radius: 25px !important;
//             }
//             .custom-input::placeholder {
//               color: rgba(255, 255, 255, 0.8);
//             }
//             .custom-input:focus {
//               box-shadow: 0 0 20px rgba(102, 126, 234, 0.4) !important;
//             }
//             .custom-label {
//               color: #667eea !important;
//               font-weight: bold !important;
//             }
//             .neon-input {
//               background: #000 !important;
//               border: 2px solid #00ff88 !important;
//               color: #00ff88 !important;
//               border-radius: 0 !important;
//               box-shadow: 0 0 10px rgba(0, 255, 136, 0.3);
//             }
//             .neon-input:focus {
//               box-shadow: 0 0 20px rgba(0, 255, 136, 0.6) !important;
//               border-color: #00ffff !important;
//             }
//             .floating-custom {
//               background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
//               color: white;
//               border: 2px solid transparent !important;
//               border-radius: 12px !important;
//             }
//             .floating-custom:focus {
//               border-color: #ffffff !important;
//               box-shadow: 0 0 20px rgba(102, 126, 234, 0.4) !important;
//             }
//             .floating-label-custom {
//               color: #667eea !important;
//               font-weight: 600 !important;
//             }
//             .floating-neon {
//               background: #000 !important;
//               border: 2px solid #00ff88 !important;
//               color: #00ff88 !important;
//               border-radius: 0 !important;
//               box-shadow: 0 0 10px rgba(0, 255, 136, 0.3);
//             }
//             .floating-neon:focus {
//               box-shadow: 0 0 20px rgba(0, 255, 136, 0.6) !important;
//               border-color: #00ffff !important;
//             }
//           `}</style>

//           <div style={{ display: 'grid', gap: '20px' }}>
//             <Input
//               label="Gradient Style"
//               placeholder="Custom gradient background"
//               inputClassName="custom-input"
//               labelClassName="custom-label"
//             />

//             <Input
//               label="Gradient Floating"
//               floatingLabel
//               inputClassName="floating-custom"
//               labelClassName="floating-label-custom"
//             />

//             <Input
//               label="Neon Style"
//               placeholder="Cyberpunk aesthetic"
//               inputClassName="neon-input"
//               labelClassName="custom-label"
//             />

//             <Input
//               label="Neon Floating"
//               floatingLabel
//               inputClassName="floating-neon"
//               labelClassName="floating-label-custom"
//             />
//           </div>
//         </section>
//       </div>
//     </div>
//   );
// };

// export default InputDemo;

{
  /* <Input 
  label="Custom Field"
  containerClassName="my-container"
  inputClassName="my-input-field"
  labelClassName="my-label"
  iconClassName="my-icons"
  errorClassName="my-errors"
/> */
}

export default Input;
