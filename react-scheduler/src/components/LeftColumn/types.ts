import { PaginatedSchedulerData, SchedulerItemClickData } from "@/types/global";

export type LeftColumnProps = {
  data: PaginatedSchedulerData;
  rows: number[];
  pageNum: number;
  pagesAmount: number;
  onLoadNext: () => void;
  onLoadPrevious: () => void;
  searchInputValue: string;
  searchPersonInputValue: string;
  onSearchInputChange: React.ChangeEventHandler<HTMLInputElement>;
  onPersonSearchInputChange: React.ChangeEventHandler<HTMLInputElement>;
  onItemClick?: (data: SchedulerItemClickData) => void;
};

export type StyledInputWrapperProps = {
  isFocused: boolean;
};
