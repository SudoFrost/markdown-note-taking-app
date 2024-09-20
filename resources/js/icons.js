import { createIcons } from 'lucide';
import {
  NotebookPen,
  Trash2,
  Pencil,
  Check,
  X
} from 'lucide';

export function loadIcons() {
  createIcons({
    icons: {
      NotebookPen,
      Trash2,
      Pencil,
      Check,
      X,
    }
  });
}
