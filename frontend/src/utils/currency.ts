export const formatIDR = (amount: number): string => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
};

export const formatIDRCompact = (amount: number): string => {
  if (amount >= 1000000000) {
    return `Rp${(amount / 1000000000).toFixed(1)}M`;
  } else if (amount >= 1000000) {
    return `Rp${(amount / 1000000).toFixed(1)}Jt`;
  } else if (amount >= 1000) {
    return `Rp${(amount / 1000).toFixed(1)}Rb`;
  }
  return formatIDR(amount);
};
