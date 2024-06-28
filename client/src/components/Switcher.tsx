'use client'

import { useState, useEffect } from 'react';
import { useTheme } from 'next-themes'

const ThemeChanger = () => {
  const [labelC, setLabelC] = useState('bg-stroke');
  const { theme, setTheme } = useTheme()

  useEffect(() => {
    if (theme === 'dark') {
      setLabelC('bg-primary');
    } else {
      setLabelC('bg-stroke');
    }
  }, [theme])

  return (
    <div>
      Theme: {labelC}
      <button onClick={() => setTheme('light')}>Light Mode</button>
      <button onClick={() => setTheme('dark')}>Dark Mode</button>
    </div>
  )
}

export default ThemeChanger;
