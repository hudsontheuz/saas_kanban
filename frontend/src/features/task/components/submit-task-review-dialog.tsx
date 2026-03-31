import { useState } from 'react';
import { Send } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';

export function SubmitTaskReviewDialog({ onSubmit }: { onSubmit: (comment: string) => void }) {
  const [comment, setComment] = useState('');

  const handleSubmit = () => {
    const value = comment.trim();
    if (!value) return;
    onSubmit(value);
    setComment('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm">
          <Send className="mr-2 h-4 w-4" /> Enviar para revisão
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Resumo da entrega</DialogTitle>
          <DialogDescription>Antes de enviar para revisão, descreva brevemente o que foi feito.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="task-delivery-comment">O que foi feito</Label>
            <textarea
              id="task-delivery-comment"
              value={comment}
              onChange={(event) => setComment(event.target.value)}
              placeholder="Ex.: implementei a validação do formulário, ajustei o fluxo e revisei as mensagens de erro."
              className="min-h-28 w-full rounded-md border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-slate-400"
            />
          </div>
          <Button className="w-full" onClick={handleSubmit} disabled={!comment.trim()}>
            Confirmar envio
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
