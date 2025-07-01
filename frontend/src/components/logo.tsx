export function Logo() {
  return (
    <div className="flex items-center gap-2">
      <svg
        width="32"
        height="32"
        viewBox="0 0 42 42"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="h-8 w-8"
      >
        <path
          d="M34.75 5.25H7.25C6.14543 5.25 5.25 6.14543 5.25 7.25V26.25C5.25 27.3546 6.14543 28.25 7.25 28.25H28L36.75 36.75V7.25C36.75 6.14543 35.8546 5.25 34.75 5.25Z"
          fill="#5D9CEC"
        />
        <circle cx="16" cy="15" r="5" fill="#4FC1E9" />
        <path
          d="M11 25C11 22.2386 13.2386 20 16 20C18.7614 20 21 22.2386 21 25H11Z"
          fill="#4FC1E9"
        />
        <g opacity="0.8">
          <circle cx="26" cy="15" r="5" fill="#AC92EC" />
          <path
            d="M21 25C21 22.2386 23.2386 20 26 20C28.7614 20 31 22.2386 31 25H21Z"
            fill="#AC92EC"
          />
        </g>
      </svg>
      <h1 className="text-2xl font-bold text-primary font-headline">
        ConnectU
      </h1>
    </div>
  );
}
